package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	envelope2 "github.com/sakemi-hiroshi/my-playground/internal/envelope"

	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 3 {
		// このサンプルは、実行時にコマンドライン引数を受け取ることを想定する
		// 集約IDとイベントの数を引数として受け取る
		fmt.Fprintln(os.Stderr, "Usage: kinesis_sample producer <aggregate-id> <count>")
		os.Exit(1)
	}
	aggregateID := os.Args[1]
	count, err := strconv.Atoi(os.Args[2])
	if err != nil || count < 1 {
		// カウントは正の整数でなければならない
		fmt.Fprintln(os.Stderr, "Count must be a positive integer")
		os.Exit(1)
	}

	ctx := context.Background()

	kinesisClient, err := newKinesisClient(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create Kinesis client: %v\n", err)
		os.Exit(1)
	}
	streamName := os.Getenv("KINESIS_STREAM_NAME")

	for i := range count {
		seqNr := uint64(i + 1)
		payload := map[string]any{
			"message": fmt.Sprintf("Event %d for aggregate %s", seqNr, aggregateID),
		}
		envelope, err := envelope2.NewEnvelope(aggregateID, seqNr, "SandboxEvent", payload)
		if err != nil {
			log.Fatalf("Failed to create envelope: %v", err)
		}
		data, err := json.Marshal(&envelope)
		if err != nil {
			log.Fatalf("Failed to marshal envelope: %v", err)
		}

		// Kinesis Streamにデータを格納する
		output, err := kinesisClient.PutRecord(ctx, &kinesis.PutRecordInput{
			StreamName:   aws.String(streamName),  // ストリーム名の指定が必要
			PartitionKey: aws.String(aggregateID), // パーティションキーは集約IDを使用して同一シャードに送る
			Data:         data,                    // データはJSONエンコードされたエンベロープ
		})
		if err != nil {
			log.Fatalf("Failed to put record: %v", err)
		}

		// 一応ログだけ書いておくか
		log.Printf("Put record to stream %s: PartitionKey=%s, SequenceNumber=%s\n", streamName, aggregateID, *output.SequenceNumber)

		time.Sleep(100 * time.Millisecond)
	}
}

func newKinesisClient(ctx context.Context) (*kinesis.Client, error) {
	endpoint := os.Getenv("KINESIS_ENDPOINT")
	region := os.Getenv("AWS_REGION")

	config, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(
			// 本来はここは環境変数とかから入れないといけないが、
			// サンプルなのでダミーを入れる
			credentials.NewStaticCredentialsProvider("dummy-access-key", "dummy-secret-key", ""),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS SDK config: %w", err)
	}

	return kinesis.NewFromConfig(config, func(o *kinesis.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	}), nil
}

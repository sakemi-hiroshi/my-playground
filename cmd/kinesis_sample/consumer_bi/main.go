package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sakemi-hiroshi/my-playground/internal/envelope"
)

func handler(_ context.Context, ev events.KinesisEvent) error {
	log.Printf("BI Consumer Received %d records\n", len(ev.Records))

	for _, record := range ev.Records {
		var env envelope.Envelope
		if err := json.Unmarshal(record.Kinesis.Data, &env); err != nil {
			log.Printf("Failed to unmarshal record: %v\n", err)
			continue
		}

		// 本来はここでprojectionをするためのロジックを動かしたりする
		log.Printf("BI Consumer Kinesis record: AggregateID=%s, SeqNr=%d, EventType=%s, Payload=%s\n",
			env.AggregateID, env.SeqNr, env.EventType, string(env.Payload))
	}

	return nil
}

func main() {
	// lambdaを動かす（お作法）
	lambda.Start(handler)
}

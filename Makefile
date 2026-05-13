.DEFAULT_GOAL := help

.PHONY: help up down create-stream build-lambda deploy-lambda redeploy-lambda build-lambda-bi deploy-lambda-bi redeploy-lambda-bi put put/multi logs logs/bi clean

STREAM_NAME ?= playground
SHARD_COUNT ?= 2
AID         ?= a-1
N           ?= 3
LAMBDA_NAME ?= kinesis-consumer
LAMBDA_ZIP  ?= /tmp/kinesis-lambda.zip

AWS_CMD = aws --endpoint-url=http://localhost:4566 --region ap-northeast-1

## ヘルプ
help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "LocalStack:"
	@echo "  up                  LocalStack を起動"
	@echo "  down                LocalStack を停止"
	@echo ""
	@echo "Kinesis Stream:"
	@echo "  create-stream       Stream を作成 (SHARD_COUNT=2)"
	@echo ""
	@echo "Producer:"
	@echo "  put AID=<id> N=<n>  指定 aggregateId を N 件 PutRecord"
	@echo "  put/multi           a-1, b-1, c-1 を各 3 件投げる（Shard 観察用）"
	@echo ""
	@echo "Lambda (consumer_projection):"
	@echo "  deploy-lambda       ビルド + Lambda 登録 + ESM 作成"
	@echo "  redeploy-lambda     コード変更後の再デプロイ"
	@echo "  logs                Lambda ログを tail"
	@echo ""
	@echo "Lambda BI (consumer_bi / Fan-out):"
	@echo "  deploy-lambda-bi    ビルド + Lambda 登録 + ESM 作成"
	@echo "  redeploy-lambda-bi  コード変更後の再デプロイ"
	@echo "  logs/bi             Lambda BI ログを tail"
	@echo ""
	@echo "Other:"
	@echo "  clean               中間ファイル削除"

## LocalStack
up:
	docker compose up -d

down:
	docker compose down

## Kinesis Stream
create-stream:
	$(AWS_CMD) kinesis create-stream \
		--stream-name $(STREAM_NAME) \
		--shard-count $(SHARD_COUNT)

## Producer
put:
	KINESIS_ENDPOINT=http://localhost:4566 \
	KINESIS_STREAM_NAME=$(STREAM_NAME) \
	AWS_REGION=ap-northeast-1 \
	go run ./cmd/kinesis_sample/producer $(AID) $(N)

# 複数の aggregateId を混ぜて投げる（Shard 振り分けを観察）
put/multi:
	KINESIS_ENDPOINT=http://localhost:4566 \
	KINESIS_STREAM_NAME=$(STREAM_NAME) \
	AWS_REGION=ap-northeast-1 \
	go run ./cmd/kinesis_sample/producer a-1 3
	KINESIS_ENDPOINT=http://localhost:4566 \
	KINESIS_STREAM_NAME=$(STREAM_NAME) \
	AWS_REGION=ap-northeast-1 \
	go run ./cmd/kinesis_sample/producer b-1 3
	KINESIS_ENDPOINT=http://localhost:4566 \
	KINESIS_STREAM_NAME=$(STREAM_NAME) \
	AWS_REGION=ap-northeast-1 \
	go run ./cmd/kinesis_sample/producer c-1 3

## Lambda
build-lambda:
	GOOS=linux GOARCH=arm64 go build -o /tmp/bootstrap ./cmd/kinesis_sample/consumer_projection
	zip -j $(LAMBDA_ZIP) /tmp/bootstrap

deploy-lambda: build-lambda
	$(AWS_CMD) lambda create-function \
		--function-name $(LAMBDA_NAME) \
		--runtime provided.al2023 \
		--handler bootstrap \
		--role arn:aws:iam::000000000000:role/lambda-role \
		--zip-file fileb://$(LAMBDA_ZIP) \
		--architectures arm64
	$(AWS_CMD) lambda create-event-source-mapping \
		--function-name $(LAMBDA_NAME) \
		--event-source-arn $$($(AWS_CMD) kinesis describe-stream --stream-name $(STREAM_NAME) --query 'StreamDescription.StreamARN' --output text) \
		--batch-size 10 \
		--starting-position LATEST

# コードを変えて再デプロイしたいとき
redeploy-lambda: build-lambda
	$(AWS_CMD) lambda update-function-code \
		--function-name $(LAMBDA_NAME) \
		--zip-file fileb://$(LAMBDA_ZIP)

## Lambda BI（Fan-out 用 2つ目の Consumer）
LAMBDA_BI_NAME ?= kinesis-consumer-bi
LAMBDA_BI_ZIP  ?= /tmp/kinesis-lambda-bi.zip

build-lambda-bi:
	GOOS=linux GOARCH=arm64 go build -o /tmp/bootstrap ./cmd/kinesis_sample/consumer_bi
	zip -j $(LAMBDA_BI_ZIP) /tmp/bootstrap

deploy-lambda-bi: build-lambda-bi
	$(AWS_CMD) lambda create-function \
		--function-name $(LAMBDA_BI_NAME) \
		--runtime provided.al2023 \
		--handler bootstrap-bi \
		--role arn:aws:iam::000000000000:role/lambda-role \
		--zip-file fileb://$(LAMBDA_BI_ZIP) \
		--architectures arm64
	$(AWS_CMD) lambda create-event-source-mapping \
		--function-name $(LAMBDA_BI_NAME) \
		--event-source-arn $$($(AWS_CMD) kinesis describe-stream --stream-name $(STREAM_NAME) --query 'StreamDescription.StreamARN' --output text) \
		--batch-size 10 \
		--starting-position LATEST

redeploy-lambda-bi: build-lambda-bi
	$(AWS_CMD) lambda update-function-code \
		--function-name $(LAMBDA_BI_NAME) \
		--zip-file fileb://$(LAMBDA_BI_ZIP)

## ログ確認
logs:
	$(AWS_CMD) logs tail /aws/lambda/$(LAMBDA_NAME) --follow

logs/bi:
	$(AWS_CMD) logs tail /aws/lambda/$(LAMBDA_BI_NAME) --follow

clean:
	rm -f /tmp/bootstrap $(LAMBDA_ZIP)

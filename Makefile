PROJECT = $(shell basename $(CURDIR))
REVISION ?= $(shell git rev-parse --short HEAD)
BRANCH ?= $(shell git branch --no-color |sort |tail -1 |cut -c 3-)
CF_TEMPLATE ?= deploy/sam.yaml
PACKAGE_TEMPLATE = deploy/package.yaml
BUCKET ?= unspecified
STACK_NAME ?= $(PROJECT)
SLACK_TOKEN ?= ChangeMe

.PHONY: clean test build package deploy slack

default: test

clean:
	-rm -rf build/*
	-rm deploy/package.yaml

test:
	go test -race -v ./...
	# sam local invoke -e testdata/events/apigateway-shoutout.json -t deploy/sam.yaml ShoutoutHandlerFunction

params:
	aws cloudformation deploy \
		--stack-name $(STACK_NAME)-params \
		--template-file ./deploy/params.yaml \
		--parameter-overrides \
			slackToken=$(SLACK_TOKEN) \
			metricNamespace=$(METRIC_NAMESPACE) \
		--output json

deleteParams:
	aws cloudformation delete-stack \
		--stack-name $(STACK_NAME)-params

bucket:
	aws s3 mb s3://$(BUCKET)

build:
	GOOS=linux GOARCH=amd64 go build -v -o ./build/handler ./cmd/handler
	GOOS=linux GOARCH=amd64 go build -v -o ./build/processor ./cmd/processor
	GOOS=linux GOARCH=amd64 go build -v -o ./build/saver ./cmd/saver
	GOOS=linux GOARCH=amd64 go build -v -o ./build/metrics ./cmd/metrics
	GOOS=linux GOARCH=amd64 go build -v -o ./build/callback ./cmd/callback

zip:
	@cd ./build && zip handler.zip handler
	@cd ./build && zip processor.zip processor
	@cd ./build && zip saver.zip saver
	@cd ./build && zip metrics.zip metrics
	@cd ./build && zip callback.zip callback

package: test build zip
	sam validate --template $(CF_TEMPLATE)
	sam package \
		--debug \
		--template-file $(CF_TEMPLATE) \
		--output-template-file $(PACKAGE_TEMPLATE) \
		--s3-bucket $(BUCKET)

deploy: clean package
	sam deploy \
		--template-file $(PACKAGE_TEMPLATE) \
		--stack-name $(STACK_NAME) \
		--capabilities CAPABILITY_IAM \
		--no-fail-on-empty-changeset

destroy:
	aws cloudformation delete-stack \
		--stack-name $(STACK_NAME)

outputs:
	aws cloudformation describe-stacks \
		--stack-name $(STACK_NAME) \
		--query 'Stacks[].Outputs' \
		--output json

describe:
	aws cloudformation describe-stacks \
		--stack-name $(STACK_NAME) \
		--output json

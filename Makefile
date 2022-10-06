PROJECT = $(shell basename $(CURDIR))
STACK_NAME ?= $(PROJECT)
STACK_NAME_PARAMS ?= $(STACK_NAME)-params
BUCKET ?= $(STACK_NAME)
METRIC_NAMESPACE ?= $(STACK_NAME)
SLACK_TOKEN ?= ChangeMe
CF_TEMPLATE ?= deploy/sam.yaml
PACKAGE_TEMPLATE = deploy/package.yaml

.PHONY: clean test build package deploy slack

default: test

clean:
	-rm -rf build/*
	-rm deploy/package.yaml

test:
	go test -race -v ./...

params:
	aws cloudformation deploy \
		--stack-name $(STACK_NAME_PARAMS) \
		--template-file ./deploy/params.yaml \
		--parameter-overrides \
			slackToken=$(SLACK_TOKEN) \
			metricNamespace=$(METRIC_NAMESPACE) \
		--output json

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

validate:
	sam validate --template $(CF_TEMPLATE)

package: test build validate zip
	sam package \
		--template-file $(CF_TEMPLATE) \
		--output-template-file $(PACKAGE_TEMPLATE) \
		--s3-bucket $(BUCKET)

deploy: clean package
	sam deploy \
		--template-file $(PACKAGE_TEMPLATE) \
		--stack-name $(STACK_NAME) \
		--parameter-overrides \
			paramsStackName=$(STACK_NAME_PARAMS) \
		--capabilities CAPABILITY_IAM \
		--no-fail-on-empty-changeset

destroy:
	aws cloudformation delete-stack \
		--stack-name $(STACK_NAME)
	aws cloudformation wait stack-delete-complete \
		--stack-name $(STACK_NAME)
	aws cloudformation delete-stack \
		--stack-name $(STACK_NAME_PARAMS)
	aws s3 rb s3://$(BUCKET) --force  

outputs:
	aws cloudformation describe-stacks \
		--stack-name $(STACK_NAME) \
		--query 'Stacks[].Outputs'

describe:
	aws cloudformation describe-stacks \
		--stack-name $(STACK_NAME) \
		--output json

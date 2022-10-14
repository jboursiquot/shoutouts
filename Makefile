PROJECT = $(shell basename $(CURDIR))
REVISION ?= $(shell git rev-parse --short HEAD)
STACK_NAME ?= $(PROJECT)
STACK_NAME_PARAMS ?= $(STACK_NAME)-params
BUCKET ?= $(STACK_NAME)
METRIC_NAMESPACE ?= $(STACK_NAME)
SLACK_TOKEN ?= ChangeMe
CF_TEMPLATE ?= deploy/sam.yaml
PACKAGE_TEMPLATE = deploy/package.yaml
AWS_ACCOUNT_ID ?= $(shell aws sts get-caller-identity --query Account --output text)
AWS_REGION ?= $(shell aws configure get region)
ECR_REGISTRY ?= $(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com

.PHONY: clean test build package deploy slack
.PHONY: ecr-login build-and-push-sanitizer-image build-and-package-sanitizing-handler
.PHONY: deploy-services-base destroy-services-base deploy-sanitizer destroy-sanitizer

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
	GOOS=linux GOARCH=amd64 go build -v -o ./build/sanitizing-handler ./cmd/sanitizing-handler

zip:
	@cd ./build && zip handler.zip handler
	@cd ./build && zip processor.zip processor
	@cd ./build && zip saver.zip saver
	@cd ./build && zip metrics.zip metrics
	@cd ./build && zip callback.zip callback
	@cd ./build && zip sanitizing-handler.zip sanitizing-handler

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

deploy-services-base:
	aws cloudformation deploy \
		--stack-name $(STACK_NAME)-services \
		--template-file ./deploy/services-base.yaml \
		--capabilities CAPABILITY_IAM \
		--output json

destroy-services-base:
	aws cloudformation delete-stack \
		--stack-name $(STACK_NAME)-services
	aws cloudformation wait stack-delete-complete \
		--stack-name $(STACK_NAME)-services

ecr-login:
	aws ecr get-login-password --region $(AWS_REGION) | docker login --username AWS --password-stdin $(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com

ECR_REPO_SANITIZER ?= $(PROJECT)/sanitizer
build-and-push-sanitizer-image: ecr-login
	GOOS=linux GOARCH=amd64 go build -v -o ./build/sanitizer ./cmd/sanitizer
	docker build -t $(ECR_REPO_SANITIZER):$(REVISION) -f Dockerfile-sanitizer .
	docker tag $(ECR_REPO_SANITIZER):$(REVISION) $(ECR_REGISTRY)/$(ECR_REPO_SANITIZER):$(REVISION)
	docker push $(ECR_REGISTRY)/$(ECR_REPO_SANITIZER):$(REVISION)

SANITIZER_CF_TEMPLATE ?= deploy/services-sanitizer.yaml 
SANITIZER_PACKAGE_TEMPLATE = deploy/services-sanitizer-package.yaml
build-and-package-sanitizing-handler:
	GOOS=linux GOARCH=amd64 go build -v -o ./build/sanitizing-handler ./cmd/sanitizing-handler
	@cd ./build && zip sanitizing-handler.zip sanitizing-handler
	sam package \
		--template-file $(SANITIZER_CF_TEMPLATE) \
		--output-template-file $(SANITIZER_PACKAGE_TEMPLATE) \
		--s3-bucket $(BUCKET)

deploy-sanitizer: build-and-push-sanitizer-image build-and-package-sanitizing-handler
	sam deploy \
		--template-file $(SANITIZER_PACKAGE_TEMPLATE)\
		--stack-name $(STACK_NAME)-sanitizer \
		--parameter-overrides \
			paramsStackName=$(STACK_NAME_PARAMS) \
			imageAndTag=$(ECR_REPO_SANITIZER):$(REVISION)

destroy-sanitizer:
	aws cloudformation delete-stack \
		--stack-name $(STACK_NAME)-sanitizer
	aws cloudformation wait stack-delete-complete \
		--stack-name $(STACK_NAME)-sanitizer
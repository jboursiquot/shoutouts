package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/jboursiquot/shoutouts"
	"github.com/sirupsen/logrus"
)

var sess *session.Session
var esqs *sqs.SQS
var ddb *dynamodb.DynamoDB
var logger *logrus.Logger

const sanitizerServiceEndpoint = "http://sanitizer.internal:8080/sanitize"

func init() {
	sess = session.Must(session.NewSession())
	esqs = sqs.New(sess)
	ddb = dynamodb.New(sess)
	logger = logrus.New()
}

func handler(req *events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	ctx := context.Background()
	return shoutouts.NewSanitizingHandler(esqs, ddb, logger, sanitizerServiceEndpoint).Handle(ctx, req)
}

func main() {
	lambda.Start(handler)
}

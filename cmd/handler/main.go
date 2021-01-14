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

func init() {
	sess = session.Must(session.NewSession())
	esqs = sqs.New(sess)
	ddb = dynamodb.New(sess)
	logger = logrus.New()
}

func handler(req *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx := context.Background()
	return shoutouts.NewHandler(esqs, ddb, logger).Handle(ctx, req)
}

func main() {
	lambda.Start(handler)
}

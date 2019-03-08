package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/jboursiquot/shoutouts/pkg/shoutouts"
)

var sess *session.Session
var esqs *sqs.SQS

func init() {
	sess = session.Must(session.NewSession())
	esqs = sqs.New(sess)
}

func handler(req *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx := context.Background()
	return shoutouts.NewHandler(esqs).Handle(ctx, req)
}

func main() {
	lambda.Start(handler)
}

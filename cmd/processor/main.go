package main

import (
	"context"

	"github.com/jboursiquot/shoutouts"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sfn"
)

var sess *session.Session
var sf *sfn.SFN

func init() {
	sess = session.Must(session.NewSession())
	sf = sfn.New(sess)
}

func handler(ctx context.Context, sqsEvent *events.SQSEvent) error {
	return shoutouts.NewSQSProcessor(sf).Process(ctx, sqsEvent)
}

func main() {
	lambda.Start(handler)
}

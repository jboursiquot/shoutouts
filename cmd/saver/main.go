package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/jboursiquot/shoutouts"
)

var sess *session.Session
var ddb *dynamodb.DynamoDB

func init() {
	sess = session.Must(session.NewSession())
	ddb = dynamodb.New(sess)
}

func handler(ctx context.Context, shoutout *shoutouts.Shoutout) (*shoutouts.Shoutout, error) {
	err := shoutouts.NewSaver(ddb).Save(ctx, shoutout)
	if err != nil {
		xray.AddError(ctx, err)
	}
	return shoutout, err
}

func main() {
	lambda.Start(handler)
}

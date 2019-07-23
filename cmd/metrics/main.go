package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/jboursiquot/shoutouts"
)

var sess *session.Session
var cw *cloudwatch.CloudWatch

func init() {
	sess = session.Must(session.NewSession())
	cw = cloudwatch.New(sess)
}

func handler(ctx context.Context, shoutout *shoutouts.Shoutout) (*shoutouts.Shoutout, error) {
	err := shoutouts.NewMetrics(cw).Capture(ctx, shoutout)
	if err != nil {
		xray.AddError(ctx, err)
	}
	return shoutout, err
}

func main() {
	lambda.Start(handler)
}

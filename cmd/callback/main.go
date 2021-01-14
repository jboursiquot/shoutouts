package main

import (
	"context"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/jboursiquot/shoutouts"
)

// var sess *session.Session
// var cw *cloudwatch.CloudWatch

func init() {
	// sess = session.Must(session.NewSession())
	// cw = cloudwatch.New(sess)
}

func handler(ctx context.Context, shoutout *shoutouts.Shoutout) (*shoutouts.Shoutout, error) {
	client := &http.Client{Timeout: time.Second * 10}
	err := shoutouts.NewCallback(client).Call(ctx, shoutout)
	if err != nil {
		xray.AddError(ctx, err)
	}
	return shoutout, err
}

func main() {
	lambda.Start(handler)
}

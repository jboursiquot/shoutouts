package shoutouts_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws/request"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"

	"github.com/aws/aws-lambda-go/events"

	"github.com/jboursiquot/shoutouts"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	cases := []struct {
		scenario string
		request  *events.APIGatewayProxyRequest
		sqs      *mockSQS
		ddb      *mockDynamoDB
	}{
		{
			scenario: "unspecified command",
			request: &events.APIGatewayProxyRequest{
				Body:       baseCommandParams("").Encode(),
				HTTPMethod: http.MethodPost,
			},
			sqs: &mockSQS{},
			ddb: &mockDynamoDB{},
		},
		{
			scenario: "unknown command",
			request: &events.APIGatewayProxyRequest{
				Body:       baseCommandParams("stuff").Encode(),
				HTTPMethod: http.MethodPost,
			},
			sqs: &mockSQS{},
			ddb: &mockDynamoDB{},
		},
		{
			scenario: "help command",
			request: &events.APIGatewayProxyRequest{
				Body:       baseCommandParams("help").Encode(),
				HTTPMethod: http.MethodPost,
			},
			sqs: &mockSQS{},
			ddb: &mockDynamoDB{},
		},
		{
			scenario: "help command (get)",
			request: &events.APIGatewayProxyRequest{
				Body:                  baseCommandParams("help").Encode(),
				QueryStringParameters: map[string]string{"token": os.Getenv("SLACK_TOKEN")},
				HTTPMethod:            http.MethodGet,
			},
			sqs: &mockSQS{},
			ddb: &mockDynamoDB{},
		},
		{
			scenario: "help usage command",
			request: &events.APIGatewayProxyRequest{
				Body:       baseCommandParams("help usage").Encode(),
				HTTPMethod: http.MethodPost,
			},
			sqs: &mockSQS{},
			ddb: &mockDynamoDB{},
		},
		{
			scenario: "help values command",
			request: &events.APIGatewayProxyRequest{
				Body:       baseCommandParams("help values").Encode(),
				HTTPMethod: http.MethodPost,
			},
			sqs: &mockSQS{},
			ddb: &mockDynamoDB{},
		},
		{
			scenario: "shoutout command",
			request: &events.APIGatewayProxyRequest{
				Body:       baseCommandParams("<@userid|johnny> it that's good thinkin'").Encode(),
				HTTPMethod: http.MethodPost,
			},
			sqs: &mockSQS{},
			ddb: &mockDynamoDB{},
		},
		{
			scenario: "shoutout command with sqs error",
			request: &events.APIGatewayProxyRequest{
				Body:       baseCommandParams("<@userid|johnny_boursiquot> tf all for the team").Encode(),
				HTTPMethod: http.MethodPost,
			},
			sqs: &mockSQS{err: errors.New("sqs api call failure")},
			ddb: &mockDynamoDB{},
		},
	}

	for _, c := range cases {
		t.Run(c.scenario, func(t *testing.T) {
			h := shoutouts.NewHandler(c.sqs, c.ddb, nullLogger())
			r, err := h.Handle(context.Background(), c.request)
			assert.NoError(t, err)
			if c.sqs.err == nil {
				assert.Equal(t, r.StatusCode, 200)
			} else {
				assert.Equal(t, r.StatusCode, 500)
			}
		})
	}
}

// baseCommandParams represents the parameters that will be sent by Slack
// to our application. Our job is to parse these params to figure out
// the user's intent. This test helper let's us try a variety of possible
// commands that our application may be asked to handle. Just pass in a
// textOverride to simulate payloads for different kinds of things a user
// might type in.
func baseCommandParams(textOverride string) *url.Values {
	p := url.Values{}
	p.Add("token", os.Getenv("SLACK_TOKEN"))
	p.Add("team_id", "T0001")
	p.Add("team_domain", "example")
	p.Add("channel_id", "C2147483705")
	p.Add("channel_name", "test")
	p.Add("user_id", "<@U030XRXJ2|jboursiquot>")
	p.Add("user_name", "johnny")
	p.Add("command", "/shoutout")
	p.Add("response_url", "https://hooks.slack.com/commands/1234/5678")
	p.Add("text", textOverride)
	return &p
}

type mockSQS struct {
	err error
}

func (m *mockSQS) SendMessageWithContext(ctx aws.Context, input *sqs.SendMessageInput, opts ...request.Option) (*sqs.SendMessageOutput, error) {
	return nil, m.err
}

package shoutouts_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jboursiquot/shoutouts"
	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-lambda-go/events"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/sfn"
)

func TestProcessor(t *testing.T) {
	cases := []struct {
		scenario string
		event    *events.SQSEvent
		sf       *mockSFN
	}{
		{
			scenario: "api error",
			event:    sampleSQSEvent(),
			sf: &mockSFN{
				nil,
				errors.New("api error"),
			},
		},
		{
			scenario: "happy path",
			event:    sampleSQSEvent(),
			sf: &mockSFN{
				&sfn.StartExecutionOutput{
					ExecutionArn: aws.String(""),
				},
				nil,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.scenario, func(t *testing.T) {
			p := shoutouts.NewSQSProcessor(c.sf)

			if c.sf.err != nil {
				assert.Error(t, p.Process(context.Background(), c.event))
			}

			if c.sf.err == nil {
				assert.NoError(t, p.Process(context.Background(), c.event))
			}
		})
	}
}

type mockSFN struct {
	out *sfn.StartExecutionOutput
	err error
}

func (m *mockSFN) StartExecutionWithContext(ctx aws.Context, input *sfn.StartExecutionInput, opts ...request.Option) (*sfn.StartExecutionOutput, error) {
	return m.out, m.err
}

func sampleSQSEvent() *events.SQSEvent {
	e := &events.SQSEvent{
		Records: []events.SQSMessage{
			{
				MessageId: "",
				Body:      "",
			},
			{
				MessageId: "",
				Body:      "",
			},
		},
	}
	return e
}

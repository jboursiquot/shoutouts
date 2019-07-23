package shoutouts

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/sfn"
)

// SFNAPI is the minimal interface needed to trigger a Step Functions State Machine.
type SFNAPI interface {
	StartExecutionWithContext(aws.Context, *sfn.StartExecutionInput, ...request.Option) (*sfn.StartExecutionOutput, error)
}

// NewSQSProcessor returns a new SQSProcessor.
func NewSQSProcessor(c SFNAPI) *SQSProcessor {
	return &SQSProcessor{sfn: c}
}

// SQSProcessor processes messages from queue.
type SQSProcessor struct {
	sfn SFNAPI
}

// Process processes individual messages.
func (p *SQSProcessor) Process(ctx context.Context, event *events.SQSEvent) error {
	for _, message := range event.Records {
		log.Printf("Processing message %s | %s", message.MessageId, message.Body)

		in := sfn.StartExecutionInput{
			Input:           aws.String(message.Body),
			StateMachineArn: aws.String(os.Getenv("STATE_MACHINE_ARN")),
		}

		out, err := p.sfn.StartExecutionWithContext(ctx, &in)
		if err != nil {
			return fmt.Errorf("Failed to start state machine execution: %s", err)
		}

		log.Printf("Started State Machine execution | ARN: %s", *out.ExecutionArn)
	}

	return nil
}

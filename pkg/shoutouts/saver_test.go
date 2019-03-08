package shoutouts_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws/request"

	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/jboursiquot/shoutouts/pkg/shoutouts"
)

func TestSaver(t *testing.T) {
	cases := []struct {
		scenario string
		shoutout *shoutouts.Shoutout
		ddb      *mockDynamoDB
	}{
		{
			scenario: "api error",
			shoutout: shoutouts.New(),
			ddb: &mockDynamoDB{
				nil,
				errors.New("api error"),
			},
		},
		{
			scenario: "happy path",
			shoutout: shoutouts.New(),
			ddb: &mockDynamoDB{
				&dynamodb.PutItemOutput{},
				nil,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.scenario, func(t *testing.T) {
			s := shoutouts.NewDynamoDBSaver(c.ddb)

			if c.ddb.err != nil {
				assert.Error(t, s.Save(context.Background(), c.shoutout))
			}

			if c.ddb.err == nil {
				assert.NoError(t, s.Save(context.Background(), c.shoutout))
			}
		})
	}
}

type mockDynamoDB struct {
	out *dynamodb.PutItemOutput
	err error
}

func (m *mockDynamoDB) PutItemWithContext(ctx aws.Context, item *dynamodb.PutItemInput, opts ...request.Option) (*dynamodb.PutItemOutput, error) {
	return m.out, m.err
}

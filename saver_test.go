package shoutouts_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws/request"

	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/jboursiquot/shoutouts"
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
				putOut: nil,
				err:    errors.New("api error"),
			},
		},
		{
			scenario: "successful put",
			shoutout: shoutouts.New(),
			ddb: &mockDynamoDB{
				putOut: &dynamodb.PutItemOutput{},
				err:    nil,
			},
		},
		{
			scenario: "successful query",
			shoutout: shoutouts.New(),
			ddb: &mockDynamoDB{
				queryOut: &dynamodb.QueryOutput{},
				err:      nil,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.scenario, func(t *testing.T) {
			s := shoutouts.NewSaver(c.ddb)

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
	putOut   *dynamodb.PutItemOutput
	queryOut *dynamodb.QueryOutput
	err      error
}

func (m *mockDynamoDB) PutItemWithContext(ctx aws.Context, item *dynamodb.PutItemInput, opts ...request.Option) (*dynamodb.PutItemOutput, error) {
	return m.putOut, m.err
}

func (m *mockDynamoDB) QueryWithContext(aws.Context, *dynamodb.QueryInput, ...request.Option) (*dynamodb.QueryOutput, error) {
	return m.queryOut, m.err
}

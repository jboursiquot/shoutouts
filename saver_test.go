package shoutouts_test

import (
	"github.com/aws/aws-sdk-go/aws/request"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

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

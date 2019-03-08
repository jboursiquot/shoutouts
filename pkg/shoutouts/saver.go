package shoutouts

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/request"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// DynamoDBAPI is the minimal interface needed to store a shoutout.
type DynamoDBAPI interface {
	PutItemWithContext(aws.Context, *dynamodb.PutItemInput, ...request.Option) (*dynamodb.PutItemOutput, error)
}

// NewDynamoDBSaver returns a new saver.
func NewDynamoDBSaver(c DynamoDBAPI) *DynamoDBSaver {
	return &DynamoDBSaver{ddb: c}
}

// DynamoDBSaver is a shoutout saver.
type DynamoDBSaver struct {
	ddb DynamoDBAPI
}

// Save saves a Shoutout
func (s *DynamoDBSaver) Save(ctx context.Context, shoutout *Shoutout) error {
	item, err := dynamodbattribute.MarshalMap(shoutout)
	if err != nil {
		return fmt.Errorf("failed to marshal shoutout for storage: %s", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(os.Getenv("TABLE_NAME")),
	}

	if _, err = s.ddb.PutItemWithContext(ctx, input); err != nil {
		return fmt.Errorf("failed to save shoutout: %s", err)
	}

	return nil
}

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

// DynamoDBQuerier is the minimal interface needed to query for shoutouts.
type DynamoDBQuerier interface {
	QueryWithContext(aws.Context, *dynamodb.QueryInput, ...request.Option) (*dynamodb.QueryOutput, error)
}

// NewLister returns a new lister.
func NewLister(c DynamoDBQuerier) *Lister {
	return &Lister{ddb: c}
}

// Lister is a shoutout lister.
type Lister struct {
	ddb DynamoDBQuerier
}

// List retrieves shoutouts for a given recipient.
func (s *Lister) List(ctx context.Context, recipientid string) ([]Shoutout, error) {
	input := &dynamodb.QueryInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		IndexName: aws.String("RecipientIDIndex"),
		KeyConditions: map[string]*dynamodb.Condition{
			"RecipientID": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(recipientid),
					},
				},
			},
		},
	}

	res, err := s.ddb.QueryWithContext(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve shoutouts: %s", err)
	}

	list := []Shoutout{}
	if err := dynamodbattribute.UnmarshalListOfMaps(res.Items, &list); err != nil {
		return nil, fmt.Errorf("failed to unmarshal shoutouts: %s", err)
	}

	return list, nil
}

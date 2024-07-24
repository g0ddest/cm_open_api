package dynamodb

import (
	"cm_open_api/internal/models"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func GetSource(region, tableName, key string) (*models.SourceResponse, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
	svc := dynamodb.New(sess)

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"mp": {
				S: aws.String(key),
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error calling GetItem: %v", err)
	}

	if result.Item == nil {
		return nil, fmt.Errorf("no item found with key: %s", key)
	}

	var sourceResponse models.SourceResponse
	err = dynamodbattribute.UnmarshalMap(result.Item, &sourceResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return &sourceResponse, nil
}

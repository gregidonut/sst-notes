package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/gregidonut/sst-notes/packages/functions/cmd/list/db"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var dynamoDbClient *dynamodb.Client

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	dynamoDbClient = dynamodb.NewFromConfig(cfg)
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	tableName := os.Getenv("NOTES_TABLE_NAME")
	params := &dynamodb.QueryInput{
		TableName:              aws.String(tableName), // Replace with your actual table name
		KeyConditionExpression: aws.String("userId = :userId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userId": &types.AttributeValueMemberS{Value: "123"},
		},
	}

	result, err := dynamoDbClient.Query(ctx, params)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: fmt.Sprintf("Error querying DynamoDB: %v", err)}, nil
	}
	var items []db.Note
	err = attributevalue.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: fmt.Sprintf("Error unmarshaling DynamoDB response: %v", err)}, nil
	}

	payload, err := json.Marshal(items)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error marshaling response"}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(payload)}, nil
}

func main() {
	lambda.Start(handler)
}

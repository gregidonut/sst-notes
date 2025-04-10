package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/gregidonut/sst-notes/packages/functions/cmd/list/db"
	"github.com/gregidonut/sst-notes/packages/functions/cmd/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userId, err := utils.GetUserId(request)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: fmt.Sprintf("Error getting userId: %v", err)}, nil
	}

	noteID := request.PathParameters["id"]
	if noteID == "" {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Missing note ID"}, nil
	}

	note := db.Note{
		UserId: userId,
		NoteId: noteID,
	}

	fmt.Printf("~noteID:%#v\n", note)

	item, err := attributevalue.MarshalMap(note)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error creating note"}, nil
	}

	params := &dynamodb.GetItemInput{
		TableName: aws.String(utils.DYNAMODB_TABLE_NAME),
		Key:       item,
	}

	dynamoDbClient, err := utils.InitDynamo()
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: fmt.Sprintf("Error initializing Dynanmo config: %v", err)}, nil
	}

	result, err := dynamoDbClient.GetItem(ctx, params)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: fmt.Sprintf("Error retrieving item: %v", err)}, nil
	}

	if result.Item == nil {
		return events.APIGatewayProxyResponse{StatusCode: 404, Body: "Item not found"}, nil
	}

	dynamoResItem, err := json.Marshal(result.Item)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error marshaling response"}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(dynamoResItem)}, nil
}

func main() {
	lambda.Start(handler)
}

package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gregidonut/sst-notes/packages/functions/cmd/list/db"
	"github.com/gregidonut/sst-notes/packages/functions/cmd/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	noteID := request.PathParameters["id"]
	if noteID == "" {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Missing note ID"}, nil
	}
	userId, err := utils.GetUserId(request)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: fmt.Sprintf("Error getting userId: %v", err)}, nil
	}

	note := db.Note{
		UserId: userId,
		NoteId: noteID,
	}

	fmt.Printf("~noteID:%#v\n", note)

	params := &dynamodb.DeleteItemInput{
		TableName: aws.String(utils.DYNAMODB_TABLE_NAME),
		Key: map[string]types.AttributeValue{
			"userId": &types.AttributeValueMemberS{Value: note.UserId},
			"noteId": &types.AttributeValueMemberS{Value: note.NoteId},
		},
	}

	dynamoDbClient, err := utils.InitDynamo()
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: fmt.Sprintf("Error initializing Dynanmo config: %v", err)}, nil
	}

	result, err := dynamoDbClient.DeleteItem(ctx, params)
	if err != nil {
		fmt.Printf("~error in delete action result:%#v\n", err.Error())
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: fmt.Sprintf("Error retrieving item: %v", err)}, nil
	}

	if result.Attributes == nil {
		return events.APIGatewayProxyResponse{StatusCode: 404, Body: "Item not found"}, nil
	}

	dynamoResItem, err := json.Marshal(result.Attributes)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error marshaling response"}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(dynamoResItem)}, nil
}

func main() {
	lambda.Start(handler)
}

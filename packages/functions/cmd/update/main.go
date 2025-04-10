package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gregidonut/sst-notes/packages/functions/cmd/list/db"
	"github.com/gregidonut/sst-notes/packages/functions/cmd/utils"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	noteID := request.PathParameters["id"]
	var requestBody db.Note
	if request.Body != "" {
		err := json.Unmarshal([]byte(request.Body), &requestBody)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid request body"}, nil
		}
	}
	userId, err := utils.GetUserId(request)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: fmt.Sprintf("Error getting userId: %v", err)}, nil
	}

	note := db.Note{
		UserId:     userId,
		NoteId:     noteID,
		Content:    requestBody.Content,
		Attachment: requestBody.Attachment,
	}

	item, err := attributevalue.MarshalMap(note)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error creating note"}, nil
	}

	updateExpression := "SET content = :content, attachment = :attachment"

	expressionValues := map[string]types.AttributeValue{
		":content":    item["content"],
		":attachment": item["attachment"],
	}

	key := map[string]types.AttributeValue{
		"userId": item["userId"],
		"noteId": item["noteId"],
	}

	params := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(utils.DYNAMODB_TABLE_NAME),
		Key:                       key,
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeValues: expressionValues,
	}

	dynamoDbClient, err := utils.InitDynamo()
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: fmt.Sprintf("Error initializing Dynanmo config: %v", err)}, nil
	}

	result, err := dynamoDbClient.UpdateItem(ctx, params)
	if err != nil {
		fmt.Printf("~error in update action result:%#v\n", err.Error())
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

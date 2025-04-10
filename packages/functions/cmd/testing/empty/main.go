package main

import (
	"context"
	"fmt"
	"os"

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
	tableName := os.Getenv("NOTES_TABLE_NAME")

	userId, err := utils.GetUserId(request)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: fmt.Sprintf("Error getting userId: %v", err)}, nil
	}
	params := &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("userId = :userId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userId": &types.AttributeValueMemberS{Value: userId},
		},
	}
	dynamoDbClient, err := utils.InitDynamo()
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: fmt.Sprintf("Error initializing Dynanmo config: %v", err)}, nil
	}

	result, err := dynamoDbClient.Query(ctx, params)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: fmt.Sprintf("Error querying DynamoDB: %v", err)}, nil
	}
	items := []db.Note{}
	err = attributevalue.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: fmt.Sprintf("Error unmarshaling DynamoDB response: %v", err)}, nil
	}

	for _, note := range items {
		fmt.Printf("~deleting %s\n", note.NoteId)
		params := &dynamodb.DeleteItemInput{
			TableName: aws.String(tableName),
			Key: map[string]types.AttributeValue{
				"userId": &types.AttributeValueMemberS{Value: note.UserId},
				"noteId": &types.AttributeValueMemberS{Value: note.NoteId},
			},
		}

		_, err := dynamoDbClient.DeleteItem(ctx, params)
		if err != nil {
			fmt.Printf("~error in delete action result:%#v\n", err.Error())
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: fmt.Sprintf("Error retrieving item: %v", err)}, nil
		}

		fmt.Printf("~finished deleting %s\n", note.NoteId)
	}

	return events.APIGatewayProxyResponse{StatusCode: 200, Body: "successfully emptied db"}, nil
}

func main() {
	lambda.Start(handler)
}

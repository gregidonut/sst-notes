package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gregidonut/sst-notes/packages/functions/cmd/list/db"
	"log"
	"os"
	"time"
)

var dynamoDbClient *dynamodb.Client

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	dynamoDbClient = dynamodb.NewFromConfig(cfg)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	tableName := os.Getenv("NOTES_TABLE_NAME")
	noteID := request.PathParameters["id"]
	var requestBody db.Note
	if request.Body != "" {
		err := json.Unmarshal([]byte(request.Body), &requestBody)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid request body"}, nil
		}
	}

	note := db.Note{
		UserId:     "123",
		NoteId:     noteID,
		Content:    requestBody.Content,
		Attachment: requestBody.Attachment,
		CreatedAt:  time.Now().String(),
	}

	item, err := attributevalue.MarshalMap(note)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error creating note"}, nil
	}

	params := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"userId": &types.AttributeValueMemberS{Value: note.UserId},
			"noteId": &types.AttributeValueMemberS{Value: note.NoteId}, // Assuming NoteId is a string
		},

		UpdateExpression:          aws.String("SET content = :content, attachment = :attachment"),
		ExpressionAttributeValues: item,
	}

	result, err := dynamoDbClient.UpdateItem(ctx, params)
	if err != nil {
		fmt.Printf("~error in update action result:%#v\n", err.Error())
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: fmt.Sprintf("Error retrieving item: %v", err)}, nil
	}

	fmt.Printf("~update result:%#v\n", result)
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

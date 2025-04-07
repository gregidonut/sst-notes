package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gregidonut/sst-notes/packages/functions/cmd/utils"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
	"github.com/gregidonut/sst-notes/packages/functions/cmd/list/db"
)

var (
	dynamoDbClient *dynamodb.Client
	tableName      = os.Getenv("NOTES_TABLE_NAME")
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}
	dynamoDbClient = dynamodb.NewFromConfig(cfg)
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var requestBody db.Note
	if event.Body != "" {
		err := json.Unmarshal([]byte(event.Body), &requestBody)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid request body"}, nil
		}
	}
	userId, err := utils.GetUserId(event)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: fmt.Sprintf("Error getting userId: %v", err)}, nil
	}

	note := db.Note{
		UserId:     userId,
		NoteId:     uuid.New().String(),
		Content:    requestBody.Content,
		Attachment: requestBody.Attachment,
		CreatedAt:  time.Now().Format(time.RFC3339Nano),
	}

	item, err := attributevalue.MarshalMap(note)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error creating note"}, nil
	}

	_, err = dynamoDbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error saving note to database"}, nil
	}

	fmt.Printf("~completed dynamoDbClient.PutItem: \n")

	responseBody, _ := json.Marshal(note)
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(responseBody)}, nil
}

func main() {
	lambda.Start(handler)
}

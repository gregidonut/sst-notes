package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/google/uuid"
	"github.com/gregidonut/sst-notes/packages/functions/cmd/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gregidonut/sst-notes/packages/functions/cmd/list/db"
)

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	notes, err := generateSeedNotes(event)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error creating the struct slice containing seed notes"}, nil
	}

	for i, note := range notes {
		item, err := attributevalue.MarshalMap(note)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error creating note"}, nil
		}

		dynamoDbClient, err := utils.InitDynamo()
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: fmt.Sprintf("Error initializing Dynanmo config: %v", err)}, nil
		}

		_, err = dynamoDbClient.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(utils.DYNAMODB_TABLE_NAME),
			Item:      item,
		})
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error saving note to database"}, nil
		}

		fmt.Printf("~completed dynamoDbClient.PutItem for note item: %d\n", i)
	}
	fmt.Printf("~completed dynamoDbClient.PutItem for all notes\n")

	responseBody, _ := json.Marshal(map[string]string{
		"status": "successfully seeded db",
	})
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(responseBody)}, nil
}

func main() {
	lambda.Start(handler)
}

func generateSeedNotes(event events.APIGatewayProxyRequest) ([]db.Note, error) {
	payload := []db.Note{}

	userId, err := utils.GetUserId(event)
	if err != nil {
		return nil, err
	}

	for _, v := range []struct {
		content     string
		attachement string
	}{
		{"pakyu", "betlog.jpg"},
		{"tanginamo", "tatlong_itlog.png"},
		{"pansit malabon", "minsan.ogg"},
		{"kahit pa kingina naman talaga", "you.mp4"},
		{"42069 tangninaka hayup ka pakshet", ""},
	} {
		payload = append(payload, db.Note{
			UserId:     userId,
			NoteId:     uuid.New().String(),
			Content:    v.content,
			Attachment: v.attachement,
			CreatedAt:  time.Now().Format(time.RFC3339Nano),
		})
	}
	return payload, nil
}

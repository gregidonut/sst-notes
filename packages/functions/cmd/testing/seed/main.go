package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/gregidonut/sst-notes/packages/functions/cmd/testing/seed/utils"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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
	notes := utils.GenerateSeedNotes()

	for i, note := range notes {
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

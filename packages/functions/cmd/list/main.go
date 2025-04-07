package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gregidonut/sst-notes/packages/functions/cmd/utils"
	"log"
	"os"
	"sort"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gregidonut/sst-notes/packages/functions/cmd/list/db"
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

	userId, err := utils.GetUserId(event)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: fmt.Sprintf("Error getting userId: %v", err)}, nil
	}

	params := &dynamodb.QueryInput{
		TableName:              aws.String(tableName), // Replace with your actual table name
		KeyConditionExpression: aws.String("userId = :userId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userId": &types.AttributeValueMemberS{Value: userId},
		},
	}

	fmt.Println("~claims check done")

	result, err := dynamoDbClient.Query(ctx, params)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: fmt.Sprintf("Error querying DynamoDB: %v", err)}, nil
	}
	var items []db.Note
	err = attributevalue.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: fmt.Sprintf("Error unmarshaling DynamoDB response: %v", err)}, nil
	}

	unorderedJson, err := json.Marshal(items)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error marshaling response"}, nil
	}

	//{{ sort
	notes := []db.Note{}
	err = json.Unmarshal(unorderedJson, &notes)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error unmarshaling json from response"}, nil
	}

	parseErr := struct {
		occured bool
		err     error
	}{
		occured: false,
		err:     nil,
	}
	sort.Slice(notes, func(i, j int) bool {
		parsedTimei, err := time.Parse(time.RFC3339Nano, notes[i].CreatedAt)
		if err != nil {
			parseErr.occured = true
			parseErr.err = err
		}
		parsedTimej, err := time.Parse(time.RFC3339Nano, notes[j].CreatedAt)
		if err != nil {
			parseErr.occured = true
			parseErr.err = err
		}
		if parseErr.occured {
			return false
		}
		return parsedTimei.After(parsedTimej)
	})
	if parseErr.occured {
		return events.APIGatewayProxyResponse{
			StatusCode: 500, Body: "Error parsing time from note struct in from dynamodb",
		}, nil
	}

	payload, err := json.Marshal(notes)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error marshaling response"}, nil
	}
	//}}

	return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(payload)}, nil
}

func main() {
	lambda.Start(handler)
}

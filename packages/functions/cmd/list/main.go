package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/gregidonut/sst-notes/packages/functions/cmd/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gregidonut/sst-notes/packages/functions/cmd/list/db"
)

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userId, err := utils.GetUserId(event)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: fmt.Sprintf("Error getting userId: %v", err)}, nil
	}

	params := &dynamodb.QueryInput{
		TableName:              aws.String(utils.DYNAMODB_TABLE_NAME),
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

	type sortableNote struct {
		note      db.Note
		createdAt time.Time
	}

	sortableNotes := []sortableNote{}
	for _, n := range items {
		parsedTime, err := time.Parse(time.RFC3339Nano, n.CreatedAt)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error parsing time from createdAt string"}, nil
		}
		sortableNotes = append(sortableNotes, sortableNote{n, parsedTime})
	}
	sort.Slice(sortableNotes, func(i, j int) bool {
		return sortableNotes[i].createdAt.After(sortableNotes[j].createdAt)
	})

	// Extract sorted notes
	sortedNotes := []db.Note{}
	for _, sn := range sortableNotes {
		sortedNotes = append(sortedNotes, sn.note)
	}

	payload, err := json.Marshal(sortedNotes)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error marshaling response"}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(payload)}, nil
}

func main() {
	lambda.Start(handler)
}

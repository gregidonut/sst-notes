package steps

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gregidonut/sst-notes/packages/functions/cmd/list/db"
)

func DeleteAll(ctx context.Context, notes []db.Note, tableName string, dynamoDbClient *dynamodb.Client) (events.APIGatewayProxyResponse, error) {

	for _, note := range notes {
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

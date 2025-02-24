package databases

import (
	"context"
	"fmt"
	"os"

	"planzoco/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func CreateEvent(event models.Event) error {
	item, err := attributevalue.MarshalMap(event)
	if err != nil {
		return err
	}

	_, err = DynamoClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("DYNAMODB_TABLE")),
		Item:      item,
	})
	return err
}

func GetEvent(eventID string) (*models.Event, error) {
	result, err := DynamoClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("DYNAMODB_TABLE")),
		Key: map[string]types.AttributeValue{
			"event_id": &types.AttributeValueMemberS{Value: eventID},
		},
	})
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, nil
	}

	var event models.Event
	err = attributevalue.UnmarshalMap(result.Item, &event)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func ListEvents() ([]models.Event, error) {
	result, err := DynamoClient.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(os.Getenv("DYNAMODB_TABLE")),
	})
	if err != nil {
		return nil, err
	}

	var events []models.Event
	err = attributevalue.UnmarshalListOfMaps(result.Items, &events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func AddOption(questionID string, option models.Option) error {
	// Get the event containing the question
	result, err := DynamoClient.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(os.Getenv("DYNAMODB_TABLE")),
		IndexName:              aws.String("EventQuestions"),
		KeyConditionExpression: aws.String("question_id = :qid"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":qid": &types.AttributeValueMemberS{Value: questionID},
		},
	})
	if err != nil {
		return err
	}

	if len(result.Items) == 0 {
		return fmt.Errorf("question not found")
	}

	// Update the question with the new option
	_, err = DynamoClient.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String(os.Getenv("DYNAMODB_TABLE")),
		Key: map[string]types.AttributeValue{
			"question_id": &types.AttributeValueMemberS{Value: questionID},
		},
		UpdateExpression: aws.String("SET options = list_append(if_not_exists(options, :empty_list), :option)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":option": &types.AttributeValueMemberL{Value: []types.AttributeValue{
				&types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
					"id":    &types.AttributeValueMemberS{Value: option.ID},
					"text":  &types.AttributeValueMemberS{Value: option.Text},
					"votes": &types.AttributeValueMemberN{Value: "0"},
				}},
			}},
			":empty_list": &types.AttributeValueMemberL{Value: []types.AttributeValue{}},
		},
	})
	return err
}

func VoteOption(optionID string) error {
	_, err := DynamoClient.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String(os.Getenv("DYNAMODB_TABLE")),
		Key: map[string]types.AttributeValue{
			"option_id": &types.AttributeValueMemberS{Value: optionID},
		},
		UpdateExpression: aws.String("ADD votes :inc"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":inc": &types.AttributeValueMemberN{Value: "1"},
		},
	})
	return err
}

func AddQuestion(eventID string, question models.Question) error {
	_, err := DynamoClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("DYNAMODB_TABLE_QUESTIONS")),
		Item: map[string]types.AttributeValue{
			"question_id": &types.AttributeValueMemberS{Value: question.ID},
			"event_id":    &types.AttributeValueMemberS{Value: eventID},
			"text":        &types.AttributeValueMemberS{Value: question.Text},
		},
	})
	return err
}

func GetQuestion(questionID string) (*models.Question, *models.Event, error) {
	result, err := DynamoClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("DYNAMODB_TABLE_QUESTIONS")),
		Key: map[string]types.AttributeValue{
			"question_id": &types.AttributeValueMemberS{Value: questionID},
		},
	})
	if err != nil {
		return nil, nil, err
	}

	var question models.Question
	if err := attributevalue.UnmarshalMap(result.Item, &question); err != nil {
		return nil, nil, err
	}

	event, err := GetEvent(question.EventID)
	if err != nil {
		return nil, nil, err
	}

	return &question, event, nil
}

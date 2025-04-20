package databases

import (
	"context"
	"fmt"

	"github.com/evoteum/planzoco/go/planzoco/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Event Operations

// CreateEvent creates a new event in DynamoDB
func CreateEvent(event models.Event) error {
	// Make sure the event uses the correct PK/SK pattern
	if event.PK == "" || event.SK == "" {
		event = models.NewEvent(event.ID, event.Name)
	}

	item, err := attributevalue.MarshalMap(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	_, err = DynamoClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(GetTableName()),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to put event in DynamoDB: %w", err)
	}

	return nil
}

// GetEvent retrieves an event by ID from DynamoDB
func GetEvent(eventID string) (*models.Event, error) {
	// Create the PK/SK for querying
	pk := string(models.EventEntity) + "#" + eventID
	sk := string(models.EventEntity) + "#" + eventID

	result, err := DynamoClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(GetTableName()),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: pk},
			"sk": &types.AttributeValueMemberS{Value: sk},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get event from DynamoDB: %w", err)
	}

	if result.Item == nil {
		return nil, nil
	}

	var event models.Event
	err = attributevalue.UnmarshalMap(result.Item, &event)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal DynamoDB result: %w", err)
	}

	// Get questions for this event
	questions, err := GetQuestionsByEventID(eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get questions for event: %w", err)
	}

	event.Questions = questions
	return &event, nil
}

// UpdateEvent updates an existing event in DynamoDB
func UpdateEvent(event models.Event) error {
	// Ensure the event uses the correct PK/SK pattern
	if event.PK == "" || event.SK == "" {
		event = models.NewEvent(event.ID, event.Name)
	}

	item, err := attributevalue.MarshalMap(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	_, err = DynamoClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(GetTableName()),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to update event in DynamoDB: %w", err)
	}

	return nil
}

// DeleteEvent deletes an event and all associated questions and options
func DeleteEvent(eventID string) error {
	// First get questions to get their IDs for deletion
	questions, err := GetQuestionsByEventID(eventID)
	if err != nil {
		return fmt.Errorf("failed to get questions to delete: %w", err)
	}

	// Delete each question and its options
	for _, question := range questions {
		if err := DeleteQuestion(question.ID); err != nil {
			return fmt.Errorf("failed to delete question %s: %w", question.ID, err)
		}
	}

	// Delete the event
	pk := string(models.EventEntity) + "#" + eventID
	sk := string(models.EventEntity) + "#" + eventID

	_, err = DynamoClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(GetTableName()),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: pk},
			"sk": &types.AttributeValueMemberS{Value: sk},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete event from DynamoDB: %w", err)
	}

	return nil
}

// ListEvents retrieves all events from DynamoDB using the GSI for entity type
func ListEvents() ([]models.Event, error) {
	// Use the EntityTypeIndex to query for events
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(GetTableName()),
		IndexName:              aws.String("EntityTypeIndex"),
		KeyConditionExpression: aws.String("entity_type = :entityType"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":entityType": &types.AttributeValueMemberS{Value: string(models.EventEntity)},
		},
	}

	result, err := DynamoClient.Query(context.TODO(), queryInput)
	if err != nil {
		return nil, fmt.Errorf("failed to query events from DynamoDB: %w", err)
	}

	var events []models.Event
	err = attributevalue.UnmarshalListOfMaps(result.Items, &events)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal DynamoDB query result: %w", err)
	}

	// Get questions for each event
	for i := range events {
		questions, err := GetQuestionsByEventID(events[i].ID)
		if err != nil {
			// Continue with empty questions if we can't get them
			events[i].Questions = []models.Question{}
			continue
		}
		events[i].Questions = questions
	}

	return events, nil
}

// Question Operations

// AddQuestion creates a new question in DynamoDB
func AddQuestion(eventID string, question models.Question) error {
	// Make sure the question uses the correct PK/SK pattern
	if question.PK == "" || question.SK == "" {
		question = models.NewQuestion(question.ID, eventID, question.Text)
	}

	item, err := attributevalue.MarshalMap(question)
	if err != nil {
		return fmt.Errorf("failed to marshal question: %w", err)
	}

	_, err = DynamoClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(GetTableName()),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to put question in DynamoDB: %w", err)
	}

	return nil
}

// GetQuestion retrieves a question by ID from DynamoDB
func GetQuestion(questionID string) (*models.Question, error) {
	// First, we need to find which event this question belongs to by querying the GSI
	// We can't directly get it because we don't know the SK (event ID)
	result, err := DynamoClient.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(GetTableName()),
		IndexName:              aws.String("EntityTypeIndex"),
		KeyConditionExpression: aws.String("entity_type = :entityType AND pk = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":entityType": &types.AttributeValueMemberS{Value: string(models.QuestionEntity)},
			":pk":         &types.AttributeValueMemberS{Value: string(models.QuestionEntity) + "#" + questionID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query question from DynamoDB: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, nil
	}

	var question models.Question
	err = attributevalue.UnmarshalMap(result.Items[0], &question)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal DynamoDB result: %w", err)
	}

	// Get options for this question
	options, err := GetOptionsByQuestionID(questionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get options for question: %w", err)
	}

	question.Options = options
	return &question, nil
}

// GetQuestionWithEvent retrieves a question with its associated event
func GetQuestionWithEvent(questionID string) (*models.Question, *models.Event, error) {
	question, err := GetQuestion(questionID)
	if err != nil {
		return nil, nil, err
	}

	if question == nil {
		return nil, nil, nil
	}

	event, err := GetEvent(question.EventID)
	if err != nil {
		return question, nil, err
	}

	return question, event, nil
}

// UpdateQuestion updates an existing question in DynamoDB
func UpdateQuestion(question models.Question) error {
	// Ensure the question uses the correct PK/SK pattern
	if question.PK == "" || question.SK == "" {
		existingQuestion, err := GetQuestion(question.ID)
		if err != nil {
			return fmt.Errorf("failed to get existing question for update: %w", err)
		}
		if existingQuestion == nil {
			return fmt.Errorf("question not found for update: %s", question.ID)
		}

		question = models.NewQuestion(question.ID, existingQuestion.EventID, question.Text)
		// Preserve options
		question.Options = existingQuestion.Options
	}

	item, err := attributevalue.MarshalMap(question)
	if err != nil {
		return fmt.Errorf("failed to marshal question: %w", err)
	}

	_, err = DynamoClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(GetTableName()),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to update question in DynamoDB: %w", err)
	}

	return nil
}

// DeleteQuestion deletes a question and all its options
func DeleteQuestion(questionID string) error {
	// First get the question to find its event ID and options
	question, err := GetQuestion(questionID)
	if err != nil {
		return fmt.Errorf("failed to get question to delete: %w", err)
	}

	if question == nil {
		return fmt.Errorf("question not found for deletion: %s", questionID)
	}

	// Delete all options for this question
	options, err := GetOptionsByQuestionID(questionID)
	if err != nil {
		return fmt.Errorf("failed to get options to delete: %w", err)
	}

	for _, option := range options {
		if err := DeleteOption(option.ID); err != nil {
			return fmt.Errorf("failed to delete option %s: %w", option.ID, err)
		}
	}

	// Delete the question using PK/SK
	pk := string(models.QuestionEntity) + "#" + questionID
	sk := string(models.EventEntity) + "#" + question.EventID

	_, err = DynamoClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(GetTableName()),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: pk},
			"sk": &types.AttributeValueMemberS{Value: sk},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete question from DynamoDB: %w", err)
	}

	return nil
}

// GetQuestionsByEventID retrieves all questions for a given event ID
func GetQuestionsByEventID(eventID string) ([]models.Question, error) {
	// Query using the EventIDIndex
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(GetTableName()),
		IndexName:              aws.String("EventIDIndex"),
		KeyConditionExpression: aws.String("event_id = :eventID"),
		FilterExpression:       aws.String("entity_type = :entityType"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":eventID":    &types.AttributeValueMemberS{Value: eventID},
			":entityType": &types.AttributeValueMemberS{Value: string(models.QuestionEntity)},
		},
	}

	result, err := DynamoClient.Query(context.TODO(), queryInput)
	if err != nil {
		return nil, fmt.Errorf("failed to query questions by event ID: %w", err)
	}

	var questions []models.Question
	if len(result.Items) == 0 {
		return questions, nil
	}

	err = attributevalue.UnmarshalListOfMaps(result.Items, &questions)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal DynamoDB query result: %w", err)
	}

	// Get options for each question
	for i := range questions {
		options, err := GetOptionsByQuestionID(questions[i].ID)
		if err != nil {
			continue
		}
		questions[i].Options = options
	}

	return questions, nil
}

// Option Operations

// AddOption creates a new option in DynamoDB
func AddOption(questionID string, option models.Option) error {
	// Make sure the option uses the correct PK/SK pattern
	if option.PK == "" || option.SK == "" {
		option = models.NewOption(option.ID, questionID, option.Text)
	}

	item, err := attributevalue.MarshalMap(option)
	if err != nil {
		return fmt.Errorf("failed to marshal option: %w", err)
	}

	_, err = DynamoClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(GetTableName()),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to put option in DynamoDB: %w", err)
	}

	return nil
}

// GetOption retrieves an option by ID from DynamoDB
func GetOption(optionID string) (*models.Option, error) {
	// First, we need to find which question this option belongs to by querying the GSI
	// We can't directly get it because we don't know the SK (question ID)
	result, err := DynamoClient.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(GetTableName()),
		IndexName:              aws.String("EntityTypeIndex"),
		KeyConditionExpression: aws.String("entity_type = :entityType AND pk = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":entityType": &types.AttributeValueMemberS{Value: string(models.OptionEntity)},
			":pk":         &types.AttributeValueMemberS{Value: string(models.OptionEntity) + "#" + optionID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query option from DynamoDB: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, nil
	}

	var option models.Option
	err = attributevalue.UnmarshalMap(result.Items[0], &option)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal DynamoDB result: %w", err)
	}

	return &option, nil
}

// UpdateOption updates an existing option in DynamoDB
func UpdateOption(option models.Option) error {
	// Ensure the option uses the correct PK/SK pattern
	if option.PK == "" || option.SK == "" {
		existingOption, err := GetOption(option.ID)
		if err != nil {
			return fmt.Errorf("failed to get existing option for update: %w", err)
		}
		if existingOption == nil {
			return fmt.Errorf("option not found for update: %s", option.ID)
		}

		option = models.NewOption(option.ID, existingOption.QuestionID, option.Text)
		// Preserve votes
		option.Votes = existingOption.Votes
	}

	item, err := attributevalue.MarshalMap(option)
	if err != nil {
		return fmt.Errorf("failed to marshal option: %w", err)
	}

	_, err = DynamoClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(GetTableName()),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to update option in DynamoDB: %w", err)
	}

	return nil
}

// DeleteOption deletes an option by ID from DynamoDB
func DeleteOption(optionID string) error {
	// First get the option to find its question ID
	option, err := GetOption(optionID)
	if err != nil {
		return fmt.Errorf("failed to get option to delete: %w", err)
	}

	if option == nil {
		return fmt.Errorf("option not found for deletion: %s", optionID)
	}

	// Delete the option using PK/SK
	pk := string(models.OptionEntity) + "#" + optionID
	sk := string(models.QuestionEntity) + "#" + option.QuestionID

	_, err = DynamoClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(GetTableName()),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: pk},
			"sk": &types.AttributeValueMemberS{Value: sk},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete option from DynamoDB: %w", err)
	}

	return nil
}

// VoteOption increments the vote count for an option
func VoteOption(optionID string) error {
	// First, get the current option
	option, err := GetOption(optionID)
	if err != nil {
		return fmt.Errorf("failed to get option to vote: %w", err)
	}
	if option == nil {
		return fmt.Errorf("option not found: %s", optionID)
	}

	// Increment vote count
	option.Votes++

	// Save back to DynamoDB
	return UpdateOption(*option)
}

// GetOptionsByQuestionID retrieves all options for a given question ID
func GetOptionsByQuestionID(questionID string) ([]models.Option, error) {
	// Query using the QuestionIDIndex
	result, err := DynamoClient.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(GetTableName()),
		IndexName:              aws.String("QuestionIDIndex"),
		KeyConditionExpression: aws.String("question_id = :questionID"),
		FilterExpression:       aws.String("entity_type = :entityType"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":questionID": &types.AttributeValueMemberS{Value: questionID},
			":entityType": &types.AttributeValueMemberS{Value: string(models.OptionEntity)},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query options by question ID: %w", err)
	}

	var options []models.Option
	err = attributevalue.UnmarshalListOfMaps(result.Items, &options)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal DynamoDB query result: %w", err)
	}

	return options, nil
}

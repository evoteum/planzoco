package models

type EntityType string

const (
	EventEntity    EntityType = "EVENT"
	QuestionEntity EntityType = "QUESTION"
	OptionEntity   EntityType = "OPTION"
)

// DynamoItem is the base structure for all items in the single DynamoDB table
type DynamoItem struct {
	PK string `json:"pk" dynamodbav:"pk"`
	SK string `json:"sk" dynamodbav:"sk"`
}

// Event represents a planning event
type Event struct {
	DynamoItem
	ID         string     `json:"id" dynamodbav:"id"`
	Name       string     `json:"name" form:"name" binding:"required" dynamodbav:"name"`
	Questions  []Question `json:"questions,omitempty" dynamodbav:"-"` // Not stored directly in the item
	EntityType EntityType `json:"-" dynamodbav:"entity_type"`
}

// NewEvent creates a new Event with the proper PK/SK pattern
func NewEvent(id string, name string) Event {
	return Event{
		DynamoItem: DynamoItem{
			PK: string(EventEntity) + "#" + id,
			SK: string(EventEntity) + "#" + id,
		},
		ID:         id,
		Name:       name,
		EntityType: EventEntity,
	}
}

// Question represents a question within an event
type Question struct {
	DynamoItem
	ID         string     `json:"id" dynamodbav:"id"`
	EventID    string     `json:"event_id" dynamodbav:"event_id"`
	Text       string     `json:"text" form:"text" binding:"required" dynamodbav:"text"`
	Options    []Option   `json:"options,omitempty" dynamodbav:"-"` // Not stored directly in the item
	EntityType EntityType `json:"-" dynamodbav:"entity_type"`
}

// NewQuestion creates a new Question with the proper PK/SK pattern
func NewQuestion(id string, eventID string, text string) Question {
	return Question{
		DynamoItem: DynamoItem{
			PK: string(QuestionEntity) + "#" + id,
			SK: string(EventEntity) + "#" + eventID,
		},
		ID:         id,
		EventID:    eventID,
		Text:       text,
		EntityType: QuestionEntity,
	}
}

func (q Question) WinningOptions() []Option {
	if len(q.Options) == 0 {
		return nil
	}

	allZero := true
	maxVotes := q.Options[0].Votes

	for _, opt := range q.Options {
		if opt.Votes > 0 {
			allZero = false
		}
		if opt.Votes > maxVotes {
			maxVotes = opt.Votes
		}
	}

	if allZero {
		return nil
	}

	// Collect all options with max votes
	var winners []Option
	for _, opt := range q.Options {
		if opt.Votes == maxVotes {
			winners = append(winners, opt)
		}
	}

	return winners
}

// Option represents an answer option for a question
type Option struct {
	DynamoItem
	ID         string     `json:"id" dynamodbav:"id"`
	QuestionID string     `json:"question_id" dynamodbav:"question_id"`
	Text       string     `json:"text" form:"text" binding:"required" dynamodbav:"text"`
	Votes      int        `json:"votes" dynamodbav:"votes"`
	EntityType EntityType `json:"-" dynamodbav:"entity_type"`
}

// NewOption creates a new Option with the proper PK/SK pattern
func NewOption(id string, questionID string, text string) Option {
	return Option{
		DynamoItem: DynamoItem{
			PK: string(OptionEntity) + "#" + id,
			SK: string(QuestionEntity) + "#" + questionID,
		},
		ID:         id,
		QuestionID: questionID,
		Text:       text,
		Votes:      0,
		EntityType: OptionEntity,
	}
}

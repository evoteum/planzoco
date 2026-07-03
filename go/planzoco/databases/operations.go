package databases

import (
	"context"
	"errors"
	"fmt"

	"github.com/evoteum/planzoco/models"

	"github.com/jackc/pgx/v5"
)

// Event Operations

// CreateEvent creates a new event in Postgres
func CreateEvent(event models.Event) error {
	_, err := Pool.Exec(context.Background(),
		"INSERT INTO events (id, name) VALUES ($1, $2)",
		event.ID, event.Name,
	)
	if err != nil {
		return fmt.Errorf("failed to insert event: %w", err)
	}

	return nil
}

// GetEvent retrieves an event by ID from Postgres
func GetEvent(eventID string) (*models.Event, error) {
	var event models.Event
	err := Pool.QueryRow(context.Background(),
		"SELECT id, name FROM events WHERE id = $1", eventID,
	).Scan(&event.ID, &event.Name)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	questions, err := GetQuestionsByEventID(eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get questions for event: %w", err)
	}

	event.Questions = questions
	return &event, nil
}

// UpdateEvent updates an existing event in Postgres
func UpdateEvent(event models.Event) error {
	tag, err := Pool.Exec(context.Background(),
		"UPDATE events SET name = $1 WHERE id = $2",
		event.Name, event.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("event not found for update: %s", event.ID)
	}

	return nil
}

// DeleteEvent deletes an event and all associated questions and options
func DeleteEvent(eventID string) error {
	tag, err := Pool.Exec(context.Background(), "DELETE FROM events WHERE id = $1", eventID)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("event not found for deletion: %s", eventID)
	}

	return nil
}

// ListEvents retrieves all events from Postgres
func ListEvents() ([]models.Event, error) {
	rows, err := Pool.Query(context.Background(), "SELECT id, name FROM events ORDER BY id")
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		if err := rows.Scan(&event.ID, &event.Name); err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

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

// AddQuestion creates a new question in Postgres
func AddQuestion(eventID string, question models.Question) error {
	_, err := Pool.Exec(context.Background(),
		"INSERT INTO questions (id, event_id, text, asked_by) VALUES ($1, $2, $3, $4)",
		question.ID, eventID, question.Text, question.AskedBy,
	)
	if err != nil {
		return fmt.Errorf("failed to insert question: %w", err)
	}

	return nil
}

// GetQuestion retrieves a question by ID from Postgres
func GetQuestion(questionID string) (*models.Question, error) {
	var question models.Question
	err := Pool.QueryRow(context.Background(),
		"SELECT id, event_id, text, asked_by FROM questions WHERE id = $1", questionID,
	).Scan(&question.ID, &question.EventID, &question.Text, &question.AskedBy)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get question: %w", err)
	}

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

// UpdateQuestion updates an existing question in Postgres
func UpdateQuestion(question models.Question) error {
	tag, err := Pool.Exec(context.Background(),
		"UPDATE questions SET text = $1 WHERE id = $2",
		question.Text, question.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update question: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("question not found for update: %s", question.ID)
	}

	return nil
}

// DeleteQuestion deletes a question and all its options
func DeleteQuestion(questionID string) error {
	tag, err := Pool.Exec(context.Background(), "DELETE FROM questions WHERE id = $1", questionID)
	if err != nil {
		return fmt.Errorf("failed to delete question: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("question not found for deletion: %s", questionID)
	}

	return nil
}

// GetQuestionsByEventID retrieves all questions for a given event ID
func GetQuestionsByEventID(eventID string) ([]models.Question, error) {
	rows, err := Pool.Query(context.Background(),
		"SELECT id, event_id, text, asked_by FROM questions WHERE event_id = $1 ORDER BY id", eventID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query questions by event ID: %w", err)
	}
	defer rows.Close()

	var questions []models.Question
	for rows.Next() {
		var question models.Question
		if err := rows.Scan(&question.ID, &question.EventID, &question.Text, &question.AskedBy); err != nil {
			return nil, fmt.Errorf("failed to scan question: %w", err)
		}
		questions = append(questions, question)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to query questions by event ID: %w", err)
	}

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

// AddOption creates a new option in Postgres
func AddOption(questionID string, option models.Option) error {
	_, err := Pool.Exec(context.Background(),
		"INSERT INTO options (id, question_id, text) VALUES ($1, $2, $3)",
		option.ID, questionID, option.Text,
	)
	if err != nil {
		return fmt.Errorf("failed to insert option: %w", err)
	}

	return nil
}

// GetOption retrieves an option by ID from Postgres, with its vote count and voters
func GetOption(optionID string) (*models.Option, error) {
	var option models.Option
	err := Pool.QueryRow(context.Background(),
		`SELECT o.id, o.question_id, o.text, COUNT(v.voter_name),
		        COALESCE(array_agg(v.voter_name ORDER BY v.voter_name) FILTER (WHERE v.voter_name IS NOT NULL), '{}')
		 FROM options o
		 LEFT JOIN votes v ON v.option_id = o.id
		 WHERE o.id = $1
		 GROUP BY o.id`, optionID,
	).Scan(&option.ID, &option.QuestionID, &option.Text, &option.Votes, &option.Voters)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get option: %w", err)
	}

	return &option, nil
}

// UpdateOption updates an existing option's text in Postgres, leaving votes untouched
func UpdateOption(option models.Option) error {
	tag, err := Pool.Exec(context.Background(),
		"UPDATE options SET text = $1 WHERE id = $2",
		option.Text, option.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update option: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("option not found for update: %s", option.ID)
	}

	return nil
}

// DeleteOption deletes an option by ID from Postgres
func DeleteOption(optionID string) error {
	tag, err := Pool.Exec(context.Background(), "DELETE FROM options WHERE id = $1", optionID)
	if err != nil {
		return fmt.Errorf("failed to delete option: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("option not found for deletion: %s", optionID)
	}

	return nil
}

// VoteOption records voterName's vote for optionID on the given question,
// replacing any previous vote they cast on that question.
func VoteOption(questionID, optionID, voterName string) error {
	_, err := Pool.Exec(context.Background(),
		`INSERT INTO votes (question_id, option_id, voter_name) VALUES ($1, $2, $3)
		 ON CONFLICT (question_id, voter_name) DO UPDATE SET option_id = EXCLUDED.option_id`,
		questionID, optionID, voterName,
	)
	if err != nil {
		return fmt.Errorf("failed to record vote: %w", err)
	}

	return nil
}

// GetOptionsByQuestionID retrieves all options for a given question ID, with vote counts and voters
func GetOptionsByQuestionID(questionID string) ([]models.Option, error) {
	rows, err := Pool.Query(context.Background(),
		`SELECT o.id, o.question_id, o.text, COUNT(v.voter_name),
		        COALESCE(array_agg(v.voter_name ORDER BY v.voter_name) FILTER (WHERE v.voter_name IS NOT NULL), '{}')
		 FROM options o
		 LEFT JOIN votes v ON v.option_id = o.id
		 WHERE o.question_id = $1
		 GROUP BY o.id
		 ORDER BY o.id`, questionID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query options by question ID: %w", err)
	}
	defer rows.Close()

	var options []models.Option
	for rows.Next() {
		var option models.Option
		if err := rows.Scan(&option.ID, &option.QuestionID, &option.Text, &option.Votes, &option.Voters); err != nil {
			return nil, fmt.Errorf("failed to scan option: %w", err)
		}
		options = append(options, option)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to query options by question ID: %w", err)
	}

	return options, nil
}

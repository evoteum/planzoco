package databases

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

const schema = `
CREATE TABLE IF NOT EXISTS events (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL CHECK (char_length(name) <= 255)
);

CREATE TABLE IF NOT EXISTS questions (
	id TEXT PRIMARY KEY,
	event_id TEXT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
	text TEXT NOT NULL CHECK (char_length(text) <= 255)
);

CREATE TABLE IF NOT EXISTS options (
	id TEXT PRIMARY KEY,
	question_id TEXT NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
	text TEXT NOT NULL CHECK (char_length(text) <= 255)
);

CREATE TABLE IF NOT EXISTS votes (
	question_id TEXT NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
	option_id TEXT NOT NULL REFERENCES options(id) ON DELETE CASCADE,
	voter_name TEXT NOT NULL CHECK (char_length(voter_name) <= 255),
	PRIMARY KEY (question_id, voter_name)
);

CREATE INDEX IF NOT EXISTS idx_questions_event_id ON questions(event_id);
CREATE INDEX IF NOT EXISTS idx_options_question_id ON options(question_id);
CREATE INDEX IF NOT EXISTS idx_votes_option_id ON votes(option_id);
`

func InitDB() error {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return fmt.Errorf("DATABASE_URL is not set")
	}

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return fmt.Errorf("unable to ping database: %w", err)
	}

	if _, err := pool.Exec(context.Background(), schema); err != nil {
		return fmt.Errorf("unable to initialize schema: %w", err)
	}

	Pool = pool
	log.Println("Postgres connection pool initialized and schema ensured")

	return nil
}

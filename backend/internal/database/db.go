package database

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Message struct {
	ID        int64     `db:"id"`
	Type      string    `db:"type"`
	Content   string    `db:"content"`
	Sender    string    `db:"sender"`
	Timestamp time.Time `db:"timestamp"`
}

type DB struct {
	*sqlx.DB
}

func NewDB(host, port, user, password, dbname string) (*DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}

	return &DB{db}, nil
}

func (db *DB) InitSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		type VARCHAR(20) NOT NULL,
		content TEXT NOT NULL,
		sender VARCHAR(255) NOT NULL,
		timestamp TIMESTAMP NOT NULL
	);`

	_, err := db.Exec(schema)
	return err
}

func (db *DB) SaveMessage(msg Message) error {
	query := `
	INSERT INTO messages (type, content, sender, timestamp)
	VALUES ($1, $2, $3, $4)`

	_, err := db.Exec(query, msg.Type, msg.Content, msg.Sender, msg.Timestamp)
	return err
}

func (db *DB) GetRecentMessages(limit int) ([]Message, error) {
	var messages []Message
	query := `
	SELECT id, type, content, sender, timestamp
	FROM messages
	ORDER BY timestamp DESC
	LIMIT $1`

	err := db.Select(&messages, query, limit)
	return messages, err
}

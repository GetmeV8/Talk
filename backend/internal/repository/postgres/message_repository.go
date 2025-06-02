package postgres

import (
	"messenger/internal/domain"

	"github.com/jmoiron/sqlx"
)

type MessageRepository struct {
	db *sqlx.DB
}

func NewMessageRepository(db *sqlx.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Save(message domain.Message) error {
	query := `
	INSERT INTO messages (type, content, sender, timestamp)
	VALUES ($1, $2, $3, $4)`

	_, err := r.db.Exec(query, message.Type, message.Content, message.Sender, message.Timestamp)
	return err
}

func (r *MessageRepository) GetRecent(limit int) ([]domain.Message, error) {
	var messages []domain.Message
	query := `
	SELECT id, type, content, sender, timestamp
	FROM messages
	ORDER BY timestamp DESC
	LIMIT $1`

	err := r.db.Select(&messages, query, limit)
	return messages, err
}

package domain

import "time"

// Message represents a chat message in the domain
type Message struct {
	ID        int64
	Type      string
	Content   string
	Sender    string
	Timestamp time.Time
}

// MessageRepository defines the interface for message storage
type MessageRepository interface {
	Save(message Message) error
	GetRecent(limit int) ([]Message, error)
}

// MessageService defines the interface for message business logic
type MessageService interface {
	SendMessage(message Message) error
	GetRecentMessages(limit int) ([]Message, error)
	BroadcastMessage(message Message) error
}

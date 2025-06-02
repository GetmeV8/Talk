package service

import (
	"messenger/internal/domain"
)

type MessageService struct {
	repo domain.MessageRepository
	// Add any additional dependencies here (e.g., event bus, cache, etc.)
}

func NewMessageService(repo domain.MessageRepository) *MessageService {
	return &MessageService{
		repo: repo,
	}
}

func (s *MessageService) SendMessage(message domain.Message) error {
	// Add any business logic here (e.g., validation, rate limiting, etc.)
	return s.repo.Save(message)
}

func (s *MessageService) GetRecentMessages(limit int) ([]domain.Message, error) {
	return s.repo.GetRecent(limit)
}

func (s *MessageService) BroadcastMessage(message domain.Message) error {
	// Add any broadcast-specific logic here
	return s.repo.Save(message)
}

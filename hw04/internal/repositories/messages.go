package repositories

import (
	"github.com/sirupsen/logrus"
	"github.com/willsem/tfs-go-hw/hw04/internal/domain"
)

const (
	sharedChatId = "shared"
)

type MessagesRepository struct {
	logger   *logrus.Logger
	messages []domain.Message
}

func NewMessagesRepository(logger *logrus.Logger) *MessagesRepository {
	return &MessagesRepository{
		logger:   logger,
		messages: make([]domain.Message, 0),
	}
}

func (repository *MessagesRepository) GetSharedMessages() ([]domain.Message, error) {
	result := make([]domain.Message, 0)
	for _, message := range repository.messages {
		if message.RecipentID == sharedChatId {
			result = append(result, message)
		}
	}
	return result, nil
}

func (repository *MessagesRepository) GetMessages(user1, user2 string) ([]domain.Message, error) {
	result := make([]domain.Message, 0)
	for _, message := range repository.messages {
		if message.SenderID == user1 && message.RecipentID == user2 ||
			message.SenderID == user2 && message.RecipentID == user1 {
			result = append(result, message)
		}
	}
	return result, nil
}

func (repository *MessagesRepository) AddMessage(message domain.Message) error {
	repository.messages = append(repository.messages, message)
	return nil
}

func (repository *MessagesRepository) GetSharedChatId() string {
	return sharedChatId
}

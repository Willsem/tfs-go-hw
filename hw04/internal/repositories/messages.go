package repositories

import (
	"github.com/sirupsen/logrus"
	"github.com/willsem/tfs-go-hw/hw04/internal/domain"
)

type MessagesRepository struct {
	logger *logrus.Logger
}

func NewMessagesRepository(logger *logrus.Logger) *MessagesRepository {
	return &MessagesRepository{
		logger: logger,
	}
}

func (repository *MessagesRepository) GetSharedMessages() ([]domain.Message, error) {
	return nil, nil
}

func (repository *MessagesRepository) GetMessages(user1, user2 string) ([]domain.Message, error) {
	return nil, nil
}

func (repository *MessagesRepository) AddMessage(message domain.Message) error {
	return nil
}

func (repository *MessagesRepository) GetSharedChatId() string {
	return ""
}

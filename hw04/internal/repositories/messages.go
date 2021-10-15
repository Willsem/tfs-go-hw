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

func (repository *MessagesRepository) GetSharedMessages(offset, count int) ([]domain.Message, error) {
	result := make([]domain.Message, 0)
	for i := offset; i < offset+count && i < len(repository.messages); i++ {
		if repository.messages[i].RecipentID == sharedChatId {
			result = append(result, repository.messages[i])
		}
	}
	return result, nil
}

func (repository *MessagesRepository) GetMessages(offset, count int, user1, user2 string) ([]domain.Message, error) {
	result := make([]domain.Message, 0)
	for i := offset; i < offset+count && i < len(repository.messages); i++ {
		if repository.messages[i].SenderID == user1 && repository.messages[i].RecipentID == user2 ||
			repository.messages[i].SenderID == user2 && repository.messages[i].RecipentID == user1 {
			result = append(result, repository.messages[i])
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

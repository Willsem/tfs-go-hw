package services

import (
	"github.com/sirupsen/logrus"
	"github.com/willsem/tfs-go-hw/hw04/internal/dto"
)

type MessagesService struct {
	messagesRepository MessagesRepository
	logger             *logrus.Logger
}

func NewMessagesService(messagesRepository MessagesRepository, logger *logrus.Logger) *MessagesService {
	return &MessagesService{
		messagesRepository: messagesRepository,
		logger:             logger,
	}
}

func (service *MessagesService) GetSharedMessages() ([]dto.Message, error) {
	return nil, nil
}

func (service *MessagesService) SendSharedMessage(sender string, message dto.MessageContent) error {
	return nil
}

func (service *MessagesService) GetPrivateMessages(user1 string, user2 string) ([]dto.Message, error) {
	return nil, nil
}

func (service *MessagesService) SendPrivateMessage(sender string, recipent string, message dto.MessageContent) error {
	return nil
}

package services

import "github.com/sirupsen/logrus"

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

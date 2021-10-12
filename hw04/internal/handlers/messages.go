package handlers

import "github.com/sirupsen/logrus"

type MessagesHandler struct {
	messagesService MessagesService
	logger          *logrus.Logger
}

func NewMessagesHandler(messagesService MessagesService, logger *logrus.Logger) *MessagesHandler {
	return &MessagesHandler{
		messagesService: messagesService,
		logger:          logger,
	}
}

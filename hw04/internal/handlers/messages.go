package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

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

func (handler *MessagesHandler) Routes() chi.Router {
	r := chi.NewRouter()

	return r
}

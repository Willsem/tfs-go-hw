package repositories

import "github.com/sirupsen/logrus"

type MessagesRepository struct {
	logger *logrus.Logger
}

func NewMessagesRepository(logger *logrus.Logger) *MessagesRepository {
	return &MessagesRepository{
		logger: logger,
	}
}

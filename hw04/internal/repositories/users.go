package repositories

import "github.com/sirupsen/logrus"

type UsersRepository struct {
	logger *logrus.Logger
}

func NewUsersRepository(logger *logrus.Logger) *UsersRepository {
	return &UsersRepository{
		logger: logger,
	}
}

package services

import "github.com/sirupsen/logrus"

type UsersService struct {
	usersRepository UsersRepository
	logger          *logrus.Logger
}

func NewUsersService(usersRepository UsersRepository, logger *logrus.Logger) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
		logger:          logger,
	}
}

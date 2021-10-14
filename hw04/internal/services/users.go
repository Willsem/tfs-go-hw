package services

import (
	"github.com/sirupsen/logrus"
	"github.com/willsem/tfs-go-hw/hw04/internal/dto"
)

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

func (service *UsersService) Register(dto.Login) (dto.User, error) {
	return dto.User{}, nil
}

func (service *UsersService) Login(dto.Login) (dto.Session, error) {
	return dto.Session{}, nil
}

func (service *UsersService) Logout(token string) error {
	return nil
}

func (service *UsersService) CheckToken(token string) (string, error) {
	return "", nil
}

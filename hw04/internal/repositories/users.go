package repositories

import (
	"github.com/sirupsen/logrus"
	"github.com/willsem/tfs-go-hw/hw04/internal/domain"
)

type UsersRepository struct {
	logger *logrus.Logger
}

func NewUsersRepository(logger *logrus.Logger) *UsersRepository {
	return &UsersRepository{
		logger: logger,
	}
}

func (repository *UsersRepository) AddUser(login string, passcode string) (domain.User, error) {
	return domain.User{}, nil
}

func (repository *UsersRepository) GetUser(login string) (domain.User, error) {
	return domain.User{}, nil
}

func (repository *UsersRepository) GetUserByID(userId string) (domain.User, error) {
	return domain.User{}, nil
}

func (repository *UsersRepository) CreateSession(userId string, token string) error {
	return nil
}

func (repository *UsersRepository) RemoveSession(token string) error {
	return nil
}

func (repository *UsersRepository) FindSession(token string) (string, error) {
	return "", nil
}

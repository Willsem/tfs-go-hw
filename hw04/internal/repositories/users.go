package repositories

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/willsem/tfs-go-hw/hw04/internal/domain"
)

type UsersRepository struct {
	logger   *logrus.Logger
	users    []domain.User
	sessions []domain.Session
}

func NewUsersRepository(logger *logrus.Logger) *UsersRepository {
	return &UsersRepository{
		logger:   logger,
		users:    make([]domain.User, 0),
		sessions: make([]domain.Session, 0),
	}
}

func (repository *UsersRepository) AddUser(login string, passcode string) (domain.User, error) {
	user := domain.User{
		ID:       uuid.New().String(),
		Login:    login,
		Passcode: passcode,
	}
	repository.users = append(repository.users, user)
	return user, nil
}

func (repository *UsersRepository) GetUser(login string) (domain.User, error) {
	for _, user := range repository.users {
		if user.Login == login {
			return user, nil
		}
	}
	return domain.User{}, fmt.Errorf("No user with that login")
}

func (repository *UsersRepository) GetUserByID(userId string) (domain.User, error) {
	for _, user := range repository.users {
		if user.ID == userId {
			return user, nil
		}
	}
	return domain.User{}, fmt.Errorf("No user with that ID")
}

func (repository *UsersRepository) CreateSession(userId string, token string) error {
	repository.sessions = append(repository.sessions, domain.Session{
		UserID: userId,
		Token:  token,
	})
	return nil
}

func (repository *UsersRepository) RemoveSession(token string) error {
	removeIndex := -1
	for i, session := range repository.sessions {
		if session.Token == token {
			removeIndex = i
			break
		}
	}

	if removeIndex == -1 {
		return fmt.Errorf("No session with that token")
	}

	repository.sessions = append(repository.sessions[:removeIndex], repository.sessions[removeIndex+1:]...)
	return nil
}

func (repository *UsersRepository) FindSession(token string) (string, error) {
	index := -1
	for i, session := range repository.sessions {
		if session.Token == token {
			index = i
			break
		}
	}

	if index == -1 {
		return "", fmt.Errorf("No session with that token")
	}

	return repository.sessions[index].UserID, nil
}

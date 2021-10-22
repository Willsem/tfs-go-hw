package repositories

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/willsem/tfs-go-hw/hw04/internal/domain"
)

type UsersRepository struct {
	logger        *logrus.Logger
	users         []domain.User
	sessions      []domain.Session
	usersMutex    *sync.RWMutex
	sessionsMutex *sync.RWMutex
}

func NewUsersRepository(logger *logrus.Logger) *UsersRepository {
	return &UsersRepository{
		logger:        logger,
		users:         make([]domain.User, 0),
		sessions:      make([]domain.Session, 0),
		usersMutex:    &sync.RWMutex{},
		sessionsMutex: &sync.RWMutex{},
	}
}

func (repository *UsersRepository) AddUser(login string, passcode string) (domain.User, error) {
	repository.usersMutex.RLock()
	for _, user := range repository.users {
		if user.Login == login {
			repository.usersMutex.RUnlock()
			return domain.User{}, fmt.Errorf("User with this login is already exist")
		}
	}
	repository.usersMutex.RUnlock()

	user := domain.User{
		ID:       uuid.New().String(),
		Login:    login,
		Passcode: passcode,
	}

	repository.usersMutex.Lock()
	repository.users = append(repository.users, user)
	repository.usersMutex.Unlock()

	return user, nil
}

func (repository *UsersRepository) GetUser(login string) (domain.User, error) {
	repository.usersMutex.RLock()
	defer repository.usersMutex.RUnlock()

	for _, user := range repository.users {
		if user.Login == login {
			return user, nil
		}
	}

	return domain.User{}, fmt.Errorf("No user with that login")
}

func (repository *UsersRepository) GetUserByID(userId string) (domain.User, error) {
	repository.usersMutex.RLock()
	defer repository.usersMutex.RUnlock()

	for _, user := range repository.users {
		if user.ID == userId {
			return user, nil
		}
	}

	return domain.User{}, fmt.Errorf("No user with that ID")
}

func (repository *UsersRepository) CreateSession(userId string, token string) error {
	repository.sessionsMutex.Lock()
	repository.sessions = append(repository.sessions, domain.Session{
		UserID: userId,
		Token:  token,
	})
	repository.sessionsMutex.Unlock()

	return nil
}

func (repository *UsersRepository) RemoveSession(token string) error {
	removeIndex := -1

	repository.sessionsMutex.RLock()
	for i, session := range repository.sessions {
		if session.Token == token {
			removeIndex = i
			break
		}
	}
	repository.sessionsMutex.RUnlock()

	if removeIndex == -1 {
		return fmt.Errorf("No session with that token")
	}

	repository.sessionsMutex.Lock()
	repository.sessions = append(repository.sessions[:removeIndex], repository.sessions[removeIndex+1:]...)
	repository.sessionsMutex.Unlock()

	return nil
}

func (repository *UsersRepository) FindSession(token string) (string, error) {
	index := -1

	repository.sessionsMutex.RLock()
	defer repository.sessionsMutex.RUnlock()

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

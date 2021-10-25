package services

import "github.com/willsem/tfs-go-hw/hw04/internal/domain"

type UsersRepository interface {
	AddUser(login string, passcode string) (domain.User, error)
	GetUser(login string) (domain.User, error)
	GetUserByID(userId string) (domain.User, error)

	CreateSession(userId string, token string) error
	RemoveSession(token string) error
	FindSession(token string) (string, error)
}

type MessagesRepository interface {
	GetSharedMessages(offset, count int) ([]domain.Message, error)
	GetMessages(offset, count int, user1, user2 string) ([]domain.Message, error)
	AddMessage(message domain.Message) error

	GetSharedChatId() string
}

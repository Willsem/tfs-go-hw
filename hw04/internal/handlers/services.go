package handlers

import "github.com/willsem/tfs-go-hw/hw04/internal/dto"

type UsersService interface {
	Register(login dto.Login) (dto.User, error)
	Login(login dto.Login) (dto.Session, error)
	Logout(token string) error
	CheckToken(token string) (string, error)
}

type MessagesService interface {
	GetSharedMessages(offset, count int) ([]dto.Message, error)
	SendSharedMessage(sender string, message dto.MessageContent) error

	GetPrivateMessages(offset, count int, user1, user2 string) ([]dto.Message, error)
	SendPrivateMessage(sender string, recipent string, message dto.MessageContent) error
}

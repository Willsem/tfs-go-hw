package handlers

import "github.com/willsem/tfs-go-hw/hw04/internal/dto"

type UsersService interface {
	Register(dto.Login) (dto.User, error)
	Login(dto.Login) (dto.Session, error)
	Logout(token string) error
	CheckToken(token string) (string, error)
}

type MessagesService interface {
	GetSharedMessages() ([]dto.Message, error)
	SendSharedMessage(sender string, message dto.MessageContent) error

	GetPrivateMessages(user1 string, user2 string) ([]dto.Message, error)
	SendPrivateMessage(sender string, recipent string, message dto.MessageContent) error
}

package services

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/willsem/tfs-go-hw/hw04/internal/domain"
	"github.com/willsem/tfs-go-hw/hw04/internal/dto"
)

type MessagesService struct {
	messagesRepository MessagesRepository
	usersRepository    UsersRepository
	logger             *logrus.Logger
}

func NewMessagesService(
	messagesRepository MessagesRepository,
	usersRepository UsersRepository,
	logger *logrus.Logger,
) *MessagesService {
	return &MessagesService{
		messagesRepository: messagesRepository,
		usersRepository:    usersRepository,
		logger:             logger,
	}
}

func (service *MessagesService) GetSharedMessages(offset, count int) ([]dto.Message, error) {
	return service.GetPrivateMessages(offset, count, "shared", "shared")
}

func (service *MessagesService) SendSharedMessage(sender string, message dto.MessageContent) error {
	return service.SendPrivateMessage(sender, service.messagesRepository.GetSharedChatId(), message)
}

func (service *MessagesService) GetPrivateMessages(offset, count int, user1, user2 string) ([]dto.Message, error) {
	messages := make([]dto.Message, 0)

	var domain []domain.Message
	var err error
	if user1 == "shared" && user2 == "shared" {
		domain, err = service.messagesRepository.GetSharedMessages(offset, count)
	} else {
		domain, err = service.messagesRepository.GetMessages(offset, count, user1, user2)
	}

	if err != nil {
		return nil, err
	}

	for _, message := range domain {
		user, err := service.usersRepository.GetUserByID(message.SenderID)
		if err != nil {
			return nil, err
		}

		messages = append(messages, dto.Message{
			User: dto.User{
				ID:    user.ID,
				Login: user.Login,
			},
			Content:  message.Content,
			DateTime: message.DateTime,
		})
	}

	return messages, nil
}

func (service *MessagesService) SendPrivateMessage(sender string, recipent string, message dto.MessageContent) error {
	messageDomain := domain.Message{
		SenderID:   sender,
		RecipentID: recipent,
		Content:    message.Content,
		DateTime:   time.Now(),
	}
	return service.messagesRepository.AddMessage(messageDomain)
}

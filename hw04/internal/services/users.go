package services

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/willsem/tfs-go-hw/hw04/internal/dto"
	"github.com/willsem/tfs-go-hw/hw04/pkg/auth"
	"golang.org/x/crypto/bcrypt"
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

func (service *UsersService) Register(login dto.Login) (dto.User, error) {
	userLogin := login.Login
	userPasscode, err := auth.HashPassword(login.Password)
	if err != nil {
		return dto.User{}, err
	}

	user, err := service.usersRepository.AddUser(userLogin, userPasscode)
	if err != nil {
		return dto.User{}, err
	}

	return dto.User{
		ID:    user.ID,
		Login: user.Login,
	}, nil
}

func (service *UsersService) Login(login dto.Login) (dto.Session, error) {
	user, err := service.usersRepository.GetUser(login.Login)
	if err != nil {
		return dto.Session{}, err
	}

	if !auth.CheckPasswordHash(login.Password, user.Passcode) {
		return dto.Session{}, fmt.Errorf("Incorrect password")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(login.Login), bcrypt.DefaultCost)
	if err != nil {
		return dto.Session{}, err
	}

	err = service.usersRepository.CreateSession(user.ID, string(hash))
	if err != nil {
		return dto.Session{}, err
	}

	return dto.Session{
		UserID: user.ID,
		Token:  string(hash),
	}, nil
}

func (service *UsersService) Logout(token string) error {
	return service.usersRepository.RemoveSession(token)
}

func (service *UsersService) CheckToken(token string) (string, error) {
	return service.usersRepository.FindSession(token)
}

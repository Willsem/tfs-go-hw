package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type UsersHandler struct {
	usersService UsersService
	logger       *logrus.Logger
}

func NewUsersHandler(usersService UsersService, logger *logrus.Logger) *UsersHandler {
	return &UsersHandler{
		usersService: usersService,
		logger:       logger,
	}
}

func (handler *UsersHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/register", handler.register)

	return r
}

func (handler *UsersHandler) register(w http.ResponseWriter, r *http.Request) {
}

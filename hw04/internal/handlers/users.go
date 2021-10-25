package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/willsem/tfs-go-hw/hw04/internal/dto"
	"github.com/willsem/tfs-go-hw/hw04/pkg/response"
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
	r.Post("/login", handler.login)
	r.Delete("/logout", handler.logout)

	return r
}

func (handler *UsersHandler) register(w http.ResponseWriter, r *http.Request) {
	login := dto.Login{}

	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		response.Respond(w, http.StatusBadRequest, dto.Error{Message: "Incorrect body"})
		return
	}
	defer r.Body.Close()

	user, err := handler.usersService.Register(login)
	if err != nil {
		response.Respond(w, http.StatusInternalServerError, dto.Error{Message: err.Error()})
		return
	}

	response.Respond(w, http.StatusOK, user)
}

func (handler *UsersHandler) login(w http.ResponseWriter, r *http.Request) {
	login := dto.Login{}

	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		response.Respond(w, http.StatusBadRequest, dto.Error{Message: "Incorrect body"})
		return
	}
	defer r.Body.Close()

	session, err := handler.usersService.Login(login)
	if err != nil {
		response.Respond(w, http.StatusInternalServerError, dto.Error{Message: err.Error()})
		return
	}

	response.Respond(w, http.StatusOK, session)
}

func (handler *UsersHandler) logout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Token")
	if token == "" {
		response.Respond(w, http.StatusUnauthorized, nil)
		return
	}

	err := handler.usersService.Logout(token)
	if err != nil {
		response.Respond(w, http.StatusInternalServerError, dto.Error{Message: err.Error()})
		return
	}

	response.Respond(w, http.StatusOK, nil)
}

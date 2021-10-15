package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/willsem/tfs-go-hw/hw04/internal/dto"
	"github.com/willsem/tfs-go-hw/hw04/pkg/response"
)

const (
	userIdCtx = "userId"
)

type MessagesHandler struct {
	messagesService MessagesService
	usersService    UsersService
	logger          *logrus.Logger
}

func NewMessagesHandler(
	messagesService MessagesService,
	usersService UsersService,
	logger *logrus.Logger,
) *MessagesHandler {
	return &MessagesHandler{
		messagesService: messagesService,
		usersService:    usersService,
		logger:          logger,
	}
}

func (handler *MessagesHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Use(handler.authMiddleware)
	r.Get("/", handler.sharedChatGet)
	r.Post("/", handler.sharedChatSendMessage)
	r.Get("/{userId}", handler.privateChatGet)
	r.Post("/{userId}", handler.privateChatSendMessage)

	return r
}

func (handler *MessagesHandler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Token")

		userId, err := handler.usersService.CheckToken(token)
		if err != nil {
			response.Respond(w, http.StatusUnauthorized, nil)
			return
		}

		ctx := context.WithValue(r.Context(), userIdCtx, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (handler *MessagesHandler) sharedChatGet(w http.ResponseWriter, r *http.Request) {
	offset := r.URL.Query().Get("offset")
	if offset == "" {
		offset = "0"
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		response.Respond(w, http.StatusBadRequest, dto.Error{Message: err.Error()})
		return
	}

	count := r.URL.Query().Get("count")
	if count == "" {
		count = "50"
	}

	countInt, err := strconv.Atoi(count)
	if err != nil {
		response.Respond(w, http.StatusBadRequest, dto.Error{Message: err.Error()})
		return
	}

	messages, err := handler.messagesService.GetSharedMessages(offsetInt, countInt)
	if err != nil {
		response.Respond(w, http.StatusInternalServerError, dto.Error{Message: err.Error()})
		return
	}

	response.Respond(w, http.StatusOK, messages)
}

func (handler *MessagesHandler) sharedChatSendMessage(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(userIdCtx).(string)
	if !ok {
		response.Respond(w, http.StatusInternalServerError, dto.Error{Message: "Cannot use auth token"})
		return
	}

	message := dto.MessageContent{}
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		response.Respond(w, http.StatusBadRequest, dto.Error{Message: "Incorrect body"})
		return
	}
	defer r.Body.Close()

	err := handler.messagesService.SendSharedMessage(userId, message)
	if err != nil {
		response.Respond(w, http.StatusInternalServerError, dto.Error{Message: err.Error()})
		return
	}

	response.Respond(w, http.StatusOK, nil)
}

func (handler *MessagesHandler) privateChatGet(w http.ResponseWriter, r *http.Request) {
	user1, ok := r.Context().Value(userIdCtx).(string)
	if !ok {
		response.Respond(w, http.StatusInternalServerError, dto.Error{Message: "Cannot use auth token"})
		return
	}

	offset := r.URL.Query().Get("offset")
	if offset == "" {
		offset = "0"
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		response.Respond(w, http.StatusBadRequest, dto.Error{Message: err.Error()})
		return
	}

	count := r.URL.Query().Get("count")
	if count == "" {
		count = "50"
	}

	countInt, err := strconv.Atoi(count)
	if err != nil {
		response.Respond(w, http.StatusBadRequest, dto.Error{Message: err.Error()})
		return
	}

	user2 := chi.URLParam(r, "userId")
	messages, err := handler.messagesService.GetPrivateMessages(offsetInt, countInt, user1, user2)
	if err != nil {
		response.Respond(w, http.StatusInternalServerError, dto.Error{Message: err.Error()})
		return
	}

	response.Respond(w, http.StatusOK, messages)
}

func (handler *MessagesHandler) privateChatSendMessage(w http.ResponseWriter, r *http.Request) {
	user1, ok := r.Context().Value(userIdCtx).(string)
	if !ok {
		response.Respond(w, http.StatusInternalServerError, dto.Error{Message: "Cannot use auth token"})
		return
	}

	user2 := chi.URLParam(r, "userId")

	message := dto.MessageContent{}
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		response.Respond(w, http.StatusBadRequest, dto.Error{Message: "Incorrect body"})
		return
	}
	defer r.Body.Close()

	err := handler.messagesService.SendPrivateMessage(user1, user2, message)
	if err != nil {
		response.Respond(w, http.StatusInternalServerError, dto.Error{Message: err.Error()})
		return
	}

	response.Respond(w, http.StatusOK, nil)
}

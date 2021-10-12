package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/willsem/tfs-go-hw/hw04/internal/handlers"
	"github.com/willsem/tfs-go-hw/hw04/internal/repositories"
	"github.com/willsem/tfs-go-hw/hw04/internal/services"
)

func main() {
	logger := logrus.New()

	usersRepository := repositories.NewUsersRepository(logger)
	messagesRepository := repositories.NewMessagesRepository(logger)

	usersService := services.NewUsersService(usersRepository, logger)
	messagesService := services.NewMessagesService(messagesRepository, logger)

	r := chi.NewRouter()

	usersHandler := handlers.NewUsersHandler(usersService, logger)
	r.Mount("/users", usersHandler.Routes())

	messagesHandler := handlers.NewMessagesHandler(messagesService, logger)
	r.Mount("/messaages", messagesHandler.Routes())

	http.ListenAndServe(":5000", r)
}

package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/willsem/tfs-go-hw/course_project/internal/dto"
	"github.com/willsem/tfs-go-hw/course_project/internal/services/tradingbot"
	"github.com/willsem/tfs-go-hw/course_project/pkg/log"
	"github.com/willsem/tfs-go-hw/course_project/pkg/response"
)

type TradingBotHandler struct {
	bot    tradingbot.TradingBot
	logger log.Logger
}

func NewTradingBotHandler(bot tradingbot.TradingBot, logger log.Logger) *TradingBotHandler {
	return &TradingBotHandler{
		bot:    bot,
		logger: logger,
	}
}

func (handler *TradingBotHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/health", handler.HealthCheck)
	r.Post("/ticker/add", handler.addTicker)
	r.Post("/ticker/remove", handler.removeTicker)

	return r
}

func (handler *TradingBotHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response.Respond(w, http.StatusOK, dto.Message{Message: "works"})
}

func (handler *TradingBotHandler) addTicker(w http.ResponseWriter, r *http.Request) {
	ticker := dto.TickerName{}
	if err := json.NewDecoder(r.Body).Decode(&ticker); err != nil {
		response.Respond(w, http.StatusBadRequest, dto.Message{Message: "Incorrect body"})
		return
	}
	defer r.Body.Close()
}

func (handler *TradingBotHandler) removeTicker(w http.ResponseWriter, r *http.Request) {
	ticker := dto.TickerName{}
	if err := json.NewDecoder(r.Body).Decode(&ticker); err != nil {
		response.Respond(w, http.StatusBadRequest, dto.Message{Message: "Incorrect body"})
		return
	}
	defer r.Body.Close()
}

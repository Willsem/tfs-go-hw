package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/willsem/tfs-go-hw/course_project/internal/services/tradingbot"
	"github.com/willsem/tfs-go-hw/course_project/pkg/log"
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
	return r
}

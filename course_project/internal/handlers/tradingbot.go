package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/willsem/tfs-go-hw/course_project/internal/dto"
	"github.com/willsem/tfs-go-hw/course_project/pkg/log"
	"github.com/willsem/tfs-go-hw/course_project/pkg/response"
)

const (
	tickerContext = "ticker"
)

type TradingBotHandler struct {
	bot    TradingBot
	logger log.Logger
}

func NewTradingBotHandler(bot TradingBot, logger log.Logger) *TradingBotHandler {
	return &TradingBotHandler{
		bot:    bot,
		logger: logger,
	}
}

func (handler *TradingBotHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/health", handler.healthCheck)

	r.Put("/continue", handler.continueWork)
	r.Put("/pause", handler.pauseWork)

	r.Get("/tickers", handler.tickers)

	r.Route("/ticker/{ticker}", func(r chi.Router) {
		r.Use(handler.tickerContext)

		r.Post("/add", handler.addTicker)
		r.Delete("/remove", handler.removeTicker)
	})

	r.Route("/setup", func(r chi.Router) {
		r.Put("/size/{size}", handler.setupSize)
	})

	return r
}

func (handler *TradingBotHandler) healthCheck(w http.ResponseWriter, r *http.Request) {
	response.Respond(w, http.StatusOK, dto.Message{Message: "Works"})
}

func (handler *TradingBotHandler) continueWork(w http.ResponseWriter, r *http.Request) {
	if handler.bot.IsWorking() {
		response.Respond(w, http.StatusBadRequest, dto.Message{Message: "Bot is already working"})
		return
	}

	if err := handler.bot.Continue(); err != nil {
		response.Respond(w, http.StatusBadRequest, dto.Message{Message: err.Error()})
		return
	}

	response.Respond(w, http.StatusOK, dto.Message{Message: "Success"})
}

func (handler *TradingBotHandler) pauseWork(w http.ResponseWriter, r *http.Request) {
	if !handler.bot.IsWorking() {
		response.Respond(w, http.StatusBadRequest, dto.Message{Message: "Bot has been already paused"})
		return
	}

	if err := handler.bot.Pause(); err != nil {
		response.Respond(w, http.StatusBadRequest, dto.Message{Message: err.Error()})
		return
	}

	response.Respond(w, http.StatusOK, dto.Message{Message: "Success"})
}

func (handler *TradingBotHandler) tickers(w http.ResponseWriter, r *http.Request) {
	tickers := handler.bot.Tickers()
	response.Respond(w, http.StatusOK, tickers)
}

func (handler *TradingBotHandler) tickerContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ticker := chi.URLParam(r, tickerContext)
		ctx := context.WithValue(r.Context(), tickerContext, ticker)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (handler *TradingBotHandler) addTicker(w http.ResponseWriter, r *http.Request) {
	ticker, ok := r.Context().Value(tickerContext).(string)
	if !ok {
		response.Respond(w, http.StatusInternalServerError, dto.Message{Message: "Internal error"})
		return
	}

	if err := handler.bot.AddTicker(ticker); err != nil {
		response.Respond(w, http.StatusBadRequest, dto.Message{Message: err.Error()})
		return
	}

	response.Respond(w, http.StatusOK, dto.Message{Message: "Success"})
}

func (handler *TradingBotHandler) removeTicker(w http.ResponseWriter, r *http.Request) {
	ticker, ok := r.Context().Value(tickerContext).(string)
	if !ok {
		response.Respond(w, http.StatusInternalServerError, dto.Message{Message: "Internal error"})
		return
	}

	if err := handler.bot.RemoveTicker(ticker); err != nil {
		response.Respond(w, http.StatusBadRequest, dto.Message{Message: err.Error()})
		return
	}

	response.Respond(w, http.StatusOK, dto.Message{Message: "Success"})
}

func (handler *TradingBotHandler) setupSize(w http.ResponseWriter, r *http.Request) {
	size, err := strconv.ParseUint(chi.URLParam(r, "size"), 10, 64)
	if err != nil {
		response.Respond(w, http.StatusBadRequest, dto.Message{Message: err.Error()})
		return
	}

	handler.bot.ChangeSize(size)
	response.Respond(w, http.StatusOK, dto.Message{Message: "Success"})
}

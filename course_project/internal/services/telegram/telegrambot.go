package telegram

import (
	"strings"

	"github.com/willsem/tfs-go-hw/course_project/pkg/log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/willsem/tfs-go-hw/course_project/internal/config"
	"github.com/willsem/tfs-go-hw/course_project/internal/repositories/applications"
)

type Bot struct {
	applicationsRepository applications.ApplicationsRepository
	botAPI                 *tgbotapi.BotAPI
	logger                 log.Logger
}

func NewBot(
	repo applications.ApplicationsRepository,
	logger log.Logger,
	config config.Telegram,
) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(config.BotKey)
	if err != nil {
		return nil, err
	}

	return &Bot{
		applicationsRepository: repo,
		logger:                 logger,
		botAPI:                 bot,
	}, nil
}

func (bot *Bot) Start() {
	go bot.listen()
}

func (bot *Bot) listen() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.botAPI.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		text := update.Message.Text
		parsed := strings.Split(text, " ")

		var reply string
		switch {
		case parsed[0] == "/all":
			apps, err := bot.applicationsRepository.GetAll()
			if err != nil {
				reply = err.Error()
				bot.logger.Error(err)
			} else {
				reply = ""
				for _, app := range apps {
					reply += app.String()
				}
			}

		case parsed[0] == "/ticker" && len(parsed) == 2:
			apps, err := bot.applicationsRepository.GetByTicker(parsed[1])
			if err != nil {
				reply = err.Error()
				bot.logger.Error(err)
			} else {
				reply = ""
				for _, app := range apps {
					reply += app.String()
				}
			}

		case parsed[0] == "/help":
			reply = "/all - Получить список всех заявок\n"
			reply += "/ticker <ticker> - Получить список заявок по тикеру"

		default:
			reply = "Неизвестная команда, используйте /help"
		}

		if reply == "" {
			reply = "По данному запросу заявки не найдены"
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		_, err := bot.botAPI.Send(msg)
		if err != nil {
			bot.logger.Error(err)
		}
	}
}

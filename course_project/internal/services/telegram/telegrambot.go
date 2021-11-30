package telegram

import (
	"strings"

	"github.com/willsem/tfs-go-hw/course_project/pkg/log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/willsem/tfs-go-hw/course_project/internal/config"
)

const (
	loggerServiceName = "[TelegramBot]"
)

type BotImpl struct {
	applicationsRepository ApplicationsRepository
	botAPI                 *tgbotapi.BotAPI
	subscribed             []int64
	logger                 log.Logger
}

func NewBot(
	repo ApplicationsRepository,
	logger log.Logger,
	config config.Telegram,
) (*BotImpl, error) {
	bot, err := tgbotapi.NewBotAPI(config.BotKey)
	if err != nil {
		return nil, err
	}

	return &BotImpl{
		applicationsRepository: repo,
		logger:                 logger,
		botAPI:                 bot,
		subscribed:             make([]int64, 0),
	}, nil
}

func (bot *BotImpl) Start() {
	go bot.listen()
}

func (bot *BotImpl) SendSubscribedMessage(message string) {
	for _, id := range bot.subscribed {
		msg := tgbotapi.NewMessage(id, message)

		_, err := bot.botAPI.Send(msg)
		if err != nil {
			bot.logger.Error(loggerServiceName, err)
		}
	}
}

func (bot *BotImpl) listen() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.botAPI.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		text := update.Message.Text
		parsed := strings.Split(text, " ")

		var reply string = ""
		switch {
		case parsed[0] == "/subscribe":
			chatID := update.Message.Chat.ID
			for _, id := range bot.subscribed {
				if chatID == id {
					reply = "Вы уже подписаны"
					break
				}
			}

			if reply == "" {
				bot.subscribed = append(bot.subscribed, update.Message.Chat.ID)
				reply = "Подписка на заявки успешно совершена"
			}

		case parsed[0] == "/unsubscribe":
			chatID := update.Message.Chat.ID
			index := -1
			for i, id := range bot.subscribed {
				if chatID == id {
					index = i
					break
				}
			}

			if index != -1 {
				bot.subscribed = append(bot.subscribed[:index], bot.subscribed[index+1:]...)
				reply = "Отписка от заявок успешно совершена"
			} else {
				reply = "Вы не подписаны"
			}

		case parsed[0] == "/all":
			apps, err := bot.applicationsRepository.GetAll()
			if err != nil {
				reply = err.Error()
				bot.logger.Error(loggerServiceName, err)
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
				bot.logger.Error(loggerServiceName, err)
			} else {
				reply = ""
				for _, app := range apps {
					reply += app.String()
				}
			}

		case parsed[0] == "/start" || parsed[0] == "/help":
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
			bot.logger.Error(loggerServiceName, err)
		}
	}
}

package telegram

type Bot interface {
	SendSubscribedMessage(message string)
}

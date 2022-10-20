package message

import "github.com/vaberof/TelegramBotUniversitySchedule/internal/infra/storage"

type MessageReceiver interface {
	GetMessage(chatId int64) (*storage.Message, error)
}

type MessageSaver interface {
	SaveMessage(chatId int64, message string)
}

type MessageReceiverSaver interface {
	MessageReceiver
	MessageSaver
}

type MessageStorage struct {
	MessageReceiverSaver
}

func NewMessageStorage() *MessageStorage {
	return &MessageStorage{
		MessageReceiverSaver: storage.NewMessageStorage(),
	}
}

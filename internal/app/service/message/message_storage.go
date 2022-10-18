package message

import "github.com/vaberof/TelegramBotUniversitySchedule/internal/infra/storage"

type MessageStorageApi interface {
	GetMessage(chatId int64) (*storage.Message, error)
	SaveMessage(chatId int64, message string)
}

type MessageStorage struct {
	MessageStorageApi
}

func NewMessageStorage() *MessageStorage {
	return &MessageStorage{
		MessageStorageApi: storage.NewMessageStorage(),
	}
}

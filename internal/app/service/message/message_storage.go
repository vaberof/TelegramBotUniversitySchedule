package message

import (
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/infra/storage/postgres/messagepg"
	"gorm.io/gorm"
)

type MessageReceiver interface {
	GetMessage(chatId int64) (*string, error)
}

type MessageSaver interface {
	SaveMessage(chatId int64, message string) error
}

type MessageReceiverSaver interface {
	MessageReceiver
	MessageSaver
}

type MessageStorage struct {
	MessageReceiverSaver
}

func NewMessageStorage(db *gorm.DB) *MessageStorage {
	return &MessageStorage{
		MessageReceiverSaver: messagepg.NewMessageStoragePostgres(db),
	}
}

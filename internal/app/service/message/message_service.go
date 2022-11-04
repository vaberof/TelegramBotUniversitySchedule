package message

import (
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/infra/storage/postgres/messagepg"
)

type MessageService struct {
	messageStorage MessageStorage
}

func NewMessageService(messageStoragePostgres *messagepg.MessageStoragePostgres) *MessageService {
	return &MessageService{
		messageStorage: messageStoragePostgres,
	}
}

func (s *MessageService) GetMessage(chatId int64) (*string, error) {
	message, err := s.messageStorage.GetMessage(chatId)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (s *MessageService) SaveMessage(chatId int64, text string) error {
	return s.messageStorage.SaveMessage(chatId, text)
}

package message

import (
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/infra/storage/postgres/messagepg"
)

type MessageStorageService struct {
	messageStorage MessageStorage
}

func NewMessageStorageService(messageStoragePostgres *messagepg.MessageStoragePostgres) *MessageStorageService {
	return &MessageStorageService{
		messageStorage: messageStoragePostgres,
	}
}

func (s *MessageStorageService) GetMessage(chatId int64) (*string, error) {
	message, err := s.messageStorage.GetMessage(chatId)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (s *MessageStorageService) SaveMessage(chatId int64, text string) error {
	return s.messageStorage.SaveMessage(chatId, text)
}

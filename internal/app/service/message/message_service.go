package message

import "github.com/vaberof/TelegramBotUniversitySchedule/internal/infra/storage"

type MessageService struct {
	messageStorage *MessageStorage
}

func NewMessageService(messageStorage *MessageStorage) *MessageService {
	return &MessageService{
		messageStorage: messageStorage,
	}
}

func (s *MessageService) GetMessage(chatId int64) (*storage.Message, error) {
	message, err := s.messageStorage.GetMessage(chatId)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (s *MessageService) SaveMessage(chatId int64, message string) {
	s.messageStorage.SaveMessage(chatId, message)
}

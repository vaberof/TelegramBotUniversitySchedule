package message

import "github.com/vaberof/TelegramBotUniversitySchedule/internal/infra/storage"

type MessageService struct {
	messageStorageApi *MessageStorage
}

func NewMessageService(messageStorage *MessageStorage) *MessageService {
	return &MessageService{
		messageStorageApi: messageStorage,
	}
}

func (s *MessageService) GetMessage(chatId int64) (*storage.Message, error) {
	message, err := s.messageStorageApi.GetMessage(chatId)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (s *MessageService) SaveMessage(chatId int64, message string) {
	s.messageStorageApi.SaveMessage(chatId, message)
}

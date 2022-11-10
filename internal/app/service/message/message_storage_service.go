package message

type MessageStorageService struct {
	messageStorage MessageStorage
}

func NewMessageStorageService(messageStorage MessageStorage) *MessageStorageService {
	return &MessageStorageService{
		messageStorage: messageStorage,
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

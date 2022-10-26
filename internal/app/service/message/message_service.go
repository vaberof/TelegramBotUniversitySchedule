package message

type MessageService struct {
	messageStorage *MessageStorage
}

func NewMessageService(messageStorage *MessageStorage) *MessageService {
	return &MessageService{
		messageStorage: messageStorage,
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

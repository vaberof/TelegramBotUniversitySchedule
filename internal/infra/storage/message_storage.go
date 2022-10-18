package storage

import "errors"

type ChatId int64
type Message string

type MessageStorage struct {
	messageStorage map[ChatId]*Message
}

func NewMessageStorage() *MessageStorage {
	return &MessageStorage{
		messageStorage: map[ChatId]*Message{},
	}
}

func (u *MessageStorage) GetMessage(chatId int64) (*Message, error) {
	if u.messageStorage[ChatId(chatId)] == nil {
		return nil, errors.New("message does not exists")
	}
	message := u.messageStorage[ChatId(chatId)]
	return message, nil
}

func (u *MessageStorage) SaveMessage(chatId int64, message string) {
	msg := Message(message)
	u.messageStorage[ChatId(chatId)] = &msg
}

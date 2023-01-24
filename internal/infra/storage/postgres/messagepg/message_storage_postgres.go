package messagepg

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type MessageStoragePostgres struct {
	db *gorm.DB
}

func NewMessageStoragePostgres(db *gorm.DB) *MessageStoragePostgres {
	return &MessageStoragePostgres{db: db}
}

func (m *MessageStoragePostgres) GetMessage(chatId int64) (*string, error) {
	message, err := m.getMessageFromDb(chatId)
	if err != nil {
		return nil, err
	}
	return &message.Text, nil
}

func (m *MessageStoragePostgres) SaveMessage(chatId int64, text string) error {
	messageFromDb, err := m.getMessageFromDb(chatId)
	if err == nil {
		return m.updateMessageInDb(messageFromDb, text)
	}

	messageToSave := &Message{
		ChatId: chatId,
		Text:   text,
	}
	return m.db.Create(&messageToSave).Error
}

func (m *MessageStoragePostgres) getMessageFromDb(chatId int64) (*Message, error) {
	var messageStorage Message

	err := m.db.Table("messages").Where("chat_id = ?", chatId).First(&messageStorage).Error
	if err != nil {
		log.Error("message not found in db")
		return nil, errors.New("message not saved")
	}
	return &messageStorage, nil
}

func (m *MessageStoragePostgres) updateMessageInDb(message *Message, text string) error {
	message.Text = text

	err := m.db.Save(&message).Error
	if err != nil {
		log.Printf("cannot update message '%v' in database, error: %v", message, err)
		return err
	}
	log.Printf("message '%v' updated", message)
	return nil
}

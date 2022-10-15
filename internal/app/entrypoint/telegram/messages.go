package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

type InputTelegramMessage struct {
	ChatId  int64
	Message string
}

// HandleNewMessage
// adds user`s chat id and his input group id to storage.MessageStorage
// and sends reply message with buttons to press to get schedule.
func (h *TelegramHandler) HandleNewMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, keyboard tgbotapi.InlineKeyboardMarkup) *InputTelegramMessage {
	responseMessage := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

	inputMsgText := responseMessage.Text

	responseMessage.ReplyMarkup = keyboard
	bot.Send(responseMessage)

	log.WithFields(log.Fields{
		"username": update.SentFrom(),
		"message":  inputMsgText,
	}).Info("User sent a message")

	return &InputTelegramMessage{
		ChatId:  responseMessage.ChatID,
		Message: inputMsgText,
	}
}

// MessageReceived checks if user sent a message.
func (h *TelegramHandler) MessageReceived(update tgbotapi.Update) bool {
	return update.Message != nil
}

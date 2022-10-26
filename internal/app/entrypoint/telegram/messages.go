package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

func (h *TelegramHandler) HandleNewMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, keyboard tgbotapi.InlineKeyboardMarkup) {
	responseMessage := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	inputMessageText := responseMessage.Text
	chatId := responseMessage.ChatID

	err := h.MessageReceiverSaver.SaveMessage(chatId, inputMessageText)
	if err != nil {
		log.Error("cannot save input message, error: ", err.Error())
		responseMessage.Text = err.Error()
		bot.Send(responseMessage)
		return
	}

	responseMessage.ReplyMarkup = keyboard
	bot.Send(responseMessage)

	log.WithFields(log.Fields{
		"username": update.SentFrom(),
		"message":  inputMessageText,
	}).Info("User sent a message")

	log.WithFields(log.Fields{
		"chatId":  chatId,
		"message": inputMessageText,
	}).Info("message is saved")
}

func (h *TelegramHandler) MessageReceived(update tgbotapi.Update) bool {
	return update.Message != nil
}

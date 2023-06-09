package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

func (h *TelegramHandler) HandleNewMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, keyboard tgbotapi.InlineKeyboardMarkup) {
	responseMessage := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

	var chatId int64
	var username string
	var inputMessageText string

	chatId = responseMessage.ChatID
	inputMessageText = responseMessage.Text

	user := update.SentFrom()
	if user != nil {
		username = user.UserName
	} else {
		logrus.Error("cannot get username from chat_id: %d", chatId)
	}

	logrus.WithFields(logrus.Fields{
		"username": username,
		"message":  inputMessageText,
	}).Info("User sent a message")

	err := h.messageService.SaveMessage(chatId, username, inputMessageText)
	if err != nil {
		logrus.Error("cannot save input message, error: ", err.Error())
		responseMessage.Text = err.Error()
		bot.Send(responseMessage)
		return
	}

	responseMessage.ReplyMarkup = keyboard
	bot.Send(responseMessage)

	logrus.WithFields(logrus.Fields{
		"chatId":   chatId,
		"username": username,
		"message":  inputMessageText,
	}).Info("message is saved")
}

func (h *TelegramHandler) MessageReceived(update tgbotapi.Update) bool {
	return update.Message != nil
}

package handler

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type TelegramHandler interface {
	HandleCommandMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update)
	CommandReceived(update tgbotapi.Update) bool
	HandleNewMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, keyboard tgbotapi.InlineKeyboardMarkup)
	MessageReceived(update tgbotapi.Update) bool
	HandleMenuButtonPress(bot *tgbotapi.BotAPI, update tgbotapi.Update, keyboard tgbotapi.InlineKeyboardMarkup)
	MenuButtonPressed(callBackQuery tgbotapi.Update) bool
}

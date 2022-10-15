package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// HandleCommandMessage handles received command and sends corresponding message to user.
func (h *TelegramHandler) HandleCommandMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	responseMessage := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	startCommandText := "Как пользоваться ботом:\n" +
		"1. Введите номер группы (БИ-11.1/БИ-11.2 и т.д.)\n" +
		"2. Выберите день, на который хотите получить расписание\n"

	defaultMessageText := "Неизвестная команда"

	switch update.Message.Command() {
	case "start":
		responseMessage.Text = startCommandText
		bot.Send(responseMessage)
	default:
		responseMessage.Text = defaultMessageText
		bot.Send(responseMessage)
	}
}

// CommandReceived checks if user sent a command.
func (h *TelegramHandler) CommandReceived(update tgbotapi.Update) bool {
	if update.Message == nil {
		return false
	}
	return update.Message.IsCommand()
}
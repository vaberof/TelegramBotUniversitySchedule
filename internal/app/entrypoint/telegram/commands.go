package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const startCommandMessageOutput string = "Как пользоваться ботом:\n" +
	"1. Введите номер группы (БИ-11.1/БИ-11.2 и т.д.)\n" +
	"2. Выберите день, на который хотите получить расписание\n"

const unknownCommandMessageOutput string = "Неизвестная команда"

func (h *TelegramHandler) HandleCommandMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	responseMessage := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	switch update.Message.Command() {
	case "start":
		responseMessage.Text = startCommandMessageOutput
		bot.Send(responseMessage)
	default:
		responseMessage.Text = unknownCommandMessageOutput
		bot.Send(responseMessage)
	}
}

func (h *TelegramHandler) CommandReceived(update tgbotapi.Update) bool {
	if update.Message == nil {
		return false
	}
	return update.Message.IsCommand()
}

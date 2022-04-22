package bot

import (
	"github.com/tg_bot_timetable/internal/controller"
	"github.com/tg_bot_timetable/internal/model"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Start starts the bot.
func Start() {

	var (
		inlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Сегодня", "Сегодня"),
				tgbotapi.NewInlineKeyboardButtonData("Завтра", "Завтра"),
			),
		)
		responseMessage tgbotapi.MessageConfig
	)

	bot := NewBot()

	channel := tgbotapi.NewUpdate(0)
	channel.Timeout = 60

	updates := bot.GetUpdatesChan(channel)

	location := model.SetLocation()
	groupStorage := model.NewGroupStorage()

	for update := range updates {

		if update.CallbackQuery != nil {

			responseCallback := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			inputCallback := responseCallback.Text
			responseCallback.ReplyMarkup = inlineKeyboard
			responseCallback.Text = *controller.HandleMessage(groupStorage, location, inputCallback, responseMessage.Text)
			bot.Send(responseCallback)
			continue
		}

		if update.Message == nil {
			continue
		}

		responseMessage = tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		responseMessage.ReplyMarkup = inlineKeyboard
		bot.Send(responseMessage)
	}
}
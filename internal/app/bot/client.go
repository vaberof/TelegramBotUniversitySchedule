package bot

import (
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/controller"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/model"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Start starts the bot.
func Start() {

	var (
		inlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Сегодня", "Сегодня"),
				tgbotapi.NewInlineKeyboardButtonData("Завтра", "Завтра"),
				tgbotapi.NewInlineKeyboardButtonData("Неделя", "Неделя"),
			),
		)
	)

	bot := NewBot()

	channel := tgbotapi.NewUpdate(0)
	channel.Timeout = 60

	updates := bot.GetUpdatesChan(channel)

	user := model.NewUser()
	groupStorage := model.NewGroupStorage()
	location := model.SetLocation()

	for update := range updates {

		if update.CallbackQuery != nil {

			responseCallback := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)

			inputCallback := responseCallback.Text
			callbackChatID := update.CallbackQuery.Message.Chat.ID

			responseCallback.ReplyMarkup = inlineKeyboard

			responseCallback.Text = *controller.HandleMessage(callbackChatID, inputCallback, user, groupStorage, location)

			bot.Send(responseCallback)
			continue
		}

		if update.Message == nil {
			continue
		}

		responseMessage := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		inputText := responseMessage.Text
		chatID := responseMessage.ChatID

		user.AddData(chatID, inputText)

		responseMessage.ReplyMarkup = inlineKeyboard
		bot.Send(responseMessage)
	}
}

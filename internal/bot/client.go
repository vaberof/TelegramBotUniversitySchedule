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
		inputText string
		chatID int64
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

			inputCallback := responseCallback.Text                      // text of the button that user pressed
			callbackChatID := update.CallbackQuery.Message.Chat.ID      // chat id where user pressed the button

			responseCallback.ReplyMarkup = inlineKeyboard

			responseCallback.Text = *controller.HandleMessage(callbackChatID, inputCallback, user, groupStorage, location)

			bot.Send(responseCallback)
			continue
		}

		if update.Message == nil {
			continue
		}

		responseMessage = tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		inputText = responseMessage.Text    // user input
		chatID = responseMessage.ChatID     // user chat id

		user.AddUser(chatID, inputText)

		responseMessage.ReplyMarkup = inlineKeyboard
		bot.Send(responseMessage)
	}
}

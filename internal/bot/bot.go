package bot

import (
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/handler"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/service"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/storage"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/integration/unisite"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Start starts the bot.
func Start() {

	botKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Сегодня", "Today"),
			tgbotapi.NewInlineKeyboardButtonData("Завтра", "Tomorrow"),
			tgbotapi.NewInlineKeyboardButtonData("Неделя", "Week"),
			tgbotapi.NewInlineKeyboardButtonData("След. неделя", "Next week"),
		),
	)

	bot := newBot()

	botUpdatesChannel := tgbotapi.NewUpdate(0)
	botUpdatesChannel.Timeout = 60

	updates := bot.GetUpdatesChan(botUpdatesChannel)

	messageStorage := storage.NewMessageStorage()
	groupStorage := storage.NewGroupStorage()

	for update := range updates {
		if menuButtonPressed(update) {
			handleMenuButtonPress(bot, update, botKeyboardMarkup, messageStorage, groupStorage)
		} else if commandReceived(update) {
			handleCommandMessage(bot, update)
		} else if messageReceived(update) {
			handleNewMessage(bot, update, botKeyboardMarkup, messageStorage)
		}
	}
}

// newBot creates bot.
func newBot() *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Bot %s is authorized.", bot.Self.UserName)

	return bot
}

// messageReceived checks if user sent a message.
func messageReceived(update tgbotapi.Update) bool {
	return update.Message != nil
}

// handleNewMessage
// adds user`s chat id and his input group id to storage.MessageStorage
// and sends reply message with buttons to press to get schedule.
func handleNewMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, keyboard tgbotapi.InlineKeyboardMarkup, msgStorage *storage.MessageStorage) {
	responseMessage := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

	inputText := responseMessage.Text
	chatID := responseMessage.ChatID

	msgStorage.AddMessageData(chatID, inputText)

	responseMessage.ReplyMarkup = keyboard
	bot.Send(responseMessage)
}

// menuButtonPressed checks if user pressed the button of replied message to him.
func menuButtonPressed(callBackQuery tgbotapi.Update) bool {
	return callBackQuery.CallbackQuery != nil
}

// handleMenuButtonPress handles pressed button value (today/tomorrow/week/next week)
// and sending a schedule for utils that user chosen.
// if user`s input group id is not exists, then sends a corresponding message.
func handleMenuButtonPress(bot *tgbotapi.BotAPI, update tgbotapi.Update, keyboard tgbotapi.InlineKeyboardMarkup, msgStorage *storage.MessageStorage, grpStorage *storage.GroupStorage) {
	responseCallback := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)

	inputCallback := responseCallback.Text
	callbackChatID := update.CallbackQuery.Message.Chat.ID
	responseCallback.ReplyMarkup = keyboard

	log.Printf("%s\n", inputCallback)

	studyGroupId, studyGroupUrl, err := handler.HandleMessage(callbackChatID, msgStorage, grpStorage)
	if err != nil {
		responseCallback.Text = err.Error()
		bot.Send(responseCallback)
		return
	}

	schedule := unisite.GetSchedule(studyGroupUrl, inputCallback)
	scheduleString := service.ScheduleToString(studyGroupId, inputCallback, schedule)

	responseCallback.Text = *scheduleString
	responseCallback.ParseMode = "markdown"
	bot.Send(responseCallback)
}

// commandReceived checks if user sent a command.
func commandReceived(update tgbotapi.Update) bool {
	return update.Message.IsCommand()
}

// handleCommandMessage handles received command and sends corresponding message to user.
func handleCommandMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	responseMessage := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	switch update.Message.Command() {
	case "start":
		responseMessage.Text = "Как пользоваться ботом:\n" +
			"1. Введите номер группы (БИ-11.1/БИ-11.2 и т.д.)\n" +
			"2. Выберите день, на который хотите получить расписание\n"
		bot.Send(responseMessage)
	default:
		responseMessage.Text = "Не знаю такой команды"
		bot.Send(responseMessage)
	}
}

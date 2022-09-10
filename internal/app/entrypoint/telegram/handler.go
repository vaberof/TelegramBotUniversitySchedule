package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"github.com/vaberof/TelegramBotUniversitySchedule/pkg/xtime"
	"github.com/vaberof/TelegramBotUniversitySchedule/pkg/xtimezone"
	"time"
)

type InputMessage struct {
	ChatId  int64
	Message string
}

// MessageReceived checks if user sent a message.
func MessageReceived(update tgbotapi.Update) bool {
	return update.Message != nil
}

// HandleNewMessage
// adds user`s chat id and his input group id to storage.MessageStorage
// and sends reply message with buttons to press to get schedule.
func HandleNewMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, keyboard tgbotapi.InlineKeyboardMarkup) *InputMessage {
	responseMessage := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

	inputText := responseMessage.Text

	responseMessage.ReplyMarkup = keyboard
	bot.Send(responseMessage)

	log.WithFields(log.Fields{
		"username": update.SentFrom(),
		"message":  inputText,
	}).Info("User sent a message")

	return &InputMessage{
		ChatId:  responseMessage.ChatID,
		Message: inputText,
	}
}

// MenuButtonPressed checks if user pressed the button of replied message to him.
func MenuButtonPressed(callBackQuery tgbotapi.Update) bool {
	return callBackQuery.CallbackQuery != nil
}

// HandleMenuButtonPress handles pressed button value (today/tomorrow/week/next week)
// and sending a schedule for date that user chosen.
// if user`s input group id is not exists, then sends a corresponding message.
func HandleMenuButtonPress(
	bot *tgbotapi.BotAPI,
	update tgbotapi.Update,
	inputMessage *InputMessage,
	keyboard tgbotapi.InlineKeyboardMarkup) {

	responseCallback := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)

	inputTelegramButton := responseCallback.Text

	log.WithFields(log.Fields{
		"username": update.SentFrom(),
		"button":   inputTelegramButton,
	}).Info("User requested a schedule")

	fmt.Println(time.Now().In(xtime.GetDefaultLocation(xtimezone.Novosibirsk)).Add(time.Hour*168), time.Now().In(xtime.GetDefaultLocation(xtimezone.Novosibirsk)).Add(time.Hour*336))
	//fromDate, toDate := parseDatesRange(inputTelegramButton)

	//	schedule, err := integration.ScheduleApi.GetSchedule(groupId, fromDate, toDate)
	//
	//scheduleString := ScheduleToString(inputGroupId, inputTelegramButton, schedule)

	responseCallback.ReplyMarkup = keyboard
	responseCallback.Text = "schedule"
	responseCallback.ParseMode = "markdown"
	bot.Send(responseCallback)
}

// CommandReceived checks if user sent a command.
func CommandReceived(update tgbotapi.Update) bool {
	return update.Message.IsCommand()
}

// HandleCommandMessage handles received command and sends corresponding message to user.
func HandleCommandMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
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

func parseDatesRange(inputTelegramButtonData string) (time.Time, time.Time) {
	switch inputTelegramButtonData {
	case "Today":
		return time.Now().In(xtime.GetDefaultLocation(xtimezone.Novosibirsk)), time.Now().In(xtime.GetDefaultLocation(xtimezone.Novosibirsk))
	case "Tomorrow":
		return time.Now().In(xtime.GetDefaultLocation(xtimezone.Novosibirsk)), time.Now().In(xtime.GetDefaultLocation(xtimezone.Novosibirsk)).Add(time.Hour * 24)
	case "Week":
		return time.Now().In(xtime.GetDefaultLocation(xtimezone.Novosibirsk)), time.Now().In(xtime.GetDefaultLocation(xtimezone.Novosibirsk)).Add(time.Hour * 168)
	default:
		return time.Now().In(xtime.GetDefaultLocation(xtimezone.Novosibirsk)).Add(time.Hour * 168), time.Now().In(xtime.GetDefaultLocation(xtimezone.Novosibirsk)).Add(time.Hour * 336)
	}
}

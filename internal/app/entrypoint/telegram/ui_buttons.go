package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/pkg/xstrconv"
	"github.com/vaberof/TelegramBotUniversitySchedule/pkg/xtime"
	"time"
)

// MenuButtonPressed checks if user pressed the button of replied message to him.
func (h *TelegramHandler) MenuButtonPressed(callBackQuery tgbotapi.Update) bool {
	return callBackQuery.CallbackQuery != nil
}

// HandleMenuButtonPress handles pressed button value (today/tomorrow/week/next week)
// and sending a schedule for date that user chosen.
// if user`s input group id is not exists, then sends a corresponding message.
func (h *TelegramHandler) HandleMenuButtonPress(
	bot *tgbotapi.BotAPI,
	update tgbotapi.Update,
	inputTelegramMessage *InputTelegramMessage,
	keyboard tgbotapi.InlineKeyboardMarkup) {

	responseCallback := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
	inputTelegramButtonName := responseCallback.Text

	responseCallback.ParseMode = "markdown"

	log.WithFields(log.Fields{
		"username": update.SentFrom(),
		"button":   inputTelegramButtonName,
	}).Info("User requested a schedule")

	if inputTelegramMessage == nil {
		responseCallback.Text = "Введите номер группы"
		bot.Send(responseCallback)
		return
	}

	groupId := inputTelegramMessage.Message
	fromDate, toDate := parseDatesRange(inputTelegramButtonName)
	//log.Printf("fromDate: %v, toDate: %v\n", fromDate, toDate)

	schedule, err := h.ScheduleReceiver.GetSchedule(groupId, fromDate, toDate)
	if err != nil {
		responseCallback.Text = err.Error()
		bot.Send(responseCallback)
		return
	}

	scheduleString := xstrconv.ScheduleToString(groupId, inputTelegramButtonName, schedule)

	//log.Printf("schedule string: %v\n", scheduleString)
	//log.Printf("schedule string length: %v\n", len(*scheduleString))

	responseCallback.Text = *scheduleString
	responseCallback.ReplyMarkup = keyboard
	_, err = bot.Send(responseCallback)
	if err != nil {
		log.Panic(err.Error())
	}
}

func parseDatesRange(inputTelegramButtonData string) (time.Time, time.Time) {
	switch inputTelegramButtonData {
	case "Today":
		today := xtime.GetDateToParse(inputTelegramButtonData)
		return today, today
	case "Tomorrow":
		tomorrow := xtime.GetDateToParse(inputTelegramButtonData)
		today := tomorrow.Add(-time.Hour * 24)
		return today, tomorrow
	case "Week":
		currentWeekDays := xtime.GetDatesToParse(inputTelegramButtonData)
		return currentWeekDays[0], currentWeekDays[6]
	default:
		nextWeekDates := xtime.GetDatesToParse(inputTelegramButtonData)
		return nextWeekDates[0], nextWeekDates[6]
	}
}

package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/pkg/xstrconv"
	"github.com/vaberof/TelegramBotUniversitySchedule/pkg/xtime"
)

const errorMessageToTelegram string = "Ошибка: невозможно получить расписание"

func (h *TelegramHandler) MenuButtonPressed(callBackQuery tgbotapi.Update) bool {
	return callBackQuery.CallbackQuery != nil
}

func (h *TelegramHandler) HandleMenuButtonPress(
	bot *tgbotapi.BotAPI,
	update tgbotapi.Update,
	keyboard tgbotapi.InlineKeyboardMarkup) {

	responseCallback := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
	inputTelegramButtonDate := responseCallback.Text

	responseCallback.ParseMode = "markdown"

	log.WithFields(log.Fields{
		"username": update.SentFrom(),
		"button":   inputTelegramButtonDate,
	}).Info("User requested a schedule")

	inputTelegramMessage, err := h.GetMessage(responseCallback.ChatID)
	if inputTelegramMessage == nil || err != nil {
		responseCallback.Text = "Введите номер группы"
		bot.Send(responseCallback)
		return
	}

	groupId := *inputTelegramMessage
	convGroupId := string(groupId)

	fromDate, toDate, err := xtime.ParseDatesRange(inputTelegramButtonDate)
	if err != nil {
		log.WithFields(log.Fields{
			"fromDate": fromDate,
			"toDate":   toDate,
			"func":     "HandleMenuButtonPress",
			"error":    err.Error(),
		}).Error("Cannot parse dates range")

		responseCallback.Text = errorMessageToTelegram
		bot.Send(responseCallback)
		return
	}

	schedule, err := h.ScheduleReceiver.GetSchedule(convGroupId, fromDate, toDate)
	if err != nil {
		log.WithFields(log.Fields{
			"schedule": schedule,
			"func":     "HandleMenuButtonPress",
			"error":    err.Error(),
		}).Error("Cannot parse dates range")

		responseCallback.Text = err.Error()
		bot.Send(responseCallback)
		return
	}

	scheduleString, err := xstrconv.ScheduleToString(convGroupId, inputTelegramButtonDate, schedule)
	if scheduleString == nil || err != nil {
		log.WithFields(log.Fields{
			"schedule string": scheduleString,
			"error":           err.Error(),
			"func":            "HandleMenuButtonPress",
		}).Error("Cannot get schedule string")

		responseCallback.Text = errorMessageToTelegram
		bot.Send(responseCallback)
		return
	}

	responseCallback.Text = *scheduleString
	responseCallback.ReplyMarkup = keyboard

	_, err = bot.Send(responseCallback)
	if err != nil {
		log.Panic(err.Error())
	}
}

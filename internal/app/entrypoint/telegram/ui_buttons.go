package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	domain "github.com/vaberof/TelegramBotUniversitySchedule/internal/domain/schedule"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/pkg/xstrconv"
	"github.com/vaberof/TelegramBotUniversitySchedule/pkg/xtime"
	"time"
)

const errorMessageToTelegram string = "Ошибка: невозможно получить расписание"

func (h *TelegramHandler) HandleMenuButtonPress(
	bot *tgbotapi.BotAPI,
	update tgbotapi.Update,
	keyboard tgbotapi.InlineKeyboardMarkup) {

	responseCallback := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
	responseCallback.ParseMode = "markdown"

	inputTelegramButtonDate := responseCallback.Text

	log.WithFields(log.Fields{"username": update.SentFrom(), "button": inputTelegramButtonDate}).
		Info("User requested a schedule")

	inputTelegramMessage := h.getMessage(responseCallback.ChatID, responseCallback, bot)
	if inputTelegramMessage == nil {
		return
	}

	fromDate, toDate, err := h.parseDatesRange(inputTelegramButtonDate, responseCallback, bot)
	if err != nil {
		return
	}

	groupId := *inputTelegramMessage

	schedule, err := h.getSchedule(groupId, fromDate, toDate, responseCallback, bot)
	if err != nil {
		return
	}

	scheduleString, err := h.scheduleToString(groupId, inputTelegramButtonDate, schedule, responseCallback, bot)
	if scheduleString == nil || err != nil {
		return
	}

	responseCallback.Text = *scheduleString
	responseCallback.ReplyMarkup = keyboard

	bot.Send(responseCallback)
}

func (h *TelegramHandler) MenuButtonPressed(callBackQuery tgbotapi.Update) bool {
	return callBackQuery.CallbackQuery != nil
}

func (h *TelegramHandler) getMessage(chatId int64, responseCallback tgbotapi.MessageConfig, bot *tgbotapi.BotAPI) *string {
	inputTelegramMessage, err := h.messageStorage.GetMessage(chatId)
	if inputTelegramMessage == nil || err != nil {
		responseCallback.Text = "Введите номер группы"
		bot.Send(responseCallback)
		return nil
	}
	return inputTelegramMessage
}

func (h *TelegramHandler) parseDatesRange(
	inputTelegramButtonDate string,
	responseCallback tgbotapi.MessageConfig,
	bot *tgbotapi.BotAPI) (time.Time, time.Time, error) {

	fromDate, toDate, err := xtime.ParseDatesRange(inputTelegramButtonDate)
	if err != nil {

		log.WithFields(log.Fields{
			"fromDate": fromDate,
			"toDate":   toDate,
			"func":     "HandleMenuButtonPress",
			"error":    err.Error()}).
			Error("Cannot parse dates range")

		responseCallback.Text = errorMessageToTelegram
		bot.Send(responseCallback)
		return time.Time{}, time.Time{}, err
	}
	return fromDate, toDate, nil
}

func (h *TelegramHandler) getSchedule(
	groupId string,
	from time.Time,
	to time.Time,
	responseCallback tgbotapi.MessageConfig,
	bot *tgbotapi.BotAPI) (*domain.Schedule, error) {

	schedule, err := h.scheduleReceiver.GetSchedule(groupId, from, to)
	if err != nil {
		log.WithFields(log.Fields{
			"schedule": schedule,
			"func":     "HandleMenuButtonPress",
			"error":    err.Error(),
		}).Error("Cannot get schedule")

		responseCallback.Text = err.Error()
		bot.Send(responseCallback)
		return nil, err
	}
	return schedule, nil
}

func (h *TelegramHandler) scheduleToString(
	groupId string,
	inputTelegramButtonDate string,
	schedule *domain.Schedule,
	responseCallback tgbotapi.MessageConfig,
	bot *tgbotapi.BotAPI) (*string, error) {

	scheduleString, err := xstrconv.ScheduleToString(groupId, inputTelegramButtonDate, schedule)
	if scheduleString == nil || err != nil {
		log.WithFields(log.Fields{
			"schedule string": scheduleString,
			"error":           err,
			"func":            "HandleMenuButtonPress",
		}).Error("Cannot get schedule string")

		responseCallback.Text = errorMessageToTelegram
		bot.Send(responseCallback)
		return nil, err
	}
	return scheduleString, nil
}

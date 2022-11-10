package telegram

import (
	domain "github.com/vaberof/TelegramBotUniversitySchedule/internal/domain/schedule"
	"time"
)

type ScheduleReceiver interface {
	GetSchedule(groupId string, from time.Time, to time.Time) (*domain.Schedule, error)
}

type MessageStorage interface {
	GetMessage(chatId int64) (*string, error)
	SaveMessage(chatId int64, text string) error
}

type TelegramHandler struct {
	scheduleReceiver ScheduleReceiver
	messageStorage   MessageStorage
}

func NewTelegramHandler(scheduleReceiver ScheduleReceiver, messageStorage MessageStorage) *TelegramHandler {
	return &TelegramHandler{
		scheduleReceiver: scheduleReceiver,
		messageStorage:   messageStorage,
	}
}

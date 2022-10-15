package telegram

import (
	domain "github.com/vaberof/TelegramBotUniversitySchedule/internal/domain/schedule"
	"time"
)

type ScheduleReceiver interface {
	GetSchedule(groupId string, from time.Time, to time.Time) (*domain.Schedule, error)
}

type TelegramHandler struct {
	ScheduleReceiver
}

func NewTelegramHandler(scheduleService *domain.ScheduleService) *TelegramHandler {
	return &TelegramHandler{ScheduleReceiver: scheduleService}
}

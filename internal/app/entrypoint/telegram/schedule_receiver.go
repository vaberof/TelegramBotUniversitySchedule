package telegram

import (
	domain "github.com/vaberof/TelegramBotUniversitySchedule/internal/domain/schedule"
	"time"
)

type ScheduleReceiver interface {
	GetSchedule(groupId string, from time.Time, to time.Time) (domain.Schedule, error)
}

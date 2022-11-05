package domain

import (
	infra "github.com/vaberof/TelegramBotUniversitySchedule/internal/infra/integration/unisite"
	"time"
)

type ScheduleApi interface {
	GetSchedule(groupId string, from time.Time, to time.Time) (*infra.GetScheduleResponse, error)
}

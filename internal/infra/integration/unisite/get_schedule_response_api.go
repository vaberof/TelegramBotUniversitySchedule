package unisite

import (
	integration "github.com/vaberof/TelegramBotUniversitySchedule/pkg/integration/unisite"
	"time"
)

type ScheduleApi interface {
	GetSchedule(groupExternalId string, from time.Time, to time.Time) (*integration.GetScheduleResponse, error)
}

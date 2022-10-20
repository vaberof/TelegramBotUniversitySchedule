package domain

import (
	infra "github.com/vaberof/TelegramBotUniversitySchedule/internal/infra/integration/unisite"
	"time"
)

type GetScheduleResponseReceiver interface {
	GetSchedule(groupId string, from time.Time, to time.Time) (*infra.GetScheduleResponse, error)
}

type ScheduleApi struct {
	GetScheduleResponseReceiver
}

func NewScheduleApi(getScheduleResponseService *infra.GetScheduleResponseService) *ScheduleApi {
	return &ScheduleApi{
		GetScheduleResponseReceiver: getScheduleResponseService,
	}
}

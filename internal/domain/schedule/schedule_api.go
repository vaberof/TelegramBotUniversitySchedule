package domain

import (
	infra "github.com/vaberof/TelegramBotUniversitySchedule/internal/infra/integration/unisite"
	"time"
)

type ScheduleApi interface {
	GetSchedule(studyGroupQueryParams string, from time.Time, to time.Time) (*infra.GetScheduleResponse, error)
}

type GetScheduleResponseApi struct {
	ScheduleApi
}

func NewGetScheduleResponseApi(scheduleApi *infra.GetScheduleResponseApiService) *GetScheduleResponseApi {
	return &GetScheduleResponseApi{
		ScheduleApi: scheduleApi,
	}
}

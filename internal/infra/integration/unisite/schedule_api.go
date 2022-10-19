package infra

import (
	"github.com/vaberof/TelegramBotUniversitySchedule/configs"
	integration "github.com/vaberof/TelegramBotUniversitySchedule/pkg/integration/unisite"
	"time"
)

type ScheduleApi interface {
	GetSchedule(studyGroupQueryParams string, from time.Time, to time.Time) (*integration.GetScheduleResponse, error)
}

type GetScheduleResponseApi struct {
	ScheduleApi
}

func NewGetScheduleResponseApi(host string, config *configs.HttpClientConfig) *GetScheduleResponseApi {
	return &GetScheduleResponseApi{
		ScheduleApi: integration.NewHttpClient(host, config),
	}
}

package infra

import (
	"github.com/vaberof/TelegramBotUniversitySchedule/configs"
	integration "github.com/vaberof/TelegramBotUniversitySchedule/pkg/integration/unisite"
	"time"
)

type ScheduleApi interface {
	GetSchedule(groupExternalId string, from time.Time, to time.Time) (*integration.GetScheduleResponse, error)
}

type GetScheduleResponseApi struct {
	ScheduleApi
}

func NewGetScheduleResponseApi(config *configs.HttpClientConfig) *GetScheduleResponseApi {
	return &GetScheduleResponseApi{
		ScheduleApi: integration.NewHttpClient(config),
	}
}

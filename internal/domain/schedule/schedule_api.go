package domain

import (
	infra "github.com/vaberof/TelegramBotUniversitySchedule/internal/infra/integration/unisite"
	"time"
)

type GetScheduleResponseReceiver interface {
	GetSchedule(groupId string, from time.Time, to time.Time) (*infra.GetScheduleResponse, error)
}

type GetScheduleResponse struct {
	GetScheduleResponseReceiver
}

func NewGetScheduleResponse(getScheduleResponseService *infra.GetScheduleResponseService) *GetScheduleResponse {
	return &GetScheduleResponse{
		GetScheduleResponseReceiver: getScheduleResponseService,
	}
}

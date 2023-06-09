package unisite

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	domain "github.com/vaberof/TelegramBotUniversitySchedule/internal/domain/schedule"
	"time"
)

type GetScheduleResponseService struct {
	scheduleApi  ScheduleApi
	groupStorage GroupStorage
}

func NewGetScheduleResponseService(scheduleApi ScheduleApi, groupStorage GroupStorage) *GetScheduleResponseService {
	return &GetScheduleResponseService{
		scheduleApi:  scheduleApi,
		groupStorage: groupStorage,
	}
}

func (s *GetScheduleResponseService) GetSchedule(groupId string, from time.Time, to time.Time) (domain.Schedule, error) {
	groupExternalId := s.groupStorage.GetGroupExternalId(groupId)
	if groupExternalId == nil {
		return nil, fmt.Errorf("Группа '%s' не найдена", groupId)
	}
	logrus.Printf("group name: %s, query params: %s", groupId, *groupExternalId)

	getScheduleResponse, err := s.scheduleApi.GetSchedule(*groupExternalId, from, to)
	if err != nil {
		return nil, err
	}
	logrus.Printf("schedule response from scheduleApi: %v", getScheduleResponse)

	if getScheduleResponse == nil || getScheduleResponse.Lessons == nil {
		return nil, errors.New("schedule api response is nil")
	}

	scheduleResponse := BuildGetScheduleResponse(getScheduleResponse)

	domainSchedule, err := BuildDomainSchedule(scheduleResponse, from, to)
	if err != nil {
		return nil, err
	}

	return domainSchedule, nil
}

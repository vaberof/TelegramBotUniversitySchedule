package domain

import (
	"time"
)

type ScheduleService struct {
	scheduleApi     ScheduleApi
	scheduleStorage ScheduleStorage
}

func NewScheduleService(scheduleApi ScheduleApi, scheduleStorage ScheduleStorage) *ScheduleService {
	return &ScheduleService{
		scheduleApi:     scheduleApi,
		scheduleStorage: scheduleStorage,
	}
}

func (s *ScheduleService) GetSchedule(groupId string, from time.Time, to time.Time) (Schedule, error) {
	return s.getScheduleImpl(groupId, from, to)
}

func (s *ScheduleService) getScheduleImpl(groupId string, from time.Time, to time.Time) (Schedule, error) {
	cachedSchedule, err := s.scheduleStorage.GetSchedule(groupId, from, to)
	if cachedSchedule == nil || err != nil {
		schedule, err := s.callScheduleApi(groupId, from, to)
		if err != nil {
			return nil, err
		}

		if err = s.cacheSchedule(groupId, schedule); err != nil {
			return nil, err
		}

		return schedule, nil
	}
	return cachedSchedule, nil
}

func (s *ScheduleService) callScheduleApi(groupId string, from time.Time, to time.Time) (Schedule, error) {
	schedule, err := s.scheduleApi.GetSchedule(groupId, from, to)
	if err != nil {
		return nil, err
	}
	return schedule, nil
}

func (s *ScheduleService) cacheSchedule(groupId string, schedule Schedule) error {
	err := s.scheduleStorage.SaveSchedule(groupId, schedule)
	if err != nil {
		return err
	}
	return nil
}

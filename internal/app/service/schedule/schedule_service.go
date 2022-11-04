package schedule

import (
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/infra/storage/postgres/schedulepg"
)

type ScheduleStorageService struct {
	scheduleStorage ScheduleStorage
}

func NewScheduleStorageService(scheduleStoragePostgres *schedulepg.ScheduleStoragePostgres) *ScheduleStorageService {
	return &ScheduleStorageService{scheduleStorage: scheduleStoragePostgres}
}

func (s *ScheduleStorageService) DeleteSchedule(groupId string, date string) error {
	return s.scheduleStorage.DeleteSchedule(groupId, date)
}

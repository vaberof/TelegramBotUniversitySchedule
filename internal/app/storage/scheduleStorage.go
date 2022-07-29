package storage

import (
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/constants"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/pkg/date"
	"time"
)

type ScheduleStorage struct {
	Schedule   map[int64][]map[string]string
	ExpireTime time.Time
}

func NewScheduleStorage() *ScheduleStorage {
	return &ScheduleStorage{
		Schedule: map[int64][]map[string]string{},
		ExpireTime: time.Date(time.Now().Year(),
			time.Now().Month(),
			time.Now().Day(),
			17, 00, 0, 0,
			time.UTC).In(date.GetDefaultLocation(constants.Location)),
	}
}

func (s *ScheduleStorage) ClearSchedule() {
	s.Schedule = map[int64][]map[string]string{}
}

func (s *ScheduleStorage) SetNewExpireTime() {
	s.ExpireTime = time.Date(time.Now().Year(),
		time.Now().Month(),
		time.Now().Day(),
		17, 0, 0, 0,
		time.UTC).In(date.GetDefaultLocation(constants.Location)).Add(24 * time.Hour)
}

// GetCachedScheduleIndex returns index array of cached schedule.
// Returns -1 if schedule not cached.
func GetCachedScheduleIndex(chatID int64, inputCallBack string, scheduleStorage *ScheduleStorage) int {
	if len(scheduleStorage.Schedule[chatID]) == 0 {
		return -1
	}

	for i := 0; i < len(scheduleStorage.Schedule[chatID]); i++ {
		for key, _ := range scheduleStorage.Schedule[chatID][i] {
			if key == inputCallBack {
				return i
			}
		}
	}
	return -1
}

func TimeExpired(currentTime time.Time, scheduleStorage *ScheduleStorage) bool {
	return currentTime.After(scheduleStorage.ExpireTime)
}

package storage

import (
	domain "github.com/vaberof/TelegramBotUniversitySchedule/internal/domain/schedule"
	"github.com/vaberof/TelegramBotUniversitySchedule/pkg/xtime"
	"github.com/vaberof/TelegramBotUniversitySchedule/pkg/xtimezone"
	"time"
)

type Schedule *domain.Schedule
type GroupId string
type Date string

type ScheduleStorage struct {
	Schedule   map[GroupId][]map[Date]*Schedule
	ExpireTime time.Time
}

func NewScheduleStorage() *ScheduleStorage {
	return &ScheduleStorage{
		Schedule: map[GroupId][]map[Date]*Schedule{},
		ExpireTime: time.Date(time.Now().Year(),
			time.Now().Month(),
			time.Now().Day(),
			17, 00, 0, 0,
			time.UTC).In(xtime.GetDefaultLocation(xtimezone.Novosibirsk)),
	}
}

func (s *ScheduleStorage) Save(groupId GroupId, date Date, schedule *Schedule) {
	inputSchedule := map[Date]*Schedule{
		date: schedule,
	}
	s.Schedule[groupId] = append(s.Schedule[groupId], inputSchedule)
}

func (s *ScheduleStorage) ClearSchedule() {
	s.Schedule = map[GroupId][]map[Date]*Schedule{}
}

func (s *ScheduleStorage) SetNewExpireTime() {
	s.ExpireTime = time.Date(time.Now().Year(),
		time.Now().Month(),
		time.Now().Day(),
		17, 0, 0, 0,
		time.UTC).In(xtime.GetDefaultLocation(xtimezone.Novosibirsk)).Add(24 * time.Hour)
}

// GetCachedScheduleIndex returns index array of cached schedule.
// Returns -1 if schedule not cached.
func GetCachedScheduleIndex(groupId GroupId, date Date, scheduleStorage *ScheduleStorage) int {
	groupIdConv := GroupId(groupId)

	if len(scheduleStorage.Schedule[groupIdConv]) == 0 {
		return -1
	}

	dateTgBtnConv := date
	for i := 0; i < len(scheduleStorage.Schedule[groupIdConv]); i++ {
		for key, _ := range scheduleStorage.Schedule[groupIdConv][i] {
			if key == dateTgBtnConv {
				return i
			}
		}
	}
	return -1
}

func TimeExpired(currentTime time.Time, scheduleStorage *ScheduleStorage) bool {
	return currentTime.After(scheduleStorage.ExpireTime)
}

package domain

import (
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/infra/storage"
	"time"
)

type ScheduleStorageApi interface {
	GetCachedLessons(groupId string, from time.Time, to time.Time) ([]*storage.Lesson, error)
	SaveLessons(groupId string, from time.Time, to time.Time, lessons []*storage.Lesson) error
}

type ScheduleStorage struct {
	ScheduleStorageApi
}

func NewScheduleStorage() *ScheduleStorage {
	return &ScheduleStorage{
		ScheduleStorageApi: storage.NewScheduleStorage(),
	}
}

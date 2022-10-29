package domain

import (
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/infra/storage/postgres/schedulepg"
	"gorm.io/gorm"
	"time"
)

type LessonsReceiver interface {
	GetLessons(groupId string, from time.Time, to time.Time) ([]*schedulepg.Lesson, error)
}

type LessonsSaver interface {
	SaveLessons(groupId string, from time.Time, to time.Time, lessons []*schedulepg.Lesson) error
}

type LessonsReceiverSaver interface {
	LessonsReceiver
	LessonsSaver
}

type ScheduleStorage struct {
	LessonsReceiverSaver
}

func NewScheduleStorage(db *gorm.DB) *ScheduleStorage {
	return &ScheduleStorage{
		LessonsReceiverSaver: schedulepg.NewScheduleStoragePostgres(db),
	}
}

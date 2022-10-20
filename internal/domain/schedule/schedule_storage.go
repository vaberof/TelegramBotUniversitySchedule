package domain

import (
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/infra/storage"
	"time"
)

type LessonsReceiver interface {
	GetLessons(groupId string, from time.Time, to time.Time) ([]*storage.Lesson, error)
}

type LessonsSaver interface {
	SaveLessons(groupId string, from time.Time, to time.Time, lessons []*storage.Lesson) error
}

type LessonsReceiverSaver interface {
	LessonsReceiver
	LessonsSaver
}

type ScheduleStorage struct {
	LessonsReceiverSaver
}

func NewScheduleStorage() *ScheduleStorage {
	return &ScheduleStorage{
		LessonsReceiverSaver: storage.NewScheduleStorage(),
	}
}

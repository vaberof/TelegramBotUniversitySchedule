package domain

import (
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/infra/storage/postgres/schedulepg"
	"time"
)

type ScheduleStorage interface {
	GetLessons(groupId string, from time.Time, to time.Time) ([]*schedulepg.Lesson, error)
	SaveSchedule(groupId string, from time.Time, to time.Time, lessons []*schedulepg.Lesson) error
}

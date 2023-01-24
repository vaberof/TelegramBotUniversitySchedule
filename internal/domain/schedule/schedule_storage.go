package domain

import (
	"time"
)

type ScheduleStorage interface {
	GetSchedule(groupId string, from time.Time, to time.Time) (Schedule, error)
	SaveSchedule(groupId string, schedule Schedule) error
}

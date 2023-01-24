package domain

import (
	"time"
)

type ScheduleApi interface {
	GetSchedule(groupId string, from time.Time, to time.Time) (Schedule, error)
}

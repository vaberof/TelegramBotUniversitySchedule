package services

import (
	"time"

	"github.com/tg_bot_timetable/internal/model"
)

// GetSchedule returns schedule with all lessons.
func GetSchedule(groupId, date, url string, location *time.Location) model.Schedule {
	return *model.ParseLessons(groupId, date, url, location)
}

// ScheduleToString converts schedule to string.
func ScheduleToString(schedule *model.Schedule) string {
	var scheduleString string

	for _, data := range schedule.Schedule {
		scheduleString += data
	}

	return scheduleString
}
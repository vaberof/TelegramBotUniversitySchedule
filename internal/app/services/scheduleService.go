package services

import (
	"time"

	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/model"
)

// GetSchedule returns schedule of type model.Schedule.
func GetSchedule(groupId, date, url string, location *time.Location) model.Schedule {
	return dateSchedule(groupId, date, url, location)
}

// ScheduleToString converts schedule of type model.Schedule to string type
// to output it to user.
func ScheduleToString(schedule *model.Schedule) string {
	var scheduleString string

	for _, data := range schedule.Schedule {
		scheduleString += data
	}

	return scheduleString
}

// dateSchedule calls ParseDayLessons or ParseWeekLessons depending on user input
func dateSchedule(groupId, date, url string, location *time.Location) model.Schedule {
	switch date {
	case "Сегодня":
		return *model.ParseDayLessons(groupId, date, url, location)
	case "Завтра":
		return *model.ParseDayLessons(groupId, date, url, location)
	case "Неделя":
		return *model.ParseWeekLessons(groupId, date, url, location)
	case "След. неделя":
		return *model.ParseWeekLessons(groupId, date, url, location)
	default:
		return model.Schedule{}
	}
}

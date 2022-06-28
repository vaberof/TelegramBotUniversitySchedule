package unisite

import (
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/model"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/pkg/date"
)

// GetSchedule returns schedule of type model.Schedule.
func GetSchedule(url, inputCallback string, parseData *model.ParseData) *model.Schedule {
	switch inputCallback {
	case date.Today, date.Tomorrow:
		return ParseDayLessons(inputCallback, url, parseData)
	default:
		return ParseWeekLessons(inputCallback, url, parseData)
	}
}

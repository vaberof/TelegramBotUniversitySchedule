package unisite

import (
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/model"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/constants"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/pkg/date"
)

// GetSchedule returns pointer to schedule of type model.Schedule.
func GetSchedule(url, inputCallback string) *model.Schedule {
	switch inputCallback {
	case constants.Today, constants.Tomorrow:
		toParseDate := date.GetParseDate(inputCallback)
		return ParseDayLessons(inputCallback, url, toParseDate)
	default:
		toParseDates := date.GetParseDates(inputCallback)
		return ParseWeekLessons(inputCallback, url, toParseDates)
	}
}

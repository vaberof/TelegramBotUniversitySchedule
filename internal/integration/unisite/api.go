package unisite

import (
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/model"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/constants"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/pkg/date"
)

// GetSchedule returns pointer to schedule of type model.Schedule
// and custom error if http request Timeout occurred, otherwise error equals to nil.
func GetSchedule(url, inputCallback string) (*model.Schedule, error) {
	switch inputCallback {
	case constants.Today, constants.Tomorrow:
		dateToParse := date.GetParseDate(inputCallback)
		return ParseDayLessons(inputCallback, url, dateToParse)
	default:
		datesToParse := date.GetParseDates(inputCallback)
		return ParseWeekLessons(inputCallback, url, datesToParse)
	}
}

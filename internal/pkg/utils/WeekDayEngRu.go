package utils

import (
	"time"

	"github.com/vaberof/TelegramBotUniversitySchedule/internal/constants"
)

// WeekDayEngRu translates day of the week from english to russian.
func WeekDayEngRu(engDay time.Weekday) string {
	switch engDay.String() {
	case constants.Monday:
		return "Понедельник"
	case constants.Tuesday:
		return "Вторник"
	case constants.Wednesday:
		return "Среда"
	case constants.Thursday:
		return "Четверг"
	case constants.Friday:
		return "Пятница"
	case constants.Saturday:
		return "Суббота"
	default:
		return "Воскресенье"
	}
}

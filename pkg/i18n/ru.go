package i18n

import (
	"time"
)

// FormatRuWeekday translates day of the week from english to russian.
func FormatRuWeekday(engDay time.Weekday) string {
	switch engDay {
	case time.Monday:
		return "Понедельник"
	case time.Tuesday:
		return "Вторник"
	case time.Wednesday:
		return "Среда"
	case time.Thursday:
		return "Четверг"
	case time.Friday:
		return "Пятница"
	case time.Saturday:
		return "Суббота"
	default:
		return "Воскресенье"
	}
}

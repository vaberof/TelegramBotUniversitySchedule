package date

import (
	"time"
)

const (
	Monday    string = "Monday"
	Tuesday          = "Tuesday"
	Wednesday        = "Wednesday"
	Thursday         = "Thursday"
	Friday           = "Friday"
	Saturday         = "Saturday"
)

const (
	Today    string = "Today"
	Tomorrow        = "Tomorrow"
	Week            = "Week"
)

// WeekDayEngRu translates day of the week from english to russian.
func WeekDayEngRu(engDay time.Weekday) string {
	switch engDay.String() {
	case Monday:
		return "Понедельник"
	case Tuesday:
		return "Вторник"
	case Wednesday:
		return "Среда"
	case Thursday:
		return "Четверг"
	case Friday:
		return "Пятница"
	case Saturday:
		return "Суббота"
	default:
		return "Воскресенье"
	}
}
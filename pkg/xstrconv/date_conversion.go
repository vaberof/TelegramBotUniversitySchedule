package xstrconv

import (
	"errors"
	"github.com/vaberof/TelegramBotUniversitySchedule/pkg/xtime"
	"github.com/vaberof/TelegramBotUniversitySchedule/pkg/xtimezone"
	"time"
)

func ConvertDateToStr(from time.Time, to time.Time) (string, error) {
	fromFormatted := from.Format("02.01")
	toFormatted := to.Format("02.01")

	Novosibirsk := xtime.GetDefaultLocation(xtimezone.Novosibirsk)

	today := xtime.GetTodayDate(Novosibirsk)
	tomorrow := xtime.GetTomorrowDate(Novosibirsk)
	week := xtime.GetWeekDatesRange()
	nextWeek := xtime.GetNextWeekDatesRange()

	// Today
	if fromFormatted == today.Format("02.01") && toFormatted == today.Format("02.01") {
		return "Today", nil
	}
	// Tomorrow
	if fromFormatted == today.Format("02.01") && toFormatted == tomorrow.Format("02.01") {
		return "Tomorrow", nil
	}
	// Week
	if fromFormatted == week[0].Format("02.01") && toFormatted == week[6].Format("02.01") {
		return "Week", nil
	}
	// Next week
	if fromFormatted == nextWeek[0].Format("02.01") && toFormatted == nextWeek[6].Format("02.01") {
		return "Next Week", nil
	}

	return "", errors.New("cannot convert date to string")
}

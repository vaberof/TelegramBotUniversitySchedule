package xtimeconv

import (
	"errors"
	"github.com/vaberof/TelegramBotUniversitySchedule/pkg/xtime"
	"github.com/vaberof/TelegramBotUniversitySchedule/pkg/xtimezone"
	"time"
)

func FromTimeToString(from time.Time, to time.Time) (string, error) {
	novosibirsk, err := xtime.GetDefaultLocation(xtimezone.Novosibirsk)
	if err != nil {
		return "", err
	}

	dateString, err := getStringPeriodOfTime(from, to, novosibirsk)
	if err != nil {
		return "", err
	}

	return dateString, nil
}

func getStringPeriodOfTime(from time.Time, to time.Time, location *time.Location) (string, error) {
	today, tomorrow, week, nextWeek, err := getDates(location)
	if err != nil {
		return "", err
	}

	if isToday(from, to, today) {
		return "Today", nil
	}

	if isTomorrow(from, to, today, tomorrow) {
		return "Tomorrow", nil
	}

	if isWeek(from, to, week[0], week[6]) {
		return "Week", nil
	}

	if isNextWeek(from, to, nextWeek[0], nextWeek[6]) {
		return "Next week", nil
	}

	return "", errors.New("cannot convert date to string")
}

func getDates(location *time.Location) (time.Time, time.Time, []time.Time, []time.Time, error) {
	today, err := xtime.GetTodayDate(location)
	if err != nil {
		return time.Time{}, time.Time{}, []time.Time{}, []time.Time{}, err
	}

	tomorrow, err := xtime.GetTomorrowDate(location)
	if err != nil {
		return time.Time{}, time.Time{}, []time.Time{}, []time.Time{}, err
	}

	week, err := xtime.GetWeekDatesRange()
	if err != nil {
		return time.Time{}, time.Time{}, []time.Time{}, []time.Time{}, err
	}

	nextWeek, err := xtime.GetNextWeekDatesRange()
	if err != nil {
		return time.Time{}, time.Time{}, []time.Time{}, []time.Time{}, err
	}

	return today, tomorrow, week, nextWeek, nil
}

func isToday(from time.Time, to time.Time, today time.Time) bool {
	return from.Format("02.01") == today.Format("02.01") &&
		to.Format("02.01") == today.Format("02.01")
}

func isTomorrow(from time.Time, to time.Time, today time.Time, tomorrow time.Time) bool {
	return from.Format("02.01") == today.Format("02.01") &&
		to.Format("02.01") == tomorrow.Format("02.01")
}

func isWeek(from time.Time, to time.Time, monday time.Time, sunday time.Time) bool {
	return from.Format("02.01") == monday.Format("02.01") &&
		to.Format("02.01") == sunday.Format("02.01")
}

func isNextWeek(from time.Time, to time.Time, nextWeekMonday time.Time, nextWeekSunday time.Time) bool {
	return from.Format("02.01") == nextWeekMonday.Format("02.01") &&
		to.Format("02.01") == nextWeekSunday.Format("02.01")
}

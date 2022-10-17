package xtime

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/vaberof/TelegramBotUniversitySchedule/pkg/xtimezone"
	"github.com/vaberof/goweekdate"
	"time"
)

func ParseDatesRange(inputTelegramButtonDate string) (time.Time, time.Time, error) {
	switch inputTelegramButtonDate {
	case "Today":
		today, err := GetDate(inputTelegramButtonDate)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		return today, today, nil
	case "Tomorrow":
		tomorrow, err := GetDate(inputTelegramButtonDate)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		today := tomorrow.Add(-time.Hour * 24)
		return today, tomorrow, nil
	case "Week":
		currentWeekDatesRange, err := GetDatesRange(inputTelegramButtonDate)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		return currentWeekDatesRange[0], currentWeekDatesRange[6], nil
	default:
		nextWeekDatesRange, err := GetDatesRange(inputTelegramButtonDate)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		return nextWeekDatesRange[0], nextWeekDatesRange[6], nil
	}
}

func GetDate(inputTelegramButtonDate string) (time.Time, error) {
	novosibirsk, err := GetDefaultLocation(xtimezone.Novosibirsk)
	if err != nil {
		return time.Time{}, err
	}
	return parseDate(inputTelegramButtonDate, novosibirsk)
}

func GetDatesRange(inputTelegramButtonDate string) ([]time.Time, error) {
	return parseDates(inputTelegramButtonDate)
}

func GetTodayDate(location *time.Location) (time.Time, error) {
	todayDate := time.Now().In(location)
	if todayDate.IsZero() {
		return time.Time{}, errors.New("error occurred while getting today`s date")
	}
	return todayDate, nil
}

func GetTomorrowDate(location *time.Location) (time.Time, error) {
	tomorrowDate := time.Now().Add(time.Hour * 24).In(location)
	if tomorrowDate.IsZero() {
		return time.Time{}, errors.New("error occurred while getting tomorrow`s date")
	}
	return tomorrowDate, nil
}

func GetWeekDatesRange() ([]time.Time, error) {
	wd := weekdate.New(time.Now(), xtimezone.Novosibirsk)
	Dates := wd.Dates(1, true)
	if len(Dates) == 0 {
		return nil, errors.New("error occurred while getting array of dates range")
	}
	return Dates, nil
}

func GetNextWeekDatesRange() ([]time.Time, error) {
	wd := weekdate.New(time.Now(), xtimezone.Novosibirsk)
	Dates := wd.Dates(2, false)
	if len(Dates) == 0 {
		return nil, errors.New("error occurred while getting array of dates range")
	}
	return Dates, nil
}

func GetDefaultLocation(timeZone string) (*time.Location, error) {
	location, err := time.LoadLocation(timeZone)
	if err != nil {
		log.WithFields(log.Fields{
			"timeZone": timeZone,
			"error":    err,
			"func":     "GetDefaultLocation",
		}).Error("Failed to load a time zone")
		return nil, err
	}
	return location, nil
}

func parseDate(inputTelegramButtonDate string, location *time.Location) (time.Time, error) {
	switch inputTelegramButtonDate {
	case Today:
		todayDate, err := GetTodayDate(location)
		if err != nil {
			return time.Time{}, err
		}
		return todayDate, nil
	default:
		tomorrowDate, err := GetTomorrowDate(location)
		if err != nil {
			return time.Time{}, err
		}
		return tomorrowDate, nil
	}
}

func parseDates(inputTelegramButtonDate string) ([]time.Time, error) {
	switch inputTelegramButtonDate {
	case Week:
		currentWeekDatesRange, err := GetWeekDatesRange()
		if err != nil {
			return nil, err
		}
		return currentWeekDatesRange, nil
	default:
		nextWeekDatesRange, err := GetNextWeekDatesRange()
		if err != nil {
			return nil, err
		}
		return nextWeekDatesRange, nil
	}
}

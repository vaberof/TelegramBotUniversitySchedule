package xtime

import (
	log "github.com/sirupsen/logrus"
	"github.com/vaberof/TelegramBotUniversitySchedule/pkg/xtimezone"
	"github.com/vaberof/goweekdate"
	"time"
)

func GetDateToParse(inputDateTelegramButton string) time.Time {
	Novosibirsk := GetDefaultLocation(xtimezone.Novosibirsk)
	switch inputDateTelegramButton {
	case Today:
		return GetTodayDate(Novosibirsk)
	default:
		return GetTomorrowDate(Novosibirsk)
	}
}

func GetDatesToParse(inputDateTelegramButton string) []time.Time {
	switch inputDateTelegramButton {
	case Week:
		return GetWeekDatesRange()
	default:
		return GetNextWeekDatesRange()
	}
}

func GetTodayDate(location *time.Location) time.Time {
	todayDate := time.Now().In(location)
	return todayDate
}

func GetTomorrowDate(location *time.Location) time.Time {
	tomorrowDate := time.Now().Add(time.Hour * 24).In(location)
	return tomorrowDate
}

func GetWeekDatesRange() []time.Time {
	wd := weekdate.New(time.Now(), xtimezone.Novosibirsk)

	Dates := wd.Dates(1, true)
	return Dates
}

func GetNextWeekDatesRange() []time.Time {
	wd := weekdate.New(time.Now(), xtimezone.Novosibirsk)

	Dates := wd.Dates(2, false)
	return Dates
}

func GetDefaultLocation(location string) *time.Location {
	loc, err := time.LoadLocation(location)
	if err != nil {
		log.WithFields(log.Fields{
			"location": location,
			"error":    err,
			"func":     "GetDefaultLocation",
		}).Fatal("Failed to load a location")
	}
	return loc
}

package xtime

import (
	log "github.com/sirupsen/logrus"
	"github.com/vaberof/TelegramBotUniversitySchedule/pkg/xtimezone"
	"github.com/vaberof/goweekdate"
	"time"
)

func GetDateToParse(inputDateTelegramButton string) time.Time {
	location := GetDefaultLocation(xtimezone.Novosibirsk)

	switch inputDateTelegramButton {
	case Today:
		return today(location)
	default:
		return tomorrow(location)
	}
}

func GetDatesToParse(inputDateTelegramButton string) []time.Time {
	switch inputDateTelegramButton {
	case Week:
		return week()
	case NextWeek:
		return nextWeek()
	}
	return nil
}

func today(location *time.Location) time.Time {
	todayDate := time.Now().In(location)
	return todayDate
}

func tomorrow(location *time.Location) time.Time {
	tomorrowDate := time.Now().Add(time.Hour * 24).In(location)
	return tomorrowDate
}

func week() []time.Time {
	wd := weekdate.New(time.Now(), xtimezone.Novosibirsk)

	Dates := wd.Dates(1, true)
	return Dates
}

func nextWeek() []time.Time {
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

func GetCurrentTime() time.Time {
	return time.Now().In(GetDefaultLocation(xtimezone.Novosibirsk))
}

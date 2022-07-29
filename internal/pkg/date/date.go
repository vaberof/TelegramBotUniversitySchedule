package date

import (
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/vaberof/TelegramBotUniversitySchedule/internal/constants"
	"github.com/vaberof/goweekdate"
)

func GetParseDate(inputCallBack string) time.Time {
	location := GetDefaultLocation(constants.Location)

	switch inputCallBack {
	case constants.Today:
		return today(location)
	default:
		return tomorrow(location)
	}
}

func GetParseDates(inputCallBack string) []time.Time {
	switch inputCallBack {
	case constants.Week:
		return week()
	default:
		return nextWeek()
	}
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
	wd := weekdate.New(time.Now(), constants.Location)

	Dates := wd.Dates(1, true)
	return Dates
}

func nextWeek() []time.Time {
	wd := weekdate.New(time.Now(), constants.Location)

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
	return time.Now().In(GetDefaultLocation(constants.Location))
}

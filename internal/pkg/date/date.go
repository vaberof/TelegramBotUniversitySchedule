package date

import (
	"log"
	"time"

	"github.com/vaberof/TelegramBotUniversitySchedule/internal/constants"
	"github.com/vaberof/goweekdate"
)

func GetParseDate(inputCallBack string) time.Time {
	loc := getDefaultLocation(constants.Location)

	switch inputCallBack {
	case constants.Today:
		return today(loc)
	default:
		return tomorrow(loc)
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

func getDefaultLocation(location string) *time.Location {
	loc, err := time.LoadLocation(location)
	if err != nil {
		log.Fatalf("failed to load a setLocation: %v", err)
	}
	return loc
}

package model

import (
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/pkg/date"
	"log"
	"time"

	"github.com/vaberof/goweekdate"
)

const Location = "Asia/Novosibirsk"

// A ParseData contains...
type ParseData struct {
	Date  time.Time
	Dates []time.Time
	Days  []time.Weekday
}

func NewParseData() *ParseData {
	return &ParseData{}
}

// GetParseData returns pointer to object of type ParseData depending on the user`s choice to manipulate it in parser.go.
func GetParseData(inputCallBack string) *ParseData {
	parseData := NewParseData()
	loc := getDefaultLocation(Location)

	switch inputCallBack {
	case date.Today:
		parseData.today(loc)
		return parseData
	case date.Tomorrow:
		parseData.tomorrow(loc)
		return parseData
	case date.Week:
		parseData.week()
		return parseData
	default:
		parseData.nextWeek()
		return parseData
	}
}

// today gets today`s date.
func (d *ParseData) today(location *time.Location) {
	todayDate := time.Now().In(location)
	d.Date = todayDate
}

// tomorrow gets next day date.
func (d *ParseData) tomorrow(location *time.Location) {
	tomorrowDate := time.Now().Add(time.Hour * 24).In(location)
	d.Date = tomorrowDate
}

// week gets dates and the names of days of the current week via weekdate package.
func (d *ParseData) week() {
	wd := weekdate.New(time.Now(), Location)

	d.Dates = wd.Dates(1, true)
	d.Days = wd.WeekDays()
}

// nextWeek gets dates and the names of days of the next week via weekdate package.
func (d *ParseData) nextWeek() {
	wd := weekdate.New(time.Now(), Location)

	d.Dates = wd.Dates(2, false)
	d.Days = wd.WeekDays()
}

// getDefaultLocation sets given location with time.LoadLocation.
// Fatales if cannot load given location.
func getDefaultLocation(location string) *time.Location {
	loc, err := time.LoadLocation(location)
	if err != nil {
		log.Fatalf("failed to load a setLocation: %v", err)
	}
	return loc
}

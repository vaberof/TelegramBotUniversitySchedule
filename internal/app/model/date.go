package model

import (
	"fmt"
	"time"

	"github.com/vaberof/goweekdate"
)

const location = "Asia/Novosibirsk"

type Date struct {
	shortDate string
	fullDate  string
	day       string

	weekShortDates []string
	weekFullLDates []string
	weekDays       []string
}

func NewDate() *Date {
	return &Date{}
}

// getDate returns object of type Date depending on the user`s choice to manipulate it in parser.
func getDate(date string, location *time.Location) *Date {
	d := NewDate()

	switch date {
	case "Сегодня":
		d.today(location)
		return d
	case "Завтра":
		d.tomorrow(location)
		return d
	case "Неделя":
		d.week()
		return d
	default:
		return nil
	}
}

// today gets today`s date and weekday.
func (d *Date) today(location *time.Location) {
	todayDate := time.Now().In(location)
	day := todayDate.Weekday()

	d.shortDate = todayDate.Format("02.01")
	d.fullDate = todayDate.Format("02.01.2006")
	d.day = dayTranslate(day.String())
}

// tomorrow gets next day date and weekday.
func (d *Date) tomorrow(location *time.Location) {
	tomorrowDate := time.Now().Add(time.Hour * 24).In(location)
	day := tomorrowDate.Weekday()

	d.shortDate = tomorrowDate.Format("02.01")
	d.fullDate = tomorrowDate.Format("02.01.2006")
	d.day = dayTranslate(day.String())
}

// week gets dates and the names of the current week days.
func (d *Date) week() {
	weekDate := weekdate.New(time.Now(), location)

	d.weekShortDates = weekDate.ShortDates(1, true)
	d.weekFullLDates = weekDate.FullDates(1, true)
	d.setWeekDays()
}

func (d *Date) setWeekDays() {
	weekDate := weekdate.New(time.Now(), location)
	wDays := weekDate.WeekDays()

	for _, day := range wDays {
		d.weekDays = append(d.weekDays, dayTranslate(day))
	}
}

// SetLocation sets "Novosibirsk" location.
func SetLocation() *time.Location {
	loc, err := time.LoadLocation(location)
	if err != nil {
		fmt.Errorf("failed to load a location: %v", err)
	}

	return loc
}

// dayTranslate translates day of the week from english to russian.
func dayTranslate(day string) string {
	switch day {
	case "Monday":
		return "Понедельник"
	case "Tuesday":
		return "Вторник"
	case "Wednesday":
		return "Среда"
	case "Thursday":
		return "Четверг"
	case "Friday":
		return "Пятница"
	case "Saturday":
		return "Суббота"
	default:
		return "Воскресенье"
	}
}

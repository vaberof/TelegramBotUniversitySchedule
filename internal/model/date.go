package model

import (
	"fmt"
	"time"
)

type Date struct {
	shortDate string
	fullDate  string
	day       string
}

func NewDate() *Date {
	return &Date{}
}

// getDate returns date depending on the user`s choice.
func getDate(date string, location *time.Location) *Date {
	d := NewDate()

	switch date {
	case "Сегодня":
		d.today(location)
		return d
	case "Завтра":
		d.tomorrow(location)
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

// tomorrow gets tomorrow`s date and weekday.
func (d *Date) tomorrow(location *time.Location) {
	tomorrowDate := time.Now().Add(time.Hour * 24).In(location)
	day := tomorrowDate.Weekday()

	d.shortDate = tomorrowDate.Format("02.01")
	d.fullDate = tomorrowDate.Format("02.01.2006")
	d.day = dayTranslate(day.String())
}

// SetLocation sets "Novosibirsk" location.
func SetLocation() *time.Location {
	location, err := time.LoadLocation("Asia/Novosibirsk")
	if err != nil {
		fmt.Errorf("failed to load a location: %v", err)
	}

	return location
}

// dayTranslate translates day of the week from english to russian.
func dayTranslate(day string) string{
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

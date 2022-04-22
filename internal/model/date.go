package model

import (
	"fmt"
	"time"
)

type Date struct {
	shortDate string
	fullDate  string
}

// NewDate returns pointer to Date structure
func NewDate() *Date {
	return &Date{}
}

// getDate returns needed date.
func getDate(date string, location *time.Location) *Date {
	d := NewDate()

	switch date {
	case "Сегодня":
		d.todayDate(location)
		return d
	case "Завтра":
		d.tomorrowDate(location)
		return d
	default:
		return nil
	}
}

// todayDate returns today`s date
func (d *Date) todayDate(location *time.Location) {
	todayDate := time.Now().In(location)
	d.shortDate = todayDate.Format("02.01")
	d.fullDate = todayDate.Format("02.01.2006")
}

// tomorrowDate returns tomorrow`s date
func (d *Date) tomorrowDate(location *time.Location) {
	tomorrowDate := time.Now().Add(time.Hour * 24).In(location)
	d.shortDate = tomorrowDate.Format("02.01")
	d.fullDate = tomorrowDate.Format("02.01.2006")
}

// SetLocation sets "Novosibirsk" location
func SetLocation() *time.Location {
	location, err := time.LoadLocation("Asia/Novosibirsk")
	if err != nil {
		fmt.Errorf("failed to load a location: %v", err)
	}

	return location
}
package model

import (
	"log"
	"time"
)

// Устанавливаем местоположение "Новосибирск"
func setLocation() *time.Location {
	location, err := time.LoadLocation("Asia/Novosibirsk")
	if err != nil {
		log.Println("Failed to load a location")
	}

	return location
}

// Получаем сегодняшнюю дату
func GetTodayDate() []string {

	var todayDate []string
	var shortTodayDate string
	var fullTodayDate string

	shortTodayDate = time.Now().In(setLocation()).Format("02.01")
	fullTodayDate = time.Now().In(setLocation()).Format("02.01.2006")

	//shortTodayDate = "17.04"
	//fullTodayDate = "17.04.2022"
	todayDate = []string{shortTodayDate, fullTodayDate}

	return todayDate
}

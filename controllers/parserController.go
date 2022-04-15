package controllers

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/tg_bot_timetable/handlers"
	"github.com/tg_bot_timetable/services"
	"log"
	"strings"
)

// Находим блок, с сегодняшней датой
func parseDate(url string) *goquery.Selection {

	var dateSelection *goquery.Selection

	document := handlers.LoadHtmlPage(url)

	document.Find("div.one_day-wrap").EachWithBreak(func(index int, tag *goquery.Selection) bool {
		// Ищем совпадение с сегодняшней датой
		everDTag := tag.Find("div.everD")
		everDTagValue := strings.ReplaceAll(everDTag.Text(), " ", "")
		// Если нашлась текущая дата
		if everDTagValue == services.GetTodayDate()[0] {
			dateSelection = tag
			return false
		}
		return true
})
	return dateSelection
}

// Проверяем, существует ли тег на сайте
func isNilSelection(selection *goquery.Selection) bool {

	if selection == nil {
		log.Printf("Ошибка, тег %v не найден", selection)
		return true
	}
	return false
}

// Ищем расписание на сегодняшний день
func parseLessons(url string) string {

	var (
		startTime       string // Начало пары
		finishTime      string // Конец пары
		lessonName      string // Название предмета
		roomNumber      string // Номер аудитории
		teacherName     string // Имя преподавателя
		lessonType      string // Тип пары (лекция/практика/лабораторная)
		responseMessage string // Сообщение пользователю
	)

	schedule := services.CreateSchedule()
	dateSelection := parseDate(url)

	// Проверяем, существует ли тег на странице
	if isNilSelection(dateSelection) {
		responseMessage = "Воскресенье - пар нет"
		return responseMessage
	}

	// Ищем все пары на сегодняшний день
	dateSelection.Find(".one_lesson").EachWithBreak(func(index int, tag *goquery.Selection) bool {
		lessonName = tag.Find(".names_of_less").Text()
		if lessonName != "" {
			startTime = tag.Find(".starting_less").Text()
			finishTime = tag.Find(".finished_less").Text()
			roomNumber = tag.Find(".kabinet_of_less").Text()
			teacherName = tag.Find(".name_of_teacher").Text()
			lessonType = tag.Find(".type_less").Text()

			// Добавляем элементы в массив
			schedule.AddLessons(
				startTime,
				finishTime,
				lessonName,
				roomNumber,
				teacherName,
				lessonType)
		}
		return true
	})

	responseMessage = schedule.GetSchedule()

	return responseMessage
}

// Получаем расписание на сегодняшний день
func getTodaySchedule(url string) string {
	return parseLessons(url)
}

package controllers

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
	"time"
)

// Делаем запрос
func makeRequest(url string) *http.Response {

	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}

	//defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Println("status code error: %d %s", res.StatusCode, res.Status)
	}
	body := res

	return body
}

// Загружаем HTML страничку
func loadHtmlPage(url string) *goquery.Document{

	doc, err := goquery.NewDocumentFromReader(makeRequest(url).Body)

	if err != nil {
		log.Println(err)
	}
	return doc
}

// Получаем сегодняшнюю дату
func getTodayDate() []string {

	var shortTodayDate string
	var fullTodayDate string

	novosibirsk, _ := time.LoadLocation("Asia/Novosibirsk")

	shortTodayDate = time.Now().In(novosibirsk).Format("02.01")
	fullTodayDate = time.Now().In(novosibirsk).Format("02.01.2006")

	//shortTodayDate = "08.04"
	//fullTodayDate = "08.04.2022"

	return []string{shortTodayDate, fullTodayDate}
}

// Проверяем во сколько начинается пара,
// чтобы правильно указывать номер пары
func lessonNum (startTime string) string {
	var numLesson string

	switch startTime {
	case "10:10":
		numLesson = "2"
	case "12:10":
		numLesson = "3"
	case "13:50":
		numLesson = "4"
	case "15:30":
		numLesson = "5"
	case "17:10":
		numLesson = "6"
	case "18:50":
		numLesson = "7"
	default:
		numLesson = "1"
	}

	return numLesson
}

// Находим блок, с сегодняшней датой
func findTodayDateTag(url string) *goquery.Selection {

	doc := loadHtmlPage(url)

	var selection *goquery.Selection

	doc.Find("div.one_day-wrap").EachWithBreak(func(index int, item *goquery.Selection) bool {
		// Ищем совпадение с сегодняшней датой
		everDTag := item.Find("div.everD")
		everDTagValue := strings.ReplaceAll(everDTag.Text(), " ", "") // убираем пробелы в дате
		// Если нашлась текущая дата
		if everDTagValue == getTodayDate()[0] {
			selection = item
			return false
		}
		return true
})
	return selection
}

// Проверяем, существует ли тег на сайте
func isNilTag(item *goquery.Selection) bool {

	if item == nil {
		log.Printf("Ошибка, тег %v не найден", item)
		return true
	}
	return false
}

// Ищем расписание на сегодняшний день
func findTodaySchedule(url string) string {

	var (
		startTime     string  // Начало пары
		finishTime    string  // Конец пары
		lessonName    string  // Название предмета
		roomNumber    string  // Номер аудитории
		teacherName   string  // Имя преподавателя
		lessonType    string  // Тип пары (лекция/практика/лабораторная)
		messageToUser string  // Сообщение пользователю
	)

	schedule := scheduleInit()
	item := findTodayDateTag(url)

	// Проверяем, существует ли тег на странице
	if isNilTag(item) {
		messageToUser = "сегодня/завтра 'воскресенье' - пар нет"
		return messageToUser
	}

	// Ищем все пары на сегодняшний день
	item.Find(".one_lesson").EachWithBreak(func(index int, item *goquery.Selection) bool {
		lessonName = item.Find(".names_of_less").Text()
		if lessonName != "" {
			startTime = item.Find(".starting_less").Text()
			finishTime = item.Find(".finished_less").Text()
			roomNumber = item.Find(".kabinet_of_less").Text()
			teacherName = item.Find(".name_of_teacher").Text()
			lessonType = item.Find(".type_less").Text()

			// Добавляем в массив, при этом убирая отступы в начале строки
			schedule.appendLessonsToArray(
				startTime,
				finishTime,
				lessonName,
				roomNumber,
				teacherName,
				lessonType)
		}
		return true
	})

	messageToUser = schedule.getScheduleFromArray()

	return messageToUser
}

// Получаем расписание на сегодняшний день
func getTodaySchedule(url string) string {
	return findTodaySchedule(url)
}

package controllers

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
)

// Делаем запрос
func DoRequest(url string) *http.Response {
	// Делаем запрос
	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}

	if res.StatusCode != 200 {
		fmt.Println("status code error: %d %s", res.StatusCode, res.Status)
	}
	return res
}

// "Сегодня" конвертируем сегодняшнюю верхнюю дату сайта в нужный формат с точкой "29.03" и "29.03.2022"
func todayDateConvert(date string) []string {

	var shortTodayDate string
	var fullTodayDate string

	for i := 0; i <= 5; i++{

		if i == 2 {
			shortTodayDate += "."
			fullTodayDate += "."
		}

		if i == 4 {
			fullTodayDate += "." + "20"
		}

		if i < 4 {
			c := fmt.Sprintf("%c", date[i])
			shortTodayDate += c
		}

		c := fmt.Sprintf("%c", date[i])
		fullTodayDate += c
	}
	shortTodayDate = "01.04"
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

// Переменные, которые используются в файлах parser.go и bot.go
var (
	DateCount int
	ResponseMassive []string
)

// Получаем расписание на сегодняшний день
func TodayShedule(url string)  {

	var (
		startTime   string  // Начало пары
		finishTime  string  // Конец пары
		nameLesson  string  // Название предмета
		kabinetNum  string  // Номер аудитории
		nameTeacher string  // Имя преподавателя
		typeLesson  string  // Тип пары (лекция/практика/лабораторная)
	)

	// Загружаем HTML страничку
	doc, err := goquery.NewDocumentFromReader(DoRequest(url).Body)
	if err != nil {
		log.Println(err)
	}

	defer DoRequest(url).Body.Close()
	if err != nil {
		log.Println(err)
	}

	// Ищем сегодняшнюю дату вверху сайта
	var date string

	doc.Find(".date_and_time").Each(func(index int, item *goquery.Selection) {
		// For each item found, get the title
		span := item.Find("span")
		today := span.Text()
		date = today
	})

	// Ищем расписание на сегодняшний день
	doc.Find("div.one_day-wrap").EachWithBreak(func(index int, item *goquery.Selection) bool{
		// Ищем совпадение с сегодняшней датой
		everDTag := item.Find("div.everD")
		value := strings.ReplaceAll(everDTag.Text(), " ", "") // убираем пробелы в дате
		// Если нашлась текущая дата
		if value == todayDateConvert(date)[0]{
			// Ищем расписание на сегодня
			item.Find(".one_lesson").EachWithBreak(func(index int, item *goquery.Selection) bool {
				nameLesson = item.Find(".names_of_less").Text()
				if nameLesson != "" {
					startTime = item.Find(".starting_less").Text()
					finishTime = item.Find(".finished_less").Text()
					kabinetNum = item.Find(".kabinet_of_less").Text()
					nameTeacher = item.Find(".name_of_teacher").Text()
					typeLesson = item.Find(".type_less").Text()

					// Добавление сегодняшней даты
					if DateCount == 0 {
						ResponseMassive = append(ResponseMassive, "[Дата]: " + todayDateConvert(date)[1] + " (Сегодня)" + "\n" )
						DateCount += 1
					}
					// Добавляем в массив пары для вывода в файл bot.go
					ResponseMassive = append(ResponseMassive,
						"\n"+ "[Пара номер]: " + lessonNum(startTime) + "\n",
						"[Начало пары]: " + startTime + "-" + finishTime + "\n",
						"[Предмет]: " + nameLesson + "\n",
						"[Аудитория]: " + kabinetNum + "\n",
						"[Преподаватель]: " + nameTeacher + "\n",
						"[Тип]: " + typeLesson + "\n")
				}
				return true
			})
		}
		return true
	})
}
package model

import (
	"log"
	"strings"
	"time"

	"github.com/tg_bot_timetable/internal/handler"

	"github.com/PuerkitoBio/goquery"
)

// parseDate finds html selection with date that user chosen.
func parseDate(date, url string, location *time.Location) *goquery.Selection {
	var dateSelection *goquery.Selection

	document := handler.LoadHtmlPage(url)

	document.Find("div.one_day-wrap").EachWithBreak(func(index int, tag *goquery.Selection) bool {
		everDTag := tag.Find("div.everD")
		everDTagValue := strings.ReplaceAll(everDTag.Text(), " ", "")

		if everDTagValue == getDate(date, location).shortDate {
			dateSelection = tag
			return false
		}
		return true
})
	return dateSelection
}

// ParseLessons finds study group`s lessons for date that user chosen,
// transforms Schedule using methods
// and returns pointer to it.
func ParseLessons(groupId, date, url string, location *time.Location) *Schedule {

	var (
		startTime   string         // Начало пары
		finishTime  string         // Конец пары
		lessonName  string         // Название предмета
		roomNumber  string         // Номер аудитории
		teacherName string         // Имя преподавателя
		lessonType  string         // Тип пары (лекция/практика/лабораторная)
	)

	dateSelection := parseDate(date, url, location)
	schedule := NewSchedule()

	if isNilSelection(dateSelection) {
		schedule.NotFound()
		schedule.AddDate(date, location)
		schedule.AddGroupId(groupId)

		return schedule
	}

	dateSelection.Find(".one_lesson").EachWithBreak(func(index int, tag *goquery.Selection) bool {
		lessonName = tag.Find(".names_of_less").Text()
		if lessonName != "" {
			startTime = tag.Find(".starting_less").Text()
			finishTime = tag.Find(".finished_less").Text()
			roomNumber = tag.Find(".kabinet_of_less").Text()
			teacherName = tag.Find(".name_of_teacher").Text()
			lessonType = tag.Find(".type_less").Text()

			schedule.AddLessons(
				startTime,
				finishTime,
				lessonName,
				roomNumber,
				strings.TrimSpace(teacherName),
				lessonType)
		}
		return true
	})

	if !schedule.ScheduleExists() {
		schedule.AddDate(date, location)
		schedule.AddGroupId(groupId)
		schedule.NoLessons()

		return schedule
	}

	schedule.AddDate(date, location)
	schedule.AddGroupId(groupId)

	return schedule
}

// isNilSelection checks if html selection exists.
func isNilSelection(selection *goquery.Selection) bool {
	if selection == nil {
		log.Printf("tag not found: %v ", selection)
		return true
	}
	return false
}
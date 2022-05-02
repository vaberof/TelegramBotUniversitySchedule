package model

import (
	"log"
	"strings"
	"time"

	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/handler"

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

// parseWeekDate finds html selection with date for the week.
func parseWeekDate(date, url string) *goquery.Selection {
	var dateSelection *goquery.Selection

	document := handler.LoadHtmlPage(url)

	document.Find("div.one_day-wrap").EachWithBreak(func(index int, tag *goquery.Selection) bool {
		everDTag := tag.Find("div.everD")
		everDTagValue := strings.ReplaceAll(everDTag.Text(), " ", "")

		if everDTagValue == date {
			dateSelection = tag
			return false
		}
		return true
	})
	return dateSelection
}

// ParseDayLessons finds study group`s lessons for date that user chosen,
// transforms Schedule using methods
// and returns pointer to it.
func ParseDayLessons(groupId, date, url string, location *time.Location) *Schedule {

	var (
		startTime   string // Начало пары
		finishTime  string // Конец пары
		lessonName  string // Название предмета
		roomNumber  string // Номер аудитории
		teacherName string // Имя преподавателя
		lessonType  string // Тип пары (лекция/практика/лабораторная)
	)

	dateSelection := parseDate(date, url, location)
	dates := getDate(date, location)

	schedule := NewSchedule()

	if isNilSelection(dateSelection) {
		schedule.NotFound()
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
		schedule.NoLessons()
	}

	schedule.AddDate(dates)
	schedule.AddGroupId(groupId)

	return schedule
}

// ParseWeekLessons finds study group`s lessons for the week,
// transforms Schedule using methods
// and returns pointer to it.
func ParseWeekLessons(groupId, date, url string, location *time.Location) *Schedule {

	var (
		startTime   string   // Начало пары
		finishTime  string   // Конец пары
		lessonName  string   // Название предмета
		roomNumber  string   // Номер аудитории
		teacherName string   // Имя преподавателя
		lessonType  string   // Тип пары (лекция/практика/лабораторная)
		lessons     []string // Для проверки на наличие пар в конкретный день во время парсинга
	)

	dates := getDate(date, location)

	schedule := NewSchedule()
	schedule.AddGroupId(groupId)

	for day := 0; day <= 6; day++ {
		lessons = []string{}

		schedule.AddWeekDate(dates, day)
		dateSelection := parseWeekDate(dates.weekShortDates[day], url)

		if isNilSelection(dateSelection) {
			schedule.NotFound()
			continue
		}

		dateSelection.Find(".one_lesson").EachWithBreak(func(index int, tag *goquery.Selection) bool {
			lessonName = tag.Find(".names_of_less").Text()
			if lessonName != "" {
				startTime = tag.Find(".starting_less").Text()
				finishTime = tag.Find(".finished_less").Text()
				roomNumber = tag.Find(".kabinet_of_less").Text()
				teacherName = tag.Find(".name_of_teacher").Text()
				lessonType = tag.Find(".type_less").Text()

				lessons = append(lessons, "lesson")

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
		if len(lessons) == 0 {
			schedule.NoLessons()
			continue
		}
	}

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

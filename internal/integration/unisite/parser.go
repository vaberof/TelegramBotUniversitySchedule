package unisite

import (
	"log"
	"strings"
	"time"

	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/model"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/pkg/http"

	"github.com/PuerkitoBio/goquery"
)

// parseDate finds html selection with date that user chosen.
func parseDate(parseData *model.ParseData, url string) *goquery.Selection {
	var dateSelection *goquery.Selection

	document := http.LoadHtmlPage(url)

	document.Find("div.one_day-wrap").EachWithBreak(func(index int, tag *goquery.Selection) bool {
		everDTag := tag.Find("div.everD")
		everDTagValue := strings.ReplaceAll(everDTag.Text(), " ", "")

		if everDTagValue == parseData.Date.Format("02.01") {
			dateSelection = tag

			return false
		}
		return true
	})

	return dateSelection
}

// parseWeekDate finds html selection with date for the week.
func parseWeekDate(parseDate time.Time, url string) *goquery.Selection {
	var dateSelection *goquery.Selection

	document := http.LoadHtmlPage(url)

	document.Find("div.one_day-wrap").EachWithBreak(func(index int, tag *goquery.Selection) bool {
		everDTag := tag.Find("div.everD")
		everDTagValue := strings.ReplaceAll(everDTag.Text(), " ", "")

		if everDTagValue == parseDate.Format("02.01") {
			dateSelection = tag
			return false
		}
		return true
	})
	return dateSelection
}

// ParseDayLessons finds study group`s lessons for date that user chosen,
// adds them to model.Schedule and returns pointer to it.
func ParseDayLessons(inputCallback, url string, parseData *model.ParseData) *model.Schedule {

	var (
		startTime   string   // Начало пары
		finishTime  string   // Конец пары
		lessonName  string   // Название предмета
		roomNumber  string   // Номер аудитории
		teacherName string   // Имя преподавателя
		lessonType  string   // Тип пары (лекция/практика/лабораторная)
		lessons     []string // check if we have lesson on certain day while parsing
	)

	dateSelection := parseDate(parseData, url)

	daySchedule := model.NewDaySchedule()
	schedule := model.Schedule{}

	if isNilSelection(dateSelection) {
		daySchedule.NotFoundSchedule("not found")
		schedule[inputCallback] = *daySchedule
		return &schedule
	}

	dateSelection.Find(".one_lesson").EachWithBreak(func(index int, tag *goquery.Selection) bool {
		lessonName = tag.Find(".names_of_less").Text()
		if lessonName != "" {
			startTime = tag.Find(".starting_less").Text()
			finishTime = tag.Find(".finished_less").Text()
			roomNumber = tag.Find(".kabinet_of_less").Text()
			teacherName = tag.Find(".name_of_teacher").Text()
			lessonType = tag.Find(".type_less").Text()

			lessons = append(lessons, "have lessons")

			daySchedule.AddLessons(
				lessonName,
				startTime,
				finishTime,
				lessonType,
				roomNumber,
				strings.TrimSpace(teacherName))
		}
		return true
	})

	if !haveLessons(lessons) {
		daySchedule.HaveNoLessons("no lessons")
		schedule[inputCallback] = *daySchedule
		return &schedule
	}

	schedule[inputCallback] = *daySchedule
	return &schedule
}

// ParseWeekLessons finds study group`s lessons for date that user chosen,
// adds them to model.Schedule and returns pointer to it.
func ParseWeekLessons(inputCallback, url string, parseData *model.ParseData) *model.Schedule {

	var (
		startTime   string   // Начало пары
		finishTime  string   // Конец пары
		lessonName  string   // Название предмета
		roomNumber  string   // Номер аудитории
		teacherName string   // Имя преподавателя
		lessonType  string   // Тип пары (лекция/практика/лабораторная)
		lessons     []string // check if we have lesson on certain day while parsing
	)

	daySchedule := model.NewDaySchedule()
	schedule := model.Schedule{}

	for day := 0; day <= 6; day++ {
		lessons = []string{}

		dateSelection := parseWeekDate(parseData.Dates[day], url)

		if isNilSelection(dateSelection) {
			daySchedule.NotFoundSchedule("not found")
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

				lessons = append(lessons, "have lessons")

				daySchedule.AddLessons(
					lessonName,
					startTime,
					finishTime,
					lessonType,
					roomNumber,
					strings.TrimSpace(teacherName))
			}
			return true
		})

		if !haveLessons(lessons) {
			daySchedule.HaveNoLessons("no lessons")
			continue
		}
	}

	schedule[inputCallback] = *daySchedule
	return &schedule
}

// haveLessons checks if we have lessons while parsing.
func haveLessons(lessons []string) bool {
	return len(lessons) != 0
}

// isNilSelection checks if html selection exists.
func isNilSelection(selection *goquery.Selection) bool {
	if selection == nil {
		log.Printf("tag not found: %v ", selection)
		return true
	}
	return false
}

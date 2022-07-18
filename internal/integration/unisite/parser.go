package unisite

import (
	log "github.com/sirupsen/logrus"
	"strings"
	"time"

	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/model"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/pkg/http"

	"github.com/PuerkitoBio/goquery"
)

// parseDate finds html selection with date that user chosen.
func parseDate(date time.Time, url string) (*goquery.Selection, error) {
	var dateSelection *goquery.Selection

	document, err := http.LoadHtmlPage(url)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"func":  "parseDate",
		}).Error("Failed while parsing")

		return dateSelection, err
	}

	document.Find("div.one_day-wrap").EachWithBreak(func(index int, tag *goquery.Selection) bool {
		everDTag := tag.Find("div.everD")
		everDTagValue := strings.ReplaceAll(everDTag.Text(), " ", "")

		if everDTagValue == date.Format("02.01") {
			dateSelection = tag

			return false
		}
		return true
	})

	return dateSelection, nil
}

// ParseDayLessons finds study group`s lessons for day that user chosen,
// adds them to model.Schedule and returns pointer to it.
func ParseDayLessons(inputCallback, url string, date time.Time) *model.Schedule {

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

	dateSelection, err := parseDate(date, url)
	if err != nil {
		daySchedule = addHttpErrorMsg(daySchedule)
		schedule[inputCallback] = *daySchedule
		return &schedule
	}

	if isNilSelection(dateSelection) {
		daySchedule = addNotFoundLessonsMsg(daySchedule)
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

	if !isHaveLessons(lessons) {
		daySchedule = addNoLessonsMsg(daySchedule)
		schedule[inputCallback] = *daySchedule
		return &schedule
	}

	schedule[inputCallback] = *daySchedule
	return &schedule
}

// ParseWeekLessons finds study group`s lessons for days that user chosen,
// adds them to model.Schedule and returns pointer to it.
func ParseWeekLessons(inputCallback, url string, dates []time.Time) *model.Schedule {

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

		dateSelection, err := parseDate(dates[day], url)
		if err != nil {
			daySchedule = addHttpErrorMsg(daySchedule)
			schedule[inputCallback] = *daySchedule
			return &schedule
		}

		if isNilSelection(dateSelection) {
			daySchedule = addNotFoundLessonsMsg(daySchedule)
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

		if !isHaveLessons(lessons) {
			daySchedule = addNoLessonsMsg(daySchedule)
			continue
		}
	}

	schedule[inputCallback] = *daySchedule
	return &schedule
}

func addNotFoundLessonsMsg(daySchedule *model.DaySchedule) *model.DaySchedule {
	*daySchedule = append(*daySchedule, model.Lesson{
		Name: "not found",
	})
	return daySchedule
}

func addNoLessonsMsg(daySchedule *model.DaySchedule) *model.DaySchedule {
	*daySchedule = append(*daySchedule, model.Lesson{
		Name: "no lessons",
	})
	return daySchedule
}

func addHttpErrorMsg(daySchedule *model.DaySchedule) *model.DaySchedule {
	*daySchedule = append(*daySchedule, model.Lesson{
		Name: "http error",
	})
	return daySchedule
}

// isHaveLessons checks if we have lessons while parsing.
func isHaveLessons(lessons []string) bool {
	return len(lessons) != 0
}

// isNilSelection checks if html selection exists.
func isNilSelection(selection *goquery.Selection) bool {
	if selection == nil {
		log.Info("Html tag not found")
		return true
	}
	return false
}

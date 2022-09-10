package integration

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/infra/storage"
	"strings"
	"time"
)

type ScheduleApi interface {
	GetSchedule(groupId string, fromDate time.Time, toDate time.Time) (*GetScheduleResponse, error)
}

type GetScheduleResponse struct {
	Lessons []*Lesson
}

type Lesson struct {
	Title           string
	StartTime       string
	FinishTime      string
	Type            string
	RoomId          string
	TeacherFullName string
}

func (s *GetScheduleResponse) addLesson(lessonName, startTime, finishTime, lessonType, roomNumber, teacherName string) {
	s.Lessons = append(s.Lessons, &Lesson{
		lessonName,
		startTime,
		finishTime,
		lessonType,
		roomNumber,
		teacherName,
	})
}

func (httpClient *HttpClient) GetSchedule(groupId string, fromDate time.Time, toDate time.Time) (*GetScheduleResponse, error) {
	return httpClient.parseLessons(groupId, fromDate, toDate)
}

func (httpClient *HttpClient) parseLessons(groupId string, fromDate time.Time, toDate time.Time) (*GetScheduleResponse, error) {

	return nil, nil
}

// parseDayLessons finds study group`s lessons for day that user chosen.
// Returns custom error if http request Timeout occurred.
func (httpClient *HttpClient) parseDayLessons(groupId string, fromDate time.Time, toDate time.Time) (*GetScheduleResponse, error) {

	var (
		startTime   string   // Начало пары
		finishTime  string   // Конец пары
		lessonName  string   // Название предмета
		roomNumber  string   // Номер аудитории
		teacherName string   // Имя преподавателя
		lessonType  string   // Тип пары (лекция/практика/лабораторная)
		lessons     []string // necessary to check if we have lesson on certain day while parsing
	)

	var scheduleResponse *GetScheduleResponse

	var dateSelection, err = httpClient.parseDate(groupId, toDate) // from, to
	if err != nil {
		httpError := fmt.Sprint("Ошибка: превышено время ожидания от сервера")
		return nil, errors.New(httpError)
	}

	if isNilSelection(dateSelection) {
		scheduleResponse = addNotFoundLessonsMsg(scheduleResponse)
		return scheduleResponse, nil
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

			scheduleResponse.addLesson(
				lessonName,
				startTime,
				finishTime,
				lessonType,
				roomNumber,
				strings.TrimSpace(teacherName))
		}
		return true
	})

	if len(lessons) == 0 {
		scheduleResponse = addNoLessonsMsg(scheduleResponse)
		return scheduleResponse, nil
	}

	return scheduleResponse, nil
}

// ParseWeekLessons finds study group`s lessons for week that user chosen.
// Returns custom error if http request Timeout occurred.
//func (httpClient *HttpClient) ParseWeekLessons(inputCallback, url string, dates []time.Time) (*GetScheduleResponse, error) {
//
//	var (
//		startTime   string   // Начало пары
//		finishTime  string   // Конец пары
//		lessonName  string   // Название предмета
//		roomNumber  string   // Номер аудитории
//		teacherName string   // Имя преподавателя
//		lessonType  string   // Тип пары (лекция/практика/лабораторная)
//		lessons     []string // necessary to check if we have lesson on certain day while parsing
//	)
//
//	var daySchedule *domain.DaySchedule
//	var schedule domain.Schedule
//
//	for day := 0; day <= 6; day++ {
//		lessons = []string{}
//
//		dateSelection, err := httpClient.parseDate(dates[day], url)
//		if err != nil {
//			httpError := fmt.Sprint("Ошибка: превышено время ожидания от сервера")
//			return &schedule, errors.New(httpError)
//		}
//
//		if isNilSelection(dateSelection) {
//			daySchedule = addNotFoundLessonsMsg(daySchedule)
//			continue
//		}
//
//		dateSelection.Find(".one_lesson").EachWithBreak(func(index int, tag *goquery.Selection) bool {
//			lessonName = tag.Find(".names_of_less").Text()
//			if lessonName != "" {
//				startTime = tag.Find(".starting_less").Text()
//				finishTime = tag.Find(".finished_less").Text()
//				roomNumber = tag.Find(".kabinet_of_less").Text()
//				teacherName = tag.Find(".name_of_teacher").Text()
//				lessonType = tag.Find(".type_less").Text()
//
//				lessons = append(lessons, "have lessons")
//
//				daySchedule.AddLesson(
//					lessonName,
//					startTime,
//					finishTime,
//					lessonType,
//					roomNumber,
//					strings.TrimSpace(teacherName))
//			}
//			return true
//		})
//
//		if !isHaveLessons(lessons) {
//			daySchedule = addNoLessonsMsg(daySchedule)
//			continue
//		}
//	}
//
//	scheduleDate := domain.Date(inputCallback)
//	schedule[scheduleDate] = *daySchedule
//	log.Printf("week schedule: %v\n", daySchedule)
//	return &schedule, nil
//}

// loadHtmlTemplate loads html page.
func (httpClient *HttpClient) loadHtmlTemplate(groupId string) (*goquery.Document, error) {
	groupStorage := storage.NewGroupStorage()
	url := groupStorage.GetStudyGroupUrl(groupId)

	response, err := httpClient.client.R().Get(*url)
	if err != nil {
		return nil, err
	}

	body := response.Body()
	rBody := bytes.NewReader(body)
	document, err := goquery.NewDocumentFromReader(rBody)

	if err != nil {
		log.WithFields(log.Fields{
			"body":  rBody,
			"error": err,
			"func":  "loadHtmlTemplate",
		}).Error("Data cannot be parsed as html")

		return nil, err
	}

	return document, nil
}

// parseDate finds html selection with date that user chosen.
func (httpClient *HttpClient) parseDate(groupId string, date time.Time) (*goquery.Selection, error) {
	var dateSelection *goquery.Selection

	document, err := httpClient.loadHtmlTemplate(groupId)
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

func addNotFoundLessonsMsg(scheduleResponse *GetScheduleResponse) *GetScheduleResponse {
	scheduleResponse.Lessons = append(scheduleResponse.Lessons, &Lesson{
		Title: "not found",
	})
	return scheduleResponse
}

func addNoLessonsMsg(scheduleResponse *GetScheduleResponse) *GetScheduleResponse {
	scheduleResponse.Lessons = append(scheduleResponse.Lessons, &Lesson{
		Title: "no lessons",
	})
	return scheduleResponse
}

// isNilSelection checks if html selection exists.
func isNilSelection(selection *goquery.Selection) bool {
	if selection == nil {
		log.Info("Html tag not found")
		return true
	}
	return false
}

package integration

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/infra/storage"
	"github.com/vaberof/TelegramBotUniversitySchedule/pkg/xstrconv"
	"io"
	"strings"
	"time"
)

type ScheduleApi interface {
	GetSchedule(group *storage.Group, from time.Time, to time.Time) (*GetScheduleResponse, error)
}

type GetScheduleResponse struct {
	Lessons []*Lesson
}

func (r *GetScheduleResponse) addLesson(title, startTime, finishTime, lessonType, roomId, teacherFullName string) {
	r.Lessons = append(r.Lessons, &Lesson{
		title,
		startTime,
		finishTime,
		lessonType,
		roomId,
		teacherFullName,
	})
}

type Lesson struct {
	Title           string
	StartTime       string
	FinishTime      string
	Type            string
	RoomId          string
	TeacherFullName string
}

func (httpClient *HttpClient) GetSchedule(group *storage.Group, from time.Time, to time.Time) (*GetScheduleResponse, error) {
	return httpClient.parseLessons(group, from, to)
}

func (httpClient *HttpClient) parseLessons(group *storage.Group, from time.Time, to time.Time) (*GetScheduleResponse, error) {
	htmlTemplate, err := httpClient.getHtmlTemplate(group.ExternalId)
	if err != nil {
		return nil, err
	}

	strDate, err := xstrconv.ConvertDateToStr(from, to)
	if err != nil {
		return nil, err
	}

	switch strDate {
	case "Today", "Tomorrow":
		return httpClient.parseDayLessons(htmlTemplate, to)
	default:
		return httpClient.parseWeekLessons(htmlTemplate, from)
	}
}

// parseDayLessons finds study group`s lessons for day that user chosen.
// Returns custom error if http request Timeout occurred.
func (httpClient *HttpClient) parseDayLessons(htmlDocument *goquery.Document, to time.Time) (*GetScheduleResponse, error) {
	var (
		startTime       string   // Начало пары
		finishTime      string   // Конец пары
		title           string   // Название предмета
		roomId          string   // Номер аудитории
		teacherFullName string   // Имя преподавателя
		lessonType      string   // Тип пары (лекция/практика/лабораторная)
		validLessons    []string // necessary to check if we have lesson on certain day while parsing
	)

	dateSelection, err := httpClient.parseDate(htmlDocument, to)
	if err != nil {
		return nil, err
	}

	var getScheduleResponse GetScheduleResponse

	if httpClient.isNilSelection(dateSelection) {
		getScheduleResponse = *addNotFoundLessonsMsg(&getScheduleResponse)
		return &getScheduleResponse, nil
	}

	dateSelection.Find(".one_lesson").EachWithBreak(func(index int, tag *goquery.Selection) bool {
		title = tag.Find(".names_of_less").Text()
		if title != "" {
			startTime = tag.Find(".starting_less").Text()
			finishTime = tag.Find(".finished_less").Text()
			roomId = tag.Find(".kabinet_of_less").Text()
			teacherFullName = tag.Find(".name_of_teacher").Text()
			lessonType = tag.Find(".type_less").Text()

			validLessons = append(validLessons, "have lessons")

			getScheduleResponse.addLesson(
				title,
				startTime,
				finishTime,
				lessonType,
				roomId,
				strings.TrimSpace(teacherFullName))
		}
		return true
	})

	if len(validLessons) == 0 {
		getScheduleResponse = *addNoLessonsMsg(&getScheduleResponse)
		return &getScheduleResponse, nil
	}

	return &getScheduleResponse, nil
}

// parseWeekLessons finds study group`s lessons for week that user chosen.
// Returns custom error if http request Timeout occurred.
func (httpClient *HttpClient) parseWeekLessons(htmlDocument *goquery.Document, from time.Time) (*GetScheduleResponse, error) {
	var (
		startTime   string   // Начало пары
		finishTime  string   // Конец пары
		title       string   // Название предмета
		roomNumber  string   // Номер аудитории
		teacherName string   // Имя преподавателя
		lessonType  string   // Тип пары (лекция/практика/лабораторная)
		lessons     []string // necessary to check if we have lesson on certain day while parsing
	)

	var getScheduleResponse GetScheduleResponse

	for weekday := 1; weekday <= 7; weekday++ {
		lessons = []string{}

		dateSelection, err := httpClient.parseDate(htmlDocument, from)
		if err != nil {
			return nil, err
		}
		from.Add(24 * time.Hour)

		if httpClient.isNilSelection(dateSelection) {
			getScheduleResponse = *addNotFoundLessonsMsg(&getScheduleResponse)
			continue
		}

		dateSelection.Find(".one_lesson").EachWithBreak(func(index int, tag *goquery.Selection) bool {
			title = tag.Find(".names_of_less").Text()
			if title != "" {
				startTime = tag.Find(".starting_less").Text()
				finishTime = tag.Find(".finished_less").Text()
				roomNumber = tag.Find(".kabinet_of_less").Text()
				teacherName = tag.Find(".name_of_teacher").Text()
				lessonType = tag.Find(".type_less").Text()

				lessons = append(lessons, "have lessons")

				getScheduleResponse.addLesson(
					title,
					startTime,
					finishTime,
					lessonType,
					roomNumber,
					strings.TrimSpace(teacherName))
			}
			return true
		})

		if len(lessons) == 0 {
			getScheduleResponse = *addNoLessonsMsg(&getScheduleResponse)
			continue
		}
	}

	return &getScheduleResponse, nil
}

func (httpClient *HttpClient) makeRequest(queryParams string) (io.Reader, error) {
	response, err := httpClient.client.R().Get(httpClient.host + queryParams)
	if err != nil {
		//log.Printf("error http response: %v", response)
		return nil, err
	}

	body := response.Body()
	rBody := bytes.NewReader(body)

	return rBody, nil
}

// createHtmlTemplate loads html page.
func (httpClient *HttpClient) createHtmlTemplate(responseBody io.Reader) (*goquery.Document, error) {
	document, err := goquery.NewDocumentFromReader(responseBody)
	if err != nil {
		log.WithFields(log.Fields{
			"responseBody": responseBody,
			"error":        err,
			"func":         "createHtmlTemplate",
		}).Error("Data cannot be parsed as html")

		return nil, err
	}

	return document, nil
}

func (httpClient *HttpClient) getHtmlTemplate(queryParams string) (*goquery.Document, error) {
	responseBody, err := httpClient.makeRequest(queryParams)
	if err != nil {
		//httpError := fmt.Sprint("Ошибка: превышено время ожидания от сервера")
		return nil, err
	}

	htmlTemplate, err := httpClient.createHtmlTemplate(responseBody)
	if err != nil {
		return nil, err
	}

	return htmlTemplate, nil
}

// parseDate finds html selection with date that user chosen.
func (httpClient *HttpClient) parseDate(htmlDocument *goquery.Document, date time.Time) (*goquery.Selection, error) {
	var dateSelection *goquery.Selection

	htmlDocument.Find("div.one_day-wrap").EachWithBreak(func(index int, tag *goquery.Selection) bool {
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

// isNilSelection checks if html selection exists.
func (httpClient *HttpClient) isNilSelection(selection *goquery.Selection) bool {
	if selection == nil {
		log.Info("Html tag not found")
		return true
	}
	return false
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

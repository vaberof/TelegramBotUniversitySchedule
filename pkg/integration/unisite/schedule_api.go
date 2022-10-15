package integration

import (
	"bytes"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
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
	return httpClient.getScheduleResponse(group, from, to)
}

func (httpClient *HttpClient) getScheduleResponse(group *storage.Group, from time.Time, to time.Time) (*GetScheduleResponse, error) {
	htmlTemplate, err := httpClient.getHtmlTemplate(group.ExternalId)
	if err != nil {
		return nil, err
	}

	return httpClient.parseLessons(htmlTemplate, from, to)
}

func (httpClient *HttpClient) parseLessons(htmlTemplate *goquery.Document, from time.Time, to time.Time) (*GetScheduleResponse, error) {
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

func (httpClient *HttpClient) parseDayLessons(htmlDocument *goquery.Document, to time.Time) (*GetScheduleResponse, error) {
	dateSelection, err := httpClient.parseDate(htmlDocument, to)
	if err != nil {
		return nil, err
	}

	if httpClient.isNilSelection(dateSelection) {
		var getScheduleResponse GetScheduleResponse
		getScheduleResponse = *addNotFoundLessonsMsg(&getScheduleResponse)
		return &getScheduleResponse, nil
	}

	var getScheduleResponse GetScheduleResponse
	var lessons []string

	httpClient.parseDateSelection(&getScheduleResponse, &lessons, dateSelection)

	if len(lessons) == 0 {
		getScheduleResponse = *addNoLessonsMsg(&getScheduleResponse)
		return &getScheduleResponse, nil
	}
	return &getScheduleResponse, nil
}

func (httpClient *HttpClient) parseWeekLessons(htmlDocument *goquery.Document, from time.Time) (*GetScheduleResponse, error) {
	var getScheduleResponse GetScheduleResponse

	for weekday := 1; weekday <= 7; weekday++ {
		var lessons []string

		dateSelection, err := httpClient.parseDate(htmlDocument, from)
		if err != nil {
			return nil, err
		}

		from.Add(24 * time.Hour)

		if httpClient.isNilSelection(dateSelection) {
			getScheduleResponse = *addNotFoundLessonsMsg(&getScheduleResponse)
			continue
		}

		httpClient.parseDateSelection(&getScheduleResponse, &lessons, dateSelection)

		if len(lessons) == 0 {
			getScheduleResponse = *addNoLessonsMsg(&getScheduleResponse)
			continue
		}
	}

	return &getScheduleResponse, nil
}

func (httpClient *HttpClient) parseDateSelection(
	getScheduleResponse *GetScheduleResponse,
	lessons *[]string,
	dateSelection *goquery.Selection) {

	var (
		startTime       string // Начало пары
		finishTime      string // Конец пары
		title           string // Название предмета
		roomId          string // Номер аудитории
		teacherFullName string // Имя преподавателя
		lessonType      string // Тип пары (лекция/практика/лабораторная)
	)

	dateSelection.Find(".one_lesson").EachWithBreak(func(index int, tag *goquery.Selection) bool {
		title = tag.Find(".names_of_less").Text()
		if title != "" {
			startTime = tag.Find(".starting_less").Text()
			finishTime = tag.Find(".finished_less").Text()
			roomId = tag.Find(".kabinet_of_less").Text()
			teacherFullName = tag.Find(".name_of_teacher").Text()
			lessonType = tag.Find(".type_less").Text()

			*lessons = append(*lessons, "have lessons")

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
}

func (httpClient *HttpClient) getHtmlTemplate(queryParams string) (*goquery.Document, error) {
	response, err := httpClient.makeRequest(queryParams)
	if err != nil {
		//httpError := fmt.Sprint("Ошибка: превышено время ожидания от сервера")
		return nil, err
	}

	responseBody, err := httpClient.getResponseBody(response)
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

func (httpClient *HttpClient) makeRequest(queryParams string) (*resty.Response, error) {
	response, err := httpClient.client.R().Get(httpClient.host + queryParams)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (httpClient *HttpClient) getResponseBody(response *resty.Response) (io.Reader, error) {
	if response == nil {
		return nil, errors.New("response is nil")
	}
	body := response.Body()
	rBody := bytes.NewReader(body)
	return rBody, nil
}

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

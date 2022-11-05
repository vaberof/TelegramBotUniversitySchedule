package infra

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/infra/storage/postgres/grouppg"
	integration "github.com/vaberof/TelegramBotUniversitySchedule/pkg/integration/unisite"
	"time"
)

type GetScheduleResponseService struct {
	scheduleApi  ScheduleApi
	groupStorage GroupStorage
}

func NewGetScheduleResponseService(httpClient *integration.HttpClient, groupStoragePostgres *grouppg.GroupStoragePostgres) *GetScheduleResponseService {
	return &GetScheduleResponseService{
		scheduleApi:  httpClient,
		groupStorage: groupStoragePostgres,
	}
}

func (s *GetScheduleResponseService) GetSchedule(groupId string, from time.Time, to time.Time) (*GetScheduleResponse, error) {
	groupExternalId := s.groupStorage.GetGroupExternalId(groupId)
	if groupExternalId == nil {
		return nil, errors.New(fmt.Sprintf("Группы '%s' не существует", groupId))
	}
	log.Printf("group name: %s, query params: %s", groupId, *groupExternalId)

	getScheduleResponse, err := s.scheduleApi.GetSchedule(*groupExternalId, from, to)
	if err != nil {
		return nil, err
	}
	log.Printf("schedule response from scheduleApi: %v", getScheduleResponse)

	if getScheduleResponse == nil || getScheduleResponse.Lessons == nil {
		return nil, errors.New("schedule api response is nil")
	}

	infraGetScheduleResponse := s.getScheduleRespToInfraSchedule(getScheduleResponse)
	return infraGetScheduleResponse, nil
}

func (s *GetScheduleResponseService) getScheduleRespToInfraSchedule(getScheduleResponse *integration.GetScheduleResponse) *GetScheduleResponse {
	var infraGetScheduleResponse GetScheduleResponse
	infraGetScheduleResponse.Lessons = s.getScheduleRespLessonsToInfraLessons(getScheduleResponse.Lessons)
	return &infraGetScheduleResponse
}

func (s *GetScheduleResponseService) getScheduleRespLessonsToInfraLessons(getScheduleRespLessons []*integration.Lesson) []*Lesson {
	var lessons []*Lesson

	for i := 0; i < len(getScheduleRespLessons); i++ {
		lesson := s.getScheduleRespLessonToInfraLesson(getScheduleRespLessons[i])
		lessons = append(lessons, lesson)
	}
	return lessons
}

func (s *GetScheduleResponseService) getScheduleRespLessonToInfraLesson(getScheduleRespLesson *integration.Lesson) *Lesson {
	var lesson Lesson

	lesson.Title = getScheduleRespLesson.Title
	lesson.StartTime = getScheduleRespLesson.StartTime
	lesson.FinishTime = getScheduleRespLesson.FinishTime
	lesson.Type = getScheduleRespLesson.Type
	lesson.RoomId = getScheduleRespLesson.RoomId
	lesson.TeacherFullName = getScheduleRespLesson.TeacherFullName

	return &lesson
}

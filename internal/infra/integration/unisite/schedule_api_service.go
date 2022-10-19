package infra

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	integration "github.com/vaberof/TelegramBotUniversitySchedule/pkg/integration/unisite"
	"time"
)

type GetScheduleResponseApiService struct {
	scheduleApi     *GetScheduleResponseApi
	groupStorageApi *GroupStorage
}

func NewGetScheduleResponseApiService(scheduleApi *GetScheduleResponseApi, groupStorage *GroupStorage) *GetScheduleResponseApiService {
	return &GetScheduleResponseApiService{
		scheduleApi:     scheduleApi,
		groupStorageApi: groupStorage,
	}
}

func (s *GetScheduleResponseApiService) GetSchedule(groupId string, from time.Time, to time.Time) (*GetScheduleResponse, error) {
	groupExternalId := s.groupStorageApi.GetGroupExternalId(groupId)
	if groupExternalId == nil {
		return nil, errors.New(fmt.Sprintf("Группы '%s' не существует", groupId))
	}
	log.Printf("group name: %s, query params: %s", groupId, *groupExternalId)

	getScheduleResponse, err := s.scheduleApi.GetSchedule(*groupExternalId, from, to)
	if err != nil {
		return nil, err
	}

	if getScheduleResponse == nil || getScheduleResponse.Lessons == nil {
		return nil, errors.New("schedule api response is nil")
	}

	infraGetScheduleResponse := s.GetScheduleRespToInfraSchedule(getScheduleResponse)
	return infraGetScheduleResponse, nil
}

func (s *GetScheduleResponseApiService) GetScheduleRespToInfraSchedule(getScheduleResponse *integration.GetScheduleResponse) *GetScheduleResponse {
	var scheduleResponse GetScheduleResponse
	scheduleResponse.Lessons = s.GetScheduleRespLessonsToInfraLessons(getScheduleResponse.Lessons)
	return &scheduleResponse
}

func (s *GetScheduleResponseApiService) GetScheduleRespLessonsToInfraLessons(getScheduleRespLessons []*integration.Lesson) []*Lesson {
	var lessons []*Lesson

	for i := 0; i < len(getScheduleRespLessons); i++ {
		lesson := s.GetScheduleRespLessonToInfraLesson(getScheduleRespLessons[i])
		lessons = append(lessons, lesson)
	}

	return lessons
}

func (s *GetScheduleResponseApiService) GetScheduleRespLessonToInfraLesson(getScheduleRespLesson *integration.Lesson) *Lesson {
	var lesson Lesson

	lesson.Title = getScheduleRespLesson.Title
	lesson.StartTime = getScheduleRespLesson.StartTime
	lesson.FinishTime = getScheduleRespLesson.FinishTime
	lesson.Type = getScheduleRespLesson.Type
	lesson.RoomId = getScheduleRespLesson.RoomId
	lesson.TeacherFullName = getScheduleRespLesson.TeacherFullName

	return &lesson
}

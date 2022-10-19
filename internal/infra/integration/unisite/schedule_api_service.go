package infra

import (
	"errors"
	integration "github.com/vaberof/TelegramBotUniversitySchedule/pkg/integration/unisite"
	"time"
)

type GetScheduleResponseApiService struct {
	scheduleApi *GetScheduleResponseApi
}

func NewGetScheduleResponseApiService(scheduleApi *GetScheduleResponseApi) *GetScheduleResponseApiService {
	return &GetScheduleResponseApiService{
		scheduleApi: scheduleApi,
	}
}

func (s *GetScheduleResponseApiService) GetSchedule(studyGroupQueryParams string, from time.Time, to time.Time) (*GetScheduleResponse, error) {
	getScheduleResponse, err := s.scheduleApi.GetSchedule(studyGroupQueryParams, from, to)
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

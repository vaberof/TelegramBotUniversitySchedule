package domain

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/infra/storage"
	integration "github.com/vaberof/TelegramBotUniversitySchedule/pkg/integration/unisite"
	"github.com/vaberof/TelegramBotUniversitySchedule/pkg/xtimeconv"
	"time"
)

type ScheduleService struct {
	scheduleStorageApi *ScheduleStorage
	groupStorageApi    *GroupStorage
	httpClient         *integration.HttpClient
}

func NewScheduleService(scheduleStorage *ScheduleStorage, groupStorage *GroupStorage, httpClient *integration.HttpClient) *ScheduleService {
	return &ScheduleService{
		scheduleStorageApi: scheduleStorage,
		groupStorageApi:    groupStorage,
		httpClient:         httpClient,
	}
}

func (s *ScheduleService) GetSchedule(groupId string, from time.Time, to time.Time) (*Schedule, error) {
	return s.getScheduleImpl(groupId, from, to)
}

func (s *ScheduleService) getScheduleImpl(groupId string, from time.Time, to time.Time) (*Schedule, error) {
	studyGroupQueryParams := s.groupStorageApi.GetStudyGroupQueryParams(groupId)
	if studyGroupQueryParams == nil {
		return nil, errors.New(fmt.Sprintf("Группы '%s' не существует", groupId))
	}
	log.Printf("group name: %s, query params: %s", groupId, *studyGroupQueryParams)

	cachedLessons, err := s.scheduleStorageApi.GetCachedLessons(groupId, from, to)
	if cachedLessons == nil || err != nil {
		getScheduleResponse, err := s.callScheduleApi(*studyGroupQueryParams, from, to)
		if err != nil {
			return nil, err
		}

		err = s.cacheLessons(groupId, getScheduleResponse.Lessons, from, to)
		if err != nil {
			return nil, err
		}

		schedule, err := s.respScheduleToDomain(getScheduleResponse, from, to)
		if err != nil {
			return nil, err
		}

		return schedule, nil
	}

	schedule, err := s.storageLessonsToDomain(cachedLessons, from, to)
	if err != nil {
		return nil, err
	}

	return schedule, nil
}

func (s *ScheduleService) cacheLessons(groupId string, lessons []*integration.Lesson, from time.Time, to time.Time) error {
	storageLessons, err := s.respLessonsToStorage(lessons)
	if err != nil {
		return err
	}

	err = s.scheduleStorageApi.SaveLessons(groupId, from, to, storageLessons)
	if err != nil {
		return err
	}

	return nil
}

func (s *ScheduleService) callScheduleApi(studyGroupQueryParams string, from time.Time, to time.Time) (*integration.GetScheduleResponse, error) {
	getScheduleResponse, err := s.callApi(s.httpClient, studyGroupQueryParams, from, to)
	log.Printf("schedule response from scheduleApi: %v", getScheduleResponse)
	if err != nil {
		return nil, err
	}

	return getScheduleResponse, nil
}

func (s *ScheduleService) callApi(scheduleApi integration.ScheduleApi, studyGroupQueryParams string, from time.Time, to time.Time) (*integration.GetScheduleResponse, error) {
	return scheduleApi.GetSchedule(studyGroupQueryParams, from, to)
}

func (s *ScheduleService) respScheduleToDomain(getScheduleResponse *integration.GetScheduleResponse, from time.Time, to time.Time) (*Schedule, error) {
	daySchedule := s.respLessonsToDomain(getScheduleResponse.Lessons)

	dateString, err := xtimeconv.FromTimeToString(from, to)
	if err != nil {
		return nil, err
	}
	log.Printf("dateString in service: %v\n", dateString)

	schedule := make(Schedule)
	schedule[Date(dateString)] = daySchedule

	return &schedule, nil
}

func (s *ScheduleService) respLessonsToDomain(respLessons []*integration.Lesson) *DaySchedule {
	var daySchedule DaySchedule

	for i := 0; i < len(respLessons); i++ {
		lesson := s.respLessonToDomain(respLessons[i])
		daySchedule = append(daySchedule, lesson)
	}

	return &daySchedule
}

func (s *ScheduleService) respLessonToDomain(respLesson *integration.Lesson) *Lesson {
	var lesson Lesson

	lesson.Title = respLesson.Title
	lesson.StartTime = respLesson.StartTime
	lesson.FinishTime = respLesson.FinishTime
	lesson.Type = respLesson.Type
	lesson.RoomId = respLesson.RoomId
	lesson.TeacherFullName = respLesson.TeacherFullName

	return &lesson
}

func (s *ScheduleService) respLessonsToStorage(respLessons []*integration.Lesson) ([]*storage.Lesson, error) {
	var lessons []*storage.Lesson

	for i := 0; i < len(respLessons); i++ {
		lesson := s.respLessonToStorage(respLessons[i])
		lessons = append(lessons, lesson)
	}
	return lessons, nil
}

func (s *ScheduleService) respLessonToStorage(respLesson *integration.Lesson) *storage.Lesson {
	var lesson storage.Lesson

	lesson.Title = respLesson.Title
	lesson.StartTime = respLesson.StartTime
	lesson.FinishTime = respLesson.FinishTime
	lesson.Type = respLesson.Type
	lesson.RoomId = respLesson.RoomId
	lesson.TeacherFullName = respLesson.TeacherFullName

	return &lesson
}

func (s *ScheduleService) storageLessonsToDomain(storageLessons []*storage.Lesson, from time.Time, to time.Time) (*Schedule, error) {
	var daySchedule DaySchedule

	for i := 0; i < len(storageLessons); i++ {
		lesson := s.storageLessonToDomain(storageLessons[i])
		daySchedule = append(daySchedule, lesson)
	}

	strDate, err := xtimeconv.FromTimeToString(from, to)
	if err != nil {
		return nil, err
	}
	log.Printf("strDate in service: %v\n", strDate)

	schedule := make(Schedule)
	schedule[Date(strDate)] = &daySchedule

	return &schedule, nil
}

func (s *ScheduleService) storageLessonToDomain(storageLesson *storage.Lesson) *Lesson {
	var lesson Lesson

	lesson.Title = storageLesson.Title
	lesson.StartTime = storageLesson.StartTime
	lesson.FinishTime = storageLesson.FinishTime
	lesson.Type = storageLesson.Type
	lesson.RoomId = storageLesson.RoomId
	lesson.TeacherFullName = storageLesson.TeacherFullName

	return &lesson
}

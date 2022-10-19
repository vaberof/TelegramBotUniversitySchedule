package domain

import (
	log "github.com/sirupsen/logrus"
	infra "github.com/vaberof/TelegramBotUniversitySchedule/internal/infra/integration/unisite"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/infra/storage"
	"github.com/vaberof/TelegramBotUniversitySchedule/pkg/xtimeconv"
	"time"
)

type ScheduleService struct {
	scheduleStorageApi *ScheduleStorage
	scheduleApi        *GetScheduleResponseApi
}

func NewScheduleService(scheduleStorage *ScheduleStorage, scheduleApi *GetScheduleResponseApi) *ScheduleService {
	return &ScheduleService{
		scheduleStorageApi: scheduleStorage,
		scheduleApi:        scheduleApi,
	}
}

func (s *ScheduleService) GetSchedule(groupId string, from time.Time, to time.Time) (*Schedule, error) {
	return s.getScheduleImpl(groupId, from, to)
}

func (s *ScheduleService) getScheduleImpl(groupId string, from time.Time, to time.Time) (*Schedule, error) {
	cachedLessons, err := s.scheduleStorageApi.GetCachedLessons(groupId, from, to)
	if cachedLessons == nil || err != nil {
		getScheduleResponse, err := s.callScheduleApi(groupId, from, to)
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

func (s *ScheduleService) cacheLessons(groupId string, lessons []*infra.Lesson, from time.Time, to time.Time) error {
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

func (s *ScheduleService) callScheduleApi(groupId string, from time.Time, to time.Time) (*infra.GetScheduleResponse, error) {
	getScheduleResponse, err := s.scheduleApi.GetSchedule(groupId, from, to)
	log.Printf("schedule response from scheduleApi: %v", getScheduleResponse)
	if err != nil {
		return nil, err
	}

	return getScheduleResponse, nil
}

func (s *ScheduleService) respScheduleToDomain(getScheduleResponse *infra.GetScheduleResponse, from time.Time, to time.Time) (*Schedule, error) {
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

func (s *ScheduleService) respLessonsToDomain(respLessons []*infra.Lesson) *DaySchedule {
	var daySchedule DaySchedule

	for i := 0; i < len(respLessons); i++ {
		lesson := s.respLessonToDomain(respLessons[i])
		daySchedule = append(daySchedule, lesson)
	}

	return &daySchedule
}

func (s *ScheduleService) respLessonToDomain(respLesson *infra.Lesson) *Lesson {
	var lesson Lesson

	lesson.Title = respLesson.Title
	lesson.StartTime = respLesson.StartTime
	lesson.FinishTime = respLesson.FinishTime
	lesson.Type = respLesson.Type
	lesson.RoomId = respLesson.RoomId
	lesson.TeacherFullName = respLesson.TeacherFullName

	return &lesson
}

func (s *ScheduleService) respLessonsToStorage(respLessons []*infra.Lesson) ([]*storage.Lesson, error) {
	var lessons []*storage.Lesson

	for i := 0; i < len(respLessons); i++ {
		lesson := s.respLessonToStorage(respLessons[i])
		lessons = append(lessons, lesson)
	}
	return lessons, nil
}

func (s *ScheduleService) respLessonToStorage(respLesson *infra.Lesson) *storage.Lesson {
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

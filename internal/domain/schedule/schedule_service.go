package domain

import (
	log "github.com/sirupsen/logrus"
	infra "github.com/vaberof/TelegramBotUniversitySchedule/internal/infra/integration/unisite"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/infra/storage/postgres/schedulepg"
	"github.com/vaberof/TelegramBotUniversitySchedule/pkg/xtimeconv"
	"time"
)

type ScheduleService struct {
	scheduleApi     *ScheduleApi
	scheduleStorage *ScheduleStorage
}

func NewScheduleService(scheduleApi *ScheduleApi, scheduleStorage *ScheduleStorage) *ScheduleService {
	return &ScheduleService{
		scheduleApi:     scheduleApi,
		scheduleStorage: scheduleStorage,
	}
}

func (s *ScheduleService) GetSchedule(groupId string, from time.Time, to time.Time) (*Schedule, error) {
	return s.getScheduleImpl(groupId, from, to)
}

func (s *ScheduleService) getScheduleImpl(groupId string, from time.Time, to time.Time) (*Schedule, error) {
	cachedLessons, err := s.scheduleStorage.GetLessons(groupId, from, to)
	if cachedLessons == nil || err != nil {
		getScheduleResponse, err := s.callScheduleApi(groupId, from, to)
		if err != nil {
			return nil, err
		}

		if err = s.cacheLessons(groupId, getScheduleResponse.Lessons, from, to); err != nil {
			return nil, err
		}

		schedule, err := s.fromGetScheduleRespToDomainSchedule(getScheduleResponse, from, to)
		if err != nil {
			return nil, err
		}

		return schedule, nil
	}

	schedule, err := s.storageLessonsToDomainSchedule(cachedLessons, from, to)
	if err != nil {
		return nil, err
	}

	return schedule, nil
}

func (s *ScheduleService) callScheduleApi(groupId string, from time.Time, to time.Time) (*infra.GetScheduleResponse, error) {
	getScheduleResponse, err := s.scheduleApi.GetSchedule(groupId, from, to)
	if err != nil {
		return nil, err
	}
	return getScheduleResponse, nil
}

func (s *ScheduleService) cacheLessons(groupId string, lessons []*infra.Lesson, from time.Time, to time.Time) error {
	storageLessons, err := s.respLessonsToStorageLessons(lessons)
	if err != nil {
		return err
	}

	err = s.scheduleStorage.SaveLessons(groupId, from, to, storageLessons)
	if err != nil {
		return err
	}

	return nil
}

func (s *ScheduleService) fromGetScheduleRespToDomainSchedule(getScheduleResponse *infra.GetScheduleResponse,
	from time.Time,
	to time.Time) (*Schedule, error) {

	daySchedule := s.respLessonsToDomainDaySchedule(getScheduleResponse.Lessons)

	dateString, err := xtimeconv.FromTimeRangeToDateString(from, to)
	if err != nil {
		return nil, err
	}
	log.Printf("dateString in service: %v\n", dateString)

	schedule := make(Schedule)
	schedule[Date(dateString)] = daySchedule

	return &schedule, nil
}

func (s *ScheduleService) respLessonsToDomainDaySchedule(respLessons []*infra.Lesson) *DaySchedule {
	var daySchedule DaySchedule

	for i := 0; i < len(respLessons); i++ {
		lesson := s.respLessonToDomainLesson(respLessons[i])
		daySchedule = append(daySchedule, lesson)
	}

	return &daySchedule
}

func (s *ScheduleService) respLessonToDomainLesson(respLesson *infra.Lesson) *Lesson {
	var lesson Lesson

	lesson.Title = respLesson.Title
	lesson.StartTime = respLesson.StartTime
	lesson.FinishTime = respLesson.FinishTime
	lesson.Type = respLesson.Type
	lesson.RoomId = respLesson.RoomId
	lesson.TeacherFullName = respLesson.TeacherFullName

	return &lesson
}

func (s *ScheduleService) respLessonsToStorageLessons(respLessons []*infra.Lesson) ([]*schedulepg.Lesson, error) {
	var lessons []*schedulepg.Lesson

	for i := 0; i < len(respLessons); i++ {
		lesson := s.respLessonToStorageLesson(respLessons[i])
		lessons = append(lessons, lesson)
	}
	return lessons, nil
}

func (s *ScheduleService) respLessonToStorageLesson(respLesson *infra.Lesson) *schedulepg.Lesson {
	var lesson schedulepg.Lesson

	lesson.Title = respLesson.Title
	lesson.StartTime = respLesson.StartTime
	lesson.FinishTime = respLesson.FinishTime
	lesson.Type = respLesson.Type
	lesson.RoomId = respLesson.RoomId
	lesson.TeacherFullName = respLesson.TeacherFullName

	return &lesson
}

func (s *ScheduleService) storageLessonsToDomainSchedule(storageLessons []*schedulepg.Lesson,
	from time.Time,
	to time.Time) (*Schedule, error) {

	daySchedule := s.storageLessonsToDomainDaySchedule(storageLessons)

	strDate, err := xtimeconv.FromTimeRangeToDateString(from, to)
	if err != nil {
		return nil, err
	}
	log.Printf("strDate in service: %v\n", strDate)

	schedule := make(Schedule)
	schedule[Date(strDate)] = daySchedule
	return &schedule, nil
}

func (s *ScheduleService) storageLessonsToDomainDaySchedule(storageLessons []*schedulepg.Lesson) *DaySchedule {
	var daySchedule DaySchedule

	for i := 0; i < len(storageLessons); i++ {
		lesson := s.storageLessonToDomainLesson(storageLessons[i])
		daySchedule = append(daySchedule, lesson)
	}
	return &daySchedule
}

func (s *ScheduleService) storageLessonToDomainLesson(storageLesson *schedulepg.Lesson) *Lesson {
	var lesson Lesson

	lesson.Title = storageLesson.Title
	lesson.StartTime = storageLesson.StartTime
	lesson.FinishTime = storageLesson.FinishTime
	lesson.Type = storageLesson.Type
	lesson.RoomId = storageLesson.RoomId
	lesson.TeacherFullName = storageLesson.TeacherFullName

	return &lesson
}
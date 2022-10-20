package storage

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/vaberof/TelegramBotUniversitySchedule/pkg/xtime"
	"github.com/vaberof/TelegramBotUniversitySchedule/pkg/xtimeconv"
	"github.com/vaberof/TelegramBotUniversitySchedule/pkg/xtimezone"
	"time"
)

type GroupId string
type Date string

type ScheduleStorage struct {
	scheduleStorage map[GroupId]map[Date]*Schedule
}

type Schedule struct {
	Lessons    []*Lesson
	ExpireTime time.Time
}

type Lesson struct {
	Title           string
	StartTime       string
	FinishTime      string
	Type            string
	RoomId          string
	TeacherFullName string
}

func NewScheduleStorage() *ScheduleStorage {
	return &ScheduleStorage{
		scheduleStorage: map[GroupId]map[Date]*Schedule{},
	}
}

func (s *ScheduleStorage) GetLessons(groupId string, from time.Time, to time.Time) ([]*Lesson, error) {
	dateString, err := xtimeconv.FromTimeToDateString(from, to)
	if err != nil {
		return nil, err
	}

	schedule := s.scheduleStorage[GroupId(groupId)][Date(dateString)]
	if schedule == nil {
		return nil, errors.New("schedule not cached yet")
	}

	if err = s.isScheduleOutdated(schedule.ExpireTime); err != nil {
		return nil, err
	}

	log.Printf("schedule sent from cache")
	log.Printf("expire time: %s", schedule.ExpireTime.Format("02.01.2006"))

	return schedule.Lessons, nil
}

func (s *ScheduleStorage) SaveLessons(groupId string, from time.Time, to time.Time, lessons []*Lesson) error {
	dateString, err := xtimeconv.FromTimeToDateString(from, to)
	if err != nil {
		return err
	}

	var schedule Schedule

	schedule.Lessons = lessons

	err = s.setExpireTime(&schedule, from, to)
	if err != nil {
		return err
	}

	if s.scheduleStorage[GroupId(groupId)] == nil {
		s.scheduleStorage[GroupId(groupId)] = make(map[Date]*Schedule)
	}

	s.scheduleStorage[GroupId(groupId)][Date(dateString)] = &schedule
	return nil
}

func (s *ScheduleStorage) setExpireTime(schedule *Schedule, from time.Time, to time.Time) error {
	dateString, err := xtimeconv.FromTimeToDateString(from, to)
	if err != nil {
		return err
	}

	s.setExpireTimeImpl(schedule, dateString, to)
	log.Printf("settled expire time: %s", schedule.ExpireTime.Format("02.01.2006"))
	return nil
}

func (s *ScheduleStorage) setExpireTimeImpl(schedule *Schedule, dateString string, date time.Time) {
	switch dateString {
	case "Today":
		tomorrowExpire := date.Add(24 * time.Hour)
		schedule.ExpireTime = tomorrowExpire
	case "Tomorrow":
		tomorrowExpire := date
		schedule.ExpireTime = tomorrowExpire
	case "Week":
		// date is equals to sunday of the current week
		nextMondayExpire := date.Add(24 * time.Hour)
		schedule.ExpireTime = nextMondayExpire
	default:
		// date is equals to sunday of the next week
		nextMondayExpire := date.Add(-6 * 24 * time.Hour)
		schedule.ExpireTime = nextMondayExpire
	}
}

func (s *ScheduleStorage) isScheduleOutdated(scheduleExpireTime time.Time) error {
	novosibirsk, err := xtime.GetDefaultLocation(xtimezone.Novosibirsk)
	if err != nil {
		return err
	}

	currentDate := time.Now().In(novosibirsk).Format("02.01")
	if currentDate == scheduleExpireTime.Format("02.01") {
		log.Printf("schedule is outdated: %s", scheduleExpireTime.Format("02.01"))
		return errors.New("scheduleExpireTime is outdated")
	}
	return nil
}

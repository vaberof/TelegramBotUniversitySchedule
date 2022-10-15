package storage

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/vaberof/TelegramBotUniversitySchedule/pkg/xstrconv"
	"github.com/vaberof/TelegramBotUniversitySchedule/pkg/xtime"
	"github.com/vaberof/TelegramBotUniversitySchedule/pkg/xtimezone"
	"time"
)

type GroupId string
type Date string

type ScheduleStorage struct {
	Schedule map[GroupId]map[Date]*Schedule
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
		Schedule: map[GroupId]map[Date]*Schedule{},
	}
}

func (s *ScheduleStorage) GetCachedLessons(groupId string, from time.Time, to time.Time) ([]*Lesson, error) {
	date, err := xstrconv.ConvertDateToStr(from, to)
	if err != nil {
		return nil, err
	}

	schedule := s.Schedule[GroupId(groupId)][Date(date)]
	if schedule == nil {
		return nil, errors.New("schedule not cached yet")
	}

	if err = s.isScheduleOutdated(schedule); err != nil {
		log.Printf("schedule outdated: %s\n", schedule.ExpireTime.Format("02.01.2006"))
		return nil, err
	}

	log.Printf("schedule settled from cache")
	return schedule.Lessons, nil
}

func (s *ScheduleStorage) SaveLessons(groupId string, from time.Time, to time.Time, lessons []*Lesson) error {
	date, err := xstrconv.ConvertDateToStr(from, to)
	if err != nil {
		return err
	}

	var schedule Schedule

	schedule.Lessons = lessons

	err = s.setExpireTime(&schedule, from, to)
	log.Printf("expire time: %s\n", schedule.ExpireTime.Format("02.01.2006"))
	if err != nil {
		return err
	}

	if s.Schedule[GroupId(groupId)] == nil {
		s.Schedule[GroupId(groupId)] = make(map[Date]*Schedule)
	}

	s.Schedule[GroupId(groupId)][Date(date)] = &schedule
	return nil
}

func (s *ScheduleStorage) setExpireTime(schedule *Schedule, from time.Time, to time.Time) error {
	date, err := xstrconv.ConvertDateToStr(from, to)
	if err != nil {
		return err
	}

	Novosibirsk := xtime.GetDefaultLocation(xtimezone.Novosibirsk)
	switch date {
	case "Today", "Tomorrow":
		schedule.ExpireTime = time.Date(time.Now().Year(),
			time.Now().Month(),
			time.Now().Day(),
			17, 00, 0, 0,
			time.UTC).In(Novosibirsk)
		return nil
	default:
		// in case when schedule is needed to current week or next week
		// 'to' is equals to sunday
		sunday := to
		schedule.ExpireTime = time.Date(sunday.Year(),
			time.Now().Month(),
			time.Now().Day(),
			17, 00, 0, 0,
			time.UTC).In(Novosibirsk)
		return nil
	}
}

func (s *ScheduleStorage) isScheduleOutdated(schedule *Schedule) error {
	Novosibirsk := xtime.GetDefaultLocation(xtimezone.Novosibirsk)

	currentTime := time.Now().In(Novosibirsk)
	if currentTime.Format("02.01") == schedule.ExpireTime.Format("02.01") {
		log.Printf("schedule expired: %s\n", schedule.ExpireTime.Format("02.01"))
		return errors.New("schedule is outdated")
	}

	return nil
}

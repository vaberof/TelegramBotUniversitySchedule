package schedulepg

import (
	"errors"
	"github.com/sirupsen/logrus"
	domain "github.com/vaberof/TelegramBotUniversitySchedule/internal/domain/schedule"
	"github.com/vaberof/TelegramBotUniversitySchedule/pkg/xtime"
	"github.com/vaberof/TelegramBotUniversitySchedule/pkg/xtimeconv"
	"github.com/vaberof/TelegramBotUniversitySchedule/pkg/xtimezone"
	"gorm.io/gorm"
	"time"
)

type ScheduleStoragePostgres struct {
	db *gorm.DB
}

func NewScheduleStoragePostgres(db *gorm.DB) *ScheduleStoragePostgres {
	return &ScheduleStoragePostgres{db: db}
}

func (s *ScheduleStoragePostgres) GetSchedule(groupId string, from time.Time, to time.Time) (domain.Schedule, error) {
	dateString, err := xtimeconv.FromTimeRangeToDateString(from, to)
	if err != nil {
		return nil, err
	}

	postgresSchedule, err := s.getSchedule(groupId, dateString)
	if err != nil {
		logrus.Error("cannot get schedule from db, error: ", err)
		return nil, err
	}

	if err = s.isScheduleOutdated(postgresSchedule.ExpireTime); err != nil {
		return nil, err
	}

	domainSchedule, err := BuildDomainSchedule(postgresSchedule, postgresSchedule.Date)
	if err != nil {
		logrus.Error("cannot build domain schedule, error: ", err)
		return nil, err
	}

	logrus.Info("schedule sent from database")
	return domainSchedule, nil
}

func (s *ScheduleStoragePostgres) SaveSchedule(groupId string, schedule domain.Schedule) error {
	var dateString string
	var daySchedule domain.DaySchedule

	for date, value := range schedule {
		dateString = string(date)
		daySchedule = *value
	}

	postgresSchedule, err := s.getSchedule(groupId, dateString)
	if err == nil {
		err = s.deleteSchedule(groupId, postgresSchedule)
		if err != nil {
			return err
		}
	}

	postgresLessons := BuildPostgresLessons(daySchedule)

	err = s.saveScheduleImpl(groupId, dateString, postgresLessons)
	if err != nil {
		return err
	}
	return nil
}

func (s *ScheduleStoragePostgres) DeleteSchedule(groupId string, date string) error {
	schedule, err := s.getSchedule(groupId, date)
	if err != nil {
		return err
	}
	return s.deleteSchedule(groupId, schedule)
}

func (s *ScheduleStoragePostgres) deleteSchedule(groupId string, schedule *Schedule) error {
	err := s.db.Select("Lessons").Where("group_id = ?", groupId).Delete(&schedule).Error
	if err != nil {
		logrus.Error("cannot delete schedule from db, error: ", err)
		return err
	}
	logrus.Info("deleted schedule from db")
	return nil
}

func (s *ScheduleStoragePostgres) getSchedule(groupId string, dateString string) (*Schedule, error) {
	var postgresSchedule Schedule

	err := s.db.Table("schedules").
		Preload("Lessons").
		Where("group_id = ? AND date = ?", groupId, dateString).
		First(&postgresSchedule).Error
	if err != nil {
		logrus.Error("schedule not found in db, error: ", err)
		return nil, errors.New("schedule not found")
	}
	return &postgresSchedule, nil
}

func (s *ScheduleStoragePostgres) saveScheduleImpl(groupId string, dateString string, lessons []*Lesson) error {
	var schedule Schedule

	schedule.GroupId = groupId
	schedule.Date = dateString
	schedule.Lessons = lessons

	err := s.setExpireTime(&schedule, dateString)
	if err != nil {
		return err
	}

	err = s.db.Create(&schedule).Error
	if err != nil {
		return err
	}

	logrus.Info("schedule cached")
	return nil
}

func (s *ScheduleStoragePostgres) setExpireTime(schedule *Schedule, date string) error {
	_, to, err := xtime.ParseDatesRange(date)

	s.setExpireTimeImpl(schedule, date, to)
	if err != nil {
		return err
	}
	logrus.Info("settled expire time: ", schedule.ExpireTime)
	return nil
}

func (s *ScheduleStoragePostgres) setExpireTimeImpl(schedule *Schedule, dateString string, requestedDate time.Time) {
	year, month, day := requestedDate.Date()
	switch dateString {
	case "Today":
		tomorrowExpireDate := time.Date(year, month, day+1, 0, 0, 0, 0, requestedDate.Location())
		schedule.ExpireTime = tomorrowExpireDate
	case "Tomorrow":
		tomorrowExpireDate := time.Date(year, month, day, 0, 0, 0, 0, requestedDate.Location())
		schedule.ExpireTime = tomorrowExpireDate
	case "Week":
		// requestedDate is equals to sunday of the current week
		nextMondayExpireDate := time.Date(year, month, day+1, 0, 0, 0, 0, requestedDate.Location())
		schedule.ExpireTime = nextMondayExpireDate
	default:
		// requestedDate is equals to sunday of the next week
		nextMondayExpireDate := time.Date(year, month, day-6, 0, 0, 0, 0, requestedDate.Location())
		schedule.ExpireTime = nextMondayExpireDate
	}
}

func (s *ScheduleStoragePostgres) isScheduleOutdated(scheduleExpireTime time.Time) error {
	novosibirsk, err := xtime.GetDefaultLocation(xtimezone.Novosibirsk)
	if err != nil {
		return err
	}

	currentTime := time.Now().In(novosibirsk)
	if currentTime.After(scheduleExpireTime) {
		logrus.Info("schedule is outdated: ", scheduleExpireTime)
		return errors.New("schedule is outdated")
	}
	return nil
}

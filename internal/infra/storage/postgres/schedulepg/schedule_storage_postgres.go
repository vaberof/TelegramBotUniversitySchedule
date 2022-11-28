package schedulepg

import (
	"errors"
	log "github.com/sirupsen/logrus"
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

func (s *ScheduleStoragePostgres) GetLessons(groupId string, from time.Time, to time.Time) ([]*Lesson, error) {
	dateString, err := xtimeconv.FromTimeRangeToDateString(from, to)
	if err != nil {
		return nil, err
	}

	scheduleFromDb, err := s.getSchedule(groupId, dateString)
	if err != nil {
		log.Error("cannot get schedule from db, error: ", err)
		return nil, err
	}

	if err = s.isScheduleOutdated(scheduleFromDb.ExpireTime); err != nil {
		return nil, err
	}

	lessonsFromDb, err := s.getLessons(scheduleFromDb.Id)
	if err != nil {
		log.Error("cannot get lessons from db, error: ", err)
		return nil, err
	}
	log.Info("lessons sent from database")
	return lessonsFromDb, nil
}

func (s *ScheduleStoragePostgres) SaveLessons(groupId string, from time.Time, to time.Time, lessons []*Lesson) error {
	dateString, err := xtimeconv.FromTimeRangeToDateString(from, to)
	if err != nil {
		return err
	}

	scheduleFromDb, err := s.getSchedule(groupId, dateString)
	if err == nil {
		err = s.updateLessons(scheduleFromDb, lessons)
		if err != nil {
			return err
		}

		err = s.deleteLessonsWithNullScheduleId()
		if err != nil {
			return err
		}
		log.Info("deleted lessons with null schedule id from db")
		return nil
	}

	err = s.saveLessonsImpl(groupId, dateString, from, to, lessons)
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

	err = s.db.Select("Lessons").Where("group_id = ?", groupId).Delete(&schedule).Error
	if err != nil {
		log.Error("cannot delete schedule from db, error: ", err)
		return err
	}
	log.Info("deleted schedule from db")
	return nil
}

func (s *ScheduleStoragePostgres) getSchedule(groupId string, dateString string) (*Schedule, error) {
	var schedule Schedule

	err := s.db.Table("schedules").Where("group_id = ? AND date = ?", groupId, dateString).First(&schedule).Error
	if err != nil {
		log.Error("schedule not found in db, error: ", err)
		return nil, errors.New("schedule not found")
	}
	return &schedule, nil
}

func (s *ScheduleStoragePostgres) getLessons(scheduleId uint) ([]*Lesson, error) {
	var lessons []*Lesson

	err := s.db.Table("lessons").Where("schedule_id = ?", scheduleId).Find(&lessons).Error
	if err != nil {
		log.Error("lessons not found in db, error: ", err)
		return nil, errors.New("cant find lessons")
	}
	return lessons, nil
}

func (s *ScheduleStoragePostgres) saveLessonsImpl(groupId string, dateString string, from time.Time, to time.Time, lessons []*Lesson) error {
	var schedule Schedule

	schedule.GroupId = groupId
	schedule.Date = dateString
	schedule.Lessons = lessons

	err := s.setExpireTime(&schedule, from, to)
	if err != nil {
		return err
	}

	err = s.db.Create(&schedule).Error
	if err != nil {
		return err
	}

	log.Info("schedule cached")
	return nil
}

func (s *ScheduleStoragePostgres) updateLessons(schedule *Schedule, lessons []*Lesson) error {
	err := s.db.Model(&schedule).Association("Lessons").Replace(lessons)
	if err != nil {
		log.Error("cannot update lessons in db, error: ", err)
		return err
	}
	log.Info("lessons updated")
	return nil
}

func (s *ScheduleStoragePostgres) deleteLessonsWithNullScheduleId() error {
	err := s.db.Table("lessons").Where("schedule_id IS NULL").Delete(&Lesson{}).Error
	if err != nil {
		log.Error("cannot delete lessons in db with null schedule id, error: ", err)
		return err
	}
	return nil
}

func (s *ScheduleStoragePostgres) setExpireTime(schedule *Schedule, from time.Time, to time.Time) error {
	dateString, err := xtimeconv.FromTimeRangeToDateString(from, to)
	if err != nil {
		return err
	}

	s.setExpireTimeImpl(schedule, dateString, to)
	log.Info("settled expire time: ", schedule.ExpireTime.Format("02.01.2006"))
	return nil
}

func (s *ScheduleStoragePostgres) setExpireTimeImpl(schedule *Schedule, dateString string, date time.Time) {
	log.Printf("DATE TIME: %v", date)
	switch dateString {
	case "Today":
		tomorrowExpireDate := date.Add(24 * time.Hour)
		schedule.ExpireTime = tomorrowExpireDate
	case "Tomorrow":
		tomorrowExpireDate := date
		schedule.ExpireTime = tomorrowExpireDate
	case "Week":
		// date is equals to sunday of the current week
		nextMondayExpireDate := date.Add(24 * time.Hour)
		schedule.ExpireTime = nextMondayExpireDate
	default:
		// date is equals to sunday of the next week
		nextMondayExpireDate := date.Add(-6 * 24 * time.Hour)
		schedule.ExpireTime = nextMondayExpireDate
	}
}

func (s *ScheduleStoragePostgres) isScheduleOutdated(scheduleExpireTime time.Time) error {
	novosibirsk, err := xtime.GetDefaultLocation(xtimezone.Novosibirsk)
	if err != nil {
		return err
	}

	currentDate := time.Now().In(novosibirsk).Format("02.01")
	if currentDate >= scheduleExpireTime.Format("02.01") {
		log.Info("schedule is outdated: ", scheduleExpireTime.Format("02.01"))
		return errors.New("schedule is outdated")
	}
	return nil
}

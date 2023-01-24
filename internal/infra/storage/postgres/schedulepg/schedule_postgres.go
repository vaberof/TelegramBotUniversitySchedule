package schedulepg

import (
	"time"
)

type Schedule struct {
	Id         uint `gorm:"primaryKey"`
	GroupId    string
	Date       string
	Lessons    []*Lesson
	ExpireTime time.Time
}

type Lesson struct {
	Id              uint `gorm:"primaryKey"`
	ScheduleId      uint
	Title           string
	StartTime       string
	FinishTime      string
	Type            string
	RoomId          string
	TeacherFullName string
}

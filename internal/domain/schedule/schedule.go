package domain

type Date string

type Schedule = map[Date]*DaySchedule
type DaySchedule []*Lesson

type Lesson struct {
	Title           string
	StartTime       string
	FinishTime      string
	Type            string
	RoomId          string
	TeacherFullName string
}

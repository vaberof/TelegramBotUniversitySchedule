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

// GetLessonNumber checks time when lesson starts
// and returns corresponding lesson`s number.
func GetLessonNumber(startTime string) int {
	switch startTime {
	case "10:10":
		return 2
	case "12:10":
		return 3
	case "13:50":
		return 4
	case "15:30":
		return 5
	case "17:10":
		return 6
	case "18:50":
		return 7
	default:
		return 1
	}
}

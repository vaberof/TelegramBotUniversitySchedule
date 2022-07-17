package model

type Schedule map[string]DaySchedule
type DaySchedule []Lesson

type Lesson struct {
	Name            string
	StartTime       string
	FinishTime      string
	Type            string
	Room            string
	TeacherFullName string
}

func NewDaySchedule() *DaySchedule {
	return &DaySchedule{}
}

func (d *DaySchedule) AddLessons(lessonName, startTime, finishTime, lessonType, roomNumber, teacherName string) {
	*d = append(*d, Lesson{
		lessonName,
		startTime,
		finishTime,
		lessonType,
		roomNumber,
		teacherName,
	})
}

// GetLessonNumber checks time when lesson starts
// and returns corresponding lesson`s number.
func GetLessonNumber(startTime string) int {
	var lessonNumber int

	switch startTime {
	case "10:10":
		lessonNumber = 2
	case "12:10":
		lessonNumber = 3
	case "13:50":
		lessonNumber = 4
	case "15:30":
		lessonNumber = 5
	case "17:10":
		lessonNumber = 6
	case "18:50":
		lessonNumber = 7
	default:
		lessonNumber = 1
	}
	return lessonNumber
}

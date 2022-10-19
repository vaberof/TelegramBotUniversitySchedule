package infra

type GetScheduleResponse struct {
	Lessons []*Lesson
}

type Lesson struct {
	Title           string
	StartTime       string
	FinishTime      string
	Type            string
	RoomId          string
	TeacherFullName string
}

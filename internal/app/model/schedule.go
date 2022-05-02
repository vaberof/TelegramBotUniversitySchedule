package model

type Schedule struct {
	Schedule []string
}

func NewSchedule() *Schedule {
	return &Schedule{}
}

// AddLessons adds lessons to schedule.
func (s *Schedule) AddLessons(startTime, finishTime, lessonName, roomNumber, teacherName, lessonType string) {
	s.Schedule = append(s.Schedule,
		"#"+*getLessonNumber(startTime)+". "+
			lessonName+" ("+
			startTime+"-"+finishTime+
			", "+lessonType+", "+
			"ауд. "+roomNumber+", "+
			"препод. "+teacherName+")"+"\n\n")
}

// AddGroupId adds study group id to schedule.
func (s *Schedule) AddGroupId(groupId string) {
	s.Schedule = append([]string{groupId + "\n"}, s.Schedule...)
}

// AddDate adds needed date to schedule while parsing on day.
func (s *Schedule) AddDate(date *Date) {
	s.Schedule = append([]string{date.fullDate + " (" + date.day + ")" + "\n\n"}, s.Schedule...)
}

// AddWeekDate adds needed date to schedule while parsing on week.
func (s *Schedule) AddWeekDate(date *Date, index int) {
	s.Schedule = append(s.Schedule, date.weekFullLDates[index]+" ("+date.weekDays[index]+")"+"\n\n")
}

// NoLessons adds "Пар нет" to schedule
// if no lessons on a certain day in non-nil selection while parsing.
func (s *Schedule) NoLessons() {
	s.Schedule = append(s.Schedule, "Пар нет"+"\n\n")
}

// NotFound adds "Расписание не найдено"
// if we catch in nil selection while parsing.
func (s *Schedule) NotFound() {
	s.Schedule = append(s.Schedule, "Расписание не найдено")
}

// ScheduleExists returns true
// if we have schedule in a certain day.
func (s *Schedule) ScheduleExists() bool {
	return len(s.Schedule) != 0
}

// getLessonNumber checks time when lesson starts
// and returns corresponding lesson`s number.
func getLessonNumber(startTime string) *string {
	var lessonNumber string

	switch startTime {
	case "10:10":
		lessonNumber = "2"
	case "12:10":
		lessonNumber = "3"
	case "13:50":
		lessonNumber = "4"
	case "15:30":
		lessonNumber = "5"
	case "17:10":
		lessonNumber = "6"
	case "18:50":
		lessonNumber = "7"
	default:
		lessonNumber = "1"
	}

	return &lessonNumber
}

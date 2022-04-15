package services

type Schedule struct {
	schedule []string
}

func CreateSchedule() *Schedule {
	return &Schedule{}
}

// Добавляем пары в массив
func (s *Schedule) AddLessons(startTime, finishTime, lessonName, roomNumber, teacherName, lessonType string) {
	s.schedule = append(s.schedule,
		"\n"+ "Пара" + "                № "+ getLessonNumber(startTime) + "\n",
		"Предмет"+ "         " + lessonName + " (" + lessonType + ")" + "\n",
		"Время"+ "             " + startTime + "-" + finishTime + "\n",
		"Аудитория"+ "     " + roomNumber + "\n",
		"Препод."+ "          " + teacherName + "\n")
}

// Возвращаем расписание пользователю
func (s *Schedule) GetSchedule() string {
	var schedule string

	if !s.todayScheduleExists(s.schedule) {
		schedule = GetTodayDate()[1] + "\n" + "Пар нет" // реализовать иначе
		return schedule
	}

	s.AddTodayDate() // реализовать иначе
	for _, text := range s.schedule {
		schedule += text
	}

	return schedule
}

// Добавляем дату в начало массива
func (s *Schedule) AddTodayDate() {
	s.schedule = append([]string{"Дата" + "                 " + GetTodayDate()[1] + "\n"}, s.schedule...)
}

// Проверяем, есть ли сегодня пары
func (s *Schedule) todayScheduleExists(schedule []string) bool {
	if len(schedule) == 0 {
		return false
	}
	return true
}

// Проверяем во сколько начинается пара,
// чтобы правильно указывать номер пары
func getLessonNumber(startTime string) string {
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

	return lessonNumber
}

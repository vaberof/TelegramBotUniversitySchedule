package controllers

type Schedule struct {
	schedule []string
}

func scheduleInit() *Schedule {
	return &Schedule{}
}

// Добавляем сегодняшнуюю дату в начало массива
func (s *Schedule) appendDateToArray() {
	s.schedule = append([]string{"Дата" + "                 " + getTodayDate()[1] + "\n"},
		s.schedule...)
}

// Добавляем пары на текущий день в массив
func (s *Schedule) appendLessonsToArray(startTime, finishTime, lessonName, roomNumber, teacherName, lessonType string) {
	s.schedule = append(s.schedule,
		"\n"+ "Пара" + "                № "+ lessonNum(startTime) + "\n",
		"Предмет"+ "         " + lessonName + " (" + lessonType + ")" + "\n",
		"Время"+ "             " + startTime + "-" + finishTime + "\n",
		"Аудитория"+ "     " + roomNumber + "\n",
		"Препод."+ "          " + teacherName + "\n")
}

// Возвращаем расписание пользователю
func (s *Schedule) getScheduleFromArray() string {
	var message string

	for _, text := range s.schedule {
		message += text
	}

	if !s.isScheduleExists(message) { // пофиксить
		message += "Пар нет"
	}

	return message
}

// Проверяем, есть ли сегодня пары
func (s *Schedule) isScheduleExists(message string) bool {

	if len(message) == 0 {            // пофиксить
		s.appendDateToArray()
		return false
	}

	return true
}

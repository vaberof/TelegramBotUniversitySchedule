package schedulepg

import domain "github.com/vaberof/TelegramBotUniversitySchedule/internal/domain/schedule"

func BuildPostgresLessons(daySchedule domain.DaySchedule) []*Lesson {
	postgresLessons := make([]*Lesson, len(daySchedule))

	for i := 0; i < len(postgresLessons); i++ {
		postgresLessons[i] = buildPostgresLesson(daySchedule[i])
	}

	return postgresLessons
}

func buildPostgresLesson(domainLesson *domain.Lesson) *Lesson {
	return &Lesson{
		Title:           domainLesson.Title,
		StartTime:       domainLesson.StartTime,
		FinishTime:      domainLesson.FinishTime,
		Type:            domainLesson.Type,
		RoomId:          domainLesson.RoomId,
		TeacherFullName: domainLesson.TeacherFullName,
	}
}

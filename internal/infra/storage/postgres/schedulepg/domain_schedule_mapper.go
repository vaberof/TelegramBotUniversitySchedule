package schedulepg

import (
	domain "github.com/vaberof/TelegramBotUniversitySchedule/internal/domain/schedule"
)

func BuildDomainSchedule(postgresLessons []*Lesson, date string) (domain.Schedule, error) {
	daySchedule := buildDomainDaySchedule(postgresLessons)

	schedule := make(domain.Schedule)
	schedule[domain.Date(date)] = daySchedule

	return schedule, nil
}

func buildDomainDaySchedule(postgresLessons []*Lesson) *domain.DaySchedule {
	daySchedule := make(domain.DaySchedule, len(postgresLessons))

	for i := 0; i < len(daySchedule); i++ {
		daySchedule[i] = buildDomainLesson(postgresLessons[i])
	}

	return &daySchedule
}

func buildDomainLesson(postgresLesson *Lesson) *domain.Lesson {
	return &domain.Lesson{
		Title:           postgresLesson.Title,
		StartTime:       postgresLesson.StartTime,
		FinishTime:      postgresLesson.FinishTime,
		Type:            postgresLesson.Type,
		RoomId:          postgresLesson.RoomId,
		TeacherFullName: postgresLesson.TeacherFullName,
	}
}

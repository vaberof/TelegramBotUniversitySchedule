package unisite

import (
	domain "github.com/vaberof/TelegramBotUniversitySchedule/internal/domain/schedule"
	"github.com/vaberof/TelegramBotUniversitySchedule/pkg/xtimeconv"
	"time"
)

func BuildDomainSchedule(scheduleResponse *GetScheduleResponse, from time.Time, to time.Time) (domain.Schedule, error) {
	daySchedule := buildDomainDaySchedule(scheduleResponse.Lessons)

	dateString, err := xtimeconv.FromTimeRangeToDateString(from, to)
	if err != nil {
		return nil, err
	}

	schedule := make(domain.Schedule)
	schedule[domain.Date(dateString)] = daySchedule

	return schedule, nil
}

func buildDomainDaySchedule(responseLessons []*Lesson) *domain.DaySchedule {
	daySchedule := make(domain.DaySchedule, len(responseLessons))

	for i := 0; i < len(daySchedule); i++ {
		daySchedule[i] = buildDomainLesson(responseLessons[i])
	}

	return &daySchedule
}

func buildDomainLesson(responseLesson *Lesson) *domain.Lesson {
	return &domain.Lesson{
		Title:           responseLesson.Title,
		StartTime:       responseLesson.StartTime,
		FinishTime:      responseLesson.FinishTime,
		Type:            responseLesson.Type,
		RoomId:          responseLesson.RoomId,
		TeacherFullName: responseLesson.TeacherFullName,
	}
}

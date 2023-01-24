package unisite

import integration "github.com/vaberof/TelegramBotUniversitySchedule/pkg/integration/unisite"

func BuildGetScheduleResponse(ApiScheduleResponse *integration.GetScheduleResponse) *GetScheduleResponse {
	lessons := buildGetScheduleResponseLessons(ApiScheduleResponse.Lessons)
	return &GetScheduleResponse{
		Lessons: lessons,
	}
}

func buildGetScheduleResponseLessons(ApiScheduleResponseLessons []*integration.Lesson) []*Lesson {
	lessons := make([]*Lesson, len(ApiScheduleResponseLessons))

	for i := 0; i < len(lessons); i++ {
		lessons[i] = buildGetScheduleResponseLesson(ApiScheduleResponseLessons[i])
	}

	return lessons
}

func buildGetScheduleResponseLesson(ApiScheduleResponseLesson *integration.Lesson) *Lesson {
	return &Lesson{
		Title:           ApiScheduleResponseLesson.Title,
		StartTime:       ApiScheduleResponseLesson.StartTime,
		FinishTime:      ApiScheduleResponseLesson.FinishTime,
		Type:            ApiScheduleResponseLesson.Type,
		RoomId:          ApiScheduleResponseLesson.RoomId,
		TeacherFullName: ApiScheduleResponseLesson.TeacherFullName,
	}
}

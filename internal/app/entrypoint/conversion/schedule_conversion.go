package conversion

import (
	"fmt"
	domain "github.com/vaberof/TelegramBotUniversitySchedule/internal/domain/schedule"
	"github.com/vaberof/TelegramBotUniversitySchedule/pkg/i18n"
	"github.com/vaberof/TelegramBotUniversitySchedule/pkg/xtime"
	"strconv"
	"time"
)

// ScheduleToString converts schedule of type model.Schedule to string type
// to output it to user.
func ScheduleToString(studyGroupId string, inputTelegramButtonDate string, schedule *domain.Schedule) *string {
	switch inputTelegramButtonDate {
	case xtime.Today, xtime.Tomorrow:
		dateToParse := xtime.GetDateToParse(inputTelegramButtonDate)
		return DayScheduleToString(studyGroupId, inputTelegramButtonDate, dateToParse, schedule)
	default:
		datesToParse := xtime.GetDatesToParse(inputTelegramButtonDate)
		return WeekScheduleToString(studyGroupId, inputTelegramButtonDate, datesToParse, schedule)
	}
}

// DayScheduleToString converts schedule of type model.Schedule to string
// if user chosen schedule on today/tomorrow.
func DayScheduleToString(
	studyGroupId string,
	inputTelegramButtonDate string,
	date time.Time,
	schedule *domain.Schedule) *string {

	var scheduleString string

	scheduleDate := domain.Date(inputTelegramButtonDate)

	scheduleString = addGroupId(studyGroupId)
	scheduleString += addDate(date)

	dereferenceSchedule := *schedule

	daySchedule := dereferenceSchedule[scheduleDate]

	for i := 0; i < len(daySchedule); i++ {
		lesson := daySchedule[i]
		if !isHaveLessons(lesson.Title) {
			scheduleString += addNoLessons()
			break
		}

		if !isFoundLessons(lesson.Title) {
			scheduleString += addNotFoundLessons()
			break
		}

		scheduleString += addLessonNumber(lesson.StartTime)
		scheduleString += addLessonName(lesson.Title)
		scheduleString += addLessonStartTime(lesson.StartTime)
		scheduleString += addLessonFinishTime(lesson.FinishTime)
		scheduleString += addLessonType(lesson.Type)
		scheduleString += addLessonRoom(lesson.RoomId)
		scheduleString += addLessonTeacherFullName(lesson.TeacherFullName)
	}

	return &scheduleString
}

// WeekScheduleToString converts schedule of type model.Schedule to string
// if user chosen schedule on week/next week.
func WeekScheduleToString(
	studyGroupId string,
	inputTelegramButtonDate string,
	dates []time.Time,
	schedule *domain.Schedule) *string {

	var scheduleString string

	scheduleDate := domain.Date(inputTelegramButtonDate)

	scheduleString = addGroupId(studyGroupId)

	dereferenceSchedule := *schedule
	daySchedule := dereferenceSchedule[scheduleDate]

	for i := 0; i < len(daySchedule); i++ {
		lesson := daySchedule[i]

		scheduleString += addDate(dates[i])

		if !isHaveLessons(lesson.Title) {
			scheduleString += addNoLessons()
			continue
		}

		if !isFoundLessons(lesson.Title) {
			scheduleString += addNotFoundLessons()
			continue
		}

		scheduleString += addLessonNumber(lesson.StartTime)
		scheduleString += addLessonName(lesson.Title)
		scheduleString += addLessonStartTime(lesson.StartTime)
		scheduleString += addLessonFinishTime(lesson.FinishTime)
		scheduleString += addLessonType(lesson.Type)
		scheduleString += addLessonRoom(lesson.RoomId)
		scheduleString += addLessonTeacherFullName(lesson.TeacherFullName)
	}

	return &scheduleString
}

func addGroupId(studyGroupId string) string {
	scheduleString := "Группа: " + fmt.Sprintf("%s\n\n", studyGroupId)
	return scheduleString
}

func addLessonNumber(lessonNumber string) string {
	scheduleString := "#" + strconv.Itoa(domain.GetLessonNumber(lessonNumber)) + ". "
	return scheduleString
}

func addLessonName(lessonName string) string {
	scheduleString := lessonName + "("
	return scheduleString
}

func addLessonStartTime(startTime string) string {
	scheduleString := startTime + "-"
	return scheduleString
}

func addLessonFinishTime(finishTime string) string {
	scheduleString := finishTime + ", "
	return scheduleString
}

func addLessonType(lessonType string) string {
	scheduleString := lessonType + ", "
	return scheduleString
}

func addLessonRoom(lessonRoom string) string {
	scheduleString := "ауд. " + lessonRoom + ", "
	return scheduleString
}

func addLessonTeacherFullName(TeacherFullName string) string {
	scheduleString := "препод. " + TeacherFullName + ")" + "\n\n"
	return scheduleString
}

func addNoLessons() string {
	scheduleString := "Пар нет\n\n"
	return scheduleString
}

func addNotFoundLessons() string {
	scheduleString := "Расписание не найдено\n\n"
	return scheduleString
}

func addDate(date time.Time) string {
	scheduleString := fmt.Sprintf("*%s*", "Дата: ") +
		fmt.Sprintf("*%s*", date.Format("02.01.2006")) +
		fmt.Sprintf(" *(%s)*\n\n", i18n.FormatRuWeekday(date.Weekday()))

	return scheduleString
}

func isHaveLessons(lessonNameField string) bool {
	return lessonNameField != "no lessons"
}

func isFoundLessons(lessonNameField string) bool {
	return lessonNameField != "not found"
}

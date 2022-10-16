package xstrconv

import (
	"fmt"
	domain "github.com/vaberof/TelegramBotUniversitySchedule/internal/domain/schedule"
	"github.com/vaberof/TelegramBotUniversitySchedule/pkg/i18n"
	"github.com/vaberof/TelegramBotUniversitySchedule/pkg/xtime"
	"strconv"
	"strings"
	"time"
)

func ScheduleToString(groupId string, inputTelegramButtonDate string, schedule *domain.Schedule) *string {
	switch inputTelegramButtonDate {
	case xtime.Today, xtime.Tomorrow:
		date := xtime.GetDateToParse(inputTelegramButtonDate)
		return DayScheduleToString(groupId, inputTelegramButtonDate, date, schedule)
	default:
		dates := xtime.GetDatesToParse(inputTelegramButtonDate)
		return WeekScheduleToString(groupId, inputTelegramButtonDate, dates, schedule)
	}
}

func DayScheduleToString(
	groupId string,
	inputTelegramButtonDate string,
	date time.Time,
	schedule *domain.Schedule) *string {

	var strSchedule string

	scheduleDate := domain.Date(inputTelegramButtonDate)

	strSchedule = addGroupId(groupId)
	strSchedule += addDate(date)

	dereferenceSchedule := *schedule

	daySchedule := dereferenceSchedule[scheduleDate]
	dereferenceDaySchedule := *daySchedule

	for i := 0; i < len(dereferenceDaySchedule); i++ {
		lesson := dereferenceDaySchedule[i]
		if !isHaveLessons(lesson.Title) {
			strSchedule += addNoLessons()
			break
		}

		if !isFoundLessons(lesson.Title) {
			strSchedule += addNotFoundLessons()
			break
		}

		strSchedule = *addLessonDataToStrSchedule(strSchedule, lesson)
	}

	return &strSchedule
}

func WeekScheduleToString(
	groupId string,
	inputTelegramButtonDate string,
	dates []time.Time,
	schedule *domain.Schedule) *string {

	var strSchedule string

	scheduleDate := domain.Date(inputTelegramButtonDate)

	strSchedule = addGroupId(groupId)
	strSchedule += addDate(dates[0])
	day := 1

	dereferenceSchedule := *schedule

	daySchedule := dereferenceSchedule[scheduleDate]
	dereferenceDaySchedule := *daySchedule

	for i := 0; i < len(dereferenceDaySchedule); i++ {
		lesson := dereferenceDaySchedule[i]
		lessonTitle := lesson.Title

		if !isHaveLessonsWhileWeekConvert(&strSchedule, lessonTitle, &day, dates) {
			continue
		}

		if !isFoundLessonsWhileWeekConvert(&strSchedule, lessonTitle, &day, dates) {
			continue
		}

		if isNextDayWhileWeekConvert(&strSchedule, lessonTitle, &day, dates) {
			continue
		}

		strSchedule = *addLessonDataToStrSchedule(strSchedule, lesson)
	}
	return &strSchedule
}

func addLessonDataToStrSchedule(strSchedule string, lesson *domain.Lesson) *string {
	strSchedule += addLessonNumber(lesson.StartTime)
	strSchedule += addLessonName(lesson.Title)
	strSchedule += addLessonStartTime(lesson.StartTime)
	strSchedule += addLessonFinishTime(lesson.FinishTime)
	strSchedule += addLessonType(lesson.Type)
	strSchedule += addLessonRoom(lesson.RoomId)
	strSchedule += addLessonTeacherFullName(lesson.TeacherFullName)
	return &strSchedule
}

func addGroupId(studyGroupId string) string {
	strSchedule := "Группа: " + fmt.Sprintf("%s\n\n", studyGroupId)
	return strSchedule
}

func addLessonNumber(lessonNumber string) string {
	strSchedule := "#" + strconv.Itoa(domain.GetLessonNumber(lessonNumber)) + ". "
	return strSchedule
}

func addLessonName(lessonName string) string {
	strSchedule := lessonName + "("
	return strSchedule
}

func addLessonStartTime(startTime string) string {
	strSchedule := startTime + "-"
	return strSchedule
}

func addLessonFinishTime(finishTime string) string {
	strSchedule := finishTime + ", "
	return strSchedule
}

func addLessonType(lessonType string) string {
	strSchedule := lessonType + ", "
	return strSchedule
}

func addLessonRoom(lessonRoom string) string {
	strSchedule := "ауд. " + lessonRoom + ", "
	return strSchedule
}

func addLessonTeacherFullName(TeacherFullName string) string {
	strSchedule := "препод. " + TeacherFullName + ")" + "\n\n"
	return strSchedule
}

func addNoLessons() string {
	strSchedule := "Пар нет\n\n"
	return strSchedule
}

func addNotFoundLessons() string {
	strSchedule := "Расписание не найдено\n\n"
	return strSchedule
}

func addDate(date time.Time) string {
	strSchedule := fmt.Sprintf("*%s*", "Дата: ") +
		fmt.Sprintf("*%s*", date.Format("02.01.2006")) +
		fmt.Sprintf(" *(%s)*\n\n", i18n.FormatRuWeekday(date.Weekday()))

	return strSchedule
}

func isHaveLessons(lessonTitleField string) bool {
	return !strings.Contains(lessonTitleField, "no lessons")
}

func isFoundLessons(lessonTitleField string) bool {
	return !strings.Contains(lessonTitleField, "not found")
}

func isNextDay(lessonTitleField string) bool {
	return strings.Contains(lessonTitleField, "next day")
}

func isHaveLessonsWhileWeekConvert(strSchedule *string, lessonTitle string, day *int, dates []time.Time) bool {
	if !isHaveLessons(lessonTitle) {
		*strSchedule += addNoLessons()
		if isNextDay(lessonTitle) {
			*strSchedule += addDate(dates[*day])
			*day++
		}
		return false
	}
	return true
}

func isFoundLessonsWhileWeekConvert(strSchedule *string, lessonTitle string, day *int, dates []time.Time) bool {
	if !isFoundLessons(lessonTitle) {
		*strSchedule += addNotFoundLessons()
		if isNextDay(lessonTitle) {
			*strSchedule += addDate(dates[*day])
			*day++
		}
		return false
	}
	return true
}

func isNextDayWhileWeekConvert(strSchedule *string, lessonTitle string, day *int, dates []time.Time) bool {
	if isNextDay(lessonTitle) {
		*strSchedule += addDate(dates[*day])
		*day++
		return true
	}
	return false
}

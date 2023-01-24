package views

import (
	"fmt"
	domain "github.com/vaberof/TelegramBotUniversitySchedule/internal/domain/schedule"
	"github.com/vaberof/TelegramBotUniversitySchedule/pkg/i18n"
	"github.com/vaberof/TelegramBotUniversitySchedule/pkg/xtime"
	"strconv"
	"strings"
	"time"
)

func ScheduleToString(groupId string, inputTelegramButtonDate string, schedule domain.Schedule) (*string, error) {
	return getScheduleString(groupId, inputTelegramButtonDate, schedule)
}

func getScheduleString(groupId string, inputTelegramButtonDate string, schedule domain.Schedule) (*string, error) {
	switch inputTelegramButtonDate {
	case xtime.Today, xtime.Tomorrow:
		date, err := xtime.GetDate(inputTelegramButtonDate)
		if err != nil {
			return nil, err
		}
		return dayScheduleToString(groupId, inputTelegramButtonDate, date, schedule), nil
	default:
		datesRange, err := xtime.GetDatesRange(inputTelegramButtonDate)
		if err != nil {
			return nil, err
		}
		return weekScheduleToString(groupId, inputTelegramButtonDate, datesRange, schedule), nil
	}
}

func dayScheduleToString(groupId string, inputTelegramButtonDate string, date time.Time, schedule domain.Schedule) *string {
	var scheduleString string

	domainDate := domain.Date(inputTelegramButtonDate)

	scheduleString = addGroupId(groupId)
	scheduleString += addDate(date)

	dereferenceSchedule := schedule

	daySchedule := dereferenceSchedule[domainDate]
	dereferenceDaySchedule := *daySchedule

	for i := 0; i < len(dereferenceDaySchedule); i++ {
		lesson := dereferenceDaySchedule[i]
		if !hasLessons(lesson.Title) {
			scheduleString += addNoLessons()
			break
		}

		if !isFoundLessons(lesson.Title) {
			scheduleString += addNotFoundLessons()
			break
		}
		scheduleString = *addLessonToScheduleString(scheduleString, lesson)
	}
	return &scheduleString
}

func weekScheduleToString(groupId string, inputTelegramButtonDate string, datesRange []time.Time, schedule domain.Schedule) *string {
	var scheduleString string

	domainDate := domain.Date(inputTelegramButtonDate)

	scheduleString = addGroupId(groupId)
	scheduleString += addDate(datesRange[0])
	day := 1

	dereferenceSchedule := schedule

	daySchedule := dereferenceSchedule[domainDate]
	dereferenceDaySchedule := *daySchedule

	for i := 0; i < len(dereferenceDaySchedule); i++ {
		lesson := dereferenceDaySchedule[i]
		lessonTitle := lesson.Title

		if !HasLessonsWhileWeekConvert(&scheduleString, lessonTitle, &day, datesRange) {
			continue
		}

		if !isFoundLessonsWhileWeekConvert(&scheduleString, lessonTitle, &day, datesRange) {
			continue
		}

		if isNextDayWhileWeekConvert(&scheduleString, lessonTitle, &day, datesRange) {
			continue
		}
		scheduleString = *addLessonToScheduleString(scheduleString, lesson)
	}
	return &scheduleString
}

func hasLessons(lessonTitle string) bool {
	return !strings.Contains(lessonTitle, "no lessons")
}

func isFoundLessons(lessonTitle string) bool {
	return !strings.Contains(lessonTitle, "not found")
}

func isNextDay(lessonTitle string) bool {
	return strings.Contains(lessonTitle, "next day")
}

func HasLessonsWhileWeekConvert(scheduleString *string, lessonTitle string, day *int, datesRange []time.Time) bool {
	if !hasLessons(lessonTitle) {
		*scheduleString += addNoLessons()
		if isNextDay(lessonTitle) {
			*scheduleString += addDate(datesRange[*day])
			*day++
		}
		return false
	}
	return true
}

func isFoundLessonsWhileWeekConvert(scheduleString *string, lessonTitle string, day *int, datesRange []time.Time) bool {
	if !isFoundLessons(lessonTitle) {
		*scheduleString += addNotFoundLessons()
		if isNextDay(lessonTitle) {
			*scheduleString += addDate(datesRange[*day])
			*day++
		}
		return false
	}
	return true
}

func isNextDayWhileWeekConvert(scheduleString *string, lessonTitle string, day *int, datesRange []time.Time) bool {
	if isNextDay(lessonTitle) {
		*scheduleString += addDate(datesRange[*day])
		*day++
		return true
	}
	return false
}

func addLessonToScheduleString(scheduleString string, lesson *domain.Lesson) *string {
	scheduleString += addLessonNumber(lesson.StartTime)
	scheduleString += addLessonName(lesson.Title)
	scheduleString += addLessonStartTime(lesson.StartTime)
	scheduleString += addLessonFinishTime(lesson.FinishTime)
	scheduleString += addLessonType(lesson.Type)
	scheduleString += addLessonRoom(lesson.RoomId)
	scheduleString += addLessonTeacherFullName(lesson.TeacherFullName)
	return &scheduleString
}

func addGroupId(studyGroupId string) string {
	strSchedule := "Группа: " + fmt.Sprintf("%s\n\n", studyGroupId)
	return strSchedule
}

func addLessonNumber(lessonNumber string) string {
	strSchedule := "#" + strconv.Itoa(getLessonNumber(lessonNumber)) + ". "
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

func getLessonNumber(startTime string) int {
	switch startTime {
	case "10:10":
		return 2
	case "12:10":
		return 3
	case "13:50":
		return 4
	case "15:30":
		return 5
	case "17:10":
		return 6
	case "18:50":
		return 7
	default:
		return 1
	}
}

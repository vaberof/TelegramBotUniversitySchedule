package service

import (
	"fmt"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/storage"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/constants"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/pkg/date"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/pkg/utils"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/model"
)

// ScheduleToString converts schedule of type model.Schedule to string type
// to output it to user.
func ScheduleToString(studyGroupId, inputCallback string, schedule *model.Schedule, cache *storage.ScheduleStorage, chatID int64) *string {
	switch inputCallback {
	case constants.Today, constants.Tomorrow:
		dateToParse := date.GetParseDate(inputCallback)
		return DayScheduleToString(studyGroupId, inputCallback, dateToParse, schedule, cache, chatID)
	default:
		datesToParse := date.GetParseDates(inputCallback)
		return WeekScheduleToString(studyGroupId, inputCallback, datesToParse, schedule, cache, chatID)
	}
}

// DayScheduleToString converts schedule of type model.Schedule to string
// if user chosen schedule on day (today/tomorrow).
func DayScheduleToString(studyGroupId, inputCallback string, date time.Time, schedule *model.Schedule, cache *storage.ScheduleStorage, chatID int64) *string {
	var scheduleString string

	// lessonStartTimeField is number of StartTime field in model.Lesson
	// and necessary to add number of each lesson while converting schedule.
	lessonStartTimeField := 1

	// lessonNameField is number of Name field in model.Lesson.
	lessonNameField := 0

	// Adding user`s input group ID.
	scheduleString = addGroupId(studyGroupId)

	// Adding date (day, month, year) and week day of bold style
	// by using ParseMode "markdown" in bot.go.
	scheduleString += addDate(date)

	dereferenceSchedule := *schedule
	daySchedule := dereferenceSchedule[inputCallback]

	for i := 0; i < len(daySchedule); i++ {
		lessonData := reflect.ValueOf(daySchedule[i])
		for j := 0; j < lessonData.NumField(); j++ {

			if !isHaveLessons(lessonData.Field(lessonNameField)) {
				scheduleString += addNoLessons()
				break
			}

			if !isFoundLessons(lessonData.Field(lessonNameField)) {
				scheduleString += addNotFoundLessons()
				break
			}

			if isHttpError(lessonData.Field(lessonNameField)) {
				scheduleString = addGroupId(studyGroupId)
				scheduleString += addHttpError()
				return &scheduleString
			}

			scheduleString += addLessonNumber(lessonData.Field(lessonStartTimeField))

			scheduleString += addLessonName(lessonData.Field(j))
			j++

			scheduleString += addLessonStartTime(lessonData.Field(j))
			j++

			scheduleString += addLessonFinishTime(lessonData.Field(j))
			j++

			scheduleString += addLessonType(lessonData.Field(j))
			j++

			scheduleString += addLessonRoom(lessonData.Field(j))
			j++

			scheduleString += addLessonTeacherFullName(lessonData.Field(j))
		}
	}

	storeScheduleString := map[string]string{
		inputCallback: scheduleString,
	}
	cache.Schedule[chatID] = append(cache.Schedule[chatID], storeScheduleString)
	log.Printf("schedule cached: chatID: %d, key: %s", chatID, inputCallback)
	return &scheduleString
}

// WeekScheduleToString converts schedule of type model.Schedule to string
// if user chosen schedule on week (current week/next week).
func WeekScheduleToString(studyGroupId, inputCallback string, dates []time.Time, schedule *model.Schedule, cache *storage.ScheduleStorage, chatID int64) *string {
	var scheduleString string

	// index starting from 0 to 6 and necessary to go through given array of dates
	// to add dates of every week day and name of the week day to scheduleString while converting schedule.
	index := 0

	// currentLessonNumber is number of certain lesson, will be compared with maxLessonNumber to
	// output correct number.
	currentLessonNumber := 0

	// maxLessonNumber necessary to compare current lesson number with next lesson number while converting schedule
	// to add correct date of each weekday and name of weekday.
	maxLessonNumber := 1000

	// lessonStartTimeField is number of StartTime field in model.Lesson
	// and necessary to add number of each lesson while converting schedule.
	lessonStartTimeField := 1

	// lessonNameField is number of Name field in model.Lesson
	lessonNameField := 0

	scheduleString = addGroupId(studyGroupId)

	dereferenceSchedule := *schedule
	daySchedule := dereferenceSchedule[inputCallback]

	for i := 0; i < len(daySchedule); i++ {
		lessonData := reflect.ValueOf(daySchedule[i])
		for j := 0; j < lessonData.NumField(); j++ {

			if !isHaveLessons(lessonData.Field(lessonNameField)) {

				scheduleString += addDates(index, dates)
				scheduleString += addNoLessons()

				currentLessonNumber = 0
				maxLessonNumber = 1000
				index++
				break
			}

			if !isFoundLessons(lessonData.Field(lessonNameField)) {

				scheduleString += addDates(index, dates)
				scheduleString += addNotFoundLessons()

				currentLessonNumber = 0
				maxLessonNumber = 1000
				index++
				break
			}

			if isHttpError(lessonData.Field(lessonNameField)) {
				scheduleString = addGroupId(studyGroupId)
				scheduleString += addHttpError()

				return &scheduleString
			}

			currentLessonNumber = model.GetLessonNumber(fmt.Sprint(lessonData.Field(lessonStartTimeField).Interface()))

			// comparing lesson numbers: if current lesson number is less than maxLessonNumber,
			// than we get lessons for a new day and add date and day of the week.
			if currentLessonNumber <= maxLessonNumber {
				scheduleString += addDates(index, dates)
				index++
			}
			maxLessonNumber = currentLessonNumber

			scheduleString += addLessonNumber(lessonData.Field(lessonStartTimeField))

			scheduleString += addLessonName(lessonData.Field(j))
			j++

			scheduleString += addLessonStartTime(lessonData.Field(j))
			j++

			scheduleString += addLessonFinishTime(lessonData.Field(j))
			j++

			scheduleString += addLessonType(lessonData.Field(j))
			j++

			scheduleString += addLessonRoom(lessonData.Field(j))
			j++

			scheduleString += addLessonTeacherFullName(lessonData.Field(j))
		}
	}

	storeScheduleString := map[string]string{
		inputCallback: scheduleString,
	}
	cache.Schedule[chatID] = append(cache.Schedule[chatID], storeScheduleString)

	return &scheduleString
}

// addGroupId adds study group id to schedule of type string while converting schedule.
func addGroupId(studyGroupId string) string {
	scheduleString := "Группа: " + fmt.Sprintf("%s\n\n", studyGroupId)
	return scheduleString
}

func addLessonNumber(lessonNumber reflect.Value) string {
	scheduleString := "#" + strconv.Itoa(model.GetLessonNumber(lessonNumber.String())) + ". "
	return scheduleString
}

func addLessonName(lessonName reflect.Value) string {
	scheduleString := lessonName.String() + "("
	return scheduleString
}

func addLessonStartTime(startTime reflect.Value) string {
	scheduleString := startTime.String() + "-"
	return scheduleString
}

func addLessonFinishTime(finishTime reflect.Value) string {
	scheduleString := finishTime.String() + ", "
	return scheduleString
}

func addLessonType(lessonType reflect.Value) string {
	scheduleString := lessonType.String() + ", "
	return scheduleString
}

func addLessonRoom(lessonRoom reflect.Value) string {
	scheduleString := "ауд. " + lessonRoom.String() + ", "
	return scheduleString
}

func addLessonTeacherFullName(TeacherFullName reflect.Value) string {
	scheduleString := "препод. " + TeacherFullName.String() + ")" + "\n\n"
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

func addHttpError() string {
	scheduleString := "Ошибка: невозможно сделать запрос на сайт\n\n"
	return scheduleString
}

// addDate adds utils and day of the week of certain day of bold style
// to schedule of type string while converting schedule on day.
func addDate(date time.Time) string {
	scheduleString := fmt.Sprintf("*%s*", "Дата: ") +
		fmt.Sprintf("*%s*", date.Format("02.01.2006")) +
		fmt.Sprintf(" *(%s)*\n\n", utils.WeekDayEngRu(date.Weekday()))

	return scheduleString
}

// addDates adds date and day of the week of certain day
// to schedule of type string while converting schedule on week.
func addDates(index int, dates []time.Time) string {
	scheduleString := fmt.Sprintf("*%s*", "Дата: ") +
		fmt.Sprintf("*%s*", dates[index].Format("02.01.2006")) +
		fmt.Sprintf(" *(%s)*\n\n", utils.WeekDayEngRu(dates[index].Weekday()))

	return scheduleString
}

// isHaveLessons returns false if we don`t have lessons for certain day
// by checking 'Name' field in model.Lesson.
func isHaveLessons(lessonNameField reflect.Value) bool {
	return lessonNameField.String() != "no lessons"
}

// isFoundLessons returns false if we didn't find lessons for certain day
// by checking 'Name' field in model.Lesson.
func isFoundLessons(lessonNameField reflect.Value) bool {
	return lessonNameField.String() != "not found"
}

// isHttpError returns true if we got http error while making request to university website
// by checking 'Name' field in model.Lesson.
func isHttpError(lessonNameField reflect.Value) bool {
	return lessonNameField.String() == "http error"
}

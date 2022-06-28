package service

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/model"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/pkg/date"
)

// ScheduleToString converts schedule of type model.Schedule to string type
// to output it to user.
func ScheduleToString(studyGroupId, inputCallback string, Date *model.ParseData, schedule *model.Schedule) *string {
	switch inputCallback {
	case date.Today, date.Tomorrow:
		return DayScheduleToString(studyGroupId, inputCallback, Date, schedule)
	default:
		return WeekScheduleToString(studyGroupId, inputCallback, Date, schedule)
	}
}

// DayScheduleToString converts schedule of type model.Schedule to string
// if user chosen schedule on day (today/tomorrow).
func DayScheduleToString(studyGroupId, inputCallback string, Date *model.ParseData, schedule *model.Schedule) *string {
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
	scheduleString += addDate(Date)

	dereferenceSchedule := *schedule
	convSchedule := dereferenceSchedule[inputCallback]

	for i := 0; i < len(convSchedule); i++ {
		lessonData := reflect.ValueOf(convSchedule[i])
		for j := 0; j < lessonData.NumField(); j++ {

			if !haveLessons(lessonData.Field(lessonNameField)) {
				scheduleString += addNoLessons()
				break
			}

			if !foundLessons(lessonData.Field(lessonNameField)) {
				scheduleString += addNotFoundLessons()
				break
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

	return &scheduleString
}

// WeekScheduleToString converts schedule of type model.Schedule to string
// if user chosen schedule on week (current week/next week).
func WeekScheduleToString(studyGroupId, inputCallback string, Date *model.ParseData, schedule *model.Schedule) *string {
	var scheduleString string

	// index starting from 0 to 6 and necessary to go through arrays of dates and week days in !model.ParseData!
	// to add date of every week day and name of the week day to scheduleString while converting schedule.
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
	convSchedule := dereferenceSchedule[inputCallback]

	for i := 0; i < len(convSchedule); i++ {
		lessonData := reflect.ValueOf(convSchedule[i])
		for j := 0; j < lessonData.NumField(); j++ {

			if !haveLessons(lessonData.Field(lessonNameField)) {

				scheduleString += addDates(index, Date)
				scheduleString += addNoLessons()

				currentLessonNumber = 0
				maxLessonNumber = 1000
				index++
				break
			}

			if !foundLessons(lessonData.Field(lessonNameField)) {

				scheduleString += addDates(index, Date)
				scheduleString += addNotFoundLessons()

				currentLessonNumber = 0
				maxLessonNumber = 1000
				index++
				break
			}

			currentLessonNumber = model.GetLessonNumber(fmt.Sprint(lessonData.Field(lessonStartTimeField).Interface()))

			// comparing lesson numbers: if next lesson number is less than current,
			// than we get lessons for a new day and add date and day of the week.
			if currentLessonNumber <= maxLessonNumber {
				scheduleString += addDates(index, Date)
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

// addDate adds date and day of the week of certain day of bold style
// to schedule of type string while converting schedule on day.
func addDate(Date *model.ParseData) string {
	scheduleString := fmt.Sprintf("*%s*", "Дата: ") +
		fmt.Sprintf("*%s*", Date.Date.Format("02.01.2006")) +
		fmt.Sprintf(" *(%s)*\n\n", date.WeekDayEngRu(Date.Date.Weekday()))

	return scheduleString
}

// addDates adds date and day of the week of certain day of bold style
// to schedule of type string while converting schedule on week.
func addDates(index int, Date *model.ParseData) string {
	scheduleString := fmt.Sprintf("*%s*", "Дата: ") +
		fmt.Sprintf("*%s*", Date.Dates[index].Format("02.01.2006")) +
		fmt.Sprintf(" *(%s)*\n\n", date.WeekDayEngRu(Date.Days[index]))

	return scheduleString
}

// haveLessons returns false if we don`t have lessons for certain day
// by checking 'Name' field in model.Lesson.
func haveLessons(lessonNameField reflect.Value) bool {
	return lessonNameField.String() != "no lessons"
}

// foundLessons returns false if we didn't find lessons for certain day
// by checking 'Name' field in model.Lesson.
func foundLessons(lessonNameField reflect.Value) bool {
	return lessonNameField.String() != "not found"
}

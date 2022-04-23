package controller

import (
	"fmt"
	"time"

	"github.com/tg_bot_timetable/internal/model"
	"github.com/tg_bot_timetable/internal/services"
)

// HandleMessage
// returns schedule for user.
// If user`s input group not exists,
// returns a corresponding message.
func HandleMessage(studyGroupStorage *model.GroupStorage, date, userText string, location *time.Location) *string {
	var response string

	studyGroupId := userText
	studyGroupUrl, exists := studyGroupStorage.StudyGroup(studyGroupId)
	if !exists {
		response = fmt.Sprintf("Группа '%s' не существует", studyGroupId)
		return &response
	}

	schedule := services.GetSchedule(studyGroupId, date, *studyGroupUrl, location)
	scheduleString := services.ScheduleToString(&schedule)

	return &scheduleString
}
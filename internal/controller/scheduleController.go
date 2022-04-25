package controller

import (
	"fmt"
	"time"

	"github.com/TelegramBotUniversitySchedule/internal/model"
	"github.com/TelegramBotUniversitySchedule/internal/services"
)

// HandleMessage
// returns schedule for user.
// If user`s input group not exists,
// returns a corresponding message.
func HandleMessage(userChatID int64, date string, user *model.User, studyGroupStorage *model.GroupStorage, location *time.Location) *string {
	var response string

	studyGroupId := user.Data[userChatID]
	studyGroupUrl, exists := studyGroupStorage.StudyGroup(studyGroupId)
	if !exists {
		response = fmt.Sprintf("Группа '%s' не существует", studyGroupId)
		return &response
	}

	schedule := services.GetSchedule(studyGroupId, date, *studyGroupUrl, location)
	scheduleString := services.ScheduleToString(&schedule)

	return &scheduleString
}

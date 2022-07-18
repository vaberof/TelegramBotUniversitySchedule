package handler

import (
	"errors"
	"fmt"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/storage"
)

// HandleMessage
// returns group id, url to parse schedule and error.
// if user`s input group id not exists, then error value becomes appropriate.
func HandleMessage(chatID int64, msgStorage *storage.MessageStorage, groupStorage *storage.GroupStorage) (string, string, error) {

	studyGroupId := msgStorage.MessageData[chatID]
	studyGroupUrl, exists := groupStorage.StudyGroup(studyGroupId)
	if !exists {
		return studyGroupId, studyGroupUrl, errors.New(fmt.Sprintf("Группы %s не существует", studyGroupId))
	}

	return studyGroupId, studyGroupUrl, nil
}

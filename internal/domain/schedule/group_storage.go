package domain

import "github.com/vaberof/TelegramBotUniversitySchedule/internal/infra/storage"

type GroupStorageApi interface {
	GetStudyGroupQueryParams(groupId string) *string
}

type GroupStorage struct {
	GroupStorageApi
}

func NewGroupStorage() *GroupStorage {
	return &GroupStorage{
		GroupStorageApi: storage.NewGroupStorage(),
	}
}

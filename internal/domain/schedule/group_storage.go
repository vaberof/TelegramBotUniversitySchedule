package domain

import "github.com/vaberof/TelegramBotUniversitySchedule/internal/infra/storage"

type GroupStorageApi interface {
	GetStudyGroup(groupId string) *storage.Group
}

type GroupStorage struct {
	GroupStorageApi
}

func NewGroupStorage() *GroupStorage {
	return &GroupStorage{
		GroupStorageApi: storage.NewGroupStorage(),
	}
}

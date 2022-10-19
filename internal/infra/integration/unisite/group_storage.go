package infra

import "github.com/vaberof/TelegramBotUniversitySchedule/internal/infra/storage"

type GroupStorageApi interface {
	GetGroupExternalId(groupId string) *string
}

type GroupStorage struct {
	GroupStorageApi
}

func NewGroupStorage() *GroupStorage {
	return &GroupStorage{
		GroupStorageApi: storage.NewGroupStorage(),
	}
}

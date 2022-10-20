package infra

import "github.com/vaberof/TelegramBotUniversitySchedule/internal/infra/storage"

type GroupExternalIdReceiver interface {
	GetGroupExternalId(groupId string) *string
}

type GroupStorage struct {
	GroupExternalIdReceiver
}

func NewGroupStorage() *GroupStorage {
	return &GroupStorage{
		GroupExternalIdReceiver: storage.NewGroupStorage(),
	}
}

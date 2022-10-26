package infra

import (
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/infra/storage/postgres/grouppg"
	"gorm.io/gorm"
)

type GroupExternalIdReceiver interface {
	GetGroupExternalId(groupId string) *string
}

type GroupStorage struct {
	GroupExternalIdReceiver
}

func NewGroupStorage(db *gorm.DB) *GroupStorage {
	return &GroupStorage{
		GroupExternalIdReceiver: grouppg.NewGroupStoragePostgres(db),
	}
}

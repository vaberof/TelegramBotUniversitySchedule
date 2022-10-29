package grouppg

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type GroupStoragePostgres struct {
	db *gorm.DB
}

func NewGroupStoragePostgres(db *gorm.DB) *GroupStoragePostgres {
	initGroups(db)
	return &GroupStoragePostgres{db: db}
}

func (g *GroupStoragePostgres) GetGroupExternalId(groupId string) *string {
	var groups []Group
	g.db.Table("groups").Find(&groups)

	for i := 0; i < len(groups); i++ {
		group := groups[i]
		if group.Id+"-"+group.Name == groupId {
			return &group.ExternalId
		}
	}
	return nil
}

func initGroups(db *gorm.DB) {
	db.Table("groups").Where("group_id > 0").Delete(&Group{})

	err := db.Create(&studyGroups).Error
	if err != nil {
		log.Fatal("cannot create studyGroups in db ", err.Error())
	}
}

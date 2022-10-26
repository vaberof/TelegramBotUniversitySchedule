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
	groups := []Group{
		{
			Id:         "БИ",
			Name:       "21.1",
			ExternalId: "?ii=1&fi=1&c=2&gn=37&",
		},
		{
			Id:         "БИ",
			Name:       "21.2",
			ExternalId: "?ii=1&fi=1&c=2&gn=38&",
		},
	}

	db.Table("groups").Where("group_id > 0").Delete(&Group{})

	err := db.Create(&groups).Error
	if err != nil {
		log.Fatal("cannot create groups in db ", err.Error())
	}
}

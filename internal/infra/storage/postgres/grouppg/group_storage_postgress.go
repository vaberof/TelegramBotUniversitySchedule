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
		// ИНСТИТУТ ИГИМ
		// ГРУППЫ БГ
		{
			Id:         "БГ",
			Name:       "21.1",
			ExternalId: "?ii=1&fi=1&c=2&gn=31&",
		},
		{
			Id:         "БГ",
			Name:       "21.2",
			ExternalId: "?ii=1&fi=1&c=2&gn=32&",
		},
		{
			Id:         "БГ",
			Name:       "22.1",
			ExternalId: "?ii=1&fi=1&c=2&gn=33&",
		},

		{
			Id:         "БГ",
			Name:       "22.1",
			ExternalId: "?ii=1&fi=1&c=2&gn=33&",
		},

		// ГРУППЫ БГД
		{
			Id:         "БГД",
			Name:       "21.1",
			ExternalId: "?ii=1&fi=1&c=2&gn=34&",
		},
		{
			Id:         "БГД",
			Name:       "21.2",
			ExternalId: "?ii=1&fi=1&c=2&gn=35&",
		},

		// ГРУППЫ БИ
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
		{
			Id:         "БИ",
			Name:       "22.1",
			ExternalId: "?ii=1&fi=1&c=2&gn=39&",
		},
		{
			Id:         "БИ",
			Name:       "22.2",
			ExternalId: "?ii=1&fi=1&c=2&gn=40&",
		},
		{
			Id:         "БИ",
			Name:       "23.1",
			ExternalId: "?ii=1&fi=1&c=2&gn=41&",
		},
		{
			Id:         "БИ",
			Name:       "23.2",
			ExternalId: "?ii=1&fi=1&c=2&gn=42&",
		},

		// ГРУППЫ БК
		{
			Id:         "БК",
			Name:       "21.1",
			ExternalId: "?ii=1&fi=1&c=2&gn=254&",
		},
		{
			Id:         "БК",
			Name:       "21.2",
			ExternalId: "?ii=1&fi=1&c=2&gn=255&",
		},
		{
			Id:         "БК",
			Name:       "22.1",
			ExternalId: "?ii=1&fi=1&c=2&gn=256&",
		},

		// ГРУППЫ БМ
		{
			Id:         "БМ",
			Name:       "21.1",
			ExternalId: "?ii=1&fi=1&c=2&gn=46&",
		},

		// ГРУППЫ МГ
		{
			Id:         "МГ",
			Name:       "21.1",
			ExternalId: "?ii=1&fi=1&c=2&gn=90&",
		},

		// ГРУППЫ МГд
		{
			Id:         "МГд",
			Name:       "21.1",
			ExternalId: "?ii=1&fi=1&c=2&gn=92&",
		},

		// ГРУППЫ МГк
		{
			Id:         "МГк",
			Name:       "21.1",
			ExternalId: "?ii=1&fi=1&c=2&gn=91&",
		},

		// ГРУППЫ МД
		{
			Id:         "МД",
			Name:       "21.1",
			ExternalId: "?ii=1&fi=1&c=2&gn=100&",
		},
		{
			Id:         "МД",
			Name:       "21.2",
			ExternalId: "?ii=1&fi=1&c=2&gn=101&",
		},
		{
			Id:         "МД",
			Name:       "22.1",
			ExternalId: "?ii=1&fi=1&c=2&gn=102&",
		},

		// ГРУППЫ МИ
		{
			Id:         "МИ",
			Name:       "21.1",
			ExternalId: "?ii=1&fi=1&c=2&gn=93&",
		},

		// ГРУППЫ ПГ
		{
			Id:         "ПГ",
			Name:       "21.1",
			ExternalId: "?ii=1&fi=1&c=2&gn=103&",
		},
		{
			Id:         "ПГ",
			Name:       "21.2",
			ExternalId: "?ii=1&fi=1&c=2&gn=104&",
		},
		{
			Id:         "ПГ",
			Name:       "22.1",
			ExternalId: "?ii=1&fi=1&c=2&gn=105&",
		},

		// ГРУППЫ ЭН
		{
			Id:         "ЭН",
			Name:       "21.1",
			ExternalId: "?ii=1&fi=1&c=2&gn=192&",
		},
	}

	db.Table("groups").Where("group_id > 0").Delete(&Group{})

	err := db.Create(&groups).Error
	if err != nil {
		log.Fatal("cannot create groups in db ", err.Error())
	}
}

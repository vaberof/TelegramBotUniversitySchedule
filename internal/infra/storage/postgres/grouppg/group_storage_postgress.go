package grouppg

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type GroupStoragePostgres struct {
	db *gorm.DB
}

func NewGroupStoragePostgres(db *gorm.DB) *GroupStoragePostgres {
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

func (g *GroupStoragePostgres) CreateGroup(id string, name string, externalId string) error {
	_, err := g.getGroup(id, name)
	if err == nil {
		log.Error("group already exists in database, error: ", err)
		return errors.New("group already exists in database")
	}

	var group Group

	group.Id = id
	group.Name = name
	group.ExternalId = externalId

	err = g.db.Table("groups").Create(&group).Error
	if err != nil {
		log.Error("cannot create group in database, error: ", err)
		return err
	}
	log.Info("created the group in database")
	return nil
}

func (g *GroupStoragePostgres) UpdateGroupExternalId(
	id string,
	name string,
	newExternalId string) error {

	group, err := g.getGroup(id, name)
	if err != nil {
		log.Error("cannot update group in database, error: ", err)
		return err
	}

	group.ExternalId = newExternalId

	err = g.db.Table("groups").Save(&group).Error
	if err != nil {
		log.Error("cannot save group in database, error: ", err)
		return err
	}
	log.Info("group updated in database")
	return nil
}

func (g *GroupStoragePostgres) DeleteGroup(id string, name string) error {
	group, err := g.getGroup(id, name)
	if err != nil {
		return err
	}

	err = g.db.Table("groups").Delete(&group).Error
	if err != nil {
		log.Error("cannot delete group from database, error ", err)
		return err
	}
	log.Info("group deleted from database")
	return nil
}

func (g *GroupStoragePostgres) getGroup(id string, name string) (*Group, error) {
	var group Group

	err := g.db.Table("groups").Where("id = ? AND name = ?", id, name).
		First(&group).Error
	if err != nil {
		log.Error("cannot find group in database, error: ", err)
		return nil, errors.New("group does not exist")
	}
	return &group, nil
}

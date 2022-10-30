package handler

type GroupStorage interface {
	CreateGroup(id string, name string, externalId string) error
	UpdateGroup(id string, name string, externalId string, newId string, newName string, newExternalId string) error
	DeleteGroup(id string, name string, externalId string) error
}

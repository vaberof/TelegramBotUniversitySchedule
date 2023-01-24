package handler

type GroupStorage interface {
	CreateGroup(id string, name string, externalId string) error
	DeleteGroup(id string, name string) error
	UpdateGroupExternalId(id string, name string, newExternalId string) error
}

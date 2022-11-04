package group

type GroupStorage interface {
	CreateGroup(id string, name string, externalId string) error
	UpdateGroupExternalId(id string, name string, newExternalId string) error
	DeleteGroup(id string, name string) error
}

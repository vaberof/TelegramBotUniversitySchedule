package group

type GroupStorageService struct {
	groupStoragePostgres *GroupStoragePostgres
}

func NewGroupStorageService(groupStoragePostgres *GroupStoragePostgres) *GroupStorageService {
	return &GroupStorageService{
		groupStoragePostgres: groupStoragePostgres,
	}
}

func (s GroupStorageService) CreateGroup(id string, name string, externalId string) error {

	return s.groupStoragePostgres.GroupStorage.CreateGroup(id, name, externalId)
}

func (s GroupStorageService) UpdateGroup(
	id string,
	name string,
	externalId string,
	newId string,
	newName string,
	newExternalId string) error {
	return s.groupStoragePostgres.GroupStorage.UpdateGroup(id, name, externalId, newId, newName, newExternalId)
}

func (s GroupStorageService) DeleteGroup(id string, name string, externalId string) error {
	return s.groupStoragePostgres.GroupStorage.DeleteGroup(id, name, externalId)
}

func (s GroupStorageService) CreateGroups(id string, name string, externalId string) error {

	return s.groupStoragePostgres.GroupStorage.CreateGroup(id, name, externalId)
}

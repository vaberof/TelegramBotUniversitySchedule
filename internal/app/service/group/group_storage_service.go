package group

type GroupStorageService struct {
	groupStoragePostgres *GroupStoragePostgres
}

func NewGroupStorageService(groupStoragePostgres *GroupStoragePostgres) *GroupStorageService {
	return &GroupStorageService{
		groupStoragePostgres: groupStoragePostgres,
	}
}

func (s *GroupStorageService) CreateGroup(id string, name string, externalId string) error {
	return s.groupStoragePostgres.GroupStorage.CreateGroup(id, name, externalId)
}

func (s *GroupStorageService) UpdateGroupExternalId(id string, name string, newExternalId string) error {
	return s.groupStoragePostgres.GroupStorage.UpdateGroupExternalId(id, name, newExternalId)
}

func (s *GroupStorageService) DeleteGroup(id string, name string) error {
	return s.groupStoragePostgres.GroupStorage.DeleteGroup(id, name)
}

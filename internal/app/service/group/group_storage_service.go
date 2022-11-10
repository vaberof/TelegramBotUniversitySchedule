package group

type GroupStorageService struct {
	groupStorage GroupStorage
}

func NewGroupStorageService(groupStorage GroupStorage) *GroupStorageService {
	return &GroupStorageService{
		groupStorage: groupStorage,
	}
}

func (s *GroupStorageService) CreateGroup(id string, name string, externalId string) error {
	return s.groupStorage.CreateGroup(id, name, externalId)
}

func (s *GroupStorageService) UpdateGroupExternalId(id string, name string, newExternalId string) error {
	return s.groupStorage.UpdateGroupExternalId(id, name, newExternalId)
}

func (s *GroupStorageService) DeleteGroup(id string, name string) error {
	return s.groupStorage.DeleteGroup(id, name)
}

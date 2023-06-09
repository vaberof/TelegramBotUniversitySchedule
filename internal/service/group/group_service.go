package group

type GroupService struct {
	groupStorage GroupStorage
}

func NewGroupService(groupStorage GroupStorage) *GroupService {
	return &GroupService{
		groupStorage: groupStorage,
	}
}

func (s *GroupService) CreateGroup(id string, name string, externalId string) error {
	return s.groupStorage.CreateGroup(id, name, externalId)
}

func (s *GroupService) UpdateGroupExternalId(id string, name string, newExternalId string) error {
	return s.groupStorage.UpdateGroupExternalId(id, name, newExternalId)
}

func (s *GroupService) DeleteGroup(id string, name string) error {
	return s.groupStorage.DeleteGroup(id, name)
}

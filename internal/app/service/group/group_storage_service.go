package group

type GroupStorageService struct {
	groupStoragePostgres *GroupStoragePostgres
}

func NewGroupStorageService(groupStoragePostgres *GroupStoragePostgres) *GroupStorageService {
	return &GroupStorageService{
		groupStoragePostgres: groupStoragePostgres,
	}
}

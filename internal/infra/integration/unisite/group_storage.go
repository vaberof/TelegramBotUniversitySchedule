package unisite

type GroupStorage interface {
	GetGroupExternalId(groupId string) *string
}

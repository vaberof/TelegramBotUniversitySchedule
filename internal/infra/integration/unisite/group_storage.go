package infra

type GroupStorage interface {
	GetGroupExternalId(groupId string) *string
}

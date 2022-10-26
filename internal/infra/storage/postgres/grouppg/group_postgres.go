package grouppg

type Group struct {
	GroupId    uint `gorm:"primarykey"`
	Id         string
	Name       string
	ExternalId string
}

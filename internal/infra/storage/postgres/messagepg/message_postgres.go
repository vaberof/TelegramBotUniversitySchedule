package messagepg

type Message struct {
	Id     uint `gorm:"primarykey"`
	ChatId int64
	Text   string
}

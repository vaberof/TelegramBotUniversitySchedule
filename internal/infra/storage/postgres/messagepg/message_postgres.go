package messagepg

type Username string

type Message struct {
	Id     uint `gorm:"primarykey"`
	ChatId int64
	From   Username
	Text   string
}

package telegram

type MessageService interface {
	GetMessage(chatId int64) (*string, error)
	SaveMessage(chatId int64, from string, text string) error
}

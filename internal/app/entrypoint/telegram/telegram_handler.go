package telegram

type TelegramHandler struct {
	scheduleReceiver ScheduleReceiver
	messageService   MessageService
}

func NewTelegramHandler(scheduleReceiver ScheduleReceiver, messageService MessageService) *TelegramHandler {
	return &TelegramHandler{
		scheduleReceiver: scheduleReceiver,
		messageService:   messageService,
	}
}

package telegram

import (
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/service/message"
	domain "github.com/vaberof/TelegramBotUniversitySchedule/internal/domain/schedule"
	"time"
)

type ScheduleReceiver interface {
	GetSchedule(groupId string, from time.Time, to time.Time) (*domain.Schedule, error)
}

type MessageReceiverSaver interface {
	GetMessage(chatId int64) (*string, error)
	SaveMessage(chatId int64, text string) error
}

type TelegramHandler struct {
	ScheduleReceiver
	MessageReceiverSaver
}

func NewTelegramHandler(scheduleService *domain.ScheduleService, messageService *message.MessageService) *TelegramHandler {
	return &TelegramHandler{
		ScheduleReceiver:     scheduleService,
		MessageReceiverSaver: messageService,
	}
}

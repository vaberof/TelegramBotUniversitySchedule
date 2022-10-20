package telegram

import (
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/service/message"
	domain "github.com/vaberof/TelegramBotUniversitySchedule/internal/domain/schedule"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/infra/storage"
	"time"
)

type ScheduleReceiver interface {
	GetSchedule(groupId string, from time.Time, to time.Time) (*domain.Schedule, error)
}

type MessageReceiverSaver interface {
	GetMessage(chatId int64) (*storage.Message, error)
	SaveMessage(chatId int64, message string)
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

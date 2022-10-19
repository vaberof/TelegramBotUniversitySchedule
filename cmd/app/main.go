package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/vaberof/TelegramBotUniversitySchedule/configs"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/entrypoint/telegram"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/service/message"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/domain/schedule"
	infra "github.com/vaberof/TelegramBotUniversitySchedule/internal/infra/integration/unisite"
	"os"
	"time"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	botKeyboardMarkup := newBotKeyboardMarkup()

	botConfig := configs.NewBotConfig(os.Getenv("token"))

	bot := newBot(botConfig)

	botUpdatesChannel := tgbotapi.NewUpdate(0)
	botUpdatesChannel.Timeout = 60

	updates := bot.GetUpdatesChan(botUpdatesChannel)

	messageStorage := message.NewMessageStorage()
	userService := message.NewMessageService(messageStorage)

	host := os.Getenv("host")
	httpClientConfig := configs.NewHttpClientConfig(3 * time.Second)

	getScheduleResponseApi := infra.NewGetScheduleResponseApi(host, httpClientConfig)
	getScheduleResponseApiService := infra.NewGetScheduleResponseApiService(getScheduleResponseApi)

	scheduleStorage := domain.NewScheduleStorage()
	groupStorage := domain.NewGroupStorage()

	scheduleApi := domain.NewGetScheduleResponseApi(getScheduleResponseApiService)
	scheduleService := domain.NewScheduleService(scheduleStorage, groupStorage, scheduleApi)

	telegramHandler := telegram.NewTelegramHandler(scheduleService, userService)

	for update := range updates {
		if telegramHandler.CommandReceived(update) {
			telegramHandler.HandleCommandMessage(bot, update)
			continue
		} else if telegramHandler.MessageReceived(update) {
			telegramHandler.HandleNewMessage(bot, update, *botKeyboardMarkup)
		} else if telegramHandler.MenuButtonPressed(update) {
			telegramHandler.HandleMenuButtonPress(bot, update, *botKeyboardMarkup)
		}
	}
}

func newBot(config *configs.BotConfig) *tgbotapi.BotAPI {
	token := config.Token
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.WithFields(log.Fields{
			"bot":   bot.Self.UserName,
			"error": err,
		}).Panic("Failed to create a bot")
	}

	log.WithFields(log.Fields{
		"bot": bot.Self.UserName,
	}).Info("Bot is authorized")

	return bot
}

func newBotKeyboardMarkup() *tgbotapi.InlineKeyboardMarkup {
	botKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Сегодня", "Today"),
			tgbotapi.NewInlineKeyboardButtonData("Завтра", "Tomorrow"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Неделя", "Week"),
			tgbotapi.NewInlineKeyboardButtonData("След. неделя", "Next week"),
		),
	)

	return &botKeyboardMarkup
}

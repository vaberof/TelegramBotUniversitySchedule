package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/vaberof/TelegramBotUniversitySchedule/configs"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/entrypoint/telegram"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/domain/schedule"
	integration "github.com/vaberof/TelegramBotUniversitySchedule/pkg/integration/unisite"
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

	host := os.Getenv("host")
	httpClientConfig := configs.NewHttpClientConfig(3 * time.Second)
	httpClient := integration.NewHttpClient(host, httpClientConfig)

	scheduleStorage := domain.NewScheduleStorage()
	groupStorage := domain.NewGroupStorage()

	scheduleService := domain.NewScheduleService(scheduleStorage, groupStorage, httpClient)
	telegramHandler := telegram.NewTelegramHandler(scheduleService)

	var inputMessage *telegram.InputTelegramMessage

	for update := range updates {
		if telegramHandler.CommandReceived(update) {
			telegramHandler.HandleCommandMessage(bot, update)
			continue
		} else if telegramHandler.MessageReceived(update) {
			inputMessage = telegramHandler.HandleNewMessage(bot, update, *botKeyboardMarkup)
		} else if telegramHandler.MenuButtonPressed(update) {
			telegramHandler.HandleMenuButtonPress(bot, update, inputMessage, *botKeyboardMarkup)
		}
	}
}

// newBot creates bot.
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
		//tgbotapi.NewInlineKeyboardRow(
		//	tgbotapi.NewInlineKeyboardButtonData("Неделя", "Week"),
		//	tgbotapi.NewInlineKeyboardButtonData("След. неделя", "Next week"),
		//),
	)

	return &botKeyboardMarkup
}

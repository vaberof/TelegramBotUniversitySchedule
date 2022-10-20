package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/vaberof/TelegramBotUniversitySchedule/configs"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/entrypoint/telegram"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/service/message"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/domain/schedule"
	infra "github.com/vaberof/TelegramBotUniversitySchedule/internal/infra/integration/unisite"
	"os"
	"time"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("failed initializating config: %s", err.Error())
	}

	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	botKeyboardMarkup := newBotKeyboardMarkup()

	botConfig := configs.NewBotConfig(os.Getenv("token"))
	bot := newBot(botConfig)

	botUpdatesChannel := tgbotapi.NewUpdate(0)
	botUpdatesChannel.Timeout = 60

	updates := bot.GetUpdatesChan(botUpdatesChannel)

	messageStorage := message.NewMessageStorage()
	messageService := message.NewMessageService(messageStorage)

	httpClientConfig := configs.NewHttpClientConfig(
		viper.GetString("server.host"),
		time.Duration(viper.GetInt("server.timeout"))*time.Second)

	scheduleApi := infra.NewGetScheduleResponseApi(httpClientConfig)
	groupStorage := infra.NewGroupStorage()
	getScheduleResponseService := infra.NewGetScheduleResponseService(scheduleApi, groupStorage)

	getScheduleResponse := domain.NewGetScheduleResponse(getScheduleResponseService)
	scheduleStorage := domain.NewScheduleStorage()
	scheduleService := domain.NewScheduleService(getScheduleResponse, scheduleStorage)

	telegramHandler := telegram.NewTelegramHandler(scheduleService, messageService)

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

func initConfig() error {
	viper.AddConfigPath("../../configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

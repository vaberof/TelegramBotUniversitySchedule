package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/vaberof/TelegramBotUniversitySchedule/configs"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/entrypoint/telegram"
	"os"
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

	// httpCfg := configs.NewHttpClientConfig(1 * time.Second)
	// httpClient := integration.NewHttpClient(httpCfg)

	var inputMessage *telegram.InputMessage

	for update := range updates {
		if telegram.MessageReceived(update) {
			inputMessage = telegram.HandleNewMessage(bot, update, *botKeyboardMarkup)
		} else if telegram.MenuButtonPressed(update) {
			telegram.HandleMenuButtonPress(bot, update, inputMessage, *botKeyboardMarkup)
		} else if telegram.CommandReceived(update) {
			telegram.HandleCommandMessage(bot, update)
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
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Неделя", "Week"),
			tgbotapi.NewInlineKeyboardButtonData("След. неделя", "Next week"),
		),
	)

	return &botKeyboardMarkup
}

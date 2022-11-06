package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/vaberof/TelegramBotUniversitySchedule/configs"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/entrypoint/telegram"
	xhttp "github.com/vaberof/TelegramBotUniversitySchedule/internal/app/http/handler"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/service/auth"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/service/group"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/service/message"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/service/schedule"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/domain/schedule"
	infra "github.com/vaberof/TelegramBotUniversitySchedule/internal/infra/integration/unisite"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/infra/storage/postgres"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/infra/storage/postgres/grouppg"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/infra/storage/postgres/messagepg"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/infra/storage/postgres/schedulepg"
	integration "github.com/vaberof/TelegramBotUniversitySchedule/pkg/integration/unisite"
	"os"
	"time"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("failed initializating config: %s", err.Error())
	}

	db, err := postgres.NewPostgresDb(&postgres.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Name:     viper.GetString("db.name"),
		User:     os.Getenv("db_username"),
		Password: os.Getenv("db_password"),
	})
	if err != nil {
		log.Fatalf("cannot connect to database %s", err.Error())
	}

	err = db.AutoMigrate(&grouppg.Group{}, &messagepg.Message{}, &schedulepg.Schedule{}, &schedulepg.Lesson{})
	if err != nil {
		log.Fatalf("cannot auto migrate models %s", err.Error())
	}

	httpClientConfig := configs.NewHttpClientConfig(
		viper.GetString("server.host"),
		time.Duration(viper.GetInt("server.timeout"))*time.Second)

	httpClient := integration.NewHttpClient(httpClientConfig)

	groupStoragePostgres := grouppg.NewGroupStoragePostgres(db)
	messageStoragePostgres := messagepg.NewMessageStoragePostgres(db)
	scheduleStoragePostgres := schedulepg.NewScheduleStoragePostgres(db)

	getScheduleResponseService := infra.NewGetScheduleResponseService(httpClient, groupStoragePostgres)
	groupStorageService := group.NewGroupStorageService(groupStoragePostgres)
	messageStorageService := message.NewMessageStorageService(messageStoragePostgres)
	scheduleStorageService := schedule.NewScheduleStorageService(scheduleStoragePostgres)
	authService := auth.NewAuthService(os.Getenv("BEARER_TOKEN"))
	scheduleService := domain.NewScheduleService(getScheduleResponseService, scheduleStoragePostgres)

	telegramHandler := telegram.NewTelegramHandler(scheduleService, messageStorageService)
	httpHandler := xhttp.NewHttpHandler(groupStorageService, scheduleStorageService, authService)

	router := httpHandler.InitRouter()
	botConfig := configs.NewBotConfig(os.Getenv("TOKEN"))
	bot := newBot(botConfig)

	botKeyboardMarkup := newBotKeyboardMarkup()

	_, err = tgbotapi.NewWebhook(os.Getenv("BASE_URL") + "/" + bot.Token)
	if err != nil {
		log.Fatalln("Problem in setting Webhook", err.Error())
	}

	updates := bot.ListenForWebhook("/" + bot.Token)

	//router.POST("/"+bot.Token, nil)

	go router.Run(":" + os.Getenv("PORT"))

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

//func webhookHandler(c *gin.Context) {
//	defer c.Request.Body.Close()
//
//	bytes, err := ioutil.ReadAll(c.Request.Body)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//
//	var update tgbotapi.Update
//	err = json.Unmarshal(bytes, &update)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//
//	log.Printf("From: %+v Text: %+v\n", update.Message.From, update.Message.Text)
//}

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
	}).Info("ssugt_timetable_bot is authorized")
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
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs/")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

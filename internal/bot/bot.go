package bot

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// NewBot creates bot.
func NewBot() *tgbotapi.BotAPI{

	config := loadConfig()

	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Авторизован на аккаунте %s", bot.Self.UserName)

	return bot
}

// loadConfig loads config from .yaml file.
func loadConfig() *Config {

	config := NewConfig()
	yamlFile, err := ioutil.ReadFile("../../configs/app.yaml")
	if err != nil{
		log.Fatalf("Error %v", err)
	}

	err = yaml.Unmarshal(yamlFile, config)
	if err != nil{
		log.Fatalf("Error %v", err)
	}

	return config
}
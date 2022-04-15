package controllers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/tg_bot_timetable/services"
	"log"
)

const TOKEN = ""

// Создаем бота
func CreateBot(token string) *tgbotapi.BotAPI{
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Авторизован на аккаунте %s", bot.Self.UserName)

	return bot
}

// Получаем сообщение от пользователя
func handleMessage(userText string) *string {
	var response string

	groups := services.CreateGroupStorage()
	studyGroupId := userText
	studyGroupUrl := groups.GetGroupUrl(studyGroupId)

	if !groups.Exists(studyGroupId) {
		log.Printf("[ERROR] Группа %s не добавлена\n", studyGroupId)

		response = fmt.Sprintf("Группа %s не существует", studyGroupId)

		return &response
	}

	schedule := getTodaySchedule(studyGroupUrl)
	response = schedule

	return &response
}

// Запускаем бота
func StartBot() {
	bot := CreateBot(TOKEN)

	channel := tgbotapi.NewUpdate(0)
	channel.Timeout = 60

	updates := bot.GetUpdatesChan(channel)

	for update := range updates {

		if update.Message == nil {
			continue
		}

		responseMessage := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		inputText := update.Message.Text
		responseText := handleMessage(inputText)
		responseMessage.Text = *responseText
		bot.Send(responseMessage)
	}
}
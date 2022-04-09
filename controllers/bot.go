package controllers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/tg_bot_timetable/models"
	"log"
)

func BotInit() {
	bot, err := tgbotapi.NewBotAPI("5265008768:AAH8cVWRC5LkYoDVKUV2QYhPxZw2d4aSNYw")
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	channel := tgbotapi.NewUpdate(0)
	channel.Timeout = 60

	updates := bot.GetUpdatesChan(channel)

	//var saved bool
	//var curGroup string

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			userMessage := update.Message.Text


			group := models.GroupInit()
			// Проверяем, существует ли группа в словаре
			if group.GroupExists(userMessage) {
				//curGroup = userMessage
				//saved = true
				// Если сущесвует, то вызываем ф-ию расписания на "Сегодня"
				makeRequest(group.GetGroup(userMessage))
				msg.Text = "Группа" + "            " + userMessage + "\n"
				//msg.Text = ""
				msg.Text += getTodaySchedule(group.GetGroup(userMessage))
				bot.Send(msg)

				// Если группы не существует -> выводим соответствующее сообщение
			} else {
				log.Printf("[ERROR] Группа %s не добавлена\n", userMessage)
				msg.Text = "Группа не найдена" +
				"\n\n" + "Если хотите добавить Вашу группу - напишите @vaberof"
				bot.Send(msg)
			}
		}
	}
}
package controllers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/tg_bot_timetable/models"
	"log"
)

func BotInit() {
	bot, err := tgbotapi.NewBotAPI("") // Убрать токен когда буду пушить на гитхаб
	if err != nil {
		log.Panic(err)
	}

	//bot.Debug = true

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

			// Проверяем, существует ли группа в словаре
			if models.GroupExists(userMessage) {
				//curGroup = userMessage
				//saved = true
				// Если сущесвует, то вызываем ф-ию расписания на "Сегодня"
				DoRequest(models.Groups[userMessage])
				TodayShedule(models.Groups[userMessage])
				msg.Text = "[Группа]: " + userMessage + "\n\n"
				// Проходимся по массиву из ф-ии и выводим в тг
				for _, text := range ResponseMassive {
					msg.Text += text
				}
				// Проверяем, если длина массива 0 -> значит пар нет, выводим сообщение
				if len(ResponseMassive) == 0 {
					msg.Text = "Пар нет"
					bot.Send(msg)
					// Если массив непустой -> пары есть,
					// выводим расписание
					// обнуляем массив и счетчик для вывода даты
				} else {
					bot.Send(msg)
					ResponseMassive = []string{}
					DateCount = 0
				}
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

package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

type WebhookHandler struct {
	bot               *tgbotapi.BotAPI
	botKeyboardMarkup *tgbotapi.InlineKeyboardMarkup
	telegramHandler   TelegramHandler
}

func NewWebhookHandler(bot *tgbotapi.BotAPI, botKeyboardMarkup *tgbotapi.InlineKeyboardMarkup, telegramHandler TelegramHandler) *WebhookHandler {
	return &WebhookHandler{
		bot:               bot,
		botKeyboardMarkup: botKeyboardMarkup,
		telegramHandler:   telegramHandler,
	}
}

func (h *WebhookHandler) HandleWebhook(c *gin.Context) {
	defer c.Request.Body.Close()

	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "cannot read a request body",
		})
		log.Error("cannot read a request body, error: ", err.Error())
		return
	}

	update, err := h.fromBytesToUpdate(bytes)
	if err != nil {
		return
	}

	h.handleUpdate(update)
}

func (h *WebhookHandler) handleUpdate(update tgbotapi.Update) {
	if h.telegramHandler.CommandReceived(update) {
		h.telegramHandler.HandleCommandMessage(h.bot, update)
	} else if h.telegramHandler.MessageReceived(update) {
		h.telegramHandler.HandleNewMessage(h.bot, update, *h.botKeyboardMarkup)
	} else if h.telegramHandler.MenuButtonPressed(update) {
		h.telegramHandler.HandleMenuButtonPress(h.bot, update, *h.botKeyboardMarkup)
	}
}

func (h *WebhookHandler) fromBytesToUpdate(bytes []byte) (tgbotapi.Update, error) {
	var update tgbotapi.Update
	err := json.Unmarshal(bytes, &update)
	if err != nil {
		log.Error("cannot unmarshal bytes to update, error: ", err.Error())
		return tgbotapi.Update{}, err
	}
	return update, nil
}

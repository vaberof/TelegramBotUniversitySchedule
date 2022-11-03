package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/http/middleware/auth"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/service/group"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/service/token"
)

type HttpHandler struct {
	GroupStorage
	TokenValidator
}

func NewHttpHandler(groupStorageService *group.GroupStorageService, tokenService *token.TokenService) *HttpHandler {
	return &HttpHandler{
		GroupStorage:   groupStorageService,
		TokenValidator: tokenService,
	}
}

func (h *HttpHandler) InitRouter() *gin.Engine {
	router := gin.Default()

	router.Group("/")
	router.Use(auth.TokenAuth(h.TokenValidator))
	{
		router.POST("/group", h.CreateGroup)
		router.PUT("/group", h.UpdateGroup)
		router.DELETE("/group", h.DeleteGroup)
	}

	return router
}

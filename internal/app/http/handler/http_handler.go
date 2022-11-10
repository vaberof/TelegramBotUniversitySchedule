package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/http/middleware"
)

type HttpHandler struct {
	groupStorage    GroupStorage
	scheduleStorage ScheduleStorage
	tokenService    TokenService
}

func NewHttpHandler(groupStorage GroupStorage, scheduleStorage ScheduleStorage, tokenService TokenService) *HttpHandler {
	return &HttpHandler{
		groupStorage:    groupStorage,
		scheduleStorage: scheduleStorage,
		tokenService:    tokenService,
	}
}

func (h *HttpHandler) InitRouter() *gin.Engine {
	router := gin.Default()

	router.Group("/")
	router.Use(middleware.Auth(h.tokenService))
	{
		router.POST("/group", h.CreateGroup)
		router.PUT("/group", h.UpdateGroup)
		router.DELETE("/group", h.DeleteGroup)

		router.DELETE("/schedule", h.DeleteSchedule)
	}

	return router
}

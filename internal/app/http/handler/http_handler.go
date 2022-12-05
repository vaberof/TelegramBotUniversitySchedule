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
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	rGroup := router.Group("/group")
	rGroup.Use(middleware.Auth(h.tokenService))
	{
		rGroup.POST("/", h.CreateGroup)
		rGroup.PUT("/", h.UpdateGroup)
		rGroup.DELETE("/", h.DeleteGroup)
	}

	rSchedule := router.Group("/schedule")
	rSchedule.Use(middleware.Auth(h.tokenService))
	{
		rSchedule.DELETE("/", h.DeleteSchedule)
	}

	return router
}

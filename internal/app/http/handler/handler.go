package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/http/middleware"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/service/auth"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/service/group"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/service/schedule"
)

type HttpHandler struct {
	GroupStorage
	ScheduleStorage
	TokenValidator
}

func NewHttpHandler(groupStorageService *group.GroupStorageService, scheduleStorage *schedule.ScheduleStorageService, authService *auth.AuthService) *HttpHandler {
	return &HttpHandler{
		GroupStorage:    groupStorageService,
		ScheduleStorage: scheduleStorage,
		TokenValidator:  authService,
	}
}

func (h *HttpHandler) InitRouter() *gin.Engine {
	router := gin.Default()

	router.Group("/")
	router.Use(middleware.Auth(h.TokenValidator))
	{
		router.POST("/group", h.CreateGroup)
		router.PUT("/group", h.UpdateGroup)
		router.DELETE("/group", h.DeleteGroup)

		router.DELETE("/schedule", h.DeleteSchedule)
	}

	return router
}

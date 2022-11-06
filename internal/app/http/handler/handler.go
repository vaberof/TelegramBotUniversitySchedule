package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/http/middleware"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/service/auth"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/service/group"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/service/schedule"
)

type HttpHandler struct {
	groupStorage    GroupStorage
	scheduleStorage ScheduleStorage
	tokenService    TokenService
}

func NewHttpHandler(groupStorageService *group.GroupStorageService, scheduleStorageService *schedule.ScheduleStorageService, authService *auth.AuthService) *HttpHandler {
	return &HttpHandler{
		groupStorage:    groupStorageService,
		scheduleStorage: scheduleStorageService,
		tokenService:    authService,
	}
}

func (h *HttpHandler) InitRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	rGroup := router.Group("/group")
	rGroup.Use(middleware.Auth(h.tokenService))
	{
		rGroup.POST("/create", h.CreateGroup)
		rGroup.PUT("/update", h.UpdateGroup)
		rGroup.DELETE("/delete", h.DeleteGroup)
	}

	rSchedule := router.Group("/schedule")
	rSchedule.Use(middleware.Auth(h.tokenService))
	{
		router.DELETE("/delete", h.DeleteSchedule)
	}

	return router
}

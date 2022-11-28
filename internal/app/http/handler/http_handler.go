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
		rSchedule.DELETE("/delete", h.DeleteSchedule)
	}

	return router
}

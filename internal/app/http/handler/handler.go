package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/app/service/group"
)

type HttpHandler struct {
	GroupStorage
}

func NewHttpHandler(groupStorageService *group.GroupStorageService) *HttpHandler {
	return &HttpHandler{GroupStorage: groupStorageService}
}

func (h *HttpHandler) InitRoutes() *gin.Engine {
	router := gin.Default()

	router.POST("/group", h.CreateGroup)
	router.PUT("/group", h.UpdateGroup)
	router.DELETE("/group", h.DeleteGroup)

	return router
}

package handler

import "github.com/vaberof/TelegramBotUniversitySchedule/internal/app/service/group"

type HttpHandler struct {
	groupStorageService *group.GroupStorageService
}

func NewHttpHandler(groupStorageService *group.GroupStorageService) *HttpHandler {
	return &HttpHandler{groupStorageService: groupStorageService}
}

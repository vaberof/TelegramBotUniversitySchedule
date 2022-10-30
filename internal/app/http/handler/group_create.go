package handler

import "github.com/gin-gonic/gin"

type CreateGroupRequestBody struct {
	Id         string `json:"id" binding:"required"`
	Name       string `json:"name" binding:"required"`
	ExternalId string `json:"external_id" binding:"required"`
}

func (h *HttpHandler) CreateGroup(c *gin.Context) {
	var groupReqBody CreateGroupRequestBody

	if err := c.Bind(&groupReqBody); err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request body",
		})
		return
	}

	err := h.GroupStorage.CreateGroup(groupReqBody.Id, groupReqBody.Name, groupReqBody.ExternalId)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "group successfully created",
	})
}

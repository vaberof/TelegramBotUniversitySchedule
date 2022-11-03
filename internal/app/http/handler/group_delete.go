package handler

import "github.com/gin-gonic/gin"

type DeleteGroupRequestBody struct {
	Id   string `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

func (h *HttpHandler) DeleteGroup(c *gin.Context) {
	var groupReqBody DeleteGroupRequestBody

	if err := c.Bind(&groupReqBody); err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request body",
		})
		return
	}

	err := h.GroupStorage.DeleteGroup(groupReqBody.Id, groupReqBody.Name)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "group successfully deleted",
	})
}

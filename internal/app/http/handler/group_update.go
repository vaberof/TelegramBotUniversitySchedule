package handler

import "github.com/gin-gonic/gin"

type UpdateGroupRequestBody struct {
	Id            string `json:"id" binding:"required"`
	Name          string `json:"name" binding:"required"`
	NewExternalId string `json:"new_external_id" binding:"required"`
}

func (h *HttpHandler) UpdateGroup(c *gin.Context) {
	var groupReqBody UpdateGroupRequestBody

	if err := c.Bind(&groupReqBody); err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request body",
		})
		return
	}

	err := h.GroupStorage.UpdateGroup(
		groupReqBody.Id,
		groupReqBody.Name,
		groupReqBody.NewExternalId)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "group successfully updated",
	})
}

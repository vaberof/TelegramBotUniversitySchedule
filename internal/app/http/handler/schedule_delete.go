package handler

import "github.com/gin-gonic/gin"

type DeleteScheduleRequestBody struct {
	GroupId string `json:"group_id"  binding:"required"`
	Date    string `json:"date"  binding:"required"`
}

func (h *HttpHandler) DeleteSchedule(c *gin.Context) {
	var scheduleReqBody DeleteScheduleRequestBody

	if err := c.Bind(&scheduleReqBody); err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request body",
		})
		return
	}

	err := h.ScheduleStorage.DeleteSchedule(scheduleReqBody.GroupId, scheduleReqBody.Date)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "schedule successfully deleted",
	})
}

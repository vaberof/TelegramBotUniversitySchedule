package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/vaberof/TelegramBotUniversitySchedule/internal/service/auth"
	"log"
)

func Auth(s auth.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		stringToken := c.GetHeader("Authorization")
		if stringToken == "" {
			c.JSON(401, gin.H{
				"error": "request does not contain an access token",
			})
			c.Abort()
			log.Print("request does not contain an access token\n")
			return
		}

		err := s.ValidateToken(stringToken)
		if err != nil {
			c.JSON(401, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

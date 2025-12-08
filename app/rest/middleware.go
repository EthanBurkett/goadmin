package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ResponseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		data, _ := c.Get("data")

		messagesVal, _ := c.Get("messages")
		messages, ok := messagesVal.([]string)
		if !ok {
			messages = []string{}
		}

		status := c.Writer.Status()

		success := status < 400

		var errObj interface{} = nil
		if !success {
			// Check if a custom error message was set
			if errorVal, exists := c.Get("error"); exists {
				errObj = errorVal
			} else {
				errObj = http.StatusText(status)
			}
		}

		c.JSON(status, ApiResponse{
			Data:     data,
			Success:  success,
			Messages: messages,
			Error:    errObj,
		})
	}
}

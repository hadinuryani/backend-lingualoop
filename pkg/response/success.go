package response

import "github.com/gin-gonic/gin"

func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	if message == "" {
		message = MsgSuccess
	}

	res := DefaultResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	c.JSON(statusCode, res)
}

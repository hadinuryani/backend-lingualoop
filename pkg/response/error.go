package response

import "github.com/gin-gonic/gin"

func Error(c *gin.Context, statusCode int, message string, errs interface{}) {
	if message == "" {
		message = MsgInternalError
	}

	res := DefaultResponse{
		Success: false,
		Message: message,
		Errors:  errs,
	}
	c.JSON(statusCode, res)
}

package response

import "github.com/gin-gonic/gin"

func SuccessWithPagination(c *gin.Context, statusCode int, message string, data interface{}, meta Meta) {
	if message == "" {
		message = MsgSuccess
	}

	res := DefaultResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	}
	c.JSON(statusCode, res)
}

package helper

import "github.com/gin-gonic/gin"

// Response is a standardized JSON response structure
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// SendSuccessResponse sends a standardized success response.
func SendSuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// SendErrorResponse sends a standardized error response.
func SendErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Response{
		Status:  "error",
		Message: message,
		Data:    nil,
	})
	c.Abort()
}

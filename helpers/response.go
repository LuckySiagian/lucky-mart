package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type APIErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}

func ResponseSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	if data == nil {
		data = map[string]interface{}{}
	}
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ResponseError(c *gin.Context, statusCode int, message string, errs interface{}) {
	if errs == nil {
		errs = []string{}
	}
	c.JSON(statusCode, APIErrorResponse{
		Success: false,
		Message: message,
		Errors:  errs,
	})
}

func SendBadRequest(c *gin.Context, message string, errs interface{}) {
	ResponseError(c, http.StatusBadRequest, message, errs)
}

func SendUnauthorized(c *gin.Context, message string) {
	ResponseError(c, http.StatusUnauthorized, message, nil)
}

func SendForbidden(c *gin.Context, message string) {
	ResponseError(c, http.StatusForbidden, message, nil)
}

func SendNotFound(c *gin.Context, message string) {
	ResponseError(c, http.StatusNotFound, message, nil)
}

func SendValidationError(c *gin.Context, errs interface{}) {
	ResponseError(c, http.StatusUnprocessableEntity, "Validation Error", errs)
}

func SendInternalServerError(c *gin.Context, message string) {
	if message == "" {
		message = "Internal Server Error"
	}
	ResponseError(c, http.StatusInternalServerError, message, nil)
}

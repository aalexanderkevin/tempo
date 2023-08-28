package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	StatusCode int    `json:"error_code"`
	Message    string `json:"error_message"`
}

type SuccessResponse struct {
	Success bool `json:"success" default:"true"`
}

var SuccessOK = SuccessResponse{
	Success: true,
}

func SendErrorResponse(c *gin.Context, err ErrorResponse, msg string) {
	c.Writer.Header().Del("content-type")
	if msg != "" {
		err.Message = msg
	}
	status := http.StatusBadRequest
	if err.StatusCode != 0 {
		status = err.StatusCode
	}
	c.JSON(status, err)
}

func WriteSuccessResponse(c *gin.Context, payload interface{}) {
	if payload != nil {
		c.JSON(http.StatusOK, payload)
		return
	}
	c.JSON(http.StatusOK, SuccessOK)

}

func WriteFailResponse(c *gin.Context, statusCode int, err error) {
	SendErrorResponse(c, ErrorResponse{
		StatusCode: statusCode,
		Message:    err.Error(),
	}, err.Error())
}

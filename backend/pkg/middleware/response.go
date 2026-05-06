package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`    // 0 for success, 1 for failure, -1 for error
	Message string      `json:"message"` // response message
	Data    interface{} `json:"data"`    // response data
}

type EmptyData struct{}
type EmptySlice []struct{}

// Succ sends a successful response with a message and data.
func Succ(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: msg,
		Data:    data,
	})
}

// SuccNoMore sends a successful response with message and no more data.
func SuccNoMore(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: msg,
		Data:    EmptyData{},
	})
}

// Fail sends a failure response with a message.
func Fail(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{
		Code:    1,
		Message: msg,
		Data:    nil,
	})
}

// Erro sends an error response with an HTTP code and a message.
func Erro(c *gin.Context, code int, msg string) {
	c.JSON(code, Response{
		Code:    -1,
		Message: msg,
		Data:    nil,
	})
}

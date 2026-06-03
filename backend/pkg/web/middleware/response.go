package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidQueryBody = errors.New("Invalid query body")
)

type Pagination[T any] struct {
	Total    int `json:"total"`
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Items    []T `json:"items"`
}

func (p *Pagination[T]) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func (p *Pagination[T]) PageTotal() int {
	return (p.Total + p.PageSize - 1) / p.PageSize
}

type Response struct {
	Code    int    `json:"code"`    // 0 for success, 1 for failure, -1 for error
	Message string `json:"message"` // response message
	Data    any    `json:"data"`    // response data
}

type EmptyData struct{}
type EmptySlice []struct{}

// Succ sends a successful response with a message and data.
func Succ(c *gin.Context, msg string, data any) {
	if data == nil {
		data = EmptyData{}
	}
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

// SuccNoList sends a successful response with message and empty list.
func SuccNoList(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: msg,
		Data:    EmptySlice{},
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

func SSEHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Next()
	}
}

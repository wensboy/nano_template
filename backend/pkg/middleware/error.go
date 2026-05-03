package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GlobalErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			// 兜底错误处理，返回 500 错误
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "服务器内部错误",
				"details": err.Error(),
			})

			c.Abort()
		}
	}
}

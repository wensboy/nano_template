package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			// 500 internal server error
			Erro(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			c.Abort()
		}
	}
}

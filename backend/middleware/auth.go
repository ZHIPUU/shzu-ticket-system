package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIKeyAuth API Key 鉴权中间件
func APIKeyAuth(expected string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-API-Key")
		if key == "" || key != expected {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error_code":    "UNAUTHORIZED",
				"error_message": "Invalid or missing API Key",
				"detail":        "X-API-Key header is required and must match server config",
			})
			return
		}
		c.Next()
	}
}

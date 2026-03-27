package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Sleep - fake network latency
func Sleep(d time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		time.Sleep(d)
		c.Next()
	}
}

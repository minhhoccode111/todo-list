package middleware

import (
	"github.com/gin-gonic/gin"
)

// _ means auto401
func Auth(_ bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

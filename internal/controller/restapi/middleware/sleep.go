package middleware

import (
	"math/rand/v2"
	"time"

	"github.com/gin-gonic/gin"
)

// Sleep - fake network latency
func Sleep() gin.HandlerFunc {
	return func(c *gin.Context) {
		//nolint:gosec,mnd // fake latency - no security implications
		time.Sleep(time.Duration(rand.IntN(5)) * time.Second)
		c.Next()
	}
}

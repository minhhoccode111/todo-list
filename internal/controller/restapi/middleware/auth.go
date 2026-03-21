package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/minhhoccode111/todo-list/internal/entity"
)

// ctxKey is a unexported type to prevent key collisions
type ctxKey string

const (
	// CtxUserIDKey is used to store and retrieve the userID from the Gin context locals.
	CtxUserIDKey ctxKey = "userID"
	// CtxUserModelKey is used to store and retrieve the user model from the Gin context locals.
	CtxUserModelKey ctxKey = "userModel"
)

// Extract token from Authorization header or query parameter
func extractToken(c *gin.Context) string {
	// Check Authorization header first
	bearerToken := c.GetHeader("Authorization")
	if len(bearerToken) > 6 && strings.ToUpper(bearerToken[0:6]) == "TOKEN " {
		return bearerToken[6:]
	}

	// Check query parameter
	token := c.Query("access_token")
	if token != "" {
		return token
	}

	return ""
}

func Auth(auto401 bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(CtxUserIDKey, 0)
		c.Set(CtxUserModelKey, &entity.User{})
		// TODO: work on this
		c.Next()
	}
}

package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/minhhoccode111/todo-list/pkg/jwt"
)

type msg struct {
	Message string `json:"message"`
}

func messageResponse(c *gin.Context, code int, s string) {
	c.AbortWithStatusJSON(code, msg{Message: s})
}

// ctxKey is a unexported type to prevent key collisions
type ctxKey string

const (
	// CtxUserIDKey is used to store and retrieve the userID from the Gin context locals.
	CtxUserIDKey ctxKey = "userID"

	JWTScheme string = "Bearer"
)

// Extract token from Authorization header or query parameter
func extractToken(c *gin.Context) string {
	// Check Authorization header first
	token := c.GetHeader("Authorization")

	l := len(JWTScheme) + 1 // +1 extra space e.g. "Bearer "

	if len(token) > l && token[0:l] == JWTScheme+" " {
		return token[l:]
	}

	// Check query parameter
	token = c.Query("access_token")
	return token
}

func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			messageResponse(c, http.StatusUnauthorized, "Unauthorized")

			return
		}

		claims, err := jwt.ValidateToken(token, secret)
		if err != nil {
			messageResponse(c, http.StatusUnauthorized, "Unauthorized")

			return
		}

		id, err := strconv.ParseInt(claims.UserID, 10, 32)
		if err != nil {
			messageResponse(c, http.StatusUnauthorized, "Unauthorized")

			return
		}

		c.Set(CtxUserIDKey, int32(id))
		c.Next()
	}
}

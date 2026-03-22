package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/minhhoccode111/todo-list/internal/controller/restapi/v1/response"
)

func messageResponse(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, response.Message{Message: msg})
}

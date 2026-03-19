package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/minhhoccode111/todo-list/internal/usecase"
	"github.com/minhhoccode111/todo-list/pkg/logger"
)

// NewV1Routes -.
func NewV1Routes(
	apiV1Group *gin.RouterGroup,
	tr usecase.Translation,
	u usecase.User,
	to usecase.Todo,
	l logger.Interface,
	v *validator.Validate,
) {
	r := &V1{tr: tr, u: u, to: to, l: l, v: v}

	translationGroup := apiV1Group.Group("/translation")
	{
		translationGroup.GET("/history", r.history)
		translationGroup.POST("/do-translate", r.doTranslate)
	}

	userGroup := apiV1Group.Group("/")
	{
		userGroup.POST("/register", r.register)
		userGroup.POST("/login", r.login)
		// TODO:
	}

	todoGroup := apiV1Group.Group("/")
	{
		// TODO:
		_ = todoGroup
	}
}

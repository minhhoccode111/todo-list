package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/minhhoccode111/todo-list/internal/usecase"
	"github.com/minhhoccode111/todo-list/pkg/logger"
)

// NewTranslationRoutes -.
func NewTranslationRoutes(
	apiV1Group *gin.RouterGroup,
	tr usecase.Translation,
	l logger.Interface,
	v *validator.Validate,
) {
	r := &V1{tr: tr, l: l, v: v}

	translationGroup := apiV1Group.Group("/translation")

	{
		translationGroup.GET("/history", r.history)
		translationGroup.POST("/do-translate", r.doTranslate)
	}
}

// NewUserRoutes -.
func NewUserRoutes(
	apiV1Group *gin.RouterGroup,
	u usecase.User,
	l logger.Interface,
	v *validator.Validate,
) {
	r := &V1{u: u, l: l, v: v}

	userGroup := apiV1Group.Group("/")
	{
	}
	_ = r
	_ = userGroup
}

// NewTodoRoutes -.
func NewTodoRoutes(
	apiV1Group *gin.RouterGroup,
	to usecase.Todo,
	l logger.Interface,
	v *validator.Validate,
) {
	r := &V1{to: to, l: l, v: v}

	todoGroup := apiV1Group.Group("/")
	{
	}
	_ = r
	_ = todoGroup
}

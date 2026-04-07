package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/minhhoccode111/todo-list/config"
	"github.com/minhhoccode111/todo-list/internal/controller/restapi/middleware"
	"github.com/minhhoccode111/todo-list/internal/usecase"
	"github.com/minhhoccode111/todo-list/pkg/logger"
)

// NewV1Routes -.
func NewV1Routes(
	apiV1Group *gin.RouterGroup,
	l logger.Interface,
	v *validator.Validate,
	cfg *config.Config,
	tr usecase.Translation,
	u usecase.User,
	to usecase.Todo,
) {
	r := &V1{
		l:   l,
		v:   v,
		cfg: cfg,
		tr:  tr,
		u:   u,
		to:  to,
	}

	translationGroup := apiV1Group.Group("/translation")
	{
		// Route-level middlewares, example
		// translationGroup.GET("/example-admin", middleware.Auth(true), middleware.AdminOnly(), r.exampleAdmin)
		translationGroup.GET("/history", r.history)
		translationGroup.POST("/do-translate", r.doTranslate)
	}

	userGroup := apiV1Group.Group("/")
	{
		userGroup.POST("/register", r.register)
		userGroup.POST("/login", r.login)
		userGroup.POST("/refresh", r.refresh)
	}

	todoGroup := apiV1Group.Group("/")
	// RouterGroup-level middlewares
	todoGroup.Use(middleware.Auth(cfg.JWT.Secret))
	{
		todoGroup.GET("/todos", r.getTodos)
		todoGroup.POST("/todos", r.createTodo)
		todoGroup.PUT("/todos/:id", r.updateTodo)
		todoGroup.DELETE("/todos/:id", r.deleteTodo)
	}
}

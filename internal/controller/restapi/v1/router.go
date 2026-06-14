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
	u usecase.User,
	to usecase.Todo,
) {
	r := &V1{
		l:   l,
		v:   v,
		cfg: cfg,
		u:   u,
		to:  to,
	}

	userGroup := apiV1Group.Group("/")
	{
		userGroup.POST("/register", r.register)
		userGroup.POST("/login", r.login)
		userGroup.POST("/refresh", r.refresh)
	}

	userAuthGroup := apiV1Group.Group("/")
	userAuthGroup.Use(middleware.Auth(cfg.JWT.Secret))
	{
		userAuthGroup.POST("/logout", r.logout)
		userAuthGroup.POST("/logout/all", r.logoutAll)
		userAuthGroup.GET("/sessions", r.listSessions)
		userAuthGroup.DELETE("/sessions/:id", r.deleteSession)
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

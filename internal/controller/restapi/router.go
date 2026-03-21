package restapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/minhhoccode111/todo-list/config"
	_ "github.com/minhhoccode111/todo-list/docs" // Swagger docs.
	"github.com/minhhoccode111/todo-list/internal/controller/restapi/middleware"
	v1 "github.com/minhhoccode111/todo-list/internal/controller/restapi/v1"
	"github.com/minhhoccode111/todo-list/internal/usecase"
	"github.com/minhhoccode111/todo-list/pkg/logger"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

// NewRouter -.
// Swagger spec:
// @title       Todo-List API
// @description A Todo-List API with Gin and Clean Architecture
// @version     1.0
// @host        localhost:8080
// @BasePath    /api/v1
func NewRouter(
	handler *gin.Engine,
	cfg *config.Config,
	tr usecase.Translation,
	u usecase.User,
	to usecase.Todo,
	l logger.Interface,
	v *validator.Validate,
) {
	// Global-level (engine-wide) middlewares
	handler.Use(
		middleware.Logger(l),
		middleware.Recovery(l),
		middleware.CORS(cfg.CORS),
		middleware.RateLimit(cfg.RateLimit),
	)

	// Prometheus metrics
	if cfg.Metrics.Enabled {
		p := ginprometheus.NewPrometheus("todo-list-gin")
		p.Use(handler)
	}

	// Swagger
	if cfg.Swagger.Enabled {
		handler.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Routers
	apiV1Group := handler.Group("/api/v1")
	{
		v1.NewV1Routes(
			apiV1Group,
			l,
			v,
			cfg,
			tr,
			u,
			to,
		)
	}
}

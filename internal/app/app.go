// Package app configures and runs application.
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/minhhoccode111/todo-list/config"
	"github.com/minhhoccode111/todo-list/internal/controller/restapi"
	"github.com/minhhoccode111/todo-list/internal/entity"
	repocache "github.com/minhhoccode111/todo-list/internal/repo/cache"
	"github.com/minhhoccode111/todo-list/internal/repo/persistent"
	"github.com/minhhoccode111/todo-list/internal/repo/webapi"
	"github.com/minhhoccode111/todo-list/internal/usecase/todo"
	"github.com/minhhoccode111/todo-list/internal/usecase/translation"
	"github.com/minhhoccode111/todo-list/internal/usecase/user"
	"github.com/minhhoccode111/todo-list/pkg/cache"
	"github.com/minhhoccode111/todo-list/pkg/httpserver"
	"github.com/minhhoccode111/todo-list/pkg/logger"
	"github.com/minhhoccode111/todo-list/pkg/postgres"
	"github.com/minhhoccode111/todo-list/pkg/validatorx"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) { //nolint: gocyclo,cyclop,funlen,gocritic,nolintlint
	l := logger.New(cfg.Log.Level)
	v := validatorx.New()

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// Cache
	translationCache, err := cache.New[string, []entity.Translation](
		cache.MaxCost(cfg.Cache.MaxCost),
		cache.TTL(cfg.Cache.TTL),
	)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - cache.New: %w", err))
	}

	userCache, err := cache.New[string, *entity.User](
		cache.MaxCost(cfg.Cache.MaxCost),
		cache.TTL(cfg.Cache.TTL),
	)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - cache.New: %w", err))
	}

	todoCache, err := cache.New[string, *entity.Todos](
		cache.MaxCost(cfg.Cache.MaxCost),
		cache.TTL(cfg.Cache.TTL),
	)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - cache.New: %w", err))
	}

	// Use-Case
	translationUseCase := translation.New(
		persistent.NewTranslationRepo(pg),
		webapi.New(),
		repocache.NewTranslationCache(translationCache),
	)
	userUseCase := user.New(
		persistent.NewUserRepo(pg),
		repocache.NewUserCache(userCache),
	)
	todoUseCase := todo.New(
		persistent.NewTodoRepo(pg),
		repocache.NewTodoCache(todoCache),
	)

	// HTTP Server
	httpServer := httpserver.New(l, httpserver.Port(cfg.HTTP.Port))
	restapi.NewRouter(httpServer.Engine, cfg, translationUseCase, userUseCase, todoUseCase, l, v)

	// Start servers
	httpServer.Start()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}

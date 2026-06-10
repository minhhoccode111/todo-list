// Package repo implements application outer layer logic. Each logic group in own file.
package repo

import (
	"context"
	"time"

	"github.com/minhhoccode111/todo-list/internal/entity"
)

//go:generate mockgen -source=contracts.go -destination=../usecase/mocks_repo_test.go -package=usecase_test

type (
	// TranslationRepo -.
	TranslationRepo interface {
		CreateHistory(context.Context, entity.Translation) error
		ReadHistory(context.Context) ([]entity.Translation, error)
	}

	// TranslationWebAPI -.
	TranslationWebAPI interface {
		Translate(entity.Translation) (entity.Translation, error)
	}

	// TranslationCache is an in-process cache for translation data.
	// Implementations must be safe for concurrent use.
	TranslationCache interface {
		// GetHistory returns the cached translation history and whether it was found.
		GetHistory(ctx context.Context) ([]entity.Translation, bool)
		// SetHistory stores the translation history in the cache.
		// Returns false if the entry was not admitted by the cache.
		SetHistory(ctx context.Context, history []entity.Translation) bool
		// InvalidateHistory removes the cached history so the next read hits the DB.
		InvalidateHistory(ctx context.Context)
	}

	// UserRepo -.
	UserRepo interface {
		CreateUser(context.Context, *entity.User) (*entity.User, error)
		ReadUserByEmail(context.Context, string) (*entity.User, error)
		ReadUserByID(context.Context, int32) (*entity.User, error)
		CreateRefreshToken(
			c context.Context,
			userID int32,
			hashed, deviceInfo string,
			expiresAt time.Time,
		) error
		ReadRefreshToken(context.Context, string) (int32, error)
		DeleteRefreshTokenByID(c context.Context, userID, id int32) error
		DeleteRefreshTokenByHash(context.Context, int32, string) error
		DeleteAllRefreshTokens(c context.Context, userID int32) error
		ListRefreshTokens(c context.Context, userID int32) ([]entity.Session, error)
	}

	// UserCache -.
	UserCache interface {
		GetUser(c context.Context, userID string) (*entity.User, bool)
		SetUser(c context.Context, userID string, u *entity.User) bool
		InvalidateUser(c context.Context, userID string)
	}

	// TodoRepo -.
	TodoRepo interface {
		CreateTodo(context.Context, *entity.Todo) (*entity.Todo, error)
		UpdateTodo(context.Context, *entity.Todo) (*entity.Todo, error)
		DeleteTodo(context.Context, int32, int32) error
		ReadTodos(context.Context, int32, int32, int32) ([]entity.Todo, int32, error)
	}

	// TodoCache -.
	TodoCache interface {
		GetTodo(c context.Context, userID, todoID string) (*entity.Todo, bool)
		SetTodo(c context.Context, userID, todoID string, t *entity.Todo) bool
		InvalidateTodo(c context.Context, userID, todoID string)
		GetTodos(c context.Context, userID, limit, offset string) (*entity.Todos, bool)
		SetTodos(c context.Context, userID, limit, offset string, t *entity.Todos) bool
		InvalidateTodos(c context.Context, userID, limit, offset string)
		InvalidateAllTodos(c context.Context)
	}
)

// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"github.com/minhhoccode111/todo-list/config"
	"github.com/minhhoccode111/todo-list/internal/entity"
)

//go:generate mockgen -source=contracts.go -destination=./mocks_usecase_test.go -package=usecase_test

type (
	// Translation -.
	Translation interface {
		Translate(context.Context, entity.Translation) (entity.Translation, error)
		ReadHistory(context.Context) (entity.TranslationHistory, error)
	}

	// User -.
	User interface {
		Register(context.Context, *entity.User, *config.JWT) (string, error)
		Login(context.Context, *entity.User, *config.JWT) (string, error)
	}

	// ITodo -.
	Todo interface {
		CreateTodo(context.Context, *entity.Todo) (*entity.Todo, error)
		UpdateTodo(context.Context, *entity.Todo) (*entity.Todo, error)
		DeleteTodo(context.Context, int32, int32) error
	}
)

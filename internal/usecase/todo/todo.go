package todo

import (
	"context"

	"github.com/minhhoccode111/todo-list/internal/entity"
	"github.com/minhhoccode111/todo-list/internal/repo"
)

// UseCase -.
type UseCase struct {
	repo  repo.TodoRepo
	cache repo.TodoCache
}

// New -.
func New(r repo.TodoRepo, c repo.TodoCache) *UseCase {
	return &UseCase{
		repo:  r,
		cache: c,
	}
}

func (uc *UseCase) CreateTodo(c context.Context, id int32, t *entity.Todo) (*entity.Todo, error) {
	return nil, nil
}

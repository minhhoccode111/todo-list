package todo

import (
	"context"
	"fmt"
	"strconv"

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

func (uc *UseCase) CreateTodo(
	c context.Context,
	t *entity.Todo,
) (*entity.Todo, error) {
	t, err := uc.repo.CreateTodo(c, t)
	if err != nil {
		return nil, fmt.Errorf("TodoUseCase - CreateTodo - uc.repo.CreateTodo: %w", err)
	}

	userID, todoID := strconv.Itoa(int(t.UserID)), strconv.Itoa(int(t.ID))
	uc.cache.SetTodo(c, userID, todoID, t)

	return t, nil
}

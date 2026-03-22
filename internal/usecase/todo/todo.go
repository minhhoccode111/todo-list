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

func (uc *UseCase) CreateTodo(c context.Context, t *entity.Todo) (*entity.Todo, error) {
	t, err := uc.repo.CreateTodo(c, t)
	if err != nil {
		return nil, fmt.Errorf("TodoUseCase - CreateTodo - uc.repo.CreateTodo: %w", err)
	}

	userID, todoID := strconv.Itoa(int(t.UserID)), strconv.Itoa(int(t.ID))
	uc.cache.SetTodo(c, userID, todoID, t)

	return t, nil
}

func (uc *UseCase) UpdateTodo(c context.Context, t *entity.Todo) (*entity.Todo, error) {
	t, err := uc.repo.UpdateTodo(c, t)
	if err != nil {
		return nil, fmt.Errorf("TodoUseCase - UpdateTodo - uc.repo.UpdateTodo: %w", err)
	}

	userID, todoID := strconv.Itoa(int(t.UserID)), strconv.Itoa(int(t.ID))
	uc.cache.SetTodo(c, userID, todoID, t)

	return t, nil
}

func (uc *UseCase) DeleteTodo(c context.Context, todoID, userID int32) error {
	err := uc.repo.DeleteTodo(c, todoID, userID)
	if err != nil {
		return fmt.Errorf("TodoUseCase - DeleteTodo - uc.repo.DeleteTodo: %w", err)
	}

	uc.cache.InvalidateTodo(c, strconv.Itoa(int(userID)), strconv.Itoa(int(todoID)))

	return nil
}

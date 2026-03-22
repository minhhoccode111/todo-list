package cache

import (
	"context"

	"github.com/minhhoccode111/todo-list/internal/entity"
	"github.com/minhhoccode111/todo-list/pkg/cache"
)

const todoKey = "todo"

// TodoCache wraps an ottercache.Cache and implements repo.TodoCache.
type TodoCache struct {
	c *cache.Cache[string, []entity.Todo]
}

// New creates a new TodoCache.
// The TTL and max-cost are configured on the ottercache.Cache itself via ottercache.TTL / ottercache.MaxCost.
func NewTodoCache(c *cache.Cache[string, []entity.Todo]) *TodoCache {
	return &TodoCache{c: c}
}

func (tc *TodoCache) GetTodo(c context.Context, userID, todoID string) (*entity.Todo, bool) {
	todos, ok := tc.c.Get(buildKey(todoKey, userID, todoID))
	if !ok || len(todos) == 0 {
		return nil, false
	}

	return &todos[0], ok
}

func (tc *TodoCache) SetTodo(c context.Context, userID, todoID string, t *entity.Todo) bool {
	if t == nil {
		return false
	}

	return tc.c.Set(buildKey(todoKey, userID, todoID), []entity.Todo{*t})
}

func (tc *TodoCache) InvalidateTodo(c context.Context, userID, todoID string) {
	tc.c.Delete(buildKey(todoKey, userID, todoID))
}

func (tc *TodoCache) GetTodos(
	c context.Context,
	userID, limit, offset string,
) ([]entity.Todo, bool) {
	return tc.c.Get(buildKey(todoKey, userID, limit, offset))
}

func (tc *TodoCache) SetTodos(
	c context.Context,
	userID, limit, offset string,
	s []entity.Todo,
) bool {
	return tc.c.Set(buildKey(todoKey, userID, limit, offset), s)
}

func (tc *TodoCache) InvalidateTodos(c context.Context, userID, limit, offset string) {
	tc.c.Delete(buildKey(todoKey, userID, limit, offset))
}

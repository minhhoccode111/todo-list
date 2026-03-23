package cache

import (
	"context"

	"github.com/minhhoccode111/todo-list/internal/entity"
	"github.com/minhhoccode111/todo-list/pkg/cache"
)

const (
	todoKey  = "todo"
	todosKey = "todos"
)

// TodoCache wraps an ottercache.Cache and implements repo.TodoCache.
type TodoCache struct {
	c *cache.Cache[string, *entity.Todos]
}

// New creates a new TodoCache.
// The TTL and max-cost are configured on the ottercache.Cache itself via ottercache.TTL / ottercache.MaxCost.
func NewTodoCache(c *cache.Cache[string, *entity.Todos]) *TodoCache {
	return &TodoCache{c: c}
}

func (tc *TodoCache) GetTodo(_ context.Context, userID, todoID string) (*entity.Todo, bool) {
	todos, ok := tc.c.Get(buildKey(todoKey, userID, todoID))
	if !ok || todos == nil || len(todos.Data) == 0 {
		return nil, false
	}

	return &todos.Data[0], ok
}

func (tc *TodoCache) SetTodo(_ context.Context, userID, todoID string, t *entity.Todo) bool {
	if t == nil {
		return false
	}

	return tc.c.Set(buildKey(todoKey, userID, todoID), &entity.Todos{Data: []entity.Todo{*t}})
}

func (tc *TodoCache) InvalidateTodo(_ context.Context, userID, todoID string) {
	tc.c.Delete(buildKey(todoKey, userID, todoID))
}

func (tc *TodoCache) GetTodos(
	_ context.Context,
	userID, limit, offset string,
) (*entity.Todos, bool) {
	return tc.c.Get(buildKey(todosKey, userID, limit, offset))
}

func (tc *TodoCache) SetTodos(
	_ context.Context,
	userID, limit, offset string,
	t *entity.Todos,
) bool {
	return tc.c.Set(buildKey(todosKey, userID, limit, offset), t)
}

func (tc *TodoCache) InvalidateTodos(_ context.Context, userID, limit, offset string) {
	tc.c.Delete(buildKey(todosKey, userID, limit, offset))
}

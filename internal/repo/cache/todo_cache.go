package cache

import (
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

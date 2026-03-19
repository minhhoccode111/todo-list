package cache

import (
	"github.com/minhhoccode111/todo-list/internal/entity"
	"github.com/minhhoccode111/todo-list/pkg/cache"
)

const userKey = "user"

// UserCache wraps an ottercache.Cache and implements repo.UserCache.
type UserCache struct {
	c *cache.Cache[string, entity.User]
}

// New creates a new UserCache.
// The TTL and max-cost are configured on the ottercache.Cache itself via ottercache.TTL / ottercache.MaxCost.
func NewUserCache(c *cache.Cache[string, entity.User]) *UserCache {
	return &UserCache{c: c}
}

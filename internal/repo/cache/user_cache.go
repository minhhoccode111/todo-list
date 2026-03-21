package cache

import (
	"context"

	"github.com/minhhoccode111/todo-list/internal/entity"
	"github.com/minhhoccode111/todo-list/pkg/cache"
)

const userKey = "user"

// UserCache wraps an ottercache.Cache and implements repo.UserCache.
type UserCache struct {
	c *cache.Cache[string, *entity.User]
}

// New creates a new UserCache.
// The TTL and max-cost are configured on the ottercache.Cache itself via ottercache.TTL / ottercache.MaxCost.
func NewUserCache(c *cache.Cache[string, *entity.User]) *UserCache {
	return &UserCache{c: c}
}

// GetUser returns the cached user by userID.
func (uc *UserCache) GetUser(_ context.Context, userID string) (*entity.User, bool) {
	return uc.c.Get(userKey + ":" + userID)
}

// SetUser stores the user in the cache by userID.
func (uc *UserCache) SetUser(_ context.Context, userID string, u *entity.User) bool {
	return uc.c.Set(userKey+":"+userID, u)
}

// InvalidateUser removes the cached user entry so the next read hits the DB.
func (uc *UserCache) InvalidateUser(_ context.Context, userID string) {
	uc.c.Delete(userKey + ":" + userID)
}

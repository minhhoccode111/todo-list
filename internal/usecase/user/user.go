package user

import (
	"github.com/minhhoccode111/todo-list/internal/repo"
)

// UseCase -.
type UseCase struct {
	repo  repo.UserRepo
	cache repo.UserCache
}

// New -.
func New(r repo.UserRepo, c repo.UserCache) *UseCase {
	return &UseCase{
		repo:  r,
		cache: c,
	}
}

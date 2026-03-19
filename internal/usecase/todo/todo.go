package todo

import (
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

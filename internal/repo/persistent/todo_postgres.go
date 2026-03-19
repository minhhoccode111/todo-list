package persistent

import (
	"github.com/minhhoccode111/todo-list/internal/repo/persistent/sqlc"
	"github.com/minhhoccode111/todo-list/pkg/postgres"
)

// TodoRepo -.
type TodoRepo struct {
	*postgres.Postgres
	queries *sqlc.Queries
}

// NewTodoRepo -.
func NewTodoRepo(pg *postgres.Postgres) *TodoRepo {
	return &TodoRepo{
		Postgres: pg,
		queries:  sqlc.New(pg.Pool),
	}
}

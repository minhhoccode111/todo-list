package persistent

import (
	"github.com/minhhoccode111/todo-list/internal/repo/persistent/sqlc"
	"github.com/minhhoccode111/todo-list/pkg/postgres"
)

// UserRepo -.
type UserRepo struct {
	*postgres.Postgres
	queries *sqlc.Queries
}

// NewUserRepo -.
func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{
		Postgres: pg,
		queries:  sqlc.New(pg.Pool),
	}
}

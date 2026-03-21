package persistent

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/minhhoccode111/todo-list/internal/entity"
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

func newUserFromDTO(u *sqlc.User) *entity.User {
	return &entity.User{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

func (ur *UserRepo) CreateUser(c context.Context, u *entity.User) (*entity.User, error) {
	sqlcUser, err := ur.queries.CreateUser(c, sqlc.CreateUserParams{
		Email:        u.Email,
		Name:         u.Name,
		PasswordHash: u.PasswordHash,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, entity.ErrDuplicateEntry
		}

		return nil, fmt.Errorf("UserRepo - CreateUser - ur.queries.CreateUser: %w", err)
	}

	return newUserFromDTO(&sqlcUser), nil
}

func (ur *UserRepo) ReadUserByID(c context.Context, id string) (entity.User, error) {
	return entity.User{}, nil
}

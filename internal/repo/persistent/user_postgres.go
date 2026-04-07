package persistent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
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

func (ur *UserRepo) ReadUserByEmail(c context.Context, email string) (*entity.User, error) {
	u, err := ur.queries.ReadUserByEmail(c, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrNoRows
		}

		return nil, fmt.Errorf("UserRepo - ReadUserByEmail - ur.queries.ReadUserByEmail: %w", err)
	}

	return newUserFromDTO(&u), nil
}

func (ur *UserRepo) ReadUserByID(c context.Context, id int32) (*entity.User, error) {
	u, err := ur.queries.ReadUserByID(c, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrNoRows
		}

		return nil, fmt.Errorf("UserRepo - ReadUserByID - ur.queries.ReadUserByID: %w", err)
	}

	return newUserFromDTO(&u), nil
}

func (ur *UserRepo) CreateRefreshToken(
	c context.Context,
	userID int32,
	hashed, deviceInfo string,
	expiresAt time.Time,
) error {
	err := ur.queries.CreateRefreshToken(c, sqlc.CreateRefreshTokenParams{
		UserID:     userID,
		TokenHash:  hashed,
		ExpiresAt:  expiresAt,
		DeviceInfo: deviceInfo,
	})
	if err != nil {
		return fmt.Errorf("UserRepo - CreateRefreshToken - ur.queries.CreateRefreshToken: %w", err)
	}

	return nil
}

func (ur *UserRepo) ReadRefreshToken(c context.Context, refresh string) (int32, error) {
	userID, err := ur.queries.ReadRefreshToken(c, refresh)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, entity.ErrNoRows
		}

		return 0, fmt.Errorf("UserRepo - ReadRefreshToken - ur.queries.ReadRefreshToken: %w", err)
	}

	return userID, nil
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

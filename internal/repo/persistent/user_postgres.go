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
	tokenHash, deviceInfo string,
	expiresAt time.Time,
) error {
	_, err := ur.queries.CreateRefreshToken(c, sqlc.CreateRefreshTokenParams{
		UserID:     userID,
		TokenHash:  tokenHash,
		ExpiresAt:  expiresAt,
		DeviceInfo: deviceInfo,
	})
	if err != nil {
		return fmt.Errorf("UserRepo - CreateRefreshToken - ur.queries.CreateRefreshToken: %w", err)
	}

	return nil
}

func (ur *UserRepo) ReadRefreshTokenByHash(c context.Context, tokenHash string) (int32, error) {
	token, err := ur.queries.ReadRefreshTokenByHash(c, tokenHash)
	if err != nil {
		return 0, fmt.Errorf(
			"UserRepo - ReadRefreshTokenByHash - ur.queries.ReadRefreshTokenByHash: %w",
			err,
		)
	}

	return token.UserID, nil
}

func (ur *UserRepo) RevokeRefreshToken(c context.Context, tokenHash string) error {
	err := ur.queries.RevokeRefreshToken(c, tokenHash)
	if err != nil {
		return fmt.Errorf("UserRepo - RevokeRefreshToken - ur.queries.RevokeRefreshToken: %w", err)
	}

	return nil
}

func (ur *UserRepo) RevokeAllUserRefreshTokens(c context.Context, userID int32) error {
	err := ur.queries.RevokeAllUserRefreshTokens(c, userID)
	if err != nil {
		return fmt.Errorf(
			"UserRepo - RevokeAllUserRefreshTokens - ur.queries.RevokeAllUserRefreshTokens: %w",
			err,
		)
	}

	return nil
}

func (ur *UserRepo) DeleteExpiredRefreshTokens(c context.Context) error {
	err := ur.queries.DeleteExpiredRefreshTokens(c)
	if err != nil {
		return fmt.Errorf(
			"UserRepo - DeleteExpiredRefreshTokens - ur.queries.DeleteExpiredRefreshTokens: %w",
			err,
		)
	}

	return nil
}

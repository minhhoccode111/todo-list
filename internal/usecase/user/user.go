package user //nolint:revive // intended

import (
	"context"
	"fmt"
	"strconv"

	"github.com/minhhoccode111/todo-list/config"
	"github.com/minhhoccode111/todo-list/internal/entity"
	"github.com/minhhoccode111/todo-list/internal/repo"
	"github.com/minhhoccode111/todo-list/pkg/jwt"
	"github.com/minhhoccode111/todo-list/pkg/password"
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

// generateToken is just a local version of jwt.GenerateToken that can take *config.JWT as an argument
// so that we don't have to pass 3 arguments at a time :)
func generateToken(userID string, cfgJWT *config.JWT) (string, error) {
	return jwt.GenerateToken(
		userID,
		cfgJWT.Secret,
		cfgJWT.Issuer,
		cfgJWT.Expiration,
	)
}

func (uc *UseCase) Register(c context.Context, u *entity.User, cfgJWT *config.JWT) (string, error) {
	hashed, err := password.HashPassword(u.PasswordHash)
	if err != nil {
		return "", fmt.Errorf(
			"UserUseCase - Register - password.HashPassword: %w",
			err,
		)
	}

	u.PasswordHash = hashed

	u, err = uc.repo.CreateUser(c, u)
	if err != nil {
		return "", fmt.Errorf(
			"UserUseCase - Register - uc.repo.CreateUser: %w",
			err,
		)
	}

	userID := strconv.Itoa(int(u.ID))

	token, err := generateToken(userID, cfgJWT)
	if err != nil {
		return "", fmt.Errorf(
			"UserUseCase - Register - generateToken: %w",
			err,
		)
	}

	uc.cache.SetUser(c, userID, u)

	return token, nil
}

func (uc *UseCase) Login(c context.Context, u *entity.User, cfgJWT *config.JWT) (string, error) {
	user, err := uc.repo.ReadUserByEmail(c, u.Email)
	if err != nil {
		return "", fmt.Errorf(
			"UserUseCase - Login - uc.repo.ReadUserByEmail: %w",
			err,
		)
	}

	if !password.ComparePassword(user.PasswordHash, u.PasswordHash) {
		return "", entity.ErrUnauthorized
	}

	userID := strconv.Itoa(int(user.ID))

	token, err := generateToken(userID, cfgJWT)
	if err != nil {
		return "", fmt.Errorf(
			"UserUseCase - Login - generateToken: %w",
			err,
		)
	}

	uc.cache.SetUser(c, userID, user)

	return token, nil
}

package user //nolint:revive // intended

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/minhhoccode111/todo-list/config"
	"github.com/minhhoccode111/todo-list/internal/entity"
	"github.com/minhhoccode111/todo-list/internal/repo"
	"github.com/minhhoccode111/todo-list/pkg/jwt"
	"github.com/minhhoccode111/todo-list/pkg/utils"
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

func (uc *UseCase) Register(
	c context.Context,
	cfg *config.Config,
	u *entity.User,
) (token, refresh string, err error) {
	hashed, err := utils.HashPassword(u.PasswordHash)
	if err != nil {
		return "", "", fmt.Errorf(
			"UserUseCase - Register - password.HashPassword: %w",
			err,
		)
	}

	u.PasswordHash = hashed

	u, err = uc.repo.CreateUser(c, u)
	if err != nil {
		return "", "", fmt.Errorf(
			"UserUseCase - Register - uc.repo.CreateUser: %w",
			err,
		)
	}

	token, err = generateToken(u.ID, &cfg.JWT)
	if err != nil {
		return "", "", fmt.Errorf(
			"UserUseCase - Register - generateToken: %w",
			err,
		)
	}

	raw, hashed, err := utils.NewRefreshToken()
	if err != nil {
		return "", "", fmt.Errorf(
			"UserUseCase - Register - utils.NewRefreshToken: %w",
			err,
		)
	}

	err = uc.repo.CreateRefreshToken(
		c,
		u.ID,
		hashed,
		"",
		time.Now().Add(cfg.RefreshToken.Expiration),
	)
	if err != nil {
		return "", "", fmt.Errorf(
			"UserUseCase - Register - uc.repo.CreateRefreshToken: %w",
			err,
		)
	}

	uc.cache.SetUser(c, strconv.Itoa(int(u.ID)), u)

	return token, raw, nil
}

func (uc *UseCase) Login(
	c context.Context,
	cfg *config.Config,
	u *entity.User,
) (token, refresh string, err error) {
	user, err := uc.repo.ReadUserByEmail(c, u.Email)
	if err != nil {
		return "", "", fmt.Errorf(
			"UserUseCase - Login - uc.repo.ReadUserByEmail: %w",
			err,
		)
	}

	if !utils.ComparePassword(user.PasswordHash, u.PasswordHash) {
		return "", "", entity.ErrUnauthorized
	}

	token, err = generateToken(user.ID, &cfg.JWT)
	if err != nil {
		return "", "", fmt.Errorf(
			"UserUseCase - Login - generateToken: %w",
			err,
		)
	}

	refresh, hashed, err := utils.NewRefreshToken()
	if err != nil {
		return "", "", fmt.Errorf(
			"UserUseCase - Login - utils.NewRefreshToken: %w",
			err,
		)
	}

	err = uc.repo.CreateRefreshToken(
		c,
		user.ID,
		hashed,
		"",
		time.Now().Add(cfg.RefreshToken.Expiration),
	)
	if err != nil {
		return "", "", fmt.Errorf(
			"UserUseCase - Login - uc.repo.CreateRefreshToken: %w",
			err,
		)
	}

	uc.cache.SetUser(c, strconv.Itoa(int(user.ID)), user)

	return token, refresh, nil
}

func (uc *UseCase) Refresh(
	c context.Context,
	cfg *config.Config,
	refresh string,
) (token, newRefresh string, err error) {
	hashed := utils.HashRefreshToken(refresh)

	userID, err := uc.repo.ReadRefreshToken(c, hashed)
	if err != nil {
		return "", "", fmt.Errorf(
			"UserUseCase - Refresh - uc.repo.ReadRefreshToken: %w",
			err,
		)
	}

	token, err = generateToken(userID, &cfg.JWT)
	if err != nil {
		return "", "", fmt.Errorf(
			"UserUseCase - Refresh - generateToken: %w",
			err,
		)
	}

	newRefresh, newRefreshhashed, err := utils.NewRefreshToken()
	if err != nil {
		return "", "", fmt.Errorf(
			"UserUseCase - Refresh - utils.NewRefreshToken: %w",
			err,
		)
	}

	err = uc.repo.CreateRefreshToken(
		c,
		userID,
		newRefreshhashed,
		"",
		time.Now().Add(cfg.RefreshToken.Expiration),
	)
	if err != nil {
		return "", "", fmt.Errorf(
			"UserUseCase - Refresh - uc.repo.CreateRefreshToken: %w",
			err,
		)
	}

	return token, newRefresh, nil
}

func (uc *UseCase) Logout(
	c context.Context,
	cfg *config.Config,
	userID, refreshTokenID int32,
	refresh string,
) error {
	// TODO: delete refresh token by hashed and id in parallel
	return nil
}

// generateToken is just a local version of jwt.GenerateToken that can take *config.JWT as an argument
// so that we don't have to pass 3 arguments at a time :)
func generateToken(userID int32, cfgJWT *config.JWT) (string, error) {
	return jwt.GenerateToken(
		strconv.Itoa(int(userID)),
		cfgJWT.Secret,
		cfgJWT.Issuer,
		cfgJWT.Expiration,
		jwt.TokenTypeAccess,
	)
}

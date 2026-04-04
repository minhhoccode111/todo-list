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
// so that we don't have to pass 3 arguments every time :)
func generateToken(userID string, cfgJWT *config.JWT, tokenType jwt.TokenType) (string, error) {
	if tokenType == jwt.TokenTypeAccess {
		return jwt.GenerateToken(
			userID,
			cfgJWT.Secret,
			cfgJWT.Issuer,
			cfgJWT.Expiration,
			tokenType,
		)
	}

	return jwt.GenerateToken(
		userID,
		cfgJWT.RefreshSecret,
		cfgJWT.RefreshIssuer,
		cfgJWT.RefreshExpiration,
		tokenType,
	)
}

func hashToken(token string) (string, error) {
	return password.HashPassword(token)
}

func (uc *UseCase) Register(
	c context.Context,
	u *entity.User,
	cfgJWT *config.JWT,
) (string, string, error) {
	hashed, err := password.HashPassword(u.PasswordHash)
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

	userID := strconv.Itoa(int(u.ID))
	userIDInt := u.ID

	accessToken, err := generateToken(userID, cfgJWT, jwt.TokenTypeAccess)
	if err != nil {
		return "", "", fmt.Errorf(
			"UserUseCase - Register - generateToken: %w",
			err,
		)
	}

	refreshToken, err := generateToken(userID, cfgJWT, jwt.TokenTypeRefresh)
	if err != nil {
		return "", "", fmt.Errorf(
			"UserUseCase - Register - generateToken: %w",
			err,
		)
	}

	refreshTokenHash, err := hashToken(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf(
			"UserUseCase - Register - hashToken: %w",
			err,
		)
	}

	expiresAt := time.Now().Add(cfgJWT.RefreshExpiration)
	err = uc.repo.CreateRefreshToken(c, userIDInt, refreshTokenHash, "", expiresAt)
	if err != nil {
		return "", "", fmt.Errorf(
			"UserUseCase - Register - uc.repo.CreateRefreshToken: %w",
			err,
		)
	}

	uc.cache.SetUser(c, userID, u)

	return accessToken, refreshToken, nil
}

func (uc *UseCase) Login(
	c context.Context,
	u *entity.User,
	cfgJWT *config.JWT,
) (string, string, error) {
	user, err := uc.repo.ReadUserByEmail(c, u.Email)
	if err != nil {
		return "", "", fmt.Errorf(
			"UserUseCase - Login - uc.repo.ReadUserByEmail: %w",
			err,
		)
	}

	if !password.ComparePassword(user.PasswordHash, u.PasswordHash) {
		return "", "", entity.ErrUnauthorized
	}

	userID := strconv.Itoa(int(user.ID))
	userIDInt := user.ID

	accessToken, err := generateToken(userID, cfgJWT, jwt.TokenTypeAccess)
	if err != nil {
		return "", "", fmt.Errorf(
			"UserUseCase - Login - generateToken: %w",
			err,
		)
	}

	refreshToken, err := generateToken(userID, cfgJWT, jwt.TokenTypeRefresh)
	if err != nil {
		return "", "", fmt.Errorf(
			"UserUseCase - Login - generateToken: %w",
			err,
		)
	}

	refreshTokenHash, err := hashToken(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf(
			"UserUseCase - Login - hashToken: %w",
			err,
		)
	}

	expiresAt := time.Now().Add(cfgJWT.RefreshExpiration)
	err = uc.repo.CreateRefreshToken(c, userIDInt, refreshTokenHash, "", expiresAt)
	if err != nil {
		return "", "", fmt.Errorf(
			"UserUseCase - Login - uc.repo.CreateRefreshToken: %w",
			err,
		)
	}

	uc.cache.SetUser(c, userID, user)

	return accessToken, refreshToken, nil
}

func (uc *UseCase) RefreshToken(
	c context.Context,
	refreshToken string,
	cfgJWT *config.JWT,
) (string, string, error) {
	refreshTokenHash, err := hashToken(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf(
			"UserUseCase - RefreshToken - hashToken: %w",
			err,
		)
	}

	userID, err := uc.repo.ReadRefreshTokenByHash(c, refreshTokenHash)
	if err != nil {
		return "", "", entity.ErrUnauthorized
	}

	err = uc.repo.RevokeRefreshToken(c, refreshTokenHash)
	if err != nil {
		return "", "", fmt.Errorf(
			"UserUseCase - RefreshToken - uc.repo.RevokeRefreshToken: %w",
			err,
		)
	}

	userIDStr := strconv.Itoa(int(userID))

	accessToken, err := generateToken(userIDStr, cfgJWT, jwt.TokenTypeAccess)
	if err != nil {
		return "", "", fmt.Errorf(
			"UserUseCase - RefreshToken - generateToken: %w",
			err,
		)
	}

	newRefreshToken, err := generateToken(userIDStr, cfgJWT, jwt.TokenTypeRefresh)
	if err != nil {
		return "", "", fmt.Errorf(
			"UserUseCase - RefreshToken - generateToken: %w",
			err,
		)
	}

	newRefreshTokenHash, err := hashToken(newRefreshToken)
	if err != nil {
		return "", "", fmt.Errorf(
			"UserUseCase - RefreshToken - hashToken: %w",
			err,
		)
	}

	expiresAt := time.Now().Add(cfgJWT.RefreshExpiration)
	err = uc.repo.CreateRefreshToken(c, userID, newRefreshTokenHash, "", expiresAt)
	if err != nil {
		return "", "", fmt.Errorf(
			"UserUseCase - RefreshToken - uc.repo.CreateRefreshToken: %w",
			err,
		)
	}

	return accessToken, newRefreshToken, nil
}

func (uc *UseCase) Logout(c context.Context, userID int32, refreshToken string) error {
	if refreshToken != "" {
		refreshTokenHash, err := hashToken(refreshToken)
		if err != nil {
			return fmt.Errorf(
				"UserUseCase - Logout - hashToken: %w",
				err,
			)
		}

		err = uc.repo.RevokeRefreshToken(c, refreshTokenHash)
		if err != nil {
			return fmt.Errorf(
				"UserUseCase - Logout - uc.repo.RevokeRefreshToken: %w",
				err,
			)
		}
	}

	err := uc.repo.RevokeAllUserRefreshTokens(c, userID)
	if err != nil {
		return fmt.Errorf(
			"UserUseCase - Logout - uc.repo.RevokeAllUserRefreshTokens: %w",
			err,
		)
	}

	return nil
}

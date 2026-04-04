package persistent

//
// import (
// 	"context"
// 	"fmt"
// 	"time"
//
// 	"github.com/jackc/pgx/v5/pgtype"
// 	"github.com/minhhoccode111/todo-list/internal/repo/persistent/sqlc"
// 	"github.com/minhhoccode111/todo-list/pkg/postgres"
// )
//
// type RefreshTokenRepo struct {
// 	*postgres.Postgres
// 	queries *sqlc.Queries
// }
//
// func NewRefreshTokenRepo(pg *postgres.Postgres) *RefreshTokenRepo {
// 	return &RefreshTokenRepo{
// 		Postgres: pg,
// 		queries:  sqlc.New(pg.Pool),
// 	}
// }
//
// func (rtr *RefreshTokenRepo) CreateRefreshToken(c context.Context, userID int32, tokenHash, deviceInfo string, expiresAt time.Time) error {
// 	_, err := rtr.queries.CreateRefreshToken(c, sqlc.CreateRefreshTokenParams{
// 		UserID:     userID,
// 		TokenHash:  tokenHash,
// 		ExpiresAt:  expiresAt,
// 		DeviceInfo: pgtype.Text{String: deviceInfo, Valid: deviceInfo != ""},
// 	})
// 	if err != nil {
// 		return fmt.Errorf("RefreshTokenRepo - CreateRefreshToken - rtr.queries.CreateRefreshToken: %w", err)
// 	}
//
// 	return nil
// }
//
// func (rtr *RefreshTokenRepo) GetRefreshTokenByHash(c context.Context, tokenHash string) (int32, error) {
// 	token, err := rtr.queries.GetRefreshTokenByHash(c, tokenHash)
// 	if err != nil {
// 		return 0, fmt.Errorf("RefreshTokenRepo - GetRefreshTokenByHash - rtr.queries.GetRefreshTokenByHash: %w", err)
// 	}
//
// 	return token.UserID, nil
// }
//
// func (rtr *RefreshTokenRepo) RevokeRefreshToken(c context.Context, tokenHash string) error {
// 	err := rtr.queries.RevokeRefreshToken(c, tokenHash)
// 	if err != nil {
// 		return fmt.Errorf("RefreshTokenRepo - RevokeRefreshToken - rtr.queries.RevokeRefreshToken: %w", err)
// 	}
//
// 	return nil
// }
//
// func (rtr *RefreshTokenRepo) RevokeAllUserRefreshTokens(c context.Context, userID int32) error {
// 	err := rtr.queries.RevokeAllUserRefreshTokens(c, userID)
// 	if err != nil {
// 		return fmt.Errorf("RefreshTokenRepo - RevokeAllUserRefreshTokens - rtr.queries.RevokeAllUserRefreshTokens: %w", err)
// 	}
//
// 	return nil
// }
//
// func (rtr *RefreshTokenRepo) DeleteExpiredRefreshTokens(c context.Context) error {
// 	err := rtr.queries.DeleteExpiredRefreshTokens(c)
// 	if err != nil {
// 		return fmt.Errorf("RefreshTokenRepo - DeleteExpiredRefreshTokens - rtr.queries.DeleteExpiredRefreshTokens: %w", err)
// 	}
//
// 	return nil
// }

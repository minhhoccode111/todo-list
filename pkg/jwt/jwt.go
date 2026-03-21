package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ClaimsJWT struct {
	UserID string
	jwt.RegisteredClaims
}

func GenerateToken(userID, secret, issuer string, expiration time.Duration) (string, error) {
	claims := ClaimsJWT{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func ValidateToken(tokenStr, secret string) (*ClaimsJWT, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&ClaimsJWT{},
		func(token *jwt.Token) (any, error) {
			// enforce algorithm
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenSignatureInvalid
			}

			return secret, nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*ClaimsJWT)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

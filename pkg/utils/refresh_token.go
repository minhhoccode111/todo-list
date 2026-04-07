package utils //nolint:revive,nolintlint // why not utils?

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

func HashRefreshToken(raw string) (hashed string) {
	hash := sha256.Sum256([]byte(raw))
	hashed = hex.EncodeToString(hash[:])

	return hashed
}

func NewRefreshToken() (raw, hashed string, err error) {
	const n = 32 // 256-bit

	b := make([]byte, n)
	if _, err = rand.Read(b); err != nil {
		return "", "", err
	}

	raw = base64.RawURLEncoding.EncodeToString(b)
	hashed = HashRefreshToken(raw)

	return raw, hashed, nil
}

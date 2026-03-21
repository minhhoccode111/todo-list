package entity

import (
	"errors"
)

var (
	ErrNoRows         = errors.New("no rows")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrDuplicateEntry = errors.New("duplicate entry")
)

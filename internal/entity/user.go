package entity

import (
	"time"
)

type User struct {
	ID           int32
	Email        string
	Name         string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

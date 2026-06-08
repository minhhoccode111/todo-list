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

type Session struct {
	ID         int32
	TokenHash  string
	DeviceInfo string
	CreatedAt  time.Time
	ExpiresAt  time.Time
	IsCurrent  bool
}

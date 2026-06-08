package response

import "time"

type Auth struct {
	Token string `json:"token" validate:"required"`
}

/*
gomodifytags -file internal/controller/restapi/v1/response/user.go -struct Session -add-tags json,validate:required -w
*/

type Session struct {
	ID        int32     `json:"id"         validate:"required"`
	Device    string    `json:"device"     validate:"required"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
	ExpiresAt time.Time `json:"expires_at" validate:"required"`
	IsCurrent bool      `json:"is_current" validate:"required"`
}

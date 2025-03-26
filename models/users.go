package model

import (
	"time"

	"github.com/google/uuid"
)

type Users struct {
	Id                  uuid.UUID  `json:"id"`
	Email               string     `json:"email" validate:"required,email"`
	Username            string     `json:"username" validate:"required"`
	Password            string     `json:"password" validate:"required"`
	Role                string     `json:"role"`
	Address             string     `json:"address"`
	CategoryPreferences []string   `json:"category_preferences"`
	CreatedAt           *time.Time `json:"created_at"`
	UpdatedAt           *time.Time `json:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at"`
}

type UserLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserLogin struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
	RefreshToken         string    `json:"refresh_token"`
	RefreshTokenExpiryAt time.Time `json:"refresh_token_expiry_at"`
	*Users
}

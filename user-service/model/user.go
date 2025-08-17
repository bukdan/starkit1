package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`            // user_id
	Username     string    `json:"username"`      // nama user
	Email        string    `json:"email"`         // email
	Phone        string    `json:"phone"`         // nomor HP
	PasswordHash string    `json:"password_hash"` // hash password
	Role         string    `json:"role"`          // user/admin
	IsActive     bool      `json:"is_active"`     // aktif/nonaktif
	IsVerified   bool      `json:"is_verified"`   // untuk OTP/email verified
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

package model

import "time"

type User struct {
	ID         string    `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Phone      *string   `json:"phone,omitempty"`
	AvatarURL  *string   `json:"avatar_url,omitempty"`
	Role       string    `json:"role"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type OTPCode struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Channel   string    `json:"channel"` // "email" or "wa"
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expires_at"`
	Used      bool      `json:"used"`
	CreatedAt time.Time `json:"created_at"`
}

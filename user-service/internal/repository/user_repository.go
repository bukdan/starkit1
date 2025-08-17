package repository

import (
	"context"
	"time"

	"user-service/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser inserts a new user, returns id
func (r *UserRepository) CreateUser(ctx context.Context, username, email, phone, passwordHash string) (string, error) {
	var id string
	row := r.db.QueryRow(ctx, `
		INSERT INTO users (username, email, phone, password_hash)
		VALUES ($1, $2, $3, $4) RETURNING id
	`, username, email, nullEmpty(phone), passwordHash)
	if err := row.Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (model.User, string, error) {
	var u model.User
	var phone, avatar *string
	var pass string
	row := r.db.QueryRow(ctx, `
		SELECT id, username, email, phone, avatar_url, role, is_verified, created_at, updated_at, password_hash
		FROM users WHERE email = $1
	`, email)
	if err := row.Scan(&u.ID, &u.Username, &u.Email, &phone, &avatar, &u.Role, &u.IsVerified, &u.CreatedAt, &u.UpdatedAt, &pass); err != nil {
		return model.User{}, "", err
	}
	u.Phone = phone
	u.AvatarURL = avatar
	return u, pass, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (model.User, error) {
	var u model.User
	var phone, avatar *string
	row := r.db.QueryRow(ctx, `
		SELECT id, username, email, phone, avatar_url, role, is_verified, created_at, updated_at
		FROM users WHERE id = $1
	`, id)
	if err := row.Scan(&u.ID, &u.Username, &u.Email, &phone, &avatar, &u.Role, &u.IsVerified, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return model.User{}, err
	}
	u.Phone = phone
	u.AvatarURL = avatar
	return u, nil
}

func (r *UserRepository) UpdateProfile(ctx context.Context, id, username, phone, avatar string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE users SET username = $1, phone = $2, avatar_url = $3, updated_at = NOW() WHERE id = $4
	`, username, nullEmpty(phone), nullEmpty(avatar), id)
	return err
}

func (r *UserRepository) SetVerified(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE users SET is_verified = true, updated_at = NOW() WHERE id = $1
	`, id)
	return err
}

func (r *UserRepository) ListUsers(ctx context.Context, limit, offset int) ([]model.User, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, username, email, phone, avatar_url, role, is_verified, created_at, updated_at
		FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]model.User, 0)
	for rows.Next() {
		var u model.User
		var phone, avatar *string
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &phone, &avatar, &u.Role, &u.IsVerified, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		u.Phone = phone
		u.AvatarURL = avatar
		out = append(out, u)
	}
	return out, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	return err
}

func (r *UserRepository) UpdateUserAdmin(ctx context.Context, id, username, email, phone, role string, isVerified bool) error {
	_, err := r.db.Exec(ctx, `
		UPDATE users SET username=$1, email=$2, phone=$3, role=$4, is_verified=$5, updated_at=NOW() WHERE id=$6
	`, username, email, nullEmpty(phone), role, isVerified, id)
	return err
}

// OTP related
func (r *UserRepository) CreateOTP(ctx context.Context, userID, channel, code string, expiresAt time.Time) (string, error) {
	var id string
	err := r.db.QueryRow(ctx, `
		INSERT INTO otp_codes (user_id, channel, code, expires_at) VALUES ($1,$2,$3,$4) RETURNING id
	`, userID, channel, code, expiresAt).Scan(&id)
	return id, err
}

func (r *UserRepository) GetValidOTP(ctx context.Context, userID, channel, code string) (model.OTPCode, error) {
	var o model.OTPCode
	row := r.db.QueryRow(ctx, `
		SELECT id, user_id, channel, code, expires_at, used, created_at
		FROM otp_codes
		WHERE user_id=$1 AND channel=$2 AND code=$3 AND used=false AND expires_at > NOW()
		ORDER BY created_at DESC LIMIT 1
	`, userID, channel, code)
	if err := row.Scan(&o.ID, &o.UserID, &o.Channel, &o.Code, &o.ExpiresAt, &o.Used, &o.CreatedAt); err != nil {
		return model.OTPCode{}, err
	}
	return o, nil
}

func (r *UserRepository) MarkOTPUsed(ctx context.Context, otpID string) error {
	_, err := r.db.Exec(ctx, `UPDATE otp_codes SET used=true WHERE id=$1`, otpID)
	return err
}

// helper to convert empty string to NULL
func nullEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

package repository

import (
	"database/sql"
	"time"
)

type OTP struct {
	ID        string
	UserID    string
	Code      string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type OTPRepository interface {
	SaveOTP(otp *OTP) error
	VerifyOTP(userID, code string) (bool, error)
	DeleteOTP(userID string) error
}

type otpRepository struct {
	db *sql.DB
}

func NewOTPRepository(db *sql.DB) OTPRepository {
	return &otpRepository{db}
}

func (r *otpRepository) SaveOTP(otp *OTP) error {
	query := `INSERT INTO otps (id, user_id, code, expires_at, created_at)
	          VALUES ($1,$2,$3,$4,NOW())`
	_, err := r.db.Exec(query, otp.ID, otp.UserID, otp.Code, otp.ExpiresAt)
	return err
}

func (r *otpRepository) VerifyOTP(userID, code string) (bool, error) {
	var dbCode string
	var expiresAt time.Time
	query := `SELECT code, expires_at FROM otps WHERE user_id=$1 ORDER BY created_at DESC LIMIT 1`
	err := r.db.QueryRow(query, userID).Scan(&dbCode, &expiresAt)
	if err != nil {
		return false, err
	}
	if dbCode == code && time.Now().Before(expiresAt) {
		return true, nil
	}
	return false, nil
}

func (r *otpRepository) DeleteOTP(userID string) error {
	_, err := r.db.Exec(`DELETE FROM otps WHERE user_id=$1`, userID)
	return err
}

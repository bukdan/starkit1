package service

import (
	"context"
	"time"

	"user-service/internal/repository"
	"user-service/internal/utils"
)

type OTPService struct {
	repo *repository.UserRepository
}

func NewOTPService(repo *repository.UserRepository) *OTPService {
	return &OTPService{repo: repo}
}

func (s *OTPService) GenerateCreateAndSend(ctx context.Context, userID, channel, contact string, expiryMinutes int) (string, error) {
	code := utils.GenerateOTP6()
	expires := time.Now().Add(time.Duration(expiryMinutes) * time.Minute)

	id, err := s.repo.CreateOTP(ctx, userID, channel, code, expires)
	if err != nil {
		return "", err
	}

	// send async
	go func() {
		if channel == "wa" {
			_ = utils.SendWhatsAppOTP(contact, code)
		} else {
			_ = utils.SendEmailOTP(contact, code)
		}
	}()

	return id, nil
}

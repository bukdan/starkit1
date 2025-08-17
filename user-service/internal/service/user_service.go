package service

import (
	"context"
	"errors"

	"user-service/internal/model"
	"user-service/internal/repository"
	"user-service/internal/utils"
)

type UserService struct {
	repo       *repository.UserRepository
	otpService *OTPService
	otpExpiry  int // minutes
}

func NewUserService(repo *repository.UserRepository, otpSvc *OTPService, otpExpiry int) *UserService {
	return &UserService{repo: repo, otpService: otpSvc, otpExpiry: otpExpiry}
}

func (s *UserService) Register(ctx context.Context, username, email, phone, password, sendVia string) (string, error) {
	hashed, err := utils.HashPassword(password)
	if err != nil {
		return "", err
	}
	uid, err := s.repo.CreateUser(ctx, username, email, phone, hashed)
	if err != nil {
		return "", err
	}

	contact := email
	if sendVia == "wa" && phone != "" {
		contact = phone
	}
	_, _ = s.otpService.GenerateCreateAndSend(ctx, uid, sendVia, contact, s.otpExpiry)
	return uid, nil
}

func (s *UserService) VerifyOTP(ctx context.Context, userID, channel, code string) error {
	otp, err := s.repo.GetValidOTP(ctx, userID, channel, code)
	if err != nil {
		return errors.New("invalid or expired otp")
	}
	if err := s.repo.MarkOTPUsed(ctx, otp.ID); err != nil {
		return err
	}
	return s.repo.SetVerified(ctx, otp.UserID)
}

func (s *UserService) Login(ctx context.Context, email, password string) (model.User, error) {
	u, passHash, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return model.User{}, err
	}
	if !utils.ComparePassword(passHash, password) {
		return model.User{}, errors.New("invalid credentials")
	}
	return u, nil
}

func (s *UserService) GetProfile(ctx context.Context, id string) (model.User, error) {
	return s.repo.GetByID(ctx, id)
}

// admin helpers omitted but can call repo methods directly

package service

import (
	"errors"
	"user-service/model"
	"user-service/repository"
	"user-service/utils"
)

type AuthService interface {
	Register(name, email, password, phone string) (*model.User, error)
	Login(email, password string) (string, error)
	LoginWithGoogle(googleID, email, name string) (string, error)
	VerifyOTP(email, otpCode string) error
}

type authService struct {
	repo       repository.UserRepository
	otpService OTPService
}

func NewAuthService(repo repository.UserRepository, otp OTPService) AuthService {
	return &authService{repo, otp}
}

// Register new user
func (s *authService) Register(name, email, password, phone string) (*model.User, error) {
	// Hash password
	hashed, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:       name,
		Email:      email,
		Password:   hashed,
		Phone:      phone,
		IsVerified: false,
		Role:       "user",
	}

	err = s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	// Kirim OTP via email/WA
	otpCode := utils.GenerateOTP()
	_ = s.otpService.SendOTP(user.Email, otpCode)

	return user, nil
}

// Login dengan email & password
func (s *authService) Login(email, password string) (string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid email or password")
	}

	// Generate JWT
	token, err := utils.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Login dengan Google ID
func (s *authService) LoginWithGoogle(googleID, email, name string) (string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		// Jika user belum ada → buat user baru
		user = &model.User{
			Name:       name,
			Email:      email,
			GoogleID:   googleID,
			IsVerified: true,
			Role:       "user",
		}
		if err := s.repo.Create(user); err != nil {
			return "", err
		}
	}

	// Generate JWT
	return utils.GenerateJWT(user.ID, user.Email, user.Role)
}

// Verifikasi OTP
func (s *authService) VerifyOTP(email, otpCode string) error {
	ok := s.otpService.VerifyOTP(email, otpCode)
	if !ok {
		return errors.New("invalid otp")
	}

	// Update user jadi verified
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return err
	}
	user.IsVerified = true
	return s.repo.Update(user)
}

package service

import (
	"fmt"
)

type OTPService interface {
	SendOTP(email, otp string) error
	VerifyOTP(email, otp string) bool
}

type otpService struct {
	storage map[string]string // simpan OTP sementara (nanti bisa pakai Redis)
}

func NewOTPService() OTPService {
	return &otpService{
		storage: make(map[string]string),
	}
}

func (o *otpService) SendOTP(email, otp string) error {
	o.storage[email] = otp
	// TODO: Integrasi dengan email/WhatsApp
	fmt.Printf("OTP untuk %s: %s\n", email, otp)
	return nil
}

func (o *otpService) VerifyOTP(email, otp string) bool {
	if val, ok := o.storage[email]; ok && val == otp {
		delete(o.storage, email) // hapus setelah dipakai
		return true
	}
	return false
}

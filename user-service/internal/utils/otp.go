package utils

import (
	"crypto/rand"
	"fmt"
	"net/smtp"
	"os"
)

// generateOTP6 returns a 6-digit numeric OTP as string
func GenerateOTP6() string {
	b := make([]byte, 3)
	_, _ = rand.Read(b)
	v := int(b[0])<<16 | int(b[1])<<8 | int(b[2])
	return fmt.Sprintf("%06d", v%1000000)
}

// A minimal email send using net/smtp (for dev). Replace with provider in production.
func SendEmailOTP(toEmail, code string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	if smtpHost == "" || smtpUser == "" {
		// not configured, no-op
		return nil
	}
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	subject := "Your verification code"
	body := "Kode verifikasi Anda: " + code
	msg := "From: " + smtpUser + "\n" +
		"To: " + toEmail + "\n" +
		"Subject: " + subject + "\n\n" + body
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, []string{toEmail}, []byte(msg))
}

// Placeholder WhatsApp sender — integrate Twilio or WhatsApp Cloud API in production
func SendWhatsAppOTP(phone, code string) error {
	_ = phone
	_ = code
	return nil
}

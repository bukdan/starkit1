package utils

import "fmt"

// Dummy OTP sender
func SendOTP(email string) {
	// real implementation bisa pakai Twilio/SMTP
	fmt.Println("📩 Sending OTP to:", email)
}

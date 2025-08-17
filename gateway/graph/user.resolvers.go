package graph

import (
	"context"
	"gateway/internal"
)

func (r *mutationResolver) Register(ctx context.Context, name string, email string, password string, phone *string) (string, error) {
	payload := map[string]interface{}{
		"name":     name,
		"email":    email,
		"password": password,
	}
	if phone != nil {
		payload["phone"] = *phone
	}

	resp, err := internal.Post("/auth/register", payload)
	if err != nil {
		return "", err
	}
	return resp["message"].(string), nil
}

func (r *mutationResolver) Login(ctx context.Context, email string, password string) (*AuthPayload, error) {
	payload := map[string]interface{}{
		"email":    email,
		"password": password,
	}
	resp, err := internal.Post("/auth/login", payload)
	if err != nil {
		return nil, err
	}
	return &AuthPayload{Token: resp["token"].(string)}, nil
}

func (r *mutationResolver) VerifyOtp(ctx context.Context, email string, code string) (string, error) {
	payload := map[string]interface{}{
		"email": email,
		"code":  code,
	}
	resp, err := internal.Post("/auth/verify-otp", payload)
	if err != nil {
		return "", err
	}
	return resp["message"].(string), nil
}

func (r *mutationResolver) LoginWithGoogle(ctx context.Context, googleID string, email string, name string) (*AuthPayload, error) {
	payload := map[string]interface{}{
		"google_id": googleID,
		"email":     email,
		"name":      name,
	}
	resp, err := internal.Post("/auth/google", payload)
	if err != nil {
		return nil, err
	}
	return &AuthPayload{Token: resp["token"].(string)}, nil
}

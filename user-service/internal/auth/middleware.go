package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"time"
)

type GoogleProfile struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Aud   string `json:"aud"`
}

// Dev helper: verify Google id_token via tokeninfo endpoint.
// Production: verify JWT signature via Google's JWKS.
func VerifyGoogleIDToken(ctx context.Context, idToken string) (*GoogleProfile, error) {
	endpoint := "https://oauth2.googleapis.com/tokeninfo"
	q := url.Values{}
	q.Set("id_token", idToken)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, endpoint+"?"+q.Encode(), nil)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("tokeninfo rejected")
	}
	var prof GoogleProfile
	if err := json.NewDecoder(resp.Body).Decode(&prof); err != nil {
		return nil, err
	}
	// optional aud check
	if cid := os.Getenv("GOOGLE_CLIENT_ID"); cid != "" {
		if prof.Aud != "" && prof.Aud != cid {
			return nil, errors.New("audience mismatch")
		}
	}
	return &prof, nil
}

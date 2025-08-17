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

// GoogleProfile minimal fields
type GoogleProfile struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Aud   string `json:"aud"`
}

// VerifyIDTokenDev: dev-friendly verification using tokeninfo endpoint.
// Production: use JWKS and verify token signature/claims properly.
func VerifyGoogleIDToken(ctx context.Context, idToken string) (*GoogleProfile, error) {
	endpoint := "https://oauth2.googleapis.com/tokeninfo"
	q := url.Values{}
	q.Set("id_token", idToken)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, endpoint+"?"+q.Encode(), nil)
	client := &http.Client{Timeout: timeSecond()}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("tokeninfo rejected")
	}
	var p GoogleProfile
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return nil, err
	}
	if cid := os.Getenv("GOOGLE_CLIENT_ID"); cid != "" && p.Aud != "" && p.Aud != cid {
		return nil, errors.New("audience mismatch")
	}
	return &p, nil
}

func timeSecond() time.Duration { return 5 * time.Second }

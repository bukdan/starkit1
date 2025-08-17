package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

func userServiceBase() string {
	if u := os.Getenv("USER_SERVICE_URL"); u != "" {
		return u
	}
	return "http://localhost:8081"
}

// postJSON posts JSON to user-service path and returns decoded body (map) and status code.
// headers parameter allows forwarding Authorization etc.
func postJSON(path string, body any, headers map[string]string) (map[string]any, int, error) {
	url := fmt.Sprintf("%s%s", userServiceBase(), path)
	b, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	var out map[string]any
	if len(data) == 0 {
		return map[string]any{"_empty": true}, resp.StatusCode, nil
	}
	if err := json.Unmarshal(data, &out); err != nil {
		// return raw text inside map
		return map[string]any{"_raw": string(data)}, resp.StatusCode, nil
	}
	return out, resp.StatusCode, nil
}

// getJSON supports GET with forwarded Authorization (for /users/me if it's a GET)
func getJSON(path string, headers map[string]string) (map[string]any, int, error) {
	url := fmt.Sprintf("%s%s", userServiceBase(), path)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	var out map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, resp.StatusCode, err
	}
	return out, resp.StatusCode, nil
}

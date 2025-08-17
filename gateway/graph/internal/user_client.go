package internal

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

var userServiceURL = os.Getenv("USER_SERVICE_URL")

func Post(path string, payload interface{}) (map[string]interface{}, error) {
	body, _ := json.Marshal(payload)
	resp, err := http.Post(userServiceURL+path, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

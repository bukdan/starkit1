package handler

import (
	"encoding/json"
	"net/http"
	"user-service/model"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	// contoh dummy user list
	users := []model.User{
		{ID: "1", Name: "Budi", Email: "budi@example.com"},
	}
	json.NewEncoder(w).Encode(users)
}

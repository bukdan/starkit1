package handler

import (
	
    "encoding/json"
    "net/http"
    "user-service/model"
    "net/http"
    "github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "users": []string{"John Doe", "Jane Smith"},
    })
}
package handler

import (
	"net/http"
	"user-service/model"
	"user-service/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService}
}

// Get profile user berdasarkan JWT (user_id)
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id") // dari middleware JWT

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Update profile user
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id") // dari middleware JWT

	var req model.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.ID = userID

	if err := h.userService.UpdateProfile(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

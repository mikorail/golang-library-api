package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"products-api-with-jwt/models"
	"products-api-with-jwt/services"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	AuthService *services.AuthService
}

// NewAuthController menginisialisasi AuthController baru
func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{AuthService: authService}
}

func (ac *AuthController) Login(c *gin.Context) {
	var input models.LoginInput
	// Step 1: Validate input fields
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Status:  "error",
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	// Validate credentials
	user, err := ac.AuthService.ValidateCredentials(input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ApiResponse{
			Status:  "error",
			Code:    http.StatusUnauthorized,
			Message: "Invalid username or password",
			Data:    nil,
		})
		return
	}

	// Step 2: Generate a new token
	expiration := time.Hour * 24 // Default to 1 day
	if input.RememberMe {
		expiration = time.Hour * 24 * 7 // 7 days
	}

	token, err := ac.AuthService.GenerateToken(uint(user.ID), user.Username, expiration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status:  "error",
			Code:    http.StatusInternalServerError,
			Message: "Could not generate token",
			Data:    nil,
		})
		return
	}

	// Step 3: Check logging history
	var recentLog models.LoggingHistory
	if err := ac.AuthService.GetRecentLoggingHistory(uint64(user.ID), &recentLog); err != nil {
		// If there's an error retrieving the logging history, just continue
		log.Printf("Could not check logging history: %v", err)
	}

	// Update user's active status
	if err := ac.AuthService.UpdateUserActiveStatus(uint(user.ID), true); err != nil {
		log.Printf("Could not update user active status: %v", err)
	}

	// If there's a recent logging history and the token is still valid
	if recentLog.ID != 0 && recentLog.ExpiredDate.After(time.Now()) {
		// Token is still valid; no need to create a new logging history
		c.JSON(http.StatusOK, models.ApiResponse{
			Status:  "success",
			Code:    http.StatusOK,
			Message: "Login successful with existing valid token.",
			Data:    gin.H{"token": token},
		})
		return
	}

	// If no logging history was found or if the previous token is expired, create a new logging history entry
	newLog := models.LoggingHistory{
		UserID:      uint(user.ID),
		JWT:         token,
		ExpiredDate: time.Now().Add(expiration),
		CreatedDate: time.Now(),
	}
	if err := ac.AuthService.CreateLoggingHistory(&newLog); err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status:  "error",
			Code:    http.StatusInternalServerError,
			Message: "Could not create logging history",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "Login successful",
		Data:    gin.H{"token": token},
	})
}

func (ac *AuthController) Logout(c *gin.Context) {
	// Dapatkan token dari header Authorization
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Status:  "error",
			Code:    http.StatusBadRequest,
			Message: "Authorization token is required",
			Data:    nil,
		})
		return
	}

	// Get the user ID from the token
	userID, err := ac.AuthService.GetUserIDFromToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ApiResponse{
			Status:  "error",
			Code:    http.StatusUnauthorized,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	// Update user's active status
	if err := ac.AuthService.UpdateUserActiveStatus(uint(userID), false); err != nil {
		log.Printf("Could not update user active status: %v", err)
		// You can choose whether to handle this error further
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Status:  "success",
		Code:    http.StatusOK,
		Message: fmt.Sprintf("Logout successful for user : %v", userID),
		Data:    nil,
	})
}

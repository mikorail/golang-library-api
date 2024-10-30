package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"products-api-with-jwt/global"
	"products-api-with-jwt/models"
	"products-api-with-jwt/services"

	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// GetJWTKey retrieves the JWT secret key from the environment variable
func GetJWTKey() []byte {
	key := os.Getenv(global.ENVSecretKey)
	if key == "" {
		panic("JWT secret key not set in environment variables")
	}
	return []byte(key)
}

// GenerateToken creates a JWT token with a dynamic expiration time (1 day or 7 days)
func GenerateToken(username string, expiration time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
		Issuer:    "your-app",
		Subject:   username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(GetJWTKey())
}

// JWTAuthMiddleware validates the JWT token in the Authorization header for each request
func JWTAuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.ApiResponse{
				Status:  "error",
				Code:    http.StatusBadRequest,
				Message: "Authorization header required",
			})
			c.Abort()
			return
		}

		// Extract token from Authorization header
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, models.ApiResponse{
				Status:  "error",
				Code:    http.StatusBadRequest,
				Message: "Bearer token required",
			})
			c.Abort()
			return
		}

		// Validate JWT token
		token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Verify signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return GetJWTKey(), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, models.ApiResponse{
				Status:  "error",
				Code:    http.StatusBadRequest,
				Message: "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Extract claims to retrieve user information if needed
		claims, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, models.ApiResponse{
				Status:  "error",
				Code:    http.StatusBadRequest,
				Message: "Invalid token claims",
			})
			c.Abort()
			return
		}

		// Check if user is active
		idCheck, err := authService.GetUserIDFromToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ApiResponse{
				Status:  "error",
				Code:    http.StatusBadRequest,
				Message: "User ID not found",
			})
			c.Abort()
			return
		}

		user, err := authService.GetUserById(int(idCheck))
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ApiResponse{
				Status:  "error",
				Code:    http.StatusBadRequest,
				Message: "User not found",
			})
			c.Abort()
			return
		}

		// Check if user is active
		if !user.Active {
			c.JSON(http.StatusUnauthorized, models.ApiResponse{
				Status:  "error",
				Code:    http.StatusBadRequest,
				Message: "You are not logged in",
			})
			c.Abort()
			return
		}

		// Store user ID in context
		c.Set("user_id", claims.Subject)
		c.Next()
	}
}

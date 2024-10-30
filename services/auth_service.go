package services

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"products-api-with-jwt/models"

	global "products-api-with-jwt/global"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AuthService struct {
	DB        *gorm.DB
	secretKey string
}

// NewAuthService menginisialisasi AuthService baru
func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{
		DB:        db,
		secretKey: os.Getenv(global.ENVSecretKey), // Ganti dengan secret key yang aman
	}
}

// GetUserByUsername retrieves a user by their username
func (s *AuthService) GetUserById(id int) (*models.User, error) {
	var user models.User
	if err := s.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err // Return error if user not found
	}
	return &user, nil
}

// ValidateCredentials memvalidasi username dan password
func (s *AuthService) ValidateCredentials(username, password string) (models.User, error) {
	var user models.User
	// Cari pengguna berdasarkan username
	if err := s.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return user, errors.New("invalid username or password")
	}

	// Verifikasi password (gunakan bcrypt atau metode hashing yang sama yang digunakan saat menyimpan)
	if err := models.CheckPasswordHash(password, user.Password); err != nil {
		return user, errors.New("invalid username or password")
	}

	return user, nil
}

// GenerateToken menghasilkan token JWT untuk pengguna
func (s *AuthService) GenerateToken(id uint, username string, expiration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  id,
		"username": username,
		"exp":      time.Now().Add(expiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *AuthService) GetRecentLoggingHistory(userID uint64, recentLog *models.LoggingHistory) error {
	return s.DB.Where("user_id = ?", userID).Order("expired_date desc").First(recentLog).Error
}

func (s *AuthService) CreateLoggingHistory(log *models.LoggingHistory) error {
	return s.DB.Create(log).Error
}

func (s *AuthService) UpdateUserActiveStatus(userID uint, active bool) error {
	result := s.DB.Model(&models.User{}).Where("id = ?", userID).Update("active", active)
	if result.Error != nil {
		return result.Error
	}

	// Check if the update affected any rows (e.g., user not found)
	if result.RowsAffected == 0 {
		return fmt.Errorf("no user found with ID %d", userID)
	}

	return nil
}

func (as *AuthService) GetUserIDFromToken(authHeader string) (uint, error) {
	// Parse the token
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	fmt.Println(tokenString)
	// Call the function to get user ID from the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the token's signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key used for signing
		return []byte(as.secretKey), nil // Replace with your actual secret key
	})

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	// Extract claims from the token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("Token is valid")

		// Assuming user ID is stored as "user_id" in the claims
		if id, ok := claims["user_id"].(float64); ok {
			return uint(id), nil
		}
	}

	return 0, fmt.Errorf("user_id claim not found in token")
}

package models

import "golang.org/x/crypto/bcrypt"

type LoginInput struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RememberMe bool   `json:"remember_me"`
}

// CheckPasswordHash memeriksa apakah password yang diberikan cocok dengan hash
func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

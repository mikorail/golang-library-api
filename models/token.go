package models

// Token represents a JWT token associated with a user
type Token struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	JWT      string `json:"jwt"`
}

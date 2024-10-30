package models

import "time"

type User struct {
	ID           int    `gorm:"primaryKey"`
	Username     string `gorm:"unique;not null"`
	Password     string `gorm:"not null"`
	BookBorrowed int
	BorrowDate   *time.Time
	Active       bool //active means active login
}

package models

import "time"

type Book struct {
	ID          int    `gorm:"primaryKey"`
	Title       string `gorm:"not null"`
	Description string
	Author      string
	Stock       int
	Borrowed    int
	CreatedAt   *time.Time
	Active      bool
}

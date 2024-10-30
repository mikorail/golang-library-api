package models

import "time"

// LoggingHistory represents a record of JWT expiration events.
type LoggingHistory struct {
	ID          uint `gorm:"primaryKey"`
	UserID      uint
	JWT         string
	ExpiredDate time.Time `gorm:"not null"`
	CreatedDate time.Time `gorm:"not null"`
}

func (LoggingHistory) TableName() string {
	return "logging_histories"
}

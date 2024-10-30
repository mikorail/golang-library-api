package config

import (
	"products-api-with-jwt/global"
	"products-api-with-jwt/models"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupDatabase initializes the PostgreSQL database
func SetupDatabase() (*gorm.DB, error) {
	global.LoadEnv() // Load environment variables

	dsn := global.GetDBConfig() // Get the database connection string
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return db, err
	}

	// Migrate tables for User and Product models
	db.AutoMigrate(&models.User{}, &models.Book{}, &models.LoggingHistory{})

	// Populate initial data
	populateInitialData(db)

	return db, nil
}

func populateInitialData(db *gorm.DB) {
	now := time.Now()

	// Check if initial data exists
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count == 0 {
		// Hash password for example users
		passwordHash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

		// Add example users
		users := []models.User{
			{Username: "admin", Password: string(passwordHash), BookBorrowed: 1, BorrowDate: &now, Active: false},
			{Username: "user1", Password: string(passwordHash), BookBorrowed: 2, BorrowDate: &now, Active: false},
			{Username: "user2", Password: string(passwordHash), Active: false},
		}
		db.Create(&users)
	}

	books := []models.Book{
		{Title: "Buku A", Description: "Deskripsi Buku A", Author: "Penulis A", Stock: 10, Borrowed: 0, CreatedAt: &now, Active: true},
		{Title: "Buku B", Description: "Deskripsi Buku B", Author: "Penulis B", Stock: 15, Borrowed: 0, CreatedAt: &now, Active: true},
		{Title: "Buku C", Description: "Deskripsi Buku C", Author: "Penulis B", Stock: 20, Borrowed: 0, CreatedAt: &now, Active: true},
	}
	db.Create(&books)

	db.Model(&models.LoggingHistory{}).Count(&count)
	if count == 0 {
		// Example data: Populate with a couple of sample records
		loggingHistoryEntries := []models.LoggingHistory{
			{UserID: 1, JWT: "aaa", ExpiredDate: time.Now().Add(time.Hour * 24), CreatedDate: time.Now().Add(time.Hour * 24)},
		}
		db.Create(&loggingHistoryEntries)
	}
}

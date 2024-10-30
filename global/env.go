package global

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const ENVSecretKey string = "SECRET_KEY"
const ENVRateLimitDur string = "RATE_LIMIT_DURATION"
const ENVRateLimitTime string = "RATE_LIMIT_TIME"

// LoadEnv loads environment variables from a .env file if it exists
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, continuing with environment variables")
	}
}

// GetDBConfig constructs the database connection string from environment variables
func GetDBConfig() string {
	host := os.Getenv("DB_URL")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")

	return "host=" + host + " user=" + user + " password=" + password +
		" dbname=" + dbname + " port=" + port + " sslmode=" + sslmode
}

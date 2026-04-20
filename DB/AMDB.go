package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres" 
	"gorm.io/gorm"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error in loading .env file")
	}
}

func InitDB() (*gorm.DB, *sql.DB) {
	// Get environment variables
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	// Connection string
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port,
	)

	// Connect to DB
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("❌ Failed to connect to database:", err)
		return nil, nil
	}

	// Get the underlying *sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println("Unable to get sql.DB from gorm.DB:", err)
		return nil, nil
	}

	fmt.Println("✅ Successfully connected to PostgreSQL!")
	return db, sqlDB
}

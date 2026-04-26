package main

import (
	"log"
	"os"

	"invoice-backend/DB"      // Your database initialization package
	"invoice-backend/Routes"  // The package where SetupRoutes lives
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	
)

func main() {
	// 1. Load environment variables (.env file)
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, using system environment variables")
	}

	// 2. Initialize Database Connection
	// This returns the GORM wrapper, so we extract the *sql.DB for our controllers
	gormDB := DB.InitDB()
	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatal("Failed to extract sql.DB from GORM:", err)
	}
	
	// Ensure the database connection closes when the application stops
	defer sqlDB.Close()

	// 3. Initialize Gin Engine
	// Use gin.Default() for built-in Logger and Recovery middleware
	r := gin.Default()

	// 4. Register Routes
	// We pass the engine and the database pointer to your routes package
	Routes.SetupRoutes(r, sqlDB)

	// 5. Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	log.Printf("Server is running on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
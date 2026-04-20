package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	
	"invoice-backend/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// 1. Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 2. Build Connection String from environment variables
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	// 3. Connect to Database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Could not connect to DB: ", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Database unreachable: ", err)
	}
	fmt.Println("✅ Database Connected Successfully")

	// 4. Setup Router
	r := gin.Default()

	// 5. Setup Routes (passing the db pool)
	routes.SetupRoutes(r, db)

	// 6. Run Server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	fmt.Printf("🚀 Server running on http://localhost:%s\n", port)
	r.Run(":" + port)
}
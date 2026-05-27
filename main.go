package main

import (
	"log"
	"os"
	"strings"

	"invoice-backend/DB"
	"invoice-backend/Routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

)

func main() {

	

	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found")
	}

	gormDB := DB.InitDB()
	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatal("Failed to extract sql.DB:", err)
	}
	defer sqlDB.Close()

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	// ✅ CORS MUST BE BEFORE ROUTES
	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {

			// localhost
			if origin == "http://localhost:5173" {
				return true
			}

			// allow all vercel domains
			if strings.Contains(origin, "vercel.app") {
				return true
			}

			return false
		},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// routes
	Routes.SetupRoutes(r, sqlDB)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
package DB 

import (
	"fmt" 
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var db *gorm.DB

// InitDB initializes the connection and returns both GORM and the raw SQL DB instance
func InitDB() *gorm.DB {
	var err error

	// Load credentials from environment variables
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode)

	// Open GORM connection with custom naming strategy
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",    
			SingularTable: true, 
		},
	})

	if err != nil {
		fmt.Printf("❌ Failed to connect to database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Database connection established successfully")
	return db
}

// GetConn provides the existing database instance to your Controllers and Services
func GetConn() *gorm.DB {
	if db == nil {
		return InitDB()
	}
	return db
}
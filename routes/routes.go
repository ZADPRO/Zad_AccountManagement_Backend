package routes

import (
	"database/sql"
	"invoice-backend/Controller"
	"invoice-backend/Middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, db *sql.DB) {
	// --- ADD CORS MIDDLEWARE HERE ---
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	v1 := r.Group("/api/v1")
	{
		// --- PUBLIC ROUTES ---
		v1.POST("/login", Controller.Login(db)) 
		v1.POST("/change-password", Controller.ChangePassword(db))
		
		drops := v1.Group("/dropdowns")
        {
            drops.GET("/states", Controller.GetStates(db))
            drops.GET("/countries", Controller.GetCountries(db))
            drops.GET("/roles", Controller.GetRoles(db))
        }
		v1.POST("/clients", Controller.CreateClient(db))
		// --- PROTECTED ROUTES ---
		protected := v1.Group("/")
		protected.Use(Middleware.AuthMiddleware())
		{
			// Client Management
			// NOTE: Changed to GetClientList to match your Controller name
			protected.GET("/clients", Controller.GetClientList(db)) 
			protected.GET("/profile", Controller.GetProfile(db))
			protected.PUT("/clients/:id", Controller.UpdateClient(db))
			protected.DELETE("/clients/:id", Controller.DeleteClient(db))
			protected.GET("/clients/:id", Controller.GetClientByID(db))

			// User Management
			protected.GET("/users", Controller.GetUserList(db))
			protected.POST("/users", Controller.CreateUser(db))
			protected.PUT("/users/:id", Controller.UpdateUser(db))
			protected.DELETE("/users/:id", Controller.DeleteUser(db)) 
			protected.GET("/users/:id", Controller.GetUserByID(db))

			// Invoice management
			protected.POST("/invoices", Controller.CreateInvoice(db))
			protected.GET("/invoices", Controller.GetInvoiceList(db))

			protected.POST("/payments", Controller.RecordPayment(db))

			// Dashboard 
			v1.GET("/dashboard/summary", Controller.GetDashboardStats(db))
			
		}
	}
}

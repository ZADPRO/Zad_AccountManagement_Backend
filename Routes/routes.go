package Routes 

import (
	"database/sql"
	"invoice-backend/Controllers"
	"invoice-backend/Middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, db *sql.DB) {
	// --- CORS MIDDLEWARE ---
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
		// --- PUBLIC ROUTES (No Token Required) ---
		v1.POST("/login", Controller.LoginHandler) 
		
		// Dropdowns (Often public or semi-public for registration/setup)
		drops := v1.Group("/dropdowns")
		{
			drops.GET("/states", Controller.GetStates(db))
			drops.GET("/countries", Controller.GetCountries(db))
			drops.GET("/roles", Controller.GetRoles(db))
		}

		// --- PROTECTED ROUTES (Requires JWT & Axios Encryption) ---
		protected := v1.Group("")
		protected.Use(Middleware.AuthMiddleware())
		{
			// Profile & Global Actions
			protected.GET("/profile", Controller.GetProfile(db))
			protected.POST("/change-password", Controller.ChangePasswordHandler(db) )
			protected.POST("/users/force-password", Controller.ForceChangePassword(db))

			// Client Management
			clientRoutes := protected.Group("/clients")
			{
				clientRoutes.GET("", Controller.GetClientList(db))
				clientRoutes.POST("", Controller.CreateClient(db))
				clientRoutes.GET("/:id", Controller.GetClientByID(db))
				clientRoutes.PUT("/:id", Controller.UpdateClient(db))
				clientRoutes.DELETE("/:id", Controller.DeleteClient(db))
			}

			// User Management
			userRoutes := protected.Group("/users")
			{
				userRoutes.GET("", Controller.GetUserList(db))
				userRoutes.POST("", Controller.CreateUser(db))
				userRoutes.GET("/:id", Controller.GetUserByID(db))
				userRoutes.PUT("/:id", Controller.UpdateUser(db))
				userRoutes.DELETE("/:id", Controller.DeleteUser(db))
			}

			// Invoice Management
			invoiceRoutes := protected.Group("/invoices")
			{
				invoiceRoutes.GET("", Controller.GetInvoiceList(db))
				invoiceRoutes.POST("", Controller.CreateInvoice(db))
				invoiceRoutes.GET("/:id", Controller.GetInvoiceByID(db))
			} 

			bankingRoutes := protected.Group("/banking")
{
    bankingRoutes.GET("/current", Controller.GetBankingInfo(db))

    // ✅ split endpoints
    bankingRoutes.POST("", Controller.CreateBanking(db))       // create
    bankingRoutes.PUT("/:id", Controller.UpdateBanking(db))    // update
	// Inside your routes setup
	bankingRoutes.DELETE("/:id", Controller.DeleteBankingInfo(db))
}

	fieldRoutes := protected.Group("/custom-fields")
    {
        fieldRoutes.GET("", Controller.GetCustomFieldList(db))
        fieldRoutes.POST("", Controller.CreateCustomField(db))
        fieldRoutes.DELETE("/:id", Controller.DeleteCustomField(db))
    }

			// Payments
			//protected.POST("/payments", Controller.RecordPayment(db))

			// Dashboard Stats
			protected.GET("/dashboard/summary", Controller.GetDashboardStats(db))
		}
	}
}
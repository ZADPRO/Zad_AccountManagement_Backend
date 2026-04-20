package Controller

import (
	"database/sql"
	"net/http"
	"invoice-backend/Service"
	"github.com/gin-gonic/gin"
)

func GetDashboardStats(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		summary, err := Service.GetDashboardSummary(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"data":   summary,
		})
	}
}
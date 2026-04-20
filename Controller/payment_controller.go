package Controller

import (
	"database/sql"
	"invoice-backend/Model"
	"invoice-backend/Query"
	"invoice-backend/Service"
	"net/http"
	"github.com/gin-gonic/gin"
)

func RecordPayment(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req Model.PaymentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid input data"})
			return
		}

		// Start Transaction to ensure Data Integrity 
		tx, err := db.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Transaction error"})
			return
		}

		// 1. Insert into transactionhistory
		_, err = tx.Exec(Query.InsertTransactionQuery, 
			req.InvoiceID, 
			req.Amount, 
			req.TransactionDate, 
			req.UpdatedBy,
		)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to log transaction"})
			return
		}

		// 2. Calculate new status based on Business Rules [cite: 44]
		newStatus, err := Service.CheckInvoiceStatus(tx, req.InvoiceID)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Status calculation failed"})
			return
		}

		// 3. Update the invoices table status
		_, err = tx.Exec(`UPDATE invoices SET paymentstatus = $1, updatedat = NOW(), updatedby = $2 WHERE invoiceid = $3`, 
			newStatus, req.UpdatedBy, req.InvoiceID)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to update invoice status"})
			return
		}

		tx.Commit()
		c.JSON(http.StatusOK, gin.H{
			"status": true, 
			"message": "Payment recorded and invoice status updated to " + newStatus,
		})
	}
}
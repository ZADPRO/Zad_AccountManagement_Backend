package Controller

import (
	"database/sql"
	"invoice-backend/Model"
	"invoice-backend/Query"
	"net/http"
	"github.com/gin-gonic/gin"
)

func CreateInvoice(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req Model.CreateInvoiceRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, Model.BaseResponse{Status: false, Message: "Invalid Data"})
			return
		}

		// 1. Begin Transaction
		tx, err := db.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, Model.BaseResponse{Status: false, Message: "Internal Server Error"})
			return
		}

		// 2. Insert Header
		var newInvoiceID int
		err = tx.QueryRow(Query.InsertInvoiceHeaderQuery, 
			req.InvoiceNumber, req.ClientID, req.InvoiceDate, req.GrandTotal, req.PaymentStatus, req.UpdatedBy,
		).Scan(&newInvoiceID)

		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, Model.BaseResponse{Status: false, Message: "Header Insert Failed: " + err.Error()})
			return
		}

		// 3. Insert Items
		for _, item := range req.Items {
			_, err = tx.Exec(Query.InsertInvoiceItemQuery, 
				newInvoiceID, item.Description, item.Quantity, item.UnitPrice, item.LineTotal, req.UpdatedBy,
			)
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, Model.BaseResponse{Status: false, Message: "Item Insert Failed"})
				return
			}
		}

		// 4. Commit
		tx.Commit()
		c.JSON(http.StatusOK, Model.BaseResponse{Status: true, Message: "Invoice Created Successfully"})
	}
} 



func GetInvoiceList(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query(Query.GetInvoiceListQuery)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Model.BaseResponse{
				Status:  false, 
				Message: "Failed to fetch invoices: " + err.Error(),
			})
			return
		}
		defer rows.Close()

		var invoices []map[string]interface{}
		for rows.Next() {
			var id int
			var invNum, clientName, date, status string
			var total float64

			// Adjust Scan variables based on your specific SELECT columns
			err := rows.Scan(&id, &invNum, &clientName, &date, &total, &status)
			if err != nil {
				continue
			}

			invoices = append(invoices, map[string]interface{}{
				"invoiceid":     id,
				"invoicenumber": invNum,
				"clientname":    clientName,
				"invoicedate":   date,
				"grandtotal":    total,
				"paymentstatus": status,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"data":   invoices,
		})
	}
}
package Controller

import (
	"database/sql"
	"encoding/json"
	"invoice-backend/Helper/HashAPI"
	"strconv"
	"invoice-backend/Models/dto"
	"invoice-backend/Query"
	"net/http"

	"github.com/gin-gonic/gin"
	"invoice-backend/Services"

)

func CreateInvoice(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getToken(c)
		var packet dto.EncryptedPacket

		// 1. Bind the Encrypted Packet
		if err := c.ShouldBindJSON(&packet); err != nil {
			c.JSON(http.StatusBadRequest, hashapi.Encrypt(gin.H{"status": false, "message": "Invalid Packet"}, true, token))
			return
		}

		// 2. Decrypt the Invoice Data
		decrypted, err := hashapi.Decrypt(packet.Data, token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, hashapi.Encrypt(gin.H{"status": false, "message": "Decryption failed"}, true, token))
			return
		}

		// Unmarshal decrypted data into the struct
		var req dto.CreateInvoiceRequest
		jsonBytes, _ := json.Marshal(decrypted)
		json.Unmarshal(jsonBytes, &req)

		// 3. Begin Transaction
		tx, err := db.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, hashapi.Encrypt(gin.H{"status": false, "message": "Transaction error"}, true, token))
			return
		}

		// 4. Insert Header
		var newInvoiceID int
		err = tx.QueryRow(Query.InsertInvoiceHeaderQuery,
			req.InvoiceNumber, req.ClientID, req.InvoiceDate, req.GrandTotal, req.PaymentStatus, req.UpdatedBy,
		).Scan(&newInvoiceID)

		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, hashapi.Encrypt(gin.H{"status": false, "message": "Header Insert Failed"}, true, token))
			return
		}

		// 5. Insert Items
		for _, item := range req.Items {
			_, err = tx.Exec(Query.InsertInvoiceItemQuery,
				newInvoiceID, item.Description, item.Quantity, item.UnitPrice, item.LineTotal, req.UpdatedBy,
			)
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, hashapi.Encrypt(gin.H{"status": false, "message": "Item insertion failed"}, true, token))
				return
			}
		}

		// 6. Commit and Send Encrypted Response
		tx.Commit()
		c.JSON(http.StatusOK, gin.H{
    "status": true,
    "data": hashapi.Encrypt(dto.BaseResponse{
        Status:  true,
        Message: "Invoice Created Successfully",
    }, true, token),
})
	}
}

func GetInvoiceList(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getToken(c)

		rows, err := db.Query(Query.GetInvoiceListQuery)
		if err != nil {
			c.JSON(http.StatusInternalServerError, hashapi.Encrypt(gin.H{"status": false, "message": err.Error()}, true, token))
			return
		}
		defer rows.Close()

		var invoices []map[string]interface{}
		for rows.Next() {
			var id int
			var invNum, clientName, date, status string
			var total float64

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

		// Encrypt the list so React can decrypt it
		c.JSON(http.StatusOK, hashapi.Encrypt(gin.H{
			"status": true,
			"data":   invoices,
		}, true, token))
	}
}
func GetInvoiceByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// 1. Get ID from URL
		idParam := c.Param("id")
		invoiceID, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid invoice ID",
			})
			return
		}

		// 2. Call service
		invoice, err := Services.GetInvoiceByID(db, invoiceID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		// 3. Return JSON
		c.JSON(http.StatusOK, invoice)
	}
}
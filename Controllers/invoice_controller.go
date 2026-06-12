package Controller

import (
	"database/sql"
	"invoice-backend/Helper/HashAPI"
	"strconv"
	"invoice-backend/Models/dto"
	"invoice-backend/Query"
	"net/http"
	"github.com/gin-gonic/gin"
	"invoice-backend/Services" 
	 
	"encoding/json"

)

func CreateInvoice(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		token := getToken(c)

		// 1. Read encrypted packet
		var packet dto.EncryptedPacket
		if err := c.ShouldBindJSON(&packet); err != nil {
			c.JSON(http.StatusBadRequest, hashapi.Encrypt(gin.H{
				"status": false,
				"message": "Invalid Packet",
			}, true, token))
			return
		}

		// 2. Decrypt request
		decrypted, err := hashapi.Decrypt(packet.Data, token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, hashapi.Encrypt(gin.H{
				"status": false,
				"message": "Decryption failed",
			}, true, token))
			return
		}

		// 3. Convert decrypted → struct
		var req dto.CreateInvoiceRequest
		jsonBytes, _ := json.Marshal(decrypted)
		if err := json.Unmarshal(jsonBytes, &req); err != nil {

		//fmt.Println("UNMARSHAL ERROR:", err)
		//fmt.Println("DECRYPTED DATA:", string(jsonBytes))

	c.JSON(http.StatusBadRequest, hashapi.Encrypt(gin.H{
		"status": false,
		"message": err.Error(),
	}, true, token))

	return
}

		// ✅ DEBUG (now this will work correctly)
		//fmt.Println("RAW REQ:", req)
		//fmt.Printf("REQ: %+v\n", req)
		//fmt.Println("IS SAVE DRAFT =", req.IsSaveDraft)
		//fmt.Println("DATE:", req.InvoiceDate)

		// 4. Call service
		newInvoiceID, err := Services.CreateFullInvoice(db, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, hashapi.Encrypt(gin.H{
				"status": false,
				"message": err.Error(),
			}, true, token))
			return
		}

		// 5. Encrypt response
		resp := hashapi.Encrypt(gin.H{
			"status": true,
			"message": "Invoice Created Successfully",
			"invoiceid": newInvoiceID,
		}, true, token)

		c.JSON(http.StatusOK, resp)
	}
}

func UpdateInvoice(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {

        invoiceID, err := strconv.Atoi(
            c.Param("id"),
        )

        if err != nil {
            c.JSON(
                http.StatusBadRequest,
                gin.H{"error": "invalid invoice id"},
            )
            return
        }

        token := getToken(c)

var packet dto.EncryptedPacket

if err := c.ShouldBindJSON(&packet); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
        "error": "Invalid Packet",
    })
    return
}

decrypted, err := hashapi.Decrypt(
    packet.Data,
    token,
)

if err != nil {
    c.JSON(http.StatusUnauthorized, gin.H{
        "error": "Decryption failed",
    })
    return
}

var req dto.CreateInvoiceRequest

jsonBytes, _ := json.Marshal(decrypted)

if err := json.Unmarshal(
    jsonBytes,
    &req,
); err != nil {

    c.JSON(http.StatusBadRequest, gin.H{
        "error": err.Error(),
    })
    return
}


        err = Services.UpdateInvoice(
            db,
            invoiceID,
            req,
        )

		

        if err != nil {
            c.JSON(
                http.StatusInternalServerError,
                gin.H{"error": err.Error()},
            )
            return
        }

        c.JSON(
            http.StatusOK,
            gin.H{
                "message": "invoice updated",
            },
        )
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
			var isSaveDraft bool
			var invoiceType string

			err := rows.Scan(&id, &invNum,&invoiceType, &clientName, &date, &total, &status,&isSaveDraft,)
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
				"issavedraft":   isSaveDraft,
				"invoicetype":   invoiceType,
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

		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"data": invoice,
		})
	}
	
}
func DeleteInvoice(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		invoiceIDParam := c.Param("id")

		invoiceID, err := strconv.Atoi(invoiceIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid invoice ID",
			})
			return
		}

		// Replace with actual logged-in admin/user ID
		adminID := 1

		err = Services.DeleteInvoice(db, invoiceID, adminID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "invoice deleted successfully",
		})
	}
}
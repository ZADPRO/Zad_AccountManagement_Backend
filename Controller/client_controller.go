package Controller

import (
	"database/sql"
	"invoice-backend/Model"
	"invoice-backend/Service"
	"net/http"
	"strconv"
	
	"github.com/gin-gonic/gin"
	
)

// GetClientList fetches all clients to display in the React table
func GetClientList(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Call Service to query Beekeeper/Postgres
		clients, err := Service.GetAllClients(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Model.BaseResponse{
				Status:  false,
				Message: "Failed to fetch clients: " + err.Error(),
			})
			return
		}

		// 2. Wrap the slice in a "clients" key to match your React needs
		// This results in: { "clients": [...] }
		c.JSON(http.StatusOK, gin.H{
			"clients": clients,
		})
	}
}
func CreateClient(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req Model.CreateClientRequest

        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
            return
        }
		
        clientID, err := Service.CreateClient(db, req)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "status":   true,
            "message":  "Client created successfully",
            "clientID": clientID,
        })
    }
}
func UpdateClient(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		clientID, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid Client ID format"})
			return
		}

		var req Model.CreateClientRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid request body"})
			return
		}

		err = Service.UpdateClient(db, clientID, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Update failed: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Client information updated successfully",
		})
	}
}

func DeleteClient(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		clientID, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, Model.BaseResponse{Status: false, Message: "Invalid Client ID"})
			return
		}

		var req struct {
			DeletedBy int `json:"deletedBy"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, Model.BaseResponse{Status: false, Message: "Admin ID (deletedBy) is required"})
			return
		}

		err = Service.DeleteClient(db, clientID, req.DeletedBy)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Model.BaseResponse{Status: false, Message: "Error: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, Model.BaseResponse{Status: true, Message: "Client deactivated successfully"})
	}
}

func GetClientByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		idParam := c.Param("id")

		clientID, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": false,
				"message": "Invalid client ID",
			})
			return
		}

		client, err := Service.GetClientByID(db, clientID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": false,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"data":   client,
		})
	}
}
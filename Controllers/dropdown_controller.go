package Controller

import (
	"database/sql"
	"invoice-backend/Helper/HashAPI" 
	"invoice-backend/Models/dto"
	"invoice-backend/Query"
	"invoice-backend/Services"
	"net/http"
	"github.com/gin-gonic/gin"
)

// GetCountries retrieves a list of countries for address forms.
// It uses the generic DropdownData service to execute a specific SQL query.
func GetCountries(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
		// Extract JWT for response encryption
        token := getToken(c)
        
		// Fetch data from the database using a predefined query constant
		data, err := Services.GetDropdownData(db, Query.GetCountriesDropdownQuery)
		if err != nil {
			c.JSON(http.StatusInternalServerError, hashapi.Encrypt(dto.DropdownResponse{
				BaseResponse: dto.BaseResponse{Status: false, Message: "Failed to fetch countries: " + err.Error()},
			}, true, token))
			return
		}
		
		// Encrypt the list of countries. Even though countries aren't sensitive,
		// the frontend expects all authenticated responses to be encrypted.
		c.JSON(http.StatusOK, hashapi.Encrypt(dto.DropdownResponse{
			BaseResponse: dto.BaseResponse{Status: true, Message: "Success"},
			Data:         data,
		}, true, token))
	}
}

// GetStates retrieves a list of states/provinces for address forms.
func GetStates(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getToken(c)

		data, err := Services.GetDropdownData(db, Query.GetStatesDropdownQuery)
		if err != nil {
			c.JSON(http.StatusInternalServerError, hashapi.Encrypt(dto.DropdownResponse{
				BaseResponse: dto.BaseResponse{Status: false, Message: "Failed to fetch states: " + err.Error()},
			}, true, token))
			return
		}

		c.JSON(http.StatusOK, hashapi.Encrypt(dto.DropdownResponse{
			BaseResponse: dto.BaseResponse{Status: true, Message: "Success"},
			Data:         data,
		}, true, token))
	}
}

// GetRoles retrieves available user roles (e.g., Admin, Staff, Viewer).
// This is used in User Management sections to assign permissions.
func GetRoles(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getToken(c)

		data, err := Services.GetDropdownData(db, Query.GetRolesDropdownQuery)
		if err != nil {
			c.JSON(http.StatusInternalServerError, hashapi.Encrypt(dto.DropdownResponse{
				BaseResponse: dto.BaseResponse{Status: false, Message: "Failed to fetch roles: " + err.Error()},
			}, true, token))
			return
		}

		c.JSON(http.StatusOK, hashapi.Encrypt(dto.DropdownResponse{
			BaseResponse: dto.BaseResponse{Status: true, Message: "Success"},
			Data:         data,
		}, true, token))
	}
}
package Controller

import (
	"database/sql"
	"invoice-backend/Helper/HashAPI" // Ensure this matches your encryption helper path
	"invoice-backend/Models/dto"
	"invoice-backend/Query"
	"invoice-backend/Services"
	"net/http"
	"fmt"
	"github.com/gin-gonic/gin"
)

// GetCountries handles the Country list for forms
func GetCountries(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := getToken(c)
        fmt.Printf("DEBUG GetCountries token: %q\n", token)  // ← add this
        fmt.Printf("DEBUG GetCountries token len: %d\n", len(token))
        // ...
    


		data, err := Services.GetDropdownData(db, Query.GetCountriesDropdownQuery)
		if err != nil {
			c.JSON(http.StatusInternalServerError, hashapi.Encrypt(dto.DropdownResponse{
				BaseResponse: dto.BaseResponse{Status: false, Message: "Failed to fetch countries: " + err.Error()},
			}, true, token))
			return
		}

		c.JSON(http.StatusOK, hashapi.Encrypt(dto.DropdownResponse{
			BaseResponse: dto.BaseResponse{Status: true, Message: "Success"},
			Data:         data,
		}, true, token))
	}
}

// GetStates handles the State list for forms
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

// GetRoles handles the Role list for User Management
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
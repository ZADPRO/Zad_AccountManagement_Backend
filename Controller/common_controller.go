package Controller

import (
	"database/sql"
	"invoice-backend/Model"
	"invoice-backend/Query"
	"invoice-backend/Service"
	"net/http"
	"github.com/gin-gonic/gin"
)

// GetCountries handles the Country list
func GetCountries(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := Service.GetDropdownData(db, Query.GetCountriesDropdownQuery)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Model.DropdownResponse{
				BaseResponse: Model.BaseResponse{Status: false, Message: err.Error()},
			})
			return
		}
		c.JSON(http.StatusOK, Model.DropdownResponse{
			BaseResponse: Model.BaseResponse{Status: true, Message: "Success"},
			Data:         data,
		})
	}
}

func GetStates(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := Service.GetDropdownData(db, Query.GetStatesDropdownQuery)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Model.DropdownResponse{
				BaseResponse: Model.BaseResponse{Status: false, Message: err.Error()},
			})
			return
		}
		c.JSON(http.StatusOK, Model.DropdownResponse{
			BaseResponse: Model.BaseResponse{Status: true, Message: "Success"},
			Data:         data,
		})
	}
}

// GetRoles handles the Role list for User Creation
func GetRoles(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := Service.GetDropdownData(db, Query.GetRolesDropdownQuery)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Model.DropdownResponse{
				BaseResponse: Model.BaseResponse{Status: false, Message: err.Error()},
			})
			return
		}
		c.JSON(http.StatusOK, Model.DropdownResponse{
			BaseResponse: Model.BaseResponse{Status: true, Message: "Success"},
			Data:         data,
		})
	}
}
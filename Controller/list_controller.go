package Controller

import (
	"invoice-backend/Model"
	"invoice-backend/Service"
	"database/sql"
	"net/http"
	"github.com/gin-gonic/gin"
)

// func GetClientList(db *sql.DB) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		data, err := Service.FetchAllClients(db)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, Model.ClientListResponse{
// 				BaseResponse: Model.BaseResponse{Status: false, Message: err.Error()},
// 			})
// 			return
// 		}
// 		c.JSON(http.StatusOK, Model.ClientListResponse{
// 			BaseResponse: Model.BaseResponse{Status: true, Message: "Clients fetched successfully"},
// 			Clients:      data,
// 		})
// 	}
// 

func GetUserList(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := Service.FetchAllUsers(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Model.UserListResponse{
				BaseResponse: Model.BaseResponse{Status: false, Message: err.Error()},
			})
			return
		}
		c.JSON(http.StatusOK, Model.UserListResponse{
			BaseResponse: Model.BaseResponse{Status: true, Message: "Users fetched successfully"},
			Users:        data,
		})
	}
}
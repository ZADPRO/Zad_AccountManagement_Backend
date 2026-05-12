package dto

import "invoice-backend/Models/internal"


type CreateUserRequest struct {
	UserCode  string `json:"userCode" binding:"required"`
	Username  string `json:"username" binding:"required"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	RoleID    int    `json:"roleId" binding:"required"`
	Email     string `json:"email"`
	internal.AuditModel
}

// CreateUserResponse returns the ID of the newly created user
type CreateUserResponse struct {
	BaseResponse
	UserID int `json:"userId,omitempty"`
	internal.AuditModel
}

// UpdateUserRequest for profile/admin updates
type UpdateUserRequest struct {
	Username  string `json:"username" binding:"required"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	RoleID    int    `json:"roleId" binding:"required"`
	Email     string `json:"email"`
}

// Specific struct for Delete to avoid binding issues
type DeleteUserRequest struct {
	AdminID int `json:"updatedBy" binding:"required"`
}
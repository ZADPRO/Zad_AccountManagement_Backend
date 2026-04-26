package responses

import (
	"invoice-backend/Models/dto"
	"invoice-backend/Models/internal"
)

// UserData matches your "active_users" view and original struct
type UserData struct {
	UserID       int    `json:"userId"`
	UserCode     string `json:"userCode"`
	Username     string `json:"username"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	RoleID       int    `json:"roleId"`
	Email        string `json:"email"`
	IsFirstLogin bool   `json:"isFirstLogin"`
	internal.AuditModel
}

type UserListResponse struct {
	dto.BaseResponse
	Users []UserData `json:"users,omitempty"`
}

type UserProfileResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	RoleName  string `json:"role_name"`
}

type ProfileResponse struct {
	Status  bool                `json:"status"`
	Message string              `json:"message"`
	User    UserProfileResponse `json:"user"`
}
package dto

// LoginRequest captures credentials from the login form
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse returns the JWT token upon successful auth
type LoginResponse struct {
	BaseResponse
	Token        string `json:"token"`
	Role         string `json:"role"`     
	Username     string `json:"name"`
	IsFirstLogin bool   `json:"isFirstLogin"`
	UserId       int    `json:"userId"`
}

// ResetFinalizeRequest for password updates
type ResetFinalizeRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}
// Create this for your "First Login" change password flow
type ChangePasswordRequest struct {
    UserID      int    `json:"userid" binding:"required"`
    NewPassword string `json:"newPassword" binding:"required,min=8"`
}
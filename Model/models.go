package Model

import (
	
	"time"
)


type AuditModel struct {
    CreatedAt time.Time  `json:"createdAt"`
    CreatedBy int32     `json:"createdBy"` // Changed to string
    UpdatedAt time.Time  `json:"updatedAt"`
    UpdatedBy int32     `json:"updatedBy"` // Changed to string
    DeletedAt *time.Time `json:"deletedAt,omitempty"`
    DeletedBy *string    `json:"deletedBy,omitempty"` // Changed to string pointer

}
// BaseResponse for consistent API output
type BaseResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

// ClientListModel matches the "active_clients" view
type ClientListModel struct {
	ClientID     int       `json:"clientId" gorm:"column:ClientID"`
	ClientCode   string    `json:"clientCode" gorm:"column:ClientCode"`
	Name         string    `json:"name" gorm:"column:Name"`
	BusinessName string    `json:"businessName" gorm:"column:BusinessName"`
	IsActive     bool      `json:"isActive" gorm:"column:IsActive"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"column:UpdatedAt"`
	AuditModel 
}

type ClientListResponse struct {
	BaseResponse
	Clients []ClientListModel `json:"clients,omitempty"`
}

// UserData matches the "active_users" view
type UserData struct {
	UserID    int    `json:"userId" gorm:"column:UserID"`
	UserCode  string `json:"userCode" gorm:"column:UserCode"`
	Username  string `json:"username" gorm:"column:Username"`
	FirstName string `json:"firstName" gorm:"column:FirstName"`
	LastName  string `json:"lastName" gorm:"column:LastName"`
	RoleID    int    `json:"roleId" gorm:"column:RoleID"` 
	Email     string `json:"email" gorm:"column:EmailID"` 
	IsFirstLogin bool `json:"isFirstLogin" gorm:"column:is_first_login"`
	AuditModel 
}

type UserListResponse struct {
	BaseResponse
	Users []UserData `json:"users,omitempty"`
}

type PaginationModel struct {
	TotalRecords int `json:"totalRecords"`
	CurrentPage  int `json:"currentPage"`
	TotalPages   int `json:"totalPages"`
	Limit        int `json:"limit"`
} 


// DropdownModel is a generic structure for Select/Dropdown menus in the UI
type DropdownModel struct {
    // Adding json tags ensures the keys are lowercase in the API response
    ID   int    `json:"id"`   
    Name string `json:"name"` 
}

type ClientResponse struct {
    // This "embeds" the Status and Message fields automatically
    BaseResponse 
    
    Data []ClientListModel `json:"data"`
    
    // Optional: add pagination if the list is long
    Pagination PaginationModel `json:"pagination,omitempty"`
}

// CreateClientRequest matches the JSON expected from your React form
type CreateClientRequest struct {
    ClientCode    string `json:"clientCode" binding:"required"`
    Name          string `json:"name" binding:"required"`
    BusinessName  string `json:"businessName" binding:"required"`
    SupplyTypeID  int    `json:"supplytypeid"` 
    ClientType    string `json:"clienttype"`   

    GSTStatus     string `json:"gststatus"`
    GSTNumber     string `json:"gstnumber"` 
    PAN           string `json:"pan"`
    IsExport      bool   `json:"isexport"`

    CountryName   string `json:"countryName" binding:"required"` // Match the capital N
    StateName     string `json:"stateName" `   // Match the capital N
    ZIP           int    `json:"zip" binding:"required"`
    Address       string `json:"registeredAddress" binding:"required"`

    PrimaryNumber string `json:"mobilenumber" binding:"required"` 
    Email         string `json:"email"` 
    
    // Ensure AuditModel fields are handled by the controller
    AuditModel 
}


// CreateClientResponse returns the ID of the newly created client
type CreateClientResponse struct {
	BaseResponse
	ClientID int `json:"clientId,omitempty"`
	AuditModel 
}

type CreateUserRequest struct {	
	UserCode  string `json:"userCode" binding:"required"`
	Username  string `json:"username" binding:"required"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	RoleID    int    `json:"roleId" binding:"required"`
	Email     string `json:"email"`
	AuditModel 
}

type CreateUserResponse struct {
	BaseResponse
	UserID int `json:"userId,omitempty"`
	AuditModel 
}


type DropdownResponse struct {
    BaseResponse                // Use embedding for consistency
    Data []DropdownModel `json:"data"`
}
// Specific struct for Delete to avoid binding issues
type DeleteUserRequest struct {
    AdminID int `json:"updatedBy" binding:"required"`
} 

type DeleteRequest struct {
    DeletedBy int `json:"deletedBy"` // Postman must send "deletedBy"
} 


// LoginRequest captures credentials from the login form
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password"` 
}

// LoginResponse returns the JWT token upon successful auth
type LoginResponse struct {
    BaseResponse
    Token string `json:"token"`
    Role  string `json:"role"` // 🔥 This matches data.role in your React code 
	Username string `json:"name"` 
	IsFirstLogin bool   `json:"isFirstLogin"` 
	UserId  int  `json:"userId"`
} 

type ClientDetailsResponse struct {
	ClientID           int    `json:"clientId"`
	ClientCode         string `json:"clientCode"`
	Name               string `json:"name"`
	BusinessName       string `json:"businessName"`
	SupplyTypeID       int    `json:"supplyTypeId"`
	IsActive           bool   `json:"isActive"`
	ClientType         string `json:"clientType"`
	Email              string `json:"email"`
	MobileNumber       string `json:"mobileNumber"`
	RegisteredAddress  string `json:"registeredAddress"`
	CountryName        string `json:"countryName"`
	StateName          string `json:"stateName"`
	ZIP                int    `json:"zip"`

	GSTNumber          string `json:"gstNumber"`
	PAN                string `json:"pan"`
	IsExport           bool   `json:"isExport"`
	GSTStatus          string `json:"gstStatus"`
}
type UpdateUserRequest struct {
	Username  string `json:"username" binding:"required"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	RoleID    int    `json:"roleId" binding:"required"`
	Email     string `json:"email"`
} 



// UserProfile represents the combined data for the Header and Profile fetch
// In Model/user_model.go
type UserProfileResponse struct {
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    RoleName  string `json:"role_name"`
}

// ProfileResponse is the standard wrapper for the profile API
type ProfileResponse struct {
    Status  bool        `json:"status"`
    Message string      `json:"message"`
    User    UserProfileResponse `json:"user"`
}

type ResetFinalizeRequest struct {
    Token       string `json:"token" binding:"required"`
    NewPassword string `json:"newPassword" binding:"required,min=6"`
}
 

type DashboardSummary struct {
	TotalClients   int     `json:"totalClients"`
	ActiveUsers    int     `json:"activeUsers"`
	TotalRevenue   float64 `json:"totalRevenue"`
	PendingAmount  float64 `json:"pendingAmount"`
	OverdueCount   int     `json:"overdueCount"`
}
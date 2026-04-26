package dto

// BaseResponse for consistent API output
type BaseResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

// PaginationModel for lists
type PaginationModel struct {
	TotalRecords int `json:"totalRecords"`
	CurrentPage  int `json:"currentPage"`
	TotalPages   int `json:"totalPages"`
	Limit        int `json:"limit"`
}

// DropdownModel is a generic structure for Select/Dropdown menus in the UI
type DropdownModel struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// DropdownResponse wrapper
type DropdownResponse struct {
	BaseResponse
	Data []DropdownModel `json:"data"`
}
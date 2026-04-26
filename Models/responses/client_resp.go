package responses

import (
	"invoice-backend/Models/dto"
	"invoice-backend/Models/internal"
)

// ClientListModel matches the "active_clients" view
type ClientListModel struct {
	ClientID     int    `json:"clientId"`
	ClientCode   string `json:"clientCode"`
	Name         string `json:"name"`
	BusinessName string `json:"businessName"`
	IsActive     bool   `json:"isActive"`
	internal.AuditModel
}

type ClientListResponse struct {
	dto.BaseResponse
	Clients []ClientListModel `json:"clients,omitempty"`
}

type ClientDetailsResponse struct {
    ClientID          int     `json:"clientId"`
    ClientCode        string  `json:"clientCode"`
    Name              string  `json:"name"`
    BusinessName      string  `json:"businessName"`
    SupplyTypeID      int     `json:"supplyTypeId"`
    IsActive          bool    `json:"isActive"`
    ClientType        string  `json:"clientType"`
    UpdatedAt         string  `json:"updatedAt"`   // add
    UpdatedBy         string  `json:"updatedBy"`   // add
    Email             string  `json:"email"`
    MobileNumber      string  `json:"mobileNumber"`

    // Registered
    RegisteredAddress string  `json:"registeredAddress"`
    CountryName       string  `json:"countryName"`
    StateName         string  `json:"stateName"`
    ZIP               int     `json:"zip"`

    // Billing
    BillingAddress     string  `json:"billingAddress"`
    BillingCountryID   int     `json:"billingCountryId"`
    BillingStateID     int     `json:"billingStateId"`
    TaxPercentage      float64 `json:"taxPercentage"`

    // Tax
    GSTNumber          string `json:"gstNumber"`
    PAN                string `json:"pan"`
    IsExport           bool   `json:"isExport"`
    GSTStatus          string `json:"gstStatus"`

    BillingCountryName string `json:"billingCountryName"` // add
    BillingStateName   string `json:"billingStateName"`   // add
}
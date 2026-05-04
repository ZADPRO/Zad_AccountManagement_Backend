package dto

import "invoice-backend/Models/internal"

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
	CountryName   string `json:"countryName"`
	StateName     string `json:"stateName"`
	ZIP           int    `json:"zip"`
	Address       string `json:"registeredAddress"`
	PrimaryNumber string `json:"mobilenumber"`
	Email         string `json:"email"`
	 
	
	// Tax + Billing
	BillingAddress   string  `json:"billingAddress"`
	BillingCountryID int     `json:"billingCountryId" binding:"required"`
	BillingStateID   *int     `json:"billingStateId"`

	TaxPercentage    float64 `json:"tax_percentage"`
	internal.AuditModel
}

// CreateClientResponse returns the ID of the newly created client
type CreateClientResponse struct {
	BaseResponse
	ClientID int `json:"clientId,omitempty"`
	internal.AuditModel
}

type DeleteRequest struct {
	DeletedBy int `json:"deletedBy"`
}
type ClientFullResponse struct {
    ClientID     int     `json:"clientid"`
    ClientCode   string  `json:"clientCode"`
    Name         string  `json:"name"`
    BusinessName string  `json:"businessName"`
	ClientType    string `json:"clienttype"`
    Email        string  `json:"email"`
    PrimaryNumber string `json:"mobilenumber"`
    Address      string  `json:"registeredAddress"`
    CountryName  string  `json:"countryName"`
    StateName    string  `json:"stateName"`
    ZIP          int     `json:"zip"`
 // Tax fields
    GSTNumber     string  `json:"gstnumber"`
    PAN           string  `json:"pan"`
    GSTStatus     string  `json:"gststatus"`
    IsExport      bool    `json:"isexport"`
    TaxPercentage float64 `json:"tax_percentage"`

    // Billing fields
    BillingAddress string `json:"billingAddress"`
    BillingCountry string `json:"billingCountry"`
    BillingState   string `json:"billingState"`
}
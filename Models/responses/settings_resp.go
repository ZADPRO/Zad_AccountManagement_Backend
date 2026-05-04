package responses

import "invoice-backend/Models/dto"


type CustomFieldData struct {
	FieldID    int    `json:"fieldId"`
	FieldLabel string `json:"fieldLabel"` // Will be decrypted before sending
	FieldType  string `json:"fieldType"`  // Will be decrypted before sending
	IsRequired bool   `json:"isRequired"`
	CreatedAt  string `json:"createdAt"` 
	CreatedBy  int    `json:"createdBy"`
}

type CustomFieldListResponse struct {
	dto.BaseResponse
	Fields []CustomFieldData `json:"fields"`
}

type BankingDetailsData struct {
    DetailsID    int    `json:"detailsId"`
    BankName     string `json:"bankName"`
    AccountNumber string `json:"accountNumber"`
    IfscCode     string `json:"ifscCode"`
    BankAddress  string `json:"bankAddress"` // Changed from businessAddress
    LogoURL      string `json:"qrCodeUrl"`
	AccountType   string `json:"accountType"` 
    SwiftCode     string `json:"swiftCode"`
} 

type BankingDetailsResponse struct {
    dto.BaseResponse
    // Change 'BankingDetailsData' to '[]BankingDetailsData'
    Data []BankingDetailsData `json:"data"` 
}

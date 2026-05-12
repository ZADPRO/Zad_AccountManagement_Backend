package responses

import "invoice-backend/Models/dto"


type CustomFieldData struct {
	FieldID    int    `json:"fieldId"`
	FieldLabel string `json:"fieldLabel"` 
	FieldType  string `json:"fieldType"`  
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
    BankAddress  string `json:"bankAddress"` 
    LogoURL      string `json:"qrCodeUrl"`
	AccountType   string `json:"accountType"` 
    SwiftCode     string `json:"swiftCode"`
} 

type BankingDetailsResponse struct {
    dto.BaseResponse
    Data []BankingDetailsData `json:"data"` 
}

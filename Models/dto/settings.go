package dto 

type SaveBankingRequest struct {
    DetailsID    int    `json:"detailsId,omitempty"` 
    BankName     string `json:"bankName"`
    AccountNumber string `json:"accountNumber"`
    IfscCode     string `json:"ifscCode"`
    BankAddress  string `json:"bankAddress"` // Changed from businessAddress
	AccountType   string `json:"accountType"` // e.g., "Savings", "Current"
    SwiftCode     string `json:"swiftCode"`
    LogoURL      string `json:"qrCodeUrl"`
    UserID int    `json:"userId"`
}

type CreateCustomFieldRequest struct {
	FieldLabel string `json:"fieldLabel"`
	FieldType  string `json:"fieldType"`
	IsRequired bool   `json:"isRequired"`
}
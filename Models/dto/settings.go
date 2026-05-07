package dto 
import "time"
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
type CustomFieldResponse struct {
    FieldID    int        `json:"fieldId"`    // Matches 'FieldID' in DB
    FieldLabel string     `json:"fieldLabel"` // Matches 'FieldLabel' in DB
    FieldType  string     `json:"fieldType"`  // Matches 'FieldType' in DB
    IsRequired bool       `json:"isRequired"` // Matches 'IsRequired' in DB
    CreatedAt  time.Time  `json:"createdAt"`
    CreatedBy  int        `json:"createdBy"`
    DeletedAt  *time.Time `json:"deletedAt,omitempty"` // Nullable for active fields
}
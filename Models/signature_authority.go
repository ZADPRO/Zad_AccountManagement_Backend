package models

type SignatureAuthority struct {
	ID            int    `json:"id" gorm:"column:id"`
	Name          string `json:"name" gorm:"column:name"`
	Designation   string `json:"designation" gorm:"column:designation"`
	ContactNumber string `json:"contactNumber" gorm:"column:contact_number"`
	Email         string `json:"email" gorm:"column:email"`
	SignatureURL  string `json:"signatureUrl" gorm:"column:signature_url"`
	IsActive      bool   `json:"isActive" gorm:"column:is_active"`
}

func (SignatureAuthority) TableName() string {
	return "signature_authorities"
}
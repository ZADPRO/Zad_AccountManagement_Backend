package models

type CompanyProfile struct {
	ID           int    `json:"id" gorm:"column:id;primaryKey"`
	CompanyName  string `json:"companyName" gorm:"column:companyname"`

	AddressLine1 string `json:"addressLine1" gorm:"column:addressline1"`
	AddressLine2 string `json:"addressLine2" gorm:"column:addressline2"`

	City         string `json:"city" gorm:"column:city"`
	State        string `json:"state" gorm:"column:state"`
	Country      string `json:"country" gorm:"column:country"`
	Pincode      string `json:"pincode" gorm:"column:pincode"`

	GSTNumber    string `json:"gstNumber" gorm:"column:gstnumber"`

	Email        string `json:"email" gorm:"column:email"`
	PhoneNumber  string `json:"phoneNumber" gorm:"column:phonenumber"`
	Website      string `json:"website" gorm:"column:website"`

	LogoUrl      string `json:"logoUrl" gorm:"column:logourl"`

	IsActive     bool   `json:"isActive" gorm:"column:isactive"`
}

func (CompanyProfile) TableName() string {
	return "companyprofile"
}
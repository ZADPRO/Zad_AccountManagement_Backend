package internal

type Client struct {
	ClientID      int    `gorm:"primaryKey;column:clientid"`
	ClientCode    string `gorm:"column:clientcode"`
	Name          string `gorm:"column:name"`           // Encrypted in DB
	BusinessName  string `gorm:"column:businessname"`   // Encrypted in DB
	SupplyTypeID  int    `gorm:"column:supplytypeid"`
	ClientType    string `gorm:"column:clienttype"`
	GSTStatus     string `gorm:"column:gststatus"`
	GSTNumber     string `gorm:"column:gstnumber"`      // Encrypted in DB
	PAN           string `gorm:"column:pan"`            // Encrypted in DB
	IsExport      bool   `gorm:"column:isexport"`
	CountryName   string `gorm:"column:countryname"`
	StateName     string `gorm:"column:statename"`
	ZIP           int    `gorm:"column:zip"`
	Address       string `gorm:"column:registeredaddress"` // Encrypted in DB
	PrimaryNumber string `gorm:"column:mobilenumber"`      // Encrypted in DB
	Email         string `gorm:"column:emailid"`           // Encrypted in DB
	IsActive      bool   `gorm:"column:isactive;default:true"`
	AuditModel           // Embeds audit fields
}


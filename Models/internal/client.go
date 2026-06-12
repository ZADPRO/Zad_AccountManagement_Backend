package internal

type Client struct {
	ClientID      int    `gorm:"primaryKey;column:clientid"`
	ClientCode    string `gorm:"column:clientcode"`
	Name          string `gorm:"column:name"`           
	BusinessName  string `gorm:"column:businessname"`   
	SupplyTypeID  int    `gorm:"column:supplytypeid"`
	ClientType    string `gorm:"column:clienttype"`
	
	GSTNumber     string `gorm:"column:gstnumber"`      
	PAN           string `gorm:"column:pan"`            
	IsExport      bool   `gorm:"column:isexport"`
	CountryName   string `gorm:"column:countryname"`
	StateName     string `gorm:"column:statename"`
	ZIP           int    `gorm:"column:zip"`
	Address       string `gorm:"column:registeredaddress"` 
	PrimaryNumber string `gorm:"column:mobilenumber"`      
	Email         string `gorm:"column:emailid"`           
	IsActive      bool   `gorm:"column:isactive;default:true"`
	AuditModel           
}


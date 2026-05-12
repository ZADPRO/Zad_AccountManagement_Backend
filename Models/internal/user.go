package internal

type User struct {
	UserID       int    `gorm:"primaryKey;column:userid"`
	UserCode     string `gorm:"column:usercode"`
	Username     string `gorm:"column:username"`
	FirstName    string `gorm:"column:firstname"` 
	LastName     string `gorm:"column:lastname"`  
	RoleID       int    `gorm:"column:roleid"`
	Email        string `gorm:"column:emailid"`   
	IsFirstLogin bool   `gorm:"column:is_first_login"`
	AuditModel          
}


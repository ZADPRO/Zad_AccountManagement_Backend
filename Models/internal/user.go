package internal

type User struct {
	UserID       int    `gorm:"primaryKey;column:userid"`
	UserCode     string `gorm:"column:usercode"`
	Username     string `gorm:"column:username"`
	FirstName    string `gorm:"column:firstname"` // Encrypted in DB
	LastName     string `gorm:"column:lastname"`  // Encrypted in DB
	RoleID       int    `gorm:"column:roleid"`
	Email        string `gorm:"column:emailid"`   // Encrypted in DB
	IsFirstLogin bool   `gorm:"column:is_first_login"`
	AuditModel          // Embeds audit fields
}


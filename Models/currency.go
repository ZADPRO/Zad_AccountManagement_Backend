package models

type Currency struct {
	ID             int    `json:"id" gorm:"column:id;primaryKey"`
	CurrencyCode   string `json:"currencyCode" gorm:"column:currencycode"`
	CurrencyName   string `json:"currencyName" gorm:"column:currencyname"`
	CurrencySymbol string `json:"currencySymbol" gorm:"column:currencysymbol"`
	IsActive       bool   `json:"isActive" gorm:"column:isactive"`
}

func (Currency) TableName() string {
	return "currencymaster"
}
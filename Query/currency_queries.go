package Query

import (
	db "invoice-backend/DB"
	models "invoice-backend/Models"
)

func GetAllCurrencies() ([]models.Currency, error) {

	var currencies []models.Currency

	err := db.GetConn().
		Order("id DESC").
		Find(&currencies).Error

	return currencies, err
}

func CreateCurrency(currency *models.Currency) error {

	return db.GetConn().
		Create(currency).Error
}

func UpdateCurrency(currency models.Currency) error {

	return db.GetConn().
		Model(&models.Currency{}).
		Where("id = ?", currency.ID).
		Updates(currency).Error
}

func DeleteCurrency(id int) error {

	return db.GetConn().
		Delete(&models.Currency{}, id).Error
}
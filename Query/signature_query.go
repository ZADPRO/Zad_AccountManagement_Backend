package Query

import (
	db "invoice-backend/DB"
	models "invoice-backend/Models"
)

func GetAllSignatureAuthorities() ([]models.SignatureAuthority, error) {

	var authorities []models.SignatureAuthority

	err := db.GetConn().
		Order("id DESC").
		Find(&authorities).Error

	return authorities, err
}

func CreateSignatureAuthority(
	authority *models.SignatureAuthority,
) error {

	return db.GetConn().
		Create(authority).Error
}

func UpdateSignatureAuthority(authority models.SignatureAuthority) error {

	return db.GetConn().
		Model(&models.SignatureAuthority{}).
		Where("id = ?", authority.ID).
		Select("name", "designation", "contact_number", "email", "signature_url").
		Updates(authority).Error
}

func DeleteSignatureAuthority(id int) error {

	return db.GetConn().
		Delete(&models.SignatureAuthority{}, id).Error
}
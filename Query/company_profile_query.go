package Query

import (
	db "invoice-backend/DB"
	models "invoice-backend/Models"
)

func GetAllCompanyProfiles() ([]models.CompanyProfile, error) {

	var profiles []models.CompanyProfile

	err := db.GetConn().
		Order("id DESC").
		Find(&profiles).Error

	return profiles, err
}

func CreateCompanyProfile(
	profile *models.CompanyProfile,
) error {

	return db.GetConn().
		Create(profile).Error
}

func UpdateCompanyProfile(
	profile models.CompanyProfile,
) error {

	return db.GetConn().
		Model(&models.CompanyProfile{}).
		Where("id = ?", profile.ID).
		Updates(profile).Error
}

func DeleteCompanyProfile(id int) error {

	return db.GetConn().
		Delete(&models.CompanyProfile{}, id).Error
}
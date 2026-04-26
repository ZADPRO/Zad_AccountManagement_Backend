package UserValidation

import (
	"fmt"
	"gorm.io/gorm"
)

// UserValidationQueryResp handles the boolean result of the check
type UserValidationQueryResp struct {
	IsValid bool `gorm:"column:is_valid"`
}

// GetUserValidation uses your specific 'active_users' view.
// It simply checks if the userid exists within that view.
var GetUserValidation = `
SELECT 
    EXISTS (
        SELECT 1 
        FROM active_users 
        WHERE userid = ?
    ) AS is_valid`

// ValidateUser checks if the user currently exists in the active_users view.
func ValidateUser(db *gorm.DB, userId interface{}) (bool, error) {
	var result UserValidationQueryResp

	// Use raw SQL to query the view
	err := db.Raw(GetUserValidation, userId).Scan(&result).Error
	if err != nil {
		return false, fmt.Errorf("database error during user validation: %v", err)
	}

	return result.IsValid, nil
}
package Services

import (
	"errors"
	"fmt"
	"invoice-backend/Query"
	"golang.org/x/crypto/bcrypt" 
	"invoice-backend/Helper/Utils"
	"gorm.io/gorm"
	"database/sql" 
)

// VerifyLogin handles the authentication logic
func VerifyLogin(db *gorm.DB, email string, password string) (int, string, string, bool, error) {
	var userID int
	var hashedPassword string
	var isActive bool
	var roleName string
	var userName string
	var isFirstLogin bool

	// 1. DB Encryption Layer:
	encEmail, err := Utils.EncryptForDB(email)
if err != nil {
	fmt.Printf("DEBUG: EncryptForDB failed: %v\n", err)
    return 0, "", "", false, errors.New("email processing failed")
}
// 2. Database Lookup:
row := db.Raw(Query.GetUserByEmailQuery, encEmail).Row()
err = row.Scan(          
    &userID,
    &hashedPassword,
    &isActive,
    &roleName,
    &userName,
    &isFirstLogin,
)
	if err != nil {
		return 0, "", "", false, errors.New("invalid email or password")
	}


	// 🚫 Check if user is active
	if !isActive {
		return 0, "", "", false, errors.New("account is deactivated")
	}
	
	
	// 4. Password Verification:
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return 0, "", "", false, errors.New("invalid email or password")
	}
	
	return userID, roleName, userName, isFirstLogin, nil
}

// ChangePassword updates user credentials and clears the first login flag
func ChangePassword(db *sql.DB, userID int, newPassword string) error {
	// Generate a new hash with a cost factor of 14.
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 14)
    if err != nil {
        
        return fmt.Errorf("failed to hash password: %w", err)
    }

    
    result, err := db.Exec(Query.ResetPasswordQuery, string(hashedPassword), userID)
    if err != nil {
        
        return fmt.Errorf("database update failed: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        
        return err
    }
    

    if rowsAffected == 0 {
        return errors.New("user not found or password remained the same")
    }

    return nil
}
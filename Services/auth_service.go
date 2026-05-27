package Services

import (
	"errors"
	"strings"
	"fmt"

	"invoice-backend/Query"
	"golang.org/x/crypto/bcrypt" 
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

	// normalize
	email = strings.TrimSpace(strings.ToLower(email))


	// DB lookup using plain searchable email
	row := db.Raw(Query.GetUserByEmailQuery, email).Row()

	err := row.Scan(
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

	

	// active check
	if !isActive {
		return 0, "", "", false, errors.New("account deactivated")
	}



	// password compare
	err = bcrypt.CompareHashAndPassword(
    []byte(hashedPassword),
    []byte(strings.TrimSpace(password)),
)



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

    
    result, err := db.Exec(Query.ResetPasswordQuery, string(hashedPassword),newPassword, userID)
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
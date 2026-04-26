package Services

import (
	"errors"
	"fmt"
	//"invoice-backend/Helper/Utils" // Updated path based on your folder structure
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

	encEmail, err := Utils.EncryptForDB(email)
if err != nil {
	fmt.Printf("DEBUG: EncryptForDB failed: %v\n", err)
    return 0, "", "", false, errors.New("email processing failed")
}
fmt.Printf("DEBUG: login encEmail=%q\n", encEmail)
fmt.Printf("DEBUG: DB stored email=%q\n", "725c8AVuczllKPhmi4ZLjVh0QgjxcBu/RQpT1CcldtjKHIfb8dGIhlBFO9/z82CnfYq2Gw==")

row := db.Raw(Query.GetUserByEmailQuery, encEmail).Row()
err = row.Scan(          // ← ✅ = not :=
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
	fmt.Printf("DEBUG: hashedPassword len=%d\n", len(hashedPassword))
fmt.Printf("DEBUG: hashedPassword=%q\n", hashedPassword)
fmt.Printf("DEBUG: isActive=%v userID=%d\n", isActive, userID)
	// 🚫 Check if user is active
	if !isActive {
		return 0, "", "", false, errors.New("account is deactivated")
	}
	// Go - temporary debug
	
	// 🔑 Compare password (bcrypt)
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return 0, "", "", false, errors.New("invalid email or password")
	}
	fmt.Printf("DEBUG: hashedPassword len=%d value=%q\n", len(hashedPassword), hashedPassword)
	return userID, roleName, userName, isFirstLogin, nil
}

// ChangePassword updates user credentials and clears the first login flag
func ChangePassword(db *sql.DB, userID int, newPassword string) error {
    fmt.Printf("DEBUG ChangePassword called: userID=%d\n", userID)
    
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 14)
    if err != nil {
        fmt.Printf("DEBUG bcrypt error: %v\n", err)
        return fmt.Errorf("failed to hash password: %w", err)
    }

    fmt.Printf("DEBUG executing query: %s\n", Query.ResetPasswordQuery)
    result, err := db.Exec(Query.ResetPasswordQuery, string(hashedPassword), userID)
    if err != nil {
        fmt.Printf("DEBUG query error: %v\n", err) // ← paste this output
        return fmt.Errorf("database update failed: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        fmt.Printf("DEBUG rowsAffected error: %v\n", err)
        return err
    }
    fmt.Printf("DEBUG rowsAffected: %d\n", rowsAffected)

    if rowsAffected == 0 {
        return errors.New("user not found or password remained the same")
    }

    return nil
}
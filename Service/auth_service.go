package Service

import (
	"database/sql"
	"errors"
	"invoice-backend/Query"
  
	"golang.org/x/crypto/bcrypt" 
    
    "invoice-backend/Utils" 
    "fmt" 

)

// Updated to return (int, string, error) -> (userID, roleName, error)

// 1. Changed parameter name from username to email for clarity
func VerifyLogin(db *sql.DB, email string, password string) (int, string, string, bool, error) {
    var userID int
    var hashedPassword string
    var isActive bool
    var roleName string 
    var userName string
    var isFirstLogin bool 

    // 1. Lowercase and Encrypt Deterministically
    encryptedEmail, err := Utils.EncryptDeterministic(email)
    if err != nil {
        return 0, "", "", false, fmt.Errorf("encryption failed: %v", err)
    }

    // 2. Query using the deterministic string
    err = db.QueryRow(Query.GetUserByEmailQuery, encryptedEmail).Scan(
        &userID, 
        &hashedPassword, 
        &isActive, 
        &roleName,
        &userName,
        &isFirstLogin, 
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return 0, "", "", false, errors.New("invalid email or password")
        }
        return 0, "", "", false, err
    }

    if !isActive {
        return 0, "", "", false, errors.New("account is deactivated")
    }

    // 3. Password check (Bcrypt remains the same)
    err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    if err != nil {
        return 0, "", "", false, errors.New("invalid email or password")
    }

    return userID, roleName, userName, isFirstLogin, nil 
}


func ChangePassword(db *sql.DB, userID int, newPassword string) error {
    // 🔐 hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    // 🔥 IMPORTANT: set is_first_login = false
    query := `
        UPDATE users
        SET password = $1,
            is_first_login = false
        WHERE userid = $2
    `

    result, err := db.Exec(query, hashedPassword, userID)
    if err != nil {
        return err
    }

    rows, _ := result.RowsAffected()
    if rows == 0 {
        return errors.New("user not found")
    }

    return nil
} 

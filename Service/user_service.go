package Service

import (
	"database/sql"
	"fmt"
	"encoding/hex"
	"crypto/rand"
	"invoice-backend/Model"
	"invoice-backend/Query"
	"invoice-backend/Utils" // 👈 Using your centralized Utils package
	"golang.org/x/crypto/bcrypt"
)

func AddNewUser(db *sql.DB, req Model.CreateUserRequest) (int, error) {
	// 1. Generate temp password
	b := make([]byte, 4)
	rand.Read(b)
	tempPassword := hex.EncodeToString(b)

	// 2. Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(tempPassword), 14)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	// 🔐 3. Encrypt sensitive data using Utils (Matching Client logic)
	encFirstName, err := Utils.Encrypt(req.FirstName)
	if err != nil {
		return 0, fmt.Errorf("first name encryption failed: %w", err)
	}

	encLastName, err := Utils.Encrypt(req.LastName)
	if err != nil {
		return 0, fmt.Errorf("last name encryption failed: %w", err)
	}

	encEmail, err := Utils.EncryptDeterministic(req.Email)
	if err != nil {
		return 0, fmt.Errorf("email encryption failed: %w", err)
	}

	var newID int
	// 4. Execute Query with encrypted values
	err = db.QueryRow(
		Query.CreateUserQuery,
		req.UserCode,
		req.Username,
		string(hashedPassword),
		encFirstName, // 🔒
		encLastName,  // 🔒
		req.RoleID,
		encEmail,     // 🔒
		true,
	).Scan(&newID)

	if err != nil {
		return 0, fmt.Errorf("database insert failed: %w", err)
	}

	// 5. Email temp password
	_ = Utils.SendWelcomeEmail(req.Email, tempPassword)

	return newID, nil
}

func UpdateUser(db *sql.DB, userID int, req Model.UpdateUserRequest) error {
	// 🔐 Encrypt before update
	encFirstName, _ := Utils.Encrypt(req.FirstName)
	encLastName, _ := Utils.Encrypt(req.LastName)
	encEmail, _ := Utils.EncryptDeterministic(req.Email)

	_, err := db.Exec(
		Query.UpdateUserQuery,
		req.Username,
		encFirstName,
		encLastName,
		req.RoleID,
		encEmail,
		userID,
	)

	return err
}

func GetUserByID(db *sql.DB, userID int) (Model.UserData, error) {
	var user Model.UserData

	err := db.QueryRow(Query.GetUserByIDQuery, userID).Scan(
		&user.UserID,
		&user.UserCode,
		&user.Username,
		&user.FirstName, // This will be the encrypted string from DB
		&user.LastName,  // This will be the encrypted string from DB
		&user.RoleID,
		&user.Email,     // This will be the encrypted string from DB
	)

	if err != nil {
		return user, err
	}

	// 🔓 Decrypt back to readable text for the frontend
	user.FirstName, _ = Utils.Decrypt(user.FirstName)
	user.LastName, _ = Utils.Decrypt(user.LastName)
	user.Email, _ = Utils.Decrypt(user.Email)

	return user, nil
}

func DeleteUser(db *sql.DB, userID int, adminID int) error { 
    // $1 = adminID (deletedby), $2 = userID
    _, err := db.Exec(Query.DeleteUserQuery, adminID, userID)
    return err
} 

func GetUserProfile(db *sql.DB, userID int) (Model.ProfileResponse, error) {
	var profile Model.ProfileResponse

	err := db.QueryRow(Query.GetUserProfileByIDQuery, userID).Scan(
		&profile.User.FirstName, // Encrypted in DB
		&profile.User.LastName,  // Encrypted in DB
		&profile.User.RoleName,
	)

	if err == nil {
		// 🔓 Decrypt for profile view
		profile.User.FirstName, _ = Utils.Decrypt(profile.User.FirstName)
		profile.User.LastName, _ = Utils.Decrypt(profile.User.LastName)
	}

	return profile, err
}



func UpdateUserPassword(db *sql.DB, userID int, newPassword string) error {
    // Hash the password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    // THE FIX: Reset is_first_login to 0 (false)
    query := `UPDATE users SET password = $1, is_first_login = 0 WHERE user_id = $2`
    
    _, err = db.Exec(query, string(hashedPassword), userID)
    return err
}
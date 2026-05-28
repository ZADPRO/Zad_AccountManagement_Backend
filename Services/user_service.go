package Services

import (
	"crypto/rand"
	"database/sql" 
	"encoding/hex"
	"fmt"

	"invoice-backend/Models/dto" 
	"invoice-backend/Models/responses" 
	"invoice-backend/Query"
	"github.com/lib/pq" 
	
	"invoice-backend/Helper/Utils"
	"golang.org/x/crypto/bcrypt"
)

// AddNewUser handles: Temp Password -> Hashing -> Static Encryption -> DB Insert
func AddNewUser(db *sql.DB, req dto.CreateUserRequest) (int, error) {
    // 1. Generate temp password
    b := make([]byte, 4)
    if _, err := rand.Read(b); err != nil {
        return 0, fmt.Errorf("failed to generate random bytes: %w", err)
    }
    tempPassword := hex.EncodeToString(b)
    fmt.Printf("DEBUG: tempPassword=%q for %q\n", tempPassword, req.Email)

    // 2. Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(tempPassword), 14)
    if err != nil {
        return 0, fmt.Errorf("failed to hash password: %w", err)
    }
	encryptedEmail, err := Utils.EncryptForDB(req.Email)
if err != nil {
    return 0, fmt.Errorf("failed to encrypt email: %w", err)
}
   
	
    // 4. Insert
    var newID int
    err = db.QueryRow(
    Query.CreateUserQuery,
    req.UserCode,
    req.Username,
    string(hashedPassword),
    tempPassword,
    req.FirstName,
    req.LastName,
    req.RoleID,
    encryptedEmail,
    req.Email,
    true,
).Scan(&newID)

    if err != nil {
   

    if pqErr, ok := err.(*pq.Error); ok {
        if pqErr.Code == "23505" {
            return 0, fmt.Errorf("email already exists")
        }
    }

    return 0, fmt.Errorf("database insert failed: %w", err)
}

   

    // 5. Send welcome email — only once, log error but don't fail
    if err := Utils.SendWelcomeEmail(req.Email, tempPassword); err != nil {
        fmt.Printf("DEBUG SendWelcomeEmail error: %v\n", err)
    } else {
        fmt.Printf("DEBUG SendWelcomeEmail: sent successfully\n")
    }

    return newID, nil
} 

// UpdateUser handles: Static Encryption -> DB Update
func UpdateUser(db *sql.DB, userID int, req dto.UpdateUserRequest) error {

	encryptedEmail, err := Utils.EncryptForDB(req.Email)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		Query.UpdateUserQuery,
		req.Username,
		req.FirstName,
		req.LastName,
		req.RoleID,
		encryptedEmail,
		req.Email,
		userID,
	)

	return err
}

// GetUserByID handles: DB Fetch -> Static Decryption
func GetUserByID(db *sql.DB, userID int) (responses.UserData, error) {
	var user responses.UserData

	err := db.QueryRow(Query.GetUserByIDQuery, userID).Scan(
		&user.UserID,
		&user.UserCode,
		&user.Username,
		&user.FirstName, // Encrypted string from DB
		&user.LastName,  // Encrypted string from DB
		&user.RoleID,
		&user.Email,     // Encrypted string from DB
	)
	
	if err != nil {
		return user, err
	}

	fmt.Println("GET USER SUCCESS:", user)
	

	return user, nil
}

// GetUserProfile retrieves the specific data needed for the header/profile page
func GetUserProfile(db *sql.DB, userID int) (responses.ProfileResponse, error) {
	var profile responses.ProfileResponse

	err := db.QueryRow(Query.GetUserProfileByIDQuery, userID).Scan(
		&profile.User.FirstName,
		&profile.User.LastName,
		&profile.User.RoleName,
	)

	if err == nil {
		// 🔓 Decrypt names for visual display
		
	}

	return profile, err
}

// FetchAllUsers retrieves the entire user list (Decrypted for the React Table)
func FetchAllUsers(db *sql.DB) ([]responses.UserData, error) {
	rows, err := db.Query(Query.GetAllUsersQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []responses.UserData
	for rows.Next() {
		var u responses.UserData
		err := rows.Scan(
	&u.UserID,
	&u.UserCode,
	&u.Username,
	&u.FirstName,
	&u.LastName,
	&u.RoleID,
	&u.Email,
)
		if err != nil {
			return nil, err
		}

		// 🔓 Decrypt each user so the table shows readable text
		
		

		users = append(users, u)
	}

	return users, nil
}

// UpdateUserPassword handles the force-reset logic
func UpdateUserPassword(db *sql.DB, userID int, newPassword string) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		Query.ResetPasswordQuery,
		string(hashedPassword),
		newPassword,
		userID,
	)

	return err
}

// DeleteUser performs a soft delete via the admin ID
func DeleteUser(db *sql.DB, userID int, adminID int) error { 
	_, err := db.Exec(Query.DeleteUserQuery, adminID, userID)
	return err
}
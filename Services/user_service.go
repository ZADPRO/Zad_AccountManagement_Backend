package Services

import (
	"crypto/rand"
	"database/sql" 
	"encoding/hex"
	"fmt"
	"invoice-backend/Helper/Utils" // Ensure this matches your package path
	"invoice-backend/Models/dto" 
	"invoice-backend/Models/responses" 
	"invoice-backend/Query"
	

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

    // 3. Encrypt for DB — handle each error explicitly
    encFirstName, err := Utils.EncryptForDB(req.FirstName)
    if err != nil {
        return 0, fmt.Errorf("firstName encryption failed: %w", err)
    }
    encLastName, err := Utils.EncryptForDB(req.LastName)
    if err != nil {
        return 0, fmt.Errorf("lastName encryption failed: %w", err)
    }
    encEmail, err := Utils.EncryptForDB(req.Email)
    if err != nil {
        return 0, fmt.Errorf("email encryption failed: %w", err)
    }

    fmt.Printf("DEBUG: inserting user code=%q username=%q roleID=%d\n", 
        req.UserCode, req.Username, req.RoleID)

    // 4. Insert
    var newID int
    err = db.QueryRow(
        Query.CreateUserQuery,
        req.UserCode,
        req.Username,
        string(hashedPassword),
        encFirstName,
        encLastName,
        req.RoleID,
        encEmail,
        true,
    ).Scan(&newID)

    if err != nil {
        fmt.Printf("DEBUG DB insert error: %v\n", err) // ← this will show the real cause
        return 0, fmt.Errorf("database insert failed: %w", err)
    }

    fmt.Printf("DEBUG: user inserted with ID=%d, sending email to %q\n", newID, req.Email)

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
	// 🔐 Encrypt updated fields before they touch the DB
	encFirstName, _ := Utils.EncryptForDB(req.FirstName)
	encLastName, _ := Utils.EncryptForDB(req.LastName)
	encEmail, _ := Utils.EncryptForDB(req.Email)

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

	// 🔓 Decrypt back to plain text so the Controller can re-encrypt it for Axios
	user.FirstName, _ = Utils.DecryptFromDB(user.FirstName)
	user.LastName, _ = Utils.DecryptFromDB(user.LastName)
	user.Email, _ = Utils.DecryptFromDB(user.Email)

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
		profile.User.FirstName, _ = Utils.DecryptFromDB(profile.User.FirstName)
		profile.User.LastName, _ = Utils.DecryptFromDB(profile.User.LastName)
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
			
		)
		if err != nil {
			return nil, err
		}

		// 🔓 Decrypt each user so the table shows readable text
		u.FirstName, _ = Utils.DecryptFromDB(u.FirstName)
		u.LastName, _ = Utils.DecryptFromDB(u.LastName)
		

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

	// Note: We also flip is_first_login to 0/false here
	query := `UPDATE users SET password = $1, is_first_login = 0, updatedat = NOW() WHERE user_id = $2`
	
	_, err = db.Exec(query, string(hashedPassword), userID)
	return err
}

// DeleteUser performs a soft delete via the admin ID
func DeleteUser(db *sql.DB, userID int, adminID int) error { 
	_, err := db.Exec(Query.DeleteUserQuery, adminID, userID)
	return err
}
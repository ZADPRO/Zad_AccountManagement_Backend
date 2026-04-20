package Service

import (
	"database/sql"
	"fmt"
	"invoice-backend/Model"
	"invoice-backend/Query"
	"invoice-backend/Utils" // Ensure this is imported
	"time"
)

// GetAllClients must start with a Capital letter to be exported!
func GetAllClients(db *sql.DB) ([]Model.ClientListModel, error) {

	query := `
	SELECT "clientid", "clientcode", "name", "businessname", "isactive"
	FROM "active_clients"
	WHERE "isactive" IS TRUE
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []Model.ClientListModel

	for rows.Next() {
		var c Model.ClientListModel

		err := rows.Scan(
			&c.ClientID,
			&c.ClientCode,
			&c.Name,
			&c.BusinessName,
			&c.IsActive,
		)
		if err != nil {
			return nil, err
		}

		clients = append(clients, c)
	}

	// Important: check for iteration errors
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return clients, nil
}
// -----------------------------
// ✅ CREATE CLIENT (MAIN LOGIC)
// -----------------------------
func CreateClient(db *sql.DB, req Model.CreateClientRequest) (int, error) {

	// Optional: set system fields here instead of controller
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()
	req.CreatedBy = 1 // Ideally from auth
	req.UpdatedBy = 2       // Ideally from auth

	// 🔐 Encrypt sensitive data
	encryptedGST, err := Utils.EncryptDeterministic(req.GSTNumber)
	if err != nil {
		return 0, fmt.Errorf("GST encryption failed: %w", err)
	}

	encryptedPAN, err := Utils.EncryptDeterministic(req.PAN)
	if err != nil {
		return 0, fmt.Errorf("PAN encryption failed: %w", err)
	}

	// 🔄 Start transaction
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	// 🔥 Safe rollback
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// -----------------------------
	// 1. Insert into clientinformation
	// -----------------------------
	var clientID int

	err = tx.QueryRow(Query.CreateClientInfoQuery,
		req.ClientCode,
		req.Name,
		req.BusinessName,
		req.SupplyTypeID,
		req.Email,
		req.PrimaryNumber,
		req.Address,
		req.CountryName,
		req.StateName,
		req.ZIP,
		req.ClientType,
		req.UpdatedBy,
	).Scan(&clientID)

	if err != nil {
		return 0, fmt.Errorf("client info insert failed: %w", err)
	}

	// -----------------------------
	// 2. Insert into clienttaxdetails
	_, err = tx.Exec(Query.CreateClientTaxQuery,
	clientID,        // $1
	encryptedGST,    // $2
	encryptedPAN,    // $3
	req.IsExport,    // $4
	req.UpdatedBy,   // $5
	req.GSTStatus,   // $6
)
	

	if err != nil {
		return 0, fmt.Errorf("tax details insert failed: %w", err)
	}

	// -----------------------------
	// ✅ Commit
	// -----------------------------
	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit failed: %w", err)
	}

	return clientID, nil
}

func UpdateClient(db *sql.DB, clientID int, req Model.CreateClientRequest) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	encryptedGST, err := Utils.EncryptDeterministic(req.GSTNumber)
	if err != nil {
		return err
	}
	encryptedPAN, err := Utils.EncryptDeterministic(req.PAN)
	if err != nil {
		return err
	}

	// 1. Update clientinformation (Merged Address updates)
	_, err = tx.Exec(Query.UpdateClientInfoQuery, 
		req.Name, req.BusinessName, req.SupplyTypeID, req.Email, 
		req.PrimaryNumber, req.Address, req.CountryName, req.StateName, req.ZIP, 
		clientID,
	)
	if err != nil {
		return err
	}

	// 2. Update clienttaxdetails
	_, err = tx.Exec(Query.UpdateClientTaxQuery, req.GSTStatus, encryptedGST, encryptedPAN, req.IsExport, clientID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func DeleteClient(db *sql.DB, clientID int, adminID int) error {
	const query = `
        UPDATE clientinformation 
        SET isactive = false, updatedat = NOW(), deletedat = NOW(), updatedby = $1 
        WHERE clientid = $2;`

	result, err := db.Exec(query, adminID, clientID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("client with ID %d not found", clientID)
	}
	return nil
}

func GetClientByID(db *sql.DB, clientID int) (Model.ClientDetailsResponse, error) {
    var client Model.ClientDetailsResponse
    
    // Temporary variables for columns we don't need in the response struct
    var updatedAt time.Time
    var updatedBy int

    err := db.QueryRow(Query.GetClientByIDQuery, clientID).Scan(
        &client.ClientID,
        &client.ClientCode,
        &client.Name,
        &client.BusinessName,
        &client.SupplyTypeID,
        &client.IsActive,
        &client.ClientType,
        &updatedAt,     // Added to match index 8
        &updatedBy,     // Added to match index 9
        &client.Email,
        &client.MobileNumber,
        &client.RegisteredAddress,
        &client.CountryName,
        &client.StateName,
        &client.ZIP,
        &client.GSTNumber,
        &client.PAN,
        &client.IsExport,
        &client.GSTStatus,
    )

    if err != nil {
        return client, err
    }

    // 🔐 Decrypt sensitive data before returning to frontend
    // Since you encrypted them in Create/Update, you MUST decrypt them here
    client.GSTNumber, _ = Utils.Decrypt(client.GSTNumber)
    client.PAN, _ = Utils.Decrypt(client.PAN)

    return client, nil
}
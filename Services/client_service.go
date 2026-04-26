package Services

import (
	"database/sql"
	"fmt"
	"invoice-backend/Helper/Utils" // Ensure this matches your path (usually Utils or Helper/Utils)
	"invoice-backend/Models/dto"
	"invoice-backend/Models/responses"
	"invoice-backend/Query"
	"time"
)

// GetAllClients retrieves all active clients
func GetAllClients(db *sql.DB) ([]responses.ClientListModel, error) {
	rows, err := db.Query(Query.GetAllClientsQuery) // Using the constant from your Query package
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []responses.ClientListModel
	for rows.Next() {
		var c responses.ClientListModel
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

	return clients, nil
}

// CreateClient handles the business logic for adding a new client
func CreateClient(db *sql.DB, req dto.CreateClientRequest, adminID int) (int, error) {
    encryptedGST, err := Utils.EncryptForDB(req.GSTNumber)
    if err != nil {
        return 0, fmt.Errorf("GST encryption failed: %w", err)
    }

    encryptedPAN, err := Utils.EncryptForDB(req.PAN)
    if err != nil {
        return 0, fmt.Errorf("PAN encryption failed: %w", err)
    }

    tx, err := db.Begin()
    if err != nil {
        return 0, err
    }

    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()

    // 3. Insert into clientinformation
    var clientID int
    err = tx.QueryRow(Query.CreateClientInfoQuery,
        req.ClientCode,    // $1
        req.Name,          // $2
        req.BusinessName,  // $3
        req.SupplyTypeID,  // $4
        req.Email,         // $5
        req.PrimaryNumber, // $6
        req.Address,       // $7
        req.CountryName,   // $8
        req.StateName,     // $9
        req.ZIP,           // $10
        req.ClientType,    // $11
        adminID,           // $12
    ).Scan(&clientID)

    if err != nil {
        return 0, fmt.Errorf("client info insert failed: %w", err)
    }

    // ✅ Convert 0 → nil for optional foreign key
    var billingStateID *int
    if req.BillingStateID != nil && *req.BillingStateID != 0 {
        billingStateID = req.BillingStateID
    }

    // 4. Insert into clienttaxdetails
    _, err = tx.Exec(Query.CreateClientTaxQuery,
        clientID,             // $1
        encryptedGST,         // $2
        encryptedPAN,         // $3
        req.TaxPercentage,    // $4
        req.IsExport,         // $5
        req.BillingAddress,   // $6
        req.BillingCountryID, // $7
        billingStateID,       // $8 ← nil when not India
        adminID,              // $9
        req.GSTStatus,        // $10
    )

    if err != nil {
        return 0, fmt.Errorf("tax details insert failed: %w", err)
    }

    if err = tx.Commit(); err != nil {
        return 0, fmt.Errorf("commit failed: %w", err)
    }

    return clientID, nil
}

// UpdateClient handles modifications to existing client data
func UpdateClient(db *sql.DB, clientID int, req dto.CreateClientRequest, adminID int) error {
    tx, err := db.Begin()
    if err != nil {
        return err
    }

    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()

    // 🔐 Encrypt sensitive fields
    encryptedGST, err := Utils.EncryptForDB(req.GSTNumber)
    if err != nil {
        return fmt.Errorf("GST encryption failed: %w", err)
    }

    encryptedPAN, err := Utils.EncryptForDB(req.PAN)
    if err != nil {
        return fmt.Errorf("PAN encryption failed: %w", err)
    }

    // ✅ Convert 0 → nil for optional billing state FK
    var billingStateID *int
    if req.BillingStateID != nil && *req.BillingStateID != 0 {
        billingStateID = req.BillingStateID
    }

    // 1. Update clientinformation
    _, err = tx.Exec(Query.UpdateClientInfoQuery,
        req.Name,          // $1  name
        req.BusinessName,  // $2  businessname
        req.SupplyTypeID,  // $3  supplytypeid
        req.Email,         // $4  email
        req.PrimaryNumber, // $5  mobilenumber
        req.Address,       // $6  registeredaddress
        req.CountryName,   // $7  countryname
        req.StateName,     // $8  statename
        req.ZIP,           // $9  zip
        req.ClientType,    // $10 clienttype
        adminID,           // $11 updatedby
        clientID,          // $12 WHERE clientid
    )
    if err != nil {
        return fmt.Errorf("client info update failed: %w", err)
    }

    // 2. Update clienttaxdetails
    _, err = tx.Exec(Query.UpdateClientTaxQuery,
        req.GSTStatus,        // $1  gststatus
        encryptedGST,         // $2  gstnumber
        encryptedPAN,         // $3  pan
        req.TaxPercentage,    // $4  taxpercentage
        req.IsExport,         // $5  isexport
        req.BillingAddress,   // $6  billingaddress
        req.BillingCountryID, // $7  billingcountryid
        billingStateID,       // $8  billingstateid (nil if not India)
        adminID,              // $9  updatedby
        clientID,             // $10 WHERE clientid
    )
    if err != nil {
        return fmt.Errorf("tax details update failed: %w", err)
    }

    return tx.Commit()
}

// GetClientByID fetches a single client and DECRYPTS the static data
func GetClientByID(db *sql.DB, clientID int) (responses.ClientDetailsResponse, error) {
	var client responses.ClientDetailsResponse
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
    &updatedAt,        // correct
    &updatedBy,        // correct
    &client.Email,
    &client.MobileNumber,
    &client.RegisteredAddress,
    &client.CountryName,
    &client.StateName,
    &client.ZIP,
    &client.BillingAddress,
    &client.BillingCountryID,
    &client.BillingStateID,
    &client.TaxPercentage,
    &client.GSTNumber,
    &client.PAN,
    &client.IsExport,
    &client.GSTStatus,
    &client.BillingCountryName,
    &client.BillingStateName,
)
	if err != nil {
		return client, err
	}

	// 🔐 6. Static Decryption
	// We decrypt the DB storage layer so we can then re-encrypt it 
	// with the dynamic token in the controller.
	client.GSTNumber, _ = Utils.DecryptFromDB(client.GSTNumber)
	client.PAN, _ = Utils.DecryptFromDB(client.PAN)

	return client, nil
}

// DeleteClient performs a soft delete
func DeleteClient(db *sql.DB, clientID int, adminID int) error {
	result, err := db.Exec(Query.DeleteClientQuery, adminID, clientID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("client with ID %d not found", clientID)
	}
	return nil
}
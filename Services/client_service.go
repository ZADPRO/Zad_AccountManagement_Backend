package Services

import (
	"database/sql"
	"fmt"
	"invoice-backend/Helper/Utils" 
	"invoice-backend/Models/dto"
	"invoice-backend/Models/responses"
	"invoice-backend/Query"
	"time" 
    "strings" 
    "strconv"
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
    var encryptedGST, encryptedPAN string
    var err error

    
    if req.GSTNumber != "" {
        encryptedGST, err = Utils.EncryptForDB(req.GSTNumber)
        if err != nil {
            return 0, fmt.Errorf("GST encryption failed: %w", err)
        }
    }

    // Only encrypt PAN if it's not empty
    if req.PAN != "" {
        encryptedPAN, err = Utils.EncryptForDB(req.PAN)
        if err != nil {
            return 0, fmt.Errorf("PAN encryption failed: %w", err)
        }
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
	req.ClientCode,
	req.Name,
	req.BusinessName,
	req.SupplyTypeID,
	toNullString(req.Email),
	req.PrimaryNumber,
	toNullString(req.Address),
	toNullString(req.CountryName),
	toNullString(req.StateName),
	toNullInt(req.ZIP),
	req.ClientType,
	adminID,
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
    var encryptedGST, encryptedPAN string

if req.GSTNumber != "" {
    encryptedGST, err = Utils.EncryptForDB(req.GSTNumber)
    if err != nil {
        return fmt.Errorf("GST encryption failed: %w", err)
    }
}

if req.PAN != "" {
    encryptedPAN, err = Utils.EncryptForDB(req.PAN)
    if err != nil {
        return fmt.Errorf("PAN encryption failed: %w", err)
    }
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

	// Nullable fields
	var email, mobileNumber, registeredAddress, countryName sql.NullString
    var stateName, billingAddress, gstNumber, pan sql.NullString
    var billingCountryName, billingStateName, gstStatus sql.NullString

    var zip sql.NullInt64
    var billingStateID sql.NullInt64
    var billingCountryID sql.NullInt64

	err := db.QueryRow(Query.GetClientByIDQuery, clientID).Scan(
	&client.ClientID,
	&client.ClientCode,
	&client.Name,
	&client.BusinessName,
	&client.SupplyTypeID,
	&client.IsActive,
	&client.ClientType,
	&updatedAt,
	&updatedBy,
	&email,
	&mobileNumber,
	&registeredAddress,
	&countryName,
	&stateName,
	&zip,
	&billingAddress,
	&billingCountryID,
	&billingStateID,
	&client.TaxPercentage,
	&gstNumber,
	&pan,
	&client.IsExport,
	&gstStatus,
	&billingCountryName,
	&billingStateName,
)
	

	if err != nil {
		return client, err
	}


if email.Valid {
	client.Email = email.String
}

if mobileNumber.Valid {
	client.MobileNumber = mobileNumber.String
}

if registeredAddress.Valid {
	client.RegisteredAddress = registeredAddress.String
}

if countryName.Valid {
	client.CountryName = countryName.String
}

if stateName.Valid {
	client.StateName = stateName.String
}

if zip.Valid {
	client.ZIP = int(zip.Int64)
}

if billingAddress.Valid {
	client.BillingAddress = billingAddress.String
}

if billingCountryID.Valid {
	client.BillingCountryID = int(billingCountryID.Int64)
}

if billingStateID.Valid {
	client.BillingStateID = int(billingStateID.Int64)
}

if gstNumber.Valid {
	client.GSTNumber = gstNumber.String
}

if pan.Valid {
	client.PAN = pan.String
}

if gstStatus.Valid {
	client.GSTStatus = gstStatus.String
}

if billingCountryName.Valid {
	client.BillingCountryName = billingCountryName.String
}

if billingStateName.Valid {
	client.BillingStateName = billingStateName.String
}

	// 🔐 Decrypt only if present
	if client.GSTNumber != "" {
		client.GSTNumber, _ = Utils.DecryptFromDB(client.GSTNumber)
	}
	if client.PAN != "" {
		client.PAN, _ = Utils.DecryptFromDB(client.PAN)
	}

	// ✅ Fix Updated fields
	client.UpdatedAt = updatedAt.Format(time.RFC3339)
	client.UpdatedBy = strconv.Itoa(updatedBy)

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


func toNullString(s string) interface{} {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
}

func toNullInt(i int) interface{} {
	if i == 0 {
		return nil
	}
	return i
}
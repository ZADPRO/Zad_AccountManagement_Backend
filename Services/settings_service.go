package Services 

import ( 
	
	"database/sql" 
	"fmt"
	"invoice-backend/Models/dto" 
	"invoice-backend/Helper/Utils" 
    "invoice-backend/Models/responses"
	
	"invoice-backend/Query"

	
) 

// SaveBankingDetails handles the Upsert logic with encryption
func CreateBankingDetails(db *sql.DB, req dto.SaveBankingRequest, adminID int) (int, error) {
	// 🔐 Encrypt
	encAcc, _ := Utils.EncryptForDB(req.AccountNumber)
	encIFSC, _ := Utils.EncryptForDB(req.IfscCode)
	encSwift, _ := Utils.EncryptForDB(req.SwiftCode)

	var newID int
	err := db.QueryRow(
		Query.InsertBankingDetailsQuery,
		req.BankName,
		encAcc,
		encIFSC,
		req.BankAddress,
		req.LogoURL,
		req.AccountType,
		encSwift,
		req.UserID, // owner
		adminID,    // CreatedBy
	).Scan(&newID)

	if err != nil {
		return 0, fmt.Errorf("insert failed: %w", err)
	}

	return newID, nil
}

func UpdateBankingDetails(db *sql.DB, req dto.SaveBankingRequest, adminID int) (int, error) {
	// 🔐 Encrypt
	encAcc, _ := Utils.EncryptForDB(req.AccountNumber)
	encIFSC, _ := Utils.EncryptForDB(req.IfscCode)
	encSwift, _ := Utils.EncryptForDB(req.SwiftCode)

	var updatedID int
	err := db.QueryRow(
		Query.UpdateBankingDetailsQuery,
		req.BankName,
		encAcc,
		encIFSC,
		req.BankAddress,
		req.LogoURL,
		req.AccountType,
		encSwift,
		adminID,        // UpdatedBy
		req.DetailsID,  // WHERE DetailsID
	).Scan(&updatedID)

	if err != nil {
		return 0, fmt.Errorf("update failed: %w", err)
	}

	return updatedID, nil
}

func GetAllBankingDetails(db *sql.DB, userID int) ([]responses.BankingDetailsData, error) {
    var banks []responses.BankingDetailsData

    rows, err := db.Query(Query.GetBankingDetailsQuery, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var b responses.BankingDetailsData
        err := rows.Scan(
            &b.DetailsID, &b.BankName, &b.AccountNumber, &b.IfscCode, 
            &b.BankAddress, &b.LogoURL, &b.AccountType, &b.SwiftCode,
        )
        if err != nil {
            return nil, err
        }

        // 🔓 Decrypt sensitive data for the UI
        b.AccountNumber, _ = Utils.DecryptFromDB(b.AccountNumber)
        b.IfscCode, _ = Utils.DecryptFromDB(b.IfscCode)
        b.SwiftCode, _ = Utils.DecryptFromDB(b.SwiftCode)

        banks = append(banks, b)
    }

    // Handle the case where the loop didn't run or encountered an error
    if err = rows.Err(); err != nil {
        return nil, err
    }

    return banks, nil
}

func SoftDeleteBankingDetails(db *sql.DB, detailsID int, userID int) error {
    _, err := db.Exec(Query.DeleteBankingDetailsQuery, userID, detailsID, userID)
    return err
}


// AddCustomField handles encryption and insertion of a new field definition
// AddCustomField saves the definition directly without DB-level encryption
// AddCustomField saves the definition directly to Postgres without DB-level encryption
func AddCustomField(db *sql.DB, req dto.CreateCustomFieldRequest, creatorID int) (int, error) {
    var newID int

    // req.FieldLabel is now passed as plain text directly to the DB
    err := db.QueryRow(
        Query.CreateCustomFieldQuery,
        req.FieldLabel, // <--- Plain Text
        req.FieldType, 
        req.IsRequired,
        creatorID,
    ).Scan(&newID)

    if err != nil {
        return 0, fmt.Errorf("database insert failed: %w", err)
    }

    return newID, nil
}

func FetchAllCustomFields(db *sql.DB) ([]responses.CustomFieldData, error) {
    rows, err := db.Query(Query.GetAllCustomFieldsQuery)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var fields []responses.CustomFieldData
    for rows.Next() {
        var f responses.CustomFieldData
        // Scan the 6 columns defined in your SQL Query
        err := rows.Scan(
            &f.FieldID, 
            &f.FieldLabel, // <--- Scanned as plain text
            &f.FieldType, 
            &f.IsRequired, 
            &f.CreatedAt, 
            &f.CreatedBy,
        )
        if err != nil {
            return nil, err
        }

        // NO decryption needed here. 
        // The Controller will handle the network encryption later.
        fields = append(fields, f)
    }
    return fields, nil
}

// RemoveCustomField performs a soft-delete by setting the DeletedAt timestamp
func RemoveCustomField(db *sql.DB, fieldID int, userID int) error {
    _, err := db.Exec(Query.DeleteCustomFieldQuery, fieldID, userID )
    if err != nil {
        return fmt.Errorf("failed to soft-delete field: %w", err)
    }
    return nil
}


package Services

import (
	"database/sql"
	"fmt"
	"invoice-backend/Models/dto"
    "invoice-backend/Helper/Utils" 
	"invoice-backend/Query"
	"time"
     "encoding/json"
)

// CreateFullInvoice handles the business logic and DB transaction
func CreateFullInvoice(db *sql.DB, req dto.CreateInvoiceRequest) (int, error) {

    

    // Optional: still ensure LineTotal exists (safe fallback)
    for i, item := range req.Items {
        if item.LineTotal == 0 {
            req.Items[i].LineTotal = float64(item.Quantity) * item.UnitPrice
        }
    }

    // ✅ Convert custom fields → JSON
    customJSON, err := json.Marshal(req.CustomValues)
    if err != nil {
        return 0, fmt.Errorf("custom json marshal failed: %w", err)
    }
    fmt.Println(customJSON);
    // ✅ Start transaction
    tx, err := db.Begin()
    if err != nil {
        return 0, err
    }

    // ✅ Insert invoice header
    var newInvoiceID int
    err = tx.QueryRow(
        Query.InsertInvoiceHeaderQuery,
        req.InvoiceNumber,
        req.ClientID,
        req.InvoiceDate,
        req.GrandTotal,
        req.PaymentStatus,
        req.UpdatedBy,
        customJSON,
        req.InvoiceDueDate,
	    req.Currency, 
        req.BankID,
        req.InvoiceType,
    ).Scan(&newInvoiceID)

    if err != nil {
        tx.Rollback()
        return 0, fmt.Errorf("header error: %w", err)
    }

    // ✅ Insert invoice items
    for _, item := range req.Items {
          itemCustomJSON, err := json.Marshal(item.CustomFieldValues)
    if err != nil {
        tx.Rollback()
        return 0, fmt.Errorf("item custom json marshal failed: %w", err)
    }
        _, err = tx.Exec(
            Query.InsertInvoiceItemQuery,
            newInvoiceID,
            item.Description,
            item.Quantity,
            item.UnitPrice,
            item.LineTotal,
            req.UpdatedBy,
            itemCustomJSON,
        )

        if err != nil {
            tx.Rollback()
            return 0, fmt.Errorf("item error: %w", err)
        }
    }

    // ✅ Commit
    if err := tx.Commit(); err != nil {
        return 0, err
    }

    return newInvoiceID, nil
}

func GetInvoiceByID(db *sql.DB, invoiceID int) (*dto.InvoiceResponse, error) {

    var inv dto.InvoiceResponse
    var clientID int
    var customJSON []byte 
    // 🔐 Decrypt Bank Details 
    

    // 1. Fetch invoice header
    err := db.QueryRow(Query.GetInvoiceByIDQuery, invoiceID).Scan(
        &inv.InvoiceID,
        &inv.InvoiceNumber,
        &inv.InvoiceDate,
        &inv.GrandTotal,
        &inv.PaymentStatus,
        &clientID,
        &customJSON,
        &inv.InvoiceDueDate,
        &inv.Currency, 
        &inv.BankID,
        &inv.InvoiceType,

        &inv.InvoiceBankName,
        &inv.InvoiceAccountNumber,
        &inv.InvoiceIFSCCode,
        &inv.InvoiceBankAddress,
        &inv.InvoiceQRCodeURL,
        &inv.InvoiceAccountType,
        &inv.InvoiceSwiftCode,
    )
    
    if err != nil {
        return nil, fmt.Errorf("invoice not found: %w", err)
    }
    if inv.InvoiceAccountNumber != "" {
    decryptedAcc, err := Utils.DecryptFromDB(inv.InvoiceAccountNumber)
    if err == nil {
        inv.InvoiceAccountNumber = decryptedAcc
    } else {
        fmt.Println("Account Number decrypt error:", err)
    }
}

if inv.InvoiceIFSCCode != "" {
    decryptedIFSC, err := Utils.DecryptFromDB(inv.InvoiceIFSCCode)
    if err == nil {
        inv.InvoiceIFSCCode = decryptedIFSC
    } else {
        fmt.Println("IFSC decrypt error:", err)
    }
}

if inv.InvoiceSwiftCode != "" {
    decryptedSwift, err := Utils.DecryptFromDB(inv.InvoiceSwiftCode)
    if err == nil {
        inv.InvoiceSwiftCode = decryptedSwift
    } else {
        fmt.Println("SWIFT decrypt error:", err)
    }
}
   if len(customJSON) > 0 {
    fmt.Println("CUSTOM JSON:", string(customJSON))
    
    inv.CustomValues = []dto.CustomFieldValue{}

  if len(customJSON) > 0 && string(customJSON) != "[]" {
    if err := json.Unmarshal(customJSON, &inv.CustomValues); err != nil {
        fmt.Println("JSON UNMARSHAL ERROR:", err)
        // Don't kill the whole request just for custom fields, just log it
    }
}
   }
    // 2. Fetch client
    var supplyTypeID, updatedBy int
    var isActive bool
    var updatedAt time.Time
    var billingCountryID, billingStateID sql.NullInt64
    var zip sql.NullInt64
    var gstNumber sql.NullString
    var pan sql.NullString
    var stateName sql.NullString
    
    err = db.QueryRow(Query.GetClientByIDQuery, clientID).Scan(
        &inv.Client.ClientID,
        &inv.Client.ClientCode,
        &inv.Client.Name,
        &inv.Client.BusinessName,
        &supplyTypeID,
        &isActive,
        &inv.Client.ClientType,
        &updatedAt,
        &updatedBy,
        &inv.Client.Email,
        &inv.Client.PrimaryNumber,
        &inv.Client.Address,
        &inv.Client.CountryName,
        &stateName,
        &zip,
        &inv.Client.BillingAddress,
        &billingCountryID,
        &billingStateID,
        &inv.Client.TaxPercentage,
        &gstNumber,
        &pan,
        &inv.Client.IsExport,
        &inv.Client.GSTStatus,
        &inv.Client.BillingCountry,
        &inv.Client.BillingState,
    )
    if err != nil {
        return nil, fmt.Errorf("client not found: %w", err)
    }
    if zip.Valid {
    inv.Client.ZIP = int(zip.Int64)
}
    if gstNumber.Valid {
    inv.Client.GSTNumber = gstNumber.String
}

if pan.Valid {
    inv.Client.PAN = pan.String
}

// ✅ Decrypt GST Number
if inv.Client.GSTNumber != "" {
    decryptedGST, decErr := Utils.DecryptFromDB(inv.Client.GSTNumber)
    if decErr == nil {
        inv.Client.GSTNumber = decryptedGST
    } else {
        fmt.Println("GST decrypt error:", decErr)
    }
}

// ✅ Decrypt PAN
if inv.Client.PAN != "" {
    decryptedPAN, decErr := Utils.DecryptFromDB(inv.Client.PAN)
    if decErr == nil {
        inv.Client.PAN = decryptedPAN
    } else {
        fmt.Println("PAN decrypt error:", decErr)
    }
}

    // 3. Fetch line items ✅ FIXED
    rows, err := db.Query(Query.GetInvoiceItemsByInvoiceIDQuery, invoiceID)
    if err != nil {
        return nil, fmt.Errorf("items fetch error: %w", err)
    }
    defer rows.Close()
    
    for rows.Next() {
        
        var item dto.InvoiceItem
        var itemCustomJSON []byte
        err := rows.Scan(
            &item.ItemID,
            &item.Description,
            &item.Quantity,
            &item.UnitPrice,
            &item.LineTotal,
            &itemCustomJSON,
        )
         if err != nil {
        return nil, fmt.Errorf("item scan error: %w", err)
    }
    if len(itemCustomJSON) > 0 && string(itemCustomJSON) != "[]" {
        err := json.Unmarshal(itemCustomJSON, &item.CustomFieldValues)
        if err != nil {
            fmt.Println("item custom JSON error:", err)
        }
    }

    inv.Items = append(inv.Items, item)
        if err != nil {
            return nil, fmt.Errorf("item scan error: %w", err)
        }
    
    }
    return &inv, nil
}
func DeleteInvoice(db *sql.DB, invoiceID int, adminID int) error {

	// 1. Soft delete invoice items
	_, err := db.Exec(Query.DeleteInvoiceItemsQuery, adminID, invoiceID)
	if err != nil {
		return err
	}

	// 2. Soft delete invoice
	result, err := db.Exec(Query.DeleteInvoiceQuery, adminID, invoiceID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("invoice with ID %d not found", invoiceID)
	}

	return nil
}
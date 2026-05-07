package Services

import (
	"database/sql"
	"fmt"
	"invoice-backend/Models/dto"
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
        string(customJSON),
    ).Scan(&newInvoiceID)

    if err != nil {
        tx.Rollback()
        return 0, fmt.Errorf("header error: %w", err)
    }

    // ✅ Insert invoice items
    for _, item := range req.Items {
        _, err = tx.Exec(
            Query.InsertInvoiceItemQuery,
            newInvoiceID,
            item.Description,
            item.Quantity,
            item.UnitPrice,
            item.LineTotal,
            req.UpdatedBy,
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
    // 1. Fetch invoice header
    err := db.QueryRow(Query.GetInvoiceByIDQuery, invoiceID).Scan(
        &inv.InvoiceID,
        &inv.InvoiceNumber,
        &inv.InvoiceDate,
        &inv.GrandTotal,
        &inv.PaymentStatus,
        &clientID,
        &customJSON, 
    )
    if err != nil {
        return nil, fmt.Errorf("invoice not found: %w", err)
    }
   if len(customJSON) > 0 {
    fmt.Println("CUSTOM JSON:", string(customJSON))

    if err := json.Unmarshal(customJSON, &inv.CustomFields); err != nil {
        return nil, fmt.Errorf("custom fields unmarshal failed: %w", err)
    }
}
    // 2. Fetch client
    var supplyTypeID, updatedBy int
    var isActive bool
    var updatedAt time.Time
    var billingCountryID, billingStateID int

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
        &inv.Client.StateName,
        &inv.Client.ZIP,
        &inv.Client.BillingAddress,
        &billingCountryID,
        &billingStateID,
        &inv.Client.TaxPercentage,
        &inv.Client.GSTNumber,
        &inv.Client.PAN,
        &inv.Client.IsExport,
        &inv.Client.GSTStatus,
        &inv.Client.BillingCountry,
        &inv.Client.BillingState,
    )
    if err != nil {
        return nil, fmt.Errorf("client not found: %w", err)
    }

    // 3. Fetch line items ✅ FIXED
    rows, err := db.Query(Query.GetInvoiceItemsByInvoiceIDQuery, invoiceID)
    if err != nil {
        return nil, fmt.Errorf("items fetch error: %w", err)
    }
    defer rows.Close()

    for rows.Next() {
        var item dto.InvoiceItem
        err := rows.Scan(
            &item.ItemID,
            &item.Description,
            &item.Quantity,
            &item.UnitPrice,
            &item.LineTotal,
        )
        if err != nil {
            return nil, fmt.Errorf("item scan error: %w", err)
        }
        inv.Items = append(inv.Items, item)
    }
    return &inv, nil
}
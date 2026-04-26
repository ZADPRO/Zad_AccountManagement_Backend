package Services

import (
	"database/sql"
	"fmt"
	"invoice-backend/Models/dto"
	"invoice-backend/Query"
)

// CreateFullInvoice handles the business logic and DB transaction
func CreateFullInvoice(db *sql.DB, req dto.CreateInvoiceRequest) (int, error) {
	// 1. Business Logic: Re-calculate totals for security
	var calculatedGrandTotal float64
	for i, item := range req.Items {
		lineTotal := float64(item.Quantity) * item.UnitPrice
		req.Items[i].LineTotal = lineTotal
		calculatedGrandTotal += lineTotal
	}
	
	// Security check: Ensure frontend total matches server total
	if calculatedGrandTotal != req.GrandTotal {
		return 0, fmt.Errorf("total mismatch: expected %v, got %v", calculatedGrandTotal, req.GrandTotal)
	}

	// 2. Start Database Transaction
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	// 3. Insert Header
	var newInvoiceID int
	err = tx.QueryRow(Query.InsertInvoiceHeaderQuery,
		req.InvoiceNumber, req.ClientID, req.InvoiceDate, req.GrandTotal, req.PaymentStatus, req.UpdatedBy,
	).Scan(&newInvoiceID)

	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("header error: %w", err)
	}

	// 4. Insert Items
	for _, item := range req.Items {
		_, err = tx.Exec(Query.InsertInvoiceItemQuery,
			newInvoiceID, item.Description, item.Quantity, item.UnitPrice, item.LineTotal, req.UpdatedBy,
		)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("item error: %w", err)
		}
	}

	// 5. Final Commit
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return newInvoiceID, nil
}
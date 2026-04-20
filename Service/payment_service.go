package Service

import (
	"database/sql"
)

// CheckInvoiceStatus determines if the sum of payments meets the grandtotal
func CheckInvoiceStatus(tx *sql.Tx, invoiceID int) (string, error) {
	var totalPaid float64
	var grandTotal float64

	// Get the grand total from the invoice header
	err := tx.QueryRow(`SELECT grandtotal FROM invoices WHERE invoiceid = $1`, invoiceID).Scan(&grandTotal)
	if err != nil {
		return "", err
	}

	// Sum all payments from transactionhistory
	err = tx.QueryRow(`SELECT COALESCE(SUM(amount), 0) FROM transactionhistory WHERE invoiceid = $1`, invoiceID).Scan(&totalPaid)
	if err != nil {
		return "", err
	}

	if totalPaid >= grandTotal {
		return "Paid", nil
	} else if totalPaid > 0 {
		return "Partially Paid", nil
	}
	return "Pending", nil
}
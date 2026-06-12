package Services

import (
	"database/sql"
	"fmt"
	"time"
)

func GetFinancialYear(invoiceDate time.Time) string {

	year := invoiceDate.Year()

	if invoiceDate.Month() >= time.April {
		return fmt.Sprintf(
			"%02d-%02d",
			year%100,
			(year+1)%100,
		)
	}

	return fmt.Sprintf(
		"%02d-%02d",
		(year-1)%100,
		year%100,
	)
}

func GenerateInvoiceNumber(
    db *sql.DB,
    invoiceDate time.Time,
) (string, error) {

    fy := GetFinancialYear(invoiceDate)

    pattern := "INV" + fy + "/%"

    var lastSeq int

   query := `
    SELECT COALESCE(
        MAX(
            CAST(
                SPLIT_PART(invoicenumber, '/', 2)
                AS INTEGER
            )
        ),
        0
    )
    FROM invoices
    WHERE invoicenumber LIKE $1
    AND issavedraft = false
    AND deletedat IS NULL
`

    err := db.QueryRow(
    query,
    pattern,
).Scan(&lastSeq)

    // First invoice in FY
    if err == sql.ErrNoRows {
        return fmt.Sprintf(
            "INV%s/001",
            fy,
        ), nil
    }

    if err != nil {
        return "", err
    }

    nextSeq := lastSeq + 1

    return fmt.Sprintf(
        "INV%s/%03d",
        fy,
        nextSeq,
    ), nil
}



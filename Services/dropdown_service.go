package Services

import (
	"database/sql"
	"invoice-backend/Helper/Utils" // Ensure this matches your package path for DecryptFromDB
	"invoice-backend/Models/dto"
)

// GetDropdownData is a generic function to fetch ID and Name from any table.
// It matches the standard *sql.DB signature used in your new routes.
func GetDropdownData(db *sql.DB, query string) ([]dto.DropdownModel, error) {
	// 1. Execute the query
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []dto.DropdownModel

	// 2. Iterate through rows
	for rows.Next() {
		var item dto.DropdownModel
		
		// Scans the two columns (usually ID and Name) into our struct
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			return nil, err
		}

		// 🔐 3. Static DB Decryption (Optional but Recommended)
		// If your DB names are stored encrypted, we unlock them here.
		// If they are plain text, this function will simply return the original string.
		decryptedName, decErr := Utils.DecryptFromDB(item.Name)
		if decErr == nil {
			item.Name = decryptedName
		}

		results = append(results, item)
	}

	// 4. Check for errors during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// 5. Handle empty results to return [] instead of null for the frontend
	if results == nil {
		results = []dto.DropdownModel{}
	}

	return results, nil
}
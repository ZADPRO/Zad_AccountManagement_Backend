package Service

import (
	"database/sql"
	"invoice-backend/Model"
)

// GetDropdownData is a generic function to fetch ID and Name from any table
func GetDropdownData(db *sql.DB, query string) ([]Model.DropdownModel, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Model.DropdownModel

	for rows.Next() {
		var item Model.DropdownModel
		// Scans the two columns (ID and Name) into our struct
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	// Handle case where no rows are returned
	if results == nil {
		results = []Model.DropdownModel{}
	}

	return results, nil
}
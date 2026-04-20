package Service

import (
	"database/sql"
	"invoice-backend/Model"
	"invoice-backend/Query"
	"invoice-backend/Utils"
)

func FetchAllClients(db *sql.DB) ([]Model.ClientListModel, error) {
	rows, err := db.Query(Query.ListActiveClients)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []Model.ClientListModel
	for rows.Next() {
		var c Model.ClientListModel
		if err := rows.Scan(&c.ClientID, &c.ClientCode, &c.Name, &c.BusinessName, &c.IsActive); err != nil {
			return nil, err
		}
		clients = append(clients, c)
	}
	return clients, nil
}

func FetchAllUsers(db *sql.DB) ([]Model.UserData, error) {
    rows, err := db.Query(Query.ListActiveUsers)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []Model.UserData
    for rows.Next() {
        var u Model.UserData
        // 1. Scan the raw (encrypted) data from the DB into the struct
        if err := rows.Scan(&u.UserID, &u.UserCode, &u.Username, &u.FirstName, &u.LastName, &u.RoleID); err != nil {
            return nil, err
        }

        // 🔓 2. Decrypt the sensitive fields before sending them to the Frontend
        // We use the blank identifier _ to ignore errors for old plain-text data
        u.FirstName, _ = Utils.Decrypt(u.FirstName)
        u.LastName, _ = Utils.Decrypt(u.LastName)
        
        // If your ListActiveUsers query also fetches email, decrypt that too:
        // u.Email, _ = Utils.Decrypt(u.Email)

        // 3. Append the "Clean" decrypted user to the list
        users = append(users, u)
    }
    return users, nil
}
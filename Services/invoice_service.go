package Services

import (
	"database/sql"
	"fmt"
	"invoice-backend/Models/dto"
    "invoice-backend/Helper/Utils" 
	"invoice-backend/Query"
	"time"
    "encoding/json"
    
    "strings"
)

// CreateFullInvoice handles the business logic and DB transaction
func CreateFullInvoice(db *sql.DB, req dto.CreateInvoiceRequest) (int, error) {

    //fmt.Println("COMPANY PROFILE ID =", req.CompanyProfileID)

    if req.IsSaveDraft {
        req.InvoiceNumber = ""
    } else {

        invoiceDate, err := time.Parse(
            "2006-01-02",
            req.InvoiceDate,
        )

        if err != nil {
            return 0, err
        }

        invoiceNumber, err := GenerateInvoiceNumber(
            db,
            invoiceDate,
        )

        if err != nil {
            return 0, err
        }

        req.InvoiceNumber = invoiceNumber
    }

   

    

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

    //fmt.Println("TAX TYPE:", req.TaxType)
    //fmt.Println("SAVING TO DB =", req.IsSaveDraft)

    err = tx.QueryRow(
        Query.InsertInvoiceHeaderQuery,
        req.InvoiceNumber,
        req.ClientID,
        req.CompanyProfileID,
        req.InvoiceDate,
        req.GrandTotal,
        req.PaymentStatus,
        req.UpdatedBy,
        string(customJSON),
        req.InvoiceDueDate,
	    req.Currency, 
        req.BankID,
        req.SignatureAuthorityID,
        req.InvoiceType,
        req.TaxType,
        req.TaxAmount,
        req.TdsAmount,
        req.IsSaveDraft,
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
            item.SACCode,
            item.Quantity,
            item.UnitPrice,
            item.LineTotal,
            req.UpdatedBy,
            string(itemCustomJSON),
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

func UpdateInvoice(
    db *sql.DB,
    invoiceID int,
    req dto.CreateInvoiceRequest,
) error {


    // Generate invoice number when draft becomes invoice
    if !req.IsSaveDraft {

    var existingInvoiceNo string

    err := db.QueryRow(
        `SELECT invoicenumber
         FROM invoices
         WHERE invoiceid = $1`,
        invoiceID,
    ).Scan(&existingInvoiceNo)

    if err != nil {
        return err
    }

    if strings.TrimSpace(existingInvoiceNo) == "" {

        invoiceDate, err := time.Parse(
            "2006-01-02",
            req.InvoiceDate,
        )

        if err != nil {
            return err
        }

        req.InvoiceNumber, err =
            GenerateInvoiceNumber(
                db,
                invoiceDate,
            )

            

        if err != nil {
            return err
        }

    } else {

        req.InvoiceNumber = existingInvoiceNo
    }
}

    customJSON, err := json.Marshal(req.CustomValues)
    if err != nil {
        return err
    }


   _, err = db.Exec(
    Query.UpdateInvoiceQuery,

    req.InvoiceNumber,      // $1
    req.ClientID,           // $2
    req.CompanyProfileID,   // $3
    req.InvoiceDate,        // $4
    req.GrandTotal,         // $5
    req.PaymentStatus,      // $6
    req.UpdatedBy,          // $7
    string(customJSON),     // $8
    req.InvoiceDueDate,     // $9
    req.Currency,           // $10
    req.BankID,             // $11
    req.SignatureAuthorityID,// $12
    req.InvoiceType,        // $13
    req.TaxType,            // $14
    req.TaxAmount,          // $15
    req.TdsAmount,          // $16
    req.IsSaveDraft,        // $17

    invoiceID,              // $18
)

    
   if err != nil {
    return err
}


return nil
}

func GetInvoiceByID(db *sql.DB, invoiceID int) (*dto.InvoiceResponse, error) {

    var inv dto.InvoiceResponse
    var clientID int
    var customJSON []byte 
    var companyProfileID sql.NullInt64
    
    // 1. Declare Nullable variables for ALL potentially missing fields
    var (
        bankID                                                        sql.NullInt64
        sigAuthID sql.NullInt64

        sigName sql.NullString
        sigRole sql.NullString
        sigContact sql.NullString
        sigEmail sql.NullString
        sigURL sql.NullString                                                    
         bankName, accNum, ifsc, bankAddr, logoUrl, accType, swiftCode sql.NullString
        
        // Fields that might be NULL in older invoices
        invDueDate, curr, invType, txType                             sql.NullString
        txAmount, tdsAmt                                              sql.NullFloat64
    )

    var (
    companyName sql.NullString
    address1 sql.NullString
    address2 sql.NullString
    city sql.NullString
    state sql.NullString
    country sql.NullString
    pincode sql.NullString
    gst sql.NullString
    companyEmail sql.NullString
    companyPhone sql.NullString
    website sql.NullString
        companyLogo sql.NullString

)

    // 2. Fetch invoice header safely
    err := db.QueryRow(Query.GetInvoiceByIDQuery, invoiceID).Scan(
        &inv.InvoiceID,
        &inv.InvoiceNumber,
        &inv.InvoiceDate,
        &inv.GrandTotal,
        &inv.PaymentStatus,
        &clientID,
        &companyProfileID,
        &customJSON,
        &invDueDate,  // Mapped to i.invoiceduedate
        &curr,        // Mapped to i.currency
        &bankID,      // Mapped to i.bankID
        &invType,     // Mapped to i.invoicetype
        &txType,      // Mapped to i.taxtype
        &txAmount,    // Mapped to i.taxamount
        &tdsAmt,      // Mapped to i.tdsamount
        &sigAuthID,
        &sigName,
        &sigRole,
        &sigContact,
        &sigEmail,
        &sigURL,
       &companyName,
&address1,
&address2,
&city,
&state,
&country,
&pincode,
&gst,
&companyEmail,
&companyPhone,
&website,
&companyLogo,
        &bankName,
        &accNum,
        &ifsc,
        &bankAddr,
        &logoUrl,
        &accType,
        &swiftCode,
    )

    //  IF IT FAILS HERE, CHECK YOUR TERMINAL
    if err != nil {
        fmt.Println("🚨 DATABASE HEADER SCAN ERROR:", err) 
        return nil, fmt.Errorf("invoice not found or scan error: %w", err)
    }

    // 3. Map valid nullable values back to struct
    if invDueDate.Valid { inv.InvoiceDueDate = invDueDate.String }
    if curr.Valid { inv.Currency = curr.String }
    if invType.Valid { inv.InvoiceType = invType.String }
    if txType.Valid { inv.TaxType = txType.String }
    if txAmount.Valid { inv.TaxAmount = txAmount.Float64 }
    if tdsAmt.Valid { inv.TdsAmount = tdsAmt.Float64 }
    if bankID.Valid { inv.BankID = int(bankID.Int64) }
    if companyProfileID.Valid { inv.CompanyProfileID = int(companyProfileID.Int64)}
    
    if sigAuthID.Valid { inv.SignatureAuthorityID = int(sigAuthID.Int64) }
    if sigName.Valid { inv.SignatureAuthorityName = sigName.String }
    if sigRole.Valid { inv.SignatureAuthorityRole = sigRole.String }
    if sigContact.Valid { inv.SignatureContactNumber = sigContact.String }
    if sigEmail.Valid { inv.SignatureEmail = sigEmail.String }
    if sigURL.Valid {inv.SignatureURL = sigURL.String}

    if bankName.Valid { inv.InvoiceBankName = bankName.String }
    if accNum.Valid { inv.InvoiceAccountNumber = accNum.String }
    if ifsc.Valid { inv.InvoiceIFSCCode = ifsc.String }
    if bankAddr.Valid { inv.InvoiceBankAddress = bankAddr.String }
    if logoUrl.Valid { inv.InvoiceQRCodeURL = logoUrl.String }
    if accType.Valid { inv.InvoiceAccountType = accType.String }
    if swiftCode.Valid { inv.InvoiceSwiftCode = swiftCode.String }
    if companyName.Valid {
    inv.CompanyName = companyName.String
}


if address1.Valid {
    inv.AddressLine1 = address1.String
}

if address2.Valid {
    inv.AddressLine2 = address2.String
}

if city.Valid {
    inv.City = city.String
}

if state.Valid {
    inv.State = state.String
}

if country.Valid {
    inv.Country = country.String
}

if pincode.Valid {
    inv.Pincode = pincode.String
}

if gst.Valid {
    inv.GSTNumber = gst.String
}

if companyEmail.Valid {
    inv.CompanyEmail = companyEmail.String
}

if companyPhone.Valid {
    inv.CompanyPhone = companyPhone.String
}

if website.Valid {
    inv.Website = website.String
}

if companyLogo.Valid {
    inv.CompanyLogoURL = companyLogo.String
}

//fmt.Println("COMPANY NAME =", inv.CompanyName)
//fmt.Println("ADDRESS 1 =", inv.AddressLine1)
//fmt.Println("CITY =", inv.City)

    // 4. Decrypt Bank Details
    if inv.InvoiceAccountNumber != "" {
        decryptedAcc, err := Utils.DecryptFromDB(inv.InvoiceAccountNumber)
        if err == nil {
            inv.InvoiceAccountNumber = decryptedAcc
        }
    }

    if inv.InvoiceIFSCCode != "" {
        decryptedIFSC, err := Utils.DecryptFromDB(inv.InvoiceIFSCCode)
        if err == nil {
            inv.InvoiceIFSCCode = decryptedIFSC
        }
    }

    if inv.InvoiceSwiftCode != "" {
        decryptedSwift, err := Utils.DecryptFromDB(inv.InvoiceSwiftCode)
        if err == nil {
            inv.InvoiceSwiftCode = decryptedSwift
        }
    }

    // 5. Handle Custom JSON
    if len(customJSON) > 0 {
        inv.CustomValues = []dto.CustomFieldValue{}
        if string(customJSON) != "[]" {
            if err := json.Unmarshal(customJSON, &inv.CustomValues); err != nil {
                fmt.Println("JSON UNMARSHAL ERROR:", err)
            }
        }
    }

    // 6. Fetch client
    var supplyTypeID, updatedBy int
    var isActive bool
    var updatedAt time.Time
    var billingCountryID, billingStateID sql.NullInt64
    var zip sql.NullInt64
    var gstNumber sql.NullString
    var pan sql.NullString
    var stateName sql.NullString
    var email, primaryNum, registeredAddr, clientCountry, billingAddr sql.NullString // Just in case
    
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
        &email,            // Was &inv.Client.Email
        &primaryNum,       // Was &inv.Client.PrimaryNumber
        &registeredAddr,   // Was &inv.Client.Address
        &clientCountry,    // Was &inv.Client.CountryName
        &stateName,
        &zip,
        &billingAddr,      // Was &inv.Client.BillingAddress
        &billingCountryID,
        &billingStateID,
        &inv.Client.TaxPercentage,
        &gstNumber,
        &pan,
        &inv.Client.IsExport,
      
        &inv.Client.BillingCountry,
        &inv.Client.BillingState,
    )
    inv.Client.SupplyTypeID = supplyTypeID

    switch supplyTypeID {
case 1:
    inv.Client.SupplyType = "B2B"
case 2:
    inv.Client.SupplyType = "B2C"
case 3:
    inv.Client.SupplyType = "C2C"
default:
    inv.Client.SupplyType = "-"
}
    
    //  IF IT FAILS HERE, CHECK YOUR TERMINAL
    if err != nil {
        fmt.Println("CLIENT SCAN ERROR:", err)
        return nil, fmt.Errorf("client not found: %w", err)
    }

    // Safely map client strings
    if email.Valid { inv.Client.Email = email.String }
    if primaryNum.Valid { inv.Client.PrimaryNumber = primaryNum.String }
    if registeredAddr.Valid { inv.Client.Address = registeredAddr.String }
    if clientCountry.Valid { inv.Client.CountryName = clientCountry.String }
    if billingAddr.Valid { inv.Client.BillingAddress = billingAddr.String }
   

    if zip.Valid { inv.Client.ZIP = int(zip.Int64) }
    if gstNumber.Valid { inv.Client.GSTNumber = gstNumber.String }
    if pan.Valid { inv.Client.PAN = pan.String }

    // Decrypt GST Number
    if inv.Client.GSTNumber != "" {
        decryptedGST, decErr := Utils.DecryptFromDB(inv.Client.GSTNumber)
        if decErr == nil {
            inv.Client.GSTNumber = decryptedGST
        }
    }

    // Decrypt PAN
    if inv.Client.PAN != "" {
        decryptedPAN, decErr := Utils.DecryptFromDB(inv.Client.PAN)
        if decErr == nil {
            inv.Client.PAN = decryptedPAN
        }
    }

    // 7. Fetch line items
    rows, err := db.Query(Query.GetInvoiceItemsByInvoiceIDQuery, invoiceID)
    if err != nil {
        fmt.Println("_ITEMS FETCH ERROR:", err)
        return nil, fmt.Errorf("items fetch error: %w", err)
    }
    defer rows.Close()
    
    for rows.Next() {
        var item dto.InvoiceItem
        var itemCustomJSON []byte
        
        err := rows.Scan(
            &item.ItemID,
            &item.Description,
            &item.SACCode,
            &item.Quantity,
            &item.UnitPrice,
            &item.LineTotal,
            &itemCustomJSON,
        )
        
        //  IF IT FAILS HERE, CHECK YOUR TERMINAL
        if err != nil {
            fmt.Println("ITEM ROW SCAN ERROR:", err)
            return nil, fmt.Errorf("item scan error: %w", err)
        }

        if len(itemCustomJSON) > 0 && string(itemCustomJSON) != "[]" {
            err := json.Unmarshal(itemCustomJSON, &item.CustomFieldValues)
            if err != nil {
                fmt.Println("item custom JSON error:", err)
            }
        }

        inv.Items = append(inv.Items, item)
    }

    //fmt.Printf("ITEMS = %+v\n", inv.Items)
    fmt.Println("SUPPLY TYPE ID =", supplyTypeID)
fmt.Println("SUPPLY TYPE =", inv.Client.SupplyType)

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


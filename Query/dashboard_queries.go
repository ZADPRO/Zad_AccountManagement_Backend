package Query

const (

    InsertInvoiceHeaderQuery = `
        INSERT INTO invoices (
            invoicenumber,
            clientid,
            invoicedate,
            grandtotal,
            paymentstatus,
            updatedat,
            updatedby,
            "CustomValues",
            invoiceduedate,
            currency,
            "bankID",
            signature_authority_id,
            invoicetype,
            taxtype,
            taxamount,
            tdsamount
        )
        VALUES (
            $1, $2, $3, $4, $5,
            NOW(), $6, $7,
            $8, $9, $10,
            $11,
            $12, $13, $14, $15
        )
        RETURNING invoiceid;
    `

    InsertInvoiceItemQuery = `
        INSERT INTO invoiceitems (
            invoiceid,
            description,
            quantity,
            unitprice,
            linetotal,
            updatedat,
            updatedby,
            custom_field_values
        )
        VALUES (
            $1, $2, $3, $4, $5,
            NOW(), $6, $7
        );
    
`


	GetInvoiceListQuery = `
        SELECT i.invoiceid, i.invoicenumber, c."name", i.invoicedate, i.grandtotal, i.paymentstatus 
        FROM invoices i
        JOIN "clientinformation" c ON i.clientid = c."clientid"
        WHERE i.deletedat IS NULL ORDER BY i.invoicedate DESC;`

	GetDashboardStatsQuery = `
        SELECT 
            (SELECT COUNT(*) FROM active_clients) as total_clients,
            (SELECT COUNT(*) FROM active_users) as active_users,
            (SELECT COALESCE(SUM(grandtotal), 0) FROM invoices WHERE paymentstatus = 'Paid') as total_revenue,
            (SELECT COALESCE(SUM(grandtotal), 0) 
            FROM invoices 
            WHERE paymentstatus = 'pending'
            AND deletedat IS NULL) as pending_amount,
            (SELECT COUNT(*) 
            FROM invoices 
            WHERE paymentstatus = 'pending'
            AND deletedat IS NULL) as overdue_count;`
 
        GetInvoiceByIDQuery = `
SELECT
    i.invoiceid,
    i.invoicenumber,
    i.invoicedate,
    i.grandtotal,
    i.paymentstatus,
    i.clientid,
    i."CustomValues",
    i.invoiceduedate,
    i.currency,
    i."bankID",
    i.invoicetype,
    i.taxtype,
    i.taxamount,
    i.tdsamount,
    i.signature_authority_id,

    sa.name AS signature_authority_name,
    sa.designation AS signature_authority_role,
    sa.contact_number,
    sa.email,

    b."BankName",
    b."AccountNumber",
    b."ifscCode",
    b."BankAddress",
    b."LogoURL",
    b."AccountType",
    b."SwiftCode"

FROM invoices i

LEFT JOIN "BankingDetails" b
ON i."bankID" = b."DetailsID"

LEFT JOIN signature_authorities sa
ON sa.id = i.signature_authority_id

WHERE i.invoiceid = $1;
`

GetInvoiceItemsByInvoiceIDQuery = `
SELECT 
    itemid,
    description,
    quantity,
    unitprice,
    linetotal,
    custom_field_values
FROM invoiceitems
WHERE invoiceid = $1;
`
 GetInvoiceCustomFieldsQuery = `
    SELECT 
     cfd."FieldLabel",
     icfv."Value"
     FROM "InvoiceCustomFieldValues" icfv
     JOIN "CustomFieldDefinitions" cfd
     ON icfv."FieldID" = cfd."FieldID"
     WHERE icfv."InvoiceID" = $1;`
     
 DeleteInvoiceItemsQuery = ` UPDATE invoiceitems 
    SET
        deletedat = NOW(),
        deletedby = $1
    WHERE invoiceid = $2; `

DeleteInvoiceQuery = `
    UPDATE invoices
    SET
        deletedat = NOW(),
        deletedby = $1
    WHERE invoiceid = $2;
`)
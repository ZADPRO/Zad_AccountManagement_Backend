package Query

const (

    InsertInvoiceHeaderQuery = `
       INSERT INTO invoices (
    invoicenumber,
    clientid,
    companyprofileid,
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
    tdsamount,
    issavedraft
)
VALUES (
    $1,  -- invoicenumber
    $2,  -- clientid
    $3,  -- companyprofileid
    $4,  -- invoicedate
    $5,  -- grandtotal
    $6,  -- paymentstatus
    NOW(),
    $7,  -- updatedby
    $8,  -- CustomValues
    $9,  -- invoiceduedate
    $10, -- currency
    $11, -- bankID
    $12, -- signature_authority_id
    $13, -- invoicetype
    $14, -- taxtype
    $15, -- taxamount
    $16, -- tdsamount
    $17  -- issavedraft
)
RETURNING invoiceid;
    `

    UpdateInvoiceQuery = `
UPDATE invoices
SET
    invoicenumber = $1,
    clientid = $2,
    companyprofileid = $3,
    invoicedate = $4,
    grandtotal = $5,
    paymentstatus = $6,
    updatedat = NOW(),
    updatedby = $7,
    "CustomValues" = $8,
    invoiceduedate = $9,
    currency = $10,
    "bankID" = $11,
    signature_authority_id = $12,
    invoicetype = $13,
    taxtype = $14,
    taxamount = $15,
    tdsamount = $16,
    issavedraft = $17
WHERE invoiceid = $18;
`

    InsertInvoiceItemQuery = `
        INSERT INTO invoiceitems (
            invoiceid,
            description,
            saccode,
            quantity,
            unitprice,
            linetotal,
            updatedat,
            updatedby,
            custom_field_values
        )
        VALUES (
            $1, $2, $3, $4, $5,  $6,
            NOW(), $7, $8
        );
    
`


	GetInvoiceListQuery = `
    SELECT
        i.invoiceid,
        i.invoicenumber,
        i.invoicetype,
        c."name",
        i.invoicedate,
        i.grandtotal,
        i.paymentstatus,
        i.issavedraft
    FROM invoices i
    JOIN "clientinformation" c
        ON i.clientid = c."clientid"
    WHERE i.deletedat IS NULL
    ORDER BY i.invoicedate DESC;
`

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
    i.companyprofileid,
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
    sa.signature_url,

    cp.companyname,
cp.addressline1,
cp.addressline2,
cp.city,
cp.state,
cp.country,
cp.pincode,
cp.gstnumber,
cp.email,
cp.phonenumber,
cp.website,
cp.logourl,

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

LEFT JOIN companyprofile cp
ON cp.id = i.companyprofileid

WHERE i.invoiceid = $1;
`

GetInvoiceItemsByInvoiceIDQuery = `
SELECT 
    itemid,
    description,
    saccode,
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
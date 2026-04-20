package Query

const (
	// Option A: Use quotes ONLY if they match your CREATE TABLE exactly
	// Option B: Remove quotes to let Postgres handle casing automatically
	ListActiveClients = `
		SELECT 
			"clientid", 
			"clientcode", 
			"name", 
			COALESCE("businessname", '') AS "businessname", 
			"isactive" 
		FROM active_clients;`

	ListActiveUsers = `
		SELECT 
			"userid", 
			"usercode", 
			"username", 
			COALESCE("firstname", '') AS "firstname", 
			COALESCE("lastname", '') AS "lastname", 
			"roleid" 
		FROM active_users;`
)  

const CreateUserQuery = `
INSERT INTO users (
    usercode,
    username,
    password,
    firstname,
    lastname,
    roleid,
    emailid,
    is_first_login
)
VALUES ($1, $2, $3, $4, $5, $6, $7,$8)
RETURNING userid;
`

// Existing dropdown queries...
// Existing dropdown queries standardized with aliases
const GetStatesDropdownQuery = `SELECT "stateid" AS id, "statename" AS name FROM "states" ORDER BY "statename" ASC;`
const GetSupplyTypesDropdownQuery = `SELECT "supplytypeid" AS id, "typename" AS name FROM "supplytypes" ORDER BY "typename" ASC;`

// New dropdown queries
const GetCountriesDropdownQuery = `SELECT "countryid" AS id, "countryname" AS name FROM "countries" ORDER BY "countryname" ASC;`
const GetRolesDropdownQuery = `SELECT "roleid" AS id, "rolename" AS name FROM "roles" ORDER BY "rolename" ASC;`



// CreateClientInfoQuery
const CreateClientInfoQuery = `
    INSERT INTO clientinformation (
    clientcode,
    name,
    businessname,
    supplytypeid,
    email,
    mobilenumber,
    registeredaddress,
    countryname,
    statename,
    zip,
    clienttype,
    updatedby,
    updatedat,
    isactive
)
VALUES (
    $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,NOW(),true
)
RETURNING clientid;
` 
const UpdateClientInfoQuery = `
    UPDATE "clientinformation" 
    SET 
        "name" = $1, 
        "businessname" = $2, 
        "supplytypeid" = $3,
        "email" = $4,
        "mobilenumber" = $5,
        "registeredaddress" = $6,
        "countryname" = $7,
        "statename" = $8,
        "zip" = $9,
        "updatedat" = NOW()
    WHERE "clientid" = $10;
`

const DeleteClientQuery = `
    UPDATE "clientinformation" 
    SET 
        "isactive" = false, 
        "updatedat" = NOW(), 
        "deletedat" = NOW()
    WHERE "clientid" = $1;`

// --- Client Tax Details ---

const CreateClientTaxQuery = `
INSERT INTO clienttaxdetails (
    clientid,
    gstnumber,
    pan,
    isexport,
    updatedat,
    updatedby,
    gststatus
)
VALUES (
    $1, $2, $3, $4, NOW(), $5, $6
);
`

const UpdateClientTaxQuery = `
    UPDATE "clienttaxdetails" 
    SET 
        "gststatus" = $1, 
        "gstnumber" = $2, 
        "pan" = $3, 
        "isexport" = $4
    WHERE "clientid" = $5;`

// --- User Queries ---

const GetUserByEmailQuery = `
    SELECT 
    u.userid, 
    u.password, 
    u.isactive, 
    r.rolename, -- This comes from the Roles table 
    u.username, 
    u.is_first_login 
FROM users u
JOIN roles r ON u.roleid = r.roleid
WHERE u.emailid = $1;`

const UpdateUserQuery = `
    UPDATE "users" 
    SET 
        "username" = $1, 
        "firstname" = $2, 
        "lastname" = $3, 
        "roleid" = $4,
        "emailid" = $5,
        "updatedat" = NOW()
    WHERE "userid" = $6;
`

const DeleteUserQuery = `
    UPDATE "users" 
    SET 
        "isactive" = false, 
        "deletedat" = NOW(),
        "deletedby" = $1
    WHERE "userid" = $2;`

// --- Invoice Queries ---

const InsertInvoiceHeaderQuery = `
    INSERT INTO invoices (invoicenumber, clientid, invoicedate, grandtotal, paymentstatus, updatedat, updatedby)
    VALUES ($1, $2, $3, $4, $5, NOW(), $6)
    RETURNING invoiceid;`

const InsertInvoiceItemQuery = `
    INSERT INTO invoiceitems (invoiceid, description, quantity, unitprice, linetotal, updatedat, updatedby)
    VALUES ($1, $2, $3, $4, $5, NOW(), $6);`

const GetInvoiceListQuery = `
    SELECT 
        i.invoiceid, 
        i.invoicenumber, 
        c."name", 
        i.invoicedate, 
        i.grandtotal, 
        i.paymentstatus 
    FROM invoices i
    JOIN "clientinformation" c ON i.clientid = c."clientid"
    WHERE i.deletedat IS NULL
    ORDER BY i.invoicedate DESC;`

const GetClientByIDQuery = `
SELECT 
    c.clientid,
    c.clientcode,
    c.name,
    c.businessname,
    c.supplytypeid,
    c.isactive,
    c.clienttype,
    c.updatedat,
    c.updatedby,
    c.email,
    c.mobilenumber,
    c.registeredaddress,
    c.countryname,
    c.statename,
    c.zip,

    t.gstnumber,
    t.pan,
    t.isexport,
    t.gststatus

FROM clientinformation c
LEFT JOIN clienttaxdetails t 
    ON c.clientid = t.clientid
WHERE c.clientid = $1;
`
const GetUserByIDQuery = `
    SELECT userid, usercode, username, firstname, lastname, roleid , emailid
    FROM active_users 
    WHERE userid = $1 LIMIT 1
`

const GetUserProfileByIDQuery = `
   -- Query to fetch display name and role name
SELECT 
    u.firstname, 
    u.lastname, 
    r.rolename 
FROM users u
JOIN roles r ON u.roleid = r.roleid
WHERE u.userid = $1;` 


const ResetPasswordQuery = `
    UPDATE users 
    SET 
        password = $1, 
        must_change_password = false 
    WHERE userid = $2;` 

const GetDashboardStatsQuery = `
    SELECT 
        (SELECT COUNT(*) FROM active_clients) as total_clients,
        (SELECT COUNT(*) FROM active_users) as active_users,
        (SELECT COALESCE(SUM(grandtotal), 0) FROM invoices WHERE paymentstatus = 'Paid') as total_revenue,
        (SELECT COALESCE(SUM(grandtotal), 0) FROM invoices WHERE paymentstatus = 'pending') as pending_amount,
        (SELECT COUNT(*) FROM invoices WHERE paymentstatus = 'pending') as overdue_count
`
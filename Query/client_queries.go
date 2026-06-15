package Query

const (
	GetAllClientsQuery = `
    SELECT "clientid", "clientcode", "name", "businessname", "isactive"
    FROM "active_clients"
    WHERE "isactive" IS TRUE
    `

	CreateClientInfoQuery = `
    INSERT INTO clientinformation (
        clientcode, name, businessname, supplytypeid, email,
        mobilenumber, registeredaddress, countryname, statename,
        zip, clienttype, updatedby, updatedat, isactive
    )
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW(), true)
    RETURNING clientid;`

	UpdateClientInfoQuery = `
    UPDATE clientinformation 
    SET 
        name             = $1,
        businessname     = $2,
        supplytypeid     = $3,
        email            = $4,
        mobilenumber     = $5,
        registeredaddress = $6,
        countryname      = $7,
        statename        = $8,
        zip              = $9,
        clienttype       = $10,
        updatedby        = $11,
        updatedat        = NOW()
    WHERE clientid = $12;`

	UpdateClientTaxQuery = `
    UPDATE clienttaxdetails 
    SET 
       
        gstnumber        = $1,
        pan              = $2,
        taxpercentage    = $3,
        isexport         = $4,
        billingaddress   = $5,
        billingcountryid = $6,
        billingstateid   = $7,
        updatedby        = $8,
        updatedat        = NOW()
    WHERE clientid = $9;`

	DeleteClientQuery = `
    UPDATE clientinformation 
    SET isactive = false, updatedat = NOW(), deletedat = NOW(), updatedby = $1 
    WHERE clientid = $2;`

	CreateClientTaxQuery = `
    INSERT INTO clienttaxdetails (
        clientid, gstnumber, pan, taxpercentage, isexport,
        billingaddress, billingcountryid, billingstateid,
        updatedby, updatedat
    )
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW());`

GetClientByIDQuery = `
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

    -- Billing + Tax
    COALESCE(t.billingaddress, ''),
    COALESCE(t.billingcountryid, 0),
    COALESCE(t.billingstateid, 0),
    COALESCE(t.taxpercentage, 0),
    COALESCE(t.gstnumber, ''),
    COALESCE(t.pan, ''),
    COALESCE(t.isexport, false),
    

    
    COALESCE(cn.countryname, ''),
    COALESCE(sn.statename, '')

FROM clientinformation c
LEFT JOIN clienttaxdetails t 
    ON c.clientid = t.clientid


LEFT JOIN countries cn 
    ON t.billingcountryid = cn.countryid

LEFT JOIN states sn 
    ON t.billingstateid = sn.stateid

WHERE c.clientid = $1;
`;
)
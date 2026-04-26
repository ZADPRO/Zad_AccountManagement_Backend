package Query

const (
	InsertInvoiceHeaderQuery = `
        INSERT INTO invoices (invoicenumber, clientid, invoicedate, grandtotal, paymentstatus, updatedat, updatedby)
        VALUES ($1, $2, $3, $4, $5, NOW(), $6)
        RETURNING invoiceid;`

	InsertInvoiceItemQuery = `
        INSERT INTO invoiceitems (invoiceid, description, quantity, unitprice, linetotal, updatedat, updatedby)
        VALUES ($1, $2, $3, $4, $5, NOW(), $6);`

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
            (SELECT COALESCE(SUM(grandtotal), 0) FROM invoices WHERE paymentstatus = 'pending') as pending_amount,
            (SELECT COUNT(*) FROM invoices WHERE paymentstatus = 'pending') as overdue_count;`
)
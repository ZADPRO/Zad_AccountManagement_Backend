package Query

// Insert into your transactionhistory table
const InsertTransactionQuery = `
    INSERT INTO transactionhistory (invoiceid, amount, transactiondate, updatedat, updatedby)
    VALUES ($1, $2, $3, NOW(), $4)`

// Update the paymentstatus in the invoices table
const UpdateInvoiceStatusQuery = `
    UPDATE invoices 
    SET paymentstatus = CASE 
        WHEN (SELECT SUM(amount) FROM transactionhistory WHERE invoiceid = $1) >= grandtotal THEN 'Paid'
        ELSE 'Partially Paid'
    END
    WHERE invoiceid = $1`
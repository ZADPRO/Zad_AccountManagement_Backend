
package dto 

type InvoiceItem struct {
	ItemID      int     `json:"itemid,omitempty"`
	InvoiceID   int     `json:"invoiceid,omitempty"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unitprice"`
	LineTotal   float64 `json:"linetotal"`
}

type CreateInvoiceRequest struct {
	InvoiceNumber string        `json:"invoicenumber"`
	ClientID      int           `json:"clientid"`
	InvoiceDate   string        `json:"invoicedate"` // YYYY-MM-DD
	GrandTotal    float64       `json:"grandtotal"`
	PaymentStatus string        `json:"paymentstatus"`
	UpdatedBy     int           `json:"updatedby"`
	Items         []InvoiceItem `json:"items"`
}
package Model

type PaymentRequest struct {
	InvoiceID       int     `json:"invoiceid"`
	Amount          float64 `json:"amount"`
	TransactionDate string  `json:"transactiondate"`
	TDSAmount       float64 `json:"tdsamount"` // captured for future tracking 
	UpdatedBy       int     `json:"updatedby"`
}

package dto 

type InvoiceItem struct {
	ItemID      int     `json:"itemid,omitempty"`
	InvoiceID   int     `json:"invoiceid,omitempty"`
	Description string  `json:"description"`
	SACCode     string  `json:"sacCode"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unitprice"`
	LineTotal   float64 `json:"linetotal"`
	CustomFieldValues []CustomFieldValue  `json:"customFieldValues"`
}
type CustomFieldValue struct {
	FieldID int    `json:"fieldId"`
	Label   string `json:"label"`
	Value   string `json:"value"`
}

type CreateInvoiceRequest struct {
	InvoiceNumber string        `json:"invoicenumber"`
	ClientID      int           `json:"clientid"`
	CompanyProfileID int `json:"companyProfileId"`
	InvoiceDate   string        `json:"invoicedate"` // YYYY-MM-DD
	GrandTotal    float64       `json:"grandtotal"`
	SignatureAuthorityID int `json:"signatureauthorityid"`
	PaymentStatus string        `json:"paymentstatus"`
	UpdatedBy     int           `json:"updatedby"`
	Items         []InvoiceItem `json:"items"`
	CustomValues []CustomFieldValue `json:"customValues"`
	InvoiceDueDate string             `json:"invoiceduedate"`
	Currency       string             `json:"currency"` 
	BankID      int    `json:"bankId"`
	InvoiceType string `json:"invoiceType"`
	TaxType string `json:"taxtype"`
	TaxAmount     float64     `json:"taxamount"`
	TdsAmount     float64     `json:"tdsamount"`
	IsSaveDraft bool `json:"isSaveDraft"`
}
type InvoiceResponse struct {
    InvoiceID     int                `json:"invoiceid"`
    InvoiceNumber string             `json:"invoicenumber"`
    InvoiceDate   string             `json:"invoicedate"`
    GrandTotal    float64            `json:"grandtotal"`
    PaymentStatus string             `json:"paymentstatus"`
    Client        ClientFullResponse `json:"client"`
    Items         []InvoiceItem      `json:"items"`
	CustomValues []CustomFieldValue `json:"customValues"`
	InvoiceDueDate string             `json:"invoiceduedate"`
	Currency       string             `json:"currency"` 
	BankID      int    `json:"bankId"`
	CompanyProfileID int    `json:"companyProfileId"`

CompanyName  string `json:"companyName"`
AddressLine1 string `json:"addressLine1"`
AddressLine2 string `json:"addressLine2"`

City         string `json:"city"`
State        string `json:"state"`
Country      string `json:"country"`
Pincode      string `json:"pincode"`

GSTNumber    string `json:"gstNumber"`

CompanyEmail string `json:"companyEmail"`
CompanyPhone string `json:"companyPhone"`
Website      string `json:"website"`

CompanyLogoURL string `json:"companyLogoUrl"`
	SignatureURL string `json:"signatureUrl"`
	SignatureAuthorityID   int    `json:"signatureAuthorityId"`
	SignatureAuthorityName string `json:"signatureAuthorityName"`
	SignatureAuthorityRole string `json:"signatureAuthorityRole"`
	SignatureContactNumber string `json:"signatureContactNumber"`
	SignatureEmail         string `json:"signatureEmail"`
	InvoiceType string `json:"invoiceType"` 
	TaxType string `json:"taxtype"`
	TaxAmount     float64     `json:"taxamount"`
	TdsAmount     float64     `json:"tdsamount"`

	InvoiceBankName      string `json:"invoiceBankName"`
	InvoiceAccountNumber string `json:"invoiceAccountNumber"`
	InvoiceIFSCCode      string `json:"invoiceIfscCode"`
	InvoiceBankAddress   string `json:"invoiceBankAddress"`
	InvoiceAccountType   string `json:"invoiceAccountType"`
	InvoiceSwiftCode     string `json:"invoiceSwiftCode"`
	InvoiceQRCodeURL     string `json:"invoiceQrCodeUrl"`
}


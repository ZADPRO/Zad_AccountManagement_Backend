package Service



func CalculateTaxType(clientState string, businessState string, isInternational bool) string {
	if isInternational {
		return "EXPORT"
	}

	if clientState == businessState {
		return "CGST_SGST"
	}
	
	return "IGST"
}
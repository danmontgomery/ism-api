package model

// ThirdPartyDistributionContract holds contractor IP rights information
// required when a 3rdPartyDistributionStatement is present.
type ThirdPartyDistributionContract struct {
	ContractNumber    string `json:"contractNumber"`
	ContractorName    string `json:"contractorName"`
	ContractorAddress string `json:"contractorAddress"`
	ExpirationDate    string `json:"expirationDate"`
}

// Complete returns true if all required fields are populated.
func (c *ThirdPartyDistributionContract) Complete() bool {
	return c.ContractNumber != "" &&
		c.ContractorName != "" &&
		c.ContractorAddress != "" &&
		c.ExpirationDate != ""
}

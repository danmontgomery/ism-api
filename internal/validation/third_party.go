package validation

import (
	"dmontgomery/ism-api/internal/model"
	"dmontgomery/ism-api/internal/refdata"
)

// ThirdPartyRule validates that a 3rdPartyDistributionContract is complete
// when a 3rdPartyDistributionStatement is present.
type ThirdPartyRule struct{}

func (r *ThirdPartyRule) Name() string { return "third_party" }

func (r *ThirdPartyRule) Applies(ism *model.ISM) bool {
	return ism.ThirdPartyDistributionStatement != ""
}

func (r *ThirdPartyRule) Validate(ism *model.ISM, _ *refdata.Registry) *ValidationResult {
	res := &ValidationResult{Valid: true}

	if ism.ThirdPartyDistributionContract == nil {
		res.AddError("3rdPartyDistributionContract", "third_party.contract_required",
			"3rdPartyDistributionContract is required when 3rdPartyDistributionStatement is present")
		return res
	}

	c := ism.ThirdPartyDistributionContract
	if c.ContractNumber == "" {
		res.AddError("3rdPartyDistributionContract.contractNumber", "third_party.missing_contract_number",
			"contractNumber is required")
	}
	if c.ContractorName == "" {
		res.AddError("3rdPartyDistributionContract.contractorName", "third_party.missing_contractor_name",
			"contractorName is required")
	}
	if c.ContractorAddress == "" {
		res.AddError("3rdPartyDistributionContract.contractorAddress", "third_party.missing_contractor_address",
			"contractorAddress is required")
	}
	if c.ExpirationDate == "" {
		res.AddError("3rdPartyDistributionContract.expirationDate", "third_party.missing_expiration_date",
			"expirationDate is required")
	}

	return res
}

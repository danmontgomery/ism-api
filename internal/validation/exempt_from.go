package validation

import (
	"expr.ai/ism-api/internal/model"
	"expr.ai/ism-api/internal/refdata"
)

// ExemptFromRule validates that every entry in exemptFrom is a known
// exemption value per CVEnumISMExemptFrom.xsd.
type ExemptFromRule struct{}

func (r *ExemptFromRule) Name() string { return "exempt_from" }

func (r *ExemptFromRule) Applies(ism *model.ISM) bool {
	return len(ism.ExemptFrom) > 0
}

func (r *ExemptFromRule) Validate(ism *model.ISM, reg *refdata.Registry) *ValidationResult {
	res := &ValidationResult{Valid: true}

	for _, code := range ism.ExemptFrom {
		if !reg.ValidExemptFrom(code) {
			res.AddError("exemptFrom", "exempt_from.invalid_code",
				"unknown exemptFrom value: "+code)
		}
	}

	return res
}

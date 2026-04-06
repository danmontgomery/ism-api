package validation

import (
	"expr.ai/ism-api/internal/model"
	"expr.ai/ism-api/internal/refdata"
)

// FGIRule validates Foreign Government Information source fields: country code
// validity for both fgiSourceOpen and fgiSourceProtected.
type FGIRule struct{}

func (r *FGIRule) Name() string { return "fgi" }

func (r *FGIRule) Applies(ism *model.ISM) bool {
	return len(ism.FGISourceOpen) > 0 || len(ism.FGISourceProtected) > 0
}

func (r *FGIRule) Validate(ism *model.ISM, reg *refdata.Registry) *ValidationResult {
	res := &ValidationResult{Valid: true}

	for _, code := range ism.FGISourceOpen {
		if !reg.ValidCountryCode(code) {
			res.AddError("fgiSourceOpen", "fgi.invalid_country_code",
				"unknown country code in fgiSourceOpen: "+code)
		}
	}

	for _, code := range ism.FGISourceProtected {
		if !reg.ValidCountryCode(code) {
			res.AddError("fgiSourceProtected", "fgi.invalid_country_code",
				"unknown country code in fgiSourceProtected: "+code)
		}
	}

	return res
}

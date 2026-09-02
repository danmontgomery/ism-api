package validation

import (
	"dmontgomery/ism-api/internal/model"
	"dmontgomery/ism-api/internal/refdata"
)

// CompliesWithRule validates that every entry in compliesWith is a known
// compliance value per CVEnumISMCompliesWith.xsd.
type CompliesWithRule struct{}

func (r *CompliesWithRule) Name() string { return "complies_with" }

func (r *CompliesWithRule) Applies(ism *model.ISM) bool {
	return len(ism.CompliesWith) > 0
}

func (r *CompliesWithRule) Validate(ism *model.ISM, reg *refdata.Registry) *ValidationResult {
	res := &ValidationResult{Valid: true}

	for _, code := range ism.CompliesWith {
		if !reg.ValidCompliesWith(code) {
			res.AddError("compliesWith", "complies_with.invalid_code",
				"unknown compliesWith value: "+code)
		}
	}

	return res
}

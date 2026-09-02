package validation

import (
	"dmontgomery/ism-api/internal/model"
	"dmontgomery/ism-api/internal/refdata"
)

// NonUSControlsRule validates that every entry in nonUSControls is a known
// non-US control marking per CVEnumISMNonUSControls.xsd.
type NonUSControlsRule struct{}

func (r *NonUSControlsRule) Name() string { return "non_us_controls" }

func (r *NonUSControlsRule) Applies(ism *model.ISM) bool {
	return len(ism.NonUSControls) > 0
}

func (r *NonUSControlsRule) Validate(ism *model.ISM, reg *refdata.Registry) *ValidationResult {
	res := &ValidationResult{Valid: true}

	for _, code := range ism.NonUSControls {
		if !reg.ValidNonUSControl(code) {
			res.AddError("nonUSControls", "non_us_controls.invalid_code",
				"unknown non-US control: "+code)
		}
	}

	return res
}

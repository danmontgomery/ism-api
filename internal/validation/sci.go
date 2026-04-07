package validation

import (
	"expr.ai/ism-api/internal/model"
	"expr.ai/ism-api/internal/refdata"
)

// SCIRule validates SCI (Sensitive Compartmented Information) controls: each
// code must be a known CVEnumISMSCIControls value and classification must be
// TOP SECRET.
type SCIRule struct{}

func (r *SCIRule) Name() string { return "sci" }

func (r *SCIRule) Applies(ism *model.ISM) bool {
	return len(ism.SCIControls) > 0
}

func (r *SCIRule) Validate(ism *model.ISM, reg *refdata.Registry) *ValidationResult {
	res := &ValidationResult{Valid: true}

	// SCI controls require TOP SECRET classification.
	if !ism.Classification.AtLeast(model.ClassificationTS) {
		res.AddError("sciControls", "sci.requires_ts",
			"SCI controls require TOP SECRET classification")
	}

	// Each control code must be in the reference data.
	for _, code := range ism.SCIControls {
		if !reg.ValidSCIControl(code) {
			res.AddError("sciControls", "sci.invalid_control",
				"unknown SCI control code: "+code)
		}
	}

	return res
}

package validation

import (
	"strings"

	"dmontgomery/ism-api/internal/model"
	"dmontgomery/ism-api/internal/refdata"
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

	// HCS/HCS-* and TK-GEOCAP require NOFORN (E4-S6.f, E4-S6.f.tk).
	hasNOFORN := false
	for _, dc := range ism.DisseminationControls {
		if dc == "NOFORN" {
			hasNOFORN = true
			break
		}
	}
	if !hasNOFORN {
		for _, code := range ism.SCIControls {
			if code == "HCS" || strings.HasPrefix(code, "HCS-") || code == "TK-GEOCAP" {
				res.AddError("sciControls", "sci.requires_noforn",
					code+" requires NOFORN in disseminationControls")
			}
		}
	}

	return res
}

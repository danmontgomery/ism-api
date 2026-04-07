package validation

import (
	"expr.ai/ism-api/internal/model"
	"expr.ai/ism-api/internal/refdata"
)

// DisseminationRule validates dissemination controls: REL requires releasableTo,
// DISPLAY ONLY requires displayOnlyTo, NOFORN/REL mutual exclusion, and
// minimum classification gates.
type DisseminationRule struct{}

func (r *DisseminationRule) Name() string { return "dissemination" }

func (r *DisseminationRule) Applies(ism *model.ISM) bool {
	return len(ism.DisseminationControls) > 0
}

func (r *DisseminationRule) Validate(ism *model.ISM, reg *refdata.Registry) *ValidationResult {
	res := &ValidationResult{Valid: true}

	controlSet := make(map[string]bool, len(ism.DisseminationControls))
	for _, dc := range ism.DisseminationControls {
		controlSet[dc] = true
	}

	// Validate each control code is known.
	for _, dc := range ism.DisseminationControls {
		if !reg.ValidDisseminationControl(dc) {
			res.AddError("disseminationControls", "dissemination.invalid_control",
				"unknown dissemination control: "+dc)
		}
	}

	// Field requirements: REL/EYES -> releasableTo, DISPLAY ONLY -> displayOnlyTo.
	for _, req := range refdata.DisseminationFieldRequirements() {
		if !controlSet[req.Control] {
			continue
		}
		switch req.Field {
		case "releasableTo":
			if len(ism.ReleasableTo) == 0 {
				res.AddError("releasableTo", "dissemination.missing_releasable_to",
					req.Control+" requires releasableTo to be populated")
			}
		case "displayOnlyTo":
			if len(ism.DisplayOnlyTo) == 0 {
				res.AddError("displayOnlyTo", "dissemination.missing_display_only_to",
					req.Control+" requires displayOnlyTo to be populated")
			}
		case "sciControls":
			if len(ism.SCIControls) == 0 {
				res.AddError("sciControls", "dissemination.missing_sci_controls",
					req.Control+" requires sciControls to be populated")
			}
		}
	}

	// Mutual exclusion checks.
	for _, pair := range refdata.ExclusiveDisseminationPairs() {
		if controlSet[pair.A] && controlSet[pair.B] {
			res.AddError("disseminationControls", "dissemination.exclusive_conflict",
				pair.A+" and "+pair.B+" cannot be used together")
		}
	}

	// Minimum classification gates.
	for _, gate := range refdata.DisseminationClassificationGates() {
		if controlSet[gate.Control] && !ism.Classification.AtLeast(gate.MinClassification) {
			res.AddError("disseminationControls", "dissemination.insufficient_classification",
				gate.Control+" requires at least "+gate.MinClassification.String()+" classification")
		}
	}

	// Validate releasableTo country codes.
	for _, code := range ism.ReleasableTo {
		if !reg.ValidCountryCode(code) {
			res.AddError("releasableTo", "dissemination.invalid_country_code",
				"unknown country code in releasableTo: "+code)
		}
	}

	// Validate displayOnlyTo country codes.
	for _, code := range ism.DisplayOnlyTo {
		if !reg.ValidCountryCode(code) {
			res.AddError("displayOnlyTo", "dissemination.invalid_country_code",
				"unknown country code in displayOnlyTo: "+code)
		}
	}

	return res
}

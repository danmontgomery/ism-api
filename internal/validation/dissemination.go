package validation

import (
	"dmontgomery/ism-api/internal/model"
	"dmontgomery/ism-api/internal/refdata"
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

	// Maximum classification gate: FOUO is U-only (CL-12, CR-11).
	if controlSet["FOUO"] && ism.Classification != model.ClassificationU {
		res.AddError("disseminationControls", "dissemination.exceeds_max_classification",
			"FOUO is only permitted at UNCLASSIFIED classification")
	}

	// CR-3: REL TO requires at least one non-USA country (E4-S10.d.5).
	if controlSet["REL"] && len(ism.ReleasableTo) > 0 {
		hasNonUSA := false
		for _, code := range ism.ReleasableTo {
			if code != "USA" {
				hasNonUSA = true
				break
			}
		}
		if !hasNonUSA {
			res.AddError("releasableTo", "dissemination.rel_to_usa_only",
				"REL TO requires at least one country besides USA per E4-S10.d.5")
		}
	}

	// CR-4: REL TO must include USA (E4-S10.d.4).
	if controlSet["REL"] && len(ism.ReleasableTo) > 0 {
		hasUSA := false
		for _, code := range ism.ReleasableTo {
			if code == "USA" {
				hasUSA = true
				break
			}
		}
		if !hasUSA {
			res.AddError("releasableTo", "dissemination.rel_to_missing_usa",
				"REL TO must include USA per E4-S10.d.4")
		}
	}

	// MX-8: US dissemination controls require USA in ownerProducer (E4-S1.d).
	if len(ism.DisseminationControls) > 0 && len(ism.OwnerProducer) > 0 {
		hasUSAOwner := false
		for _, owner := range ism.OwnerProducer {
			if owner == "USA" {
				hasUSAOwner = true
				break
			}
		}
		if !hasUSAOwner {
			res.AddError("disseminationControls", "dissemination.non_us_owner",
				"US dissemination controls require USA in ownerProducer per E4-S1.d")
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

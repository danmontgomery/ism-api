package validation

import (
	"expr.ai/ism-api/internal/model"
	"expr.ai/ism-api/internal/refdata"
)

// FGIRule validates Foreign Government Information source fields: country code
// validity, minimum classification gate (E4-S9.b), and FGI/REL TO compatibility
// (E4-S9.l).
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

	// E4-S9.b: FGI requiring protection must be classified at least CONFIDENTIAL.
	if !ism.Classification.AtLeast(model.ClassificationC) {
		res.AddError("classification", "fgi.insufficient_classification",
			"FGI requires at least CONFIDENTIAL classification per E4-S9.b")
	}

	// E4-S9.l: REL TO with FGI — FGI source countries must be in the REL TO list.
	hasREL := false
	for _, ctrl := range ism.DisseminationControls {
		if ctrl == "REL" {
			hasREL = true
			break
		}
	}
	if hasREL && len(ism.ReleasableTo) > 0 {
		relSet := make(map[string]bool, len(ism.ReleasableTo))
		for _, c := range ism.ReleasableTo {
			relSet[c] = true
		}
		for _, src := range ism.FGISourceOpen {
			if !relSet[src] {
				res.AddError("fgiSourceOpen", "fgi.rel_to_incompatible",
					"FGI source "+src+" not in releasableTo list per E4-S9.l")
			}
		}
		for _, src := range ism.FGISourceProtected {
			if !relSet[src] {
				res.AddError("fgiSourceProtected", "fgi.rel_to_incompatible",
					"FGI source "+src+" not in releasableTo list per E4-S9.l")
			}
		}
	}

	return res
}

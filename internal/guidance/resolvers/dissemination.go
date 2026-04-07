package resolvers

import (
	"expr.ai/ism-api/internal/guidance"
	"expr.ai/ism-api/internal/model"
	"expr.ai/ism-api/internal/refdata"
)

// DisseminationResolver provides guidance for disseminationControls,
// releasableTo, and displayOnlyTo fields.
type DisseminationResolver struct{}

func (r *DisseminationResolver) Fields() []string {
	return []string{"disseminationControls", "releasableTo", "displayOnlyTo", "sciControls"}
}

func (r *DisseminationResolver) Resolve(ism *model.ISM, reg *refdata.Registry) []guidance.FieldGuidance {
	var results []guidance.FieldGuidance

	// disseminationControls — always available, filtered by classification level.
	var allowed []guidance.AllowedValue
	for _, dc := range reg.DisseminationControls {
		if dc.MinClassification != "" && !ism.Classification.AtLeast(dc.MinClassification) {
			continue
		}
		allowed = append(allowed, guidance.AllowedValue{Code: dc.Code, Label: dc.Label})
	}
	results = append(results, guidance.FieldGuidance{
		Field:         "disseminationControls",
		Status:        guidance.StatusAvailable,
		AllowedValues: allowed,
	})

	// releasableTo — required when REL or EYES is selected.
	if hasControl(ism.DisseminationControls, "REL") || hasControl(ism.DisseminationControls, "EYES") {
		results = append(results, guidance.FieldGuidance{
			Field:    "releasableTo",
			Status:   guidance.StatusRequired,
			Required: true,
		})
	} else {
		results = append(results, guidance.FieldGuidance{
			Field:      "releasableTo",
			Status:     guidance.StatusNotApplicable,
			RequiredIf: "REL or EYES dissemination control is selected",
			Reason:     "Only applicable when REL or EYES dissemination control is set",
		})
	}

	// displayOnlyTo — required when DISPLAY ONLY is selected.
	if hasControl(ism.DisseminationControls, "DISPLAY ONLY") {
		results = append(results, guidance.FieldGuidance{
			Field:    "displayOnlyTo",
			Status:   guidance.StatusRequired,
			Required: true,
		})
	} else {
		results = append(results, guidance.FieldGuidance{
			Field:      "displayOnlyTo",
			Status:     guidance.StatusNotApplicable,
			RequiredIf: "DISPLAY ONLY dissemination control is selected",
			Reason:     "Only applicable when DISPLAY ONLY dissemination control is set",
		})
	}

	// sciControls — available when SCI dissemination control is selected.
	if hasControl(ism.DisseminationControls, "SCI") {
		var sciAllowed []guidance.AllowedValue
		for _, sc := range reg.SCIControls {
			sciAllowed = append(sciAllowed, guidance.AllowedValue{Code: sc.Code, Label: sc.Label})
		}
		results = append(results, guidance.FieldGuidance{
			Field:         "sciControls",
			Status:        guidance.StatusAvailable,
			AllowedValues: sciAllowed,
		})
	} else {
		results = append(results, guidance.FieldGuidance{
			Field:      "sciControls",
			Status:     guidance.StatusNotApplicable,
			RequiredIf: "SCI dissemination control is selected",
			Reason:     "Only applicable when SCI dissemination control is set",
		})
	}

	return results
}

func hasControl(controls []string, code string) bool {
	for _, c := range controls {
		if c == code {
			return true
		}
	}
	return false
}

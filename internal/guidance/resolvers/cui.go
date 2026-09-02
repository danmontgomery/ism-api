package resolvers

import (
	"dmontgomery/ism-api/internal/guidance"
	"dmontgomery/ism-api/internal/model"
	"dmontgomery/ism-api/internal/refdata"
)

// CUIResolver provides guidance for CUI-specific fields: categoryMarkings,
// controlledByName, controlledByOffice, and poc.
type CUIResolver struct{}

func (r *CUIResolver) Fields() []string {
	return []string{"categoryMarkings", "controlledByName", "controlledByOffice", "poc"}
}

func (r *CUIResolver) Resolve(ism *model.ISM, reg *refdata.Registry) []guidance.FieldGuidance {
	if ism.Classification != model.ClassificationCUI {
		return notApplicableAll(r.Fields(), "Only applicable for CUI classification")
	}

	cats := make([]guidance.AllowedValue, len(reg.CUICategories))
	for i, c := range reg.CUICategories {
		cats[i] = guidance.AllowedValue{Code: c.Code, Label: c.Label, Type: c.Type}
	}

	return []guidance.FieldGuidance{
		{
			Field:         "categoryMarkings",
			Status:        guidance.StatusAvailable,
			RequiredIf:    "CUI Specified marking is intended",
			AllowedValues: cats,
		},
		{
			Field:  "controlledByName",
			Status: guidance.StatusAvailable,
		},
		{
			Field:  "controlledByOffice",
			Status: guidance.StatusAvailable,
		},
		{
			Field:  "poc",
			Status: guidance.StatusAvailable,
		},
	}
}

// notApplicableAll returns not_applicable guidance for each field with the given reason.
func notApplicableAll(fields []string, reason string) []guidance.FieldGuidance {
	results := make([]guidance.FieldGuidance, len(fields))
	for i, f := range fields {
		results[i] = guidance.FieldGuidance{
			Field:  f,
			Status: guidance.StatusNotApplicable,
			Reason: reason,
		}
	}
	return results
}

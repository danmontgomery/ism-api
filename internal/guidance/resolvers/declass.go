package resolvers

import (
	"expr.ai/ism-api/internal/guidance"
	"expr.ai/ism-api/internal/model"
	"expr.ai/ism-api/internal/refdata"
)

// DeclassResolver provides guidance for declassification fields:
// declassDate, declassEvent, and declassException.
type DeclassResolver struct{}

func (r *DeclassResolver) Fields() []string {
	return []string{"declassDate", "declassEvent", "declassException"}
}

func (r *DeclassResolver) Resolve(ism *model.ISM, reg *refdata.Registry) []guidance.FieldGuidance {
	if ism.Classification != model.ClassificationC && ism.Classification != model.ClassificationS {
		return notApplicableAll(r.Fields(), "Only applicable for Confidential or Secret classifications")
	}

	var results []guidance.FieldGuidance

	// declassDate and declassEvent are mutually exclusive.
	if ism.DeclassEvent != "" {
		results = append(results, guidance.FieldGuidance{
			Field:  "declassDate",
			Status: guidance.StatusLocked,
			Reason: "Cannot set declassDate when declassEvent is set (mutually exclusive)",
		})
	} else {
		results = append(results, guidance.FieldGuidance{
			Field:  "declassDate",
			Status: guidance.StatusAvailable,
		})
	}

	if ism.DeclassDate != "" {
		results = append(results, guidance.FieldGuidance{
			Field:  "declassEvent",
			Status: guidance.StatusLocked,
			Reason: "Cannot set declassEvent when declassDate is set (mutually exclusive)",
		})
	} else {
		results = append(results, guidance.FieldGuidance{
			Field:  "declassEvent",
			Status: guidance.StatusAvailable,
		})
	}

	// declassException — always available for classified, with allowed values.
	exceptions := make([]guidance.AllowedValue, len(reg.DeclassExceptions))
	for i, e := range reg.DeclassExceptions {
		exceptions[i] = guidance.AllowedValue{Code: e.Code, Label: e.Label}
	}
	results = append(results, guidance.FieldGuidance{
		Field:         "declassException",
		Status:        guidance.StatusAvailable,
		AllowedValues: exceptions,
	})

	return results
}

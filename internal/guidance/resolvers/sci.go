package resolvers

import (
	"expr.ai/ism-api/internal/guidance"
	"expr.ai/ism-api/internal/model"
	"expr.ai/ism-api/internal/refdata"
)

// SCIResolver provides guidance for the sciControls field.
// SCI controls are only applicable at TOP SECRET classification.
type SCIResolver struct{}

func (r *SCIResolver) Fields() []string {
	return []string{"sciControls"}
}

func (r *SCIResolver) Resolve(ism *model.ISM, reg *refdata.Registry) []guidance.FieldGuidance {
	if !ism.Classification.AtLeast(model.ClassificationTS) {
		return []guidance.FieldGuidance{{
			Field:  "sciControls",
			Status: guidance.StatusNotApplicable,
			Reason: "SCI controls require TOP SECRET classification",
		}}
	}

	var allowed []guidance.AllowedValue
	for _, sc := range reg.SCIControls {
		allowed = append(allowed, guidance.AllowedValue{
			Code:  sc.Code,
			Label: sc.Label,
			Type:  sc.Category,
		})
	}

	return []guidance.FieldGuidance{{
		Field:         "sciControls",
		Status:        guidance.StatusAvailable,
		AllowedValues: allowed,
	}}
}

package resolvers

import (
	"github.com/danielmontgomery/ism-api/internal/guidance"
	"github.com/danielmontgomery/ism-api/internal/model"
	"github.com/danielmontgomery/ism-api/internal/refdata"
)

// AuthorityResolver provides guidance for classification authority fields:
// classifiedBy, classificationReason, derivativelyClassifiedBy, derivedFrom,
// and compilationReason.
type AuthorityResolver struct{}

func (r *AuthorityResolver) Fields() []string {
	return []string{
		"classifiedBy",
		"classificationReason",
		"derivativelyClassifiedBy",
		"derivedFrom",
		"compilationReason",
	}
}

func (r *AuthorityResolver) Resolve(ism *model.ISM, reg *refdata.Registry) []guidance.FieldGuidance {
	if ism.Classification != model.ClassificationC && ism.Classification != model.ClassificationS {
		return notApplicableAll(r.Fields(), "Only applicable for Confidential or Secret classifications")
	}

	var results []guidance.FieldGuidance

	// classifiedBy — available for original classification.
	results = append(results, guidance.FieldGuidance{
		Field:  "classifiedBy",
		Status: guidance.StatusAvailable,
	})

	// classificationReason — required when classifiedBy is set.
	if ism.ClassifiedBy != "" {
		results = append(results, guidance.FieldGuidance{
			Field:    "classificationReason",
			Status:   guidance.StatusRequired,
			Required: true,
		})
	} else {
		results = append(results, guidance.FieldGuidance{
			Field:      "classificationReason",
			Status:     guidance.StatusAvailable,
			RequiredIf: "classifiedBy is set",
		})
	}

	// derivativelyClassifiedBy — available for derivative classification.
	results = append(results, guidance.FieldGuidance{
		Field:  "derivativelyClassifiedBy",
		Status: guidance.StatusAvailable,
	})

	// derivedFrom — required when derivativelyClassifiedBy is set.
	if ism.DerivativelyClassifiedBy != "" {
		results = append(results, guidance.FieldGuidance{
			Field:    "derivedFrom",
			Status:   guidance.StatusRequired,
			Required: true,
		})
	} else {
		results = append(results, guidance.FieldGuidance{
			Field:      "derivedFrom",
			Status:     guidance.StatusAvailable,
			RequiredIf: "derivativelyClassifiedBy is set",
		})
	}

	// compilationReason — always available for classified docs.
	results = append(results, guidance.FieldGuidance{
		Field:  "compilationReason",
		Status: guidance.StatusAvailable,
	})

	return results
}

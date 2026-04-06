package resolvers

import (
	"github.com/danielmontgomery/ism-api/internal/guidance"
	"github.com/danielmontgomery/ism-api/internal/model"
	"github.com/danielmontgomery/ism-api/internal/refdata"
)

// ClassificationResolver provides guidance for classification, ownerProducer,
// joint, and version fields.
type ClassificationResolver struct{}

func (r *ClassificationResolver) Fields() []string {
	return []string{"classification", "ownerProducer", "joint", "version"}
}

func (r *ClassificationResolver) Resolve(ism *model.ISM, reg *refdata.Registry) []guidance.FieldGuidance {
	var results []guidance.FieldGuidance

	// classification — always required.
	classValues := make([]guidance.AllowedValue, len(reg.Classifications))
	for i, c := range reg.Classifications {
		classValues[i] = guidance.AllowedValue{Code: string(c.Code), Label: c.Label}
	}
	results = append(results, guidance.FieldGuidance{
		Field:         "classification",
		Status:        guidance.StatusRequired,
		Required:      true,
		AllowedValues: classValues,
	})

	// ownerProducer — required for C/S, available for CUI, not_applicable for U.
	switch ism.Classification {
	case model.ClassificationC, model.ClassificationS:
		results = append(results, guidance.FieldGuidance{
			Field:    "ownerProducer",
			Status:   guidance.StatusRequired,
			Required: true,
		})
	case model.ClassificationCUI:
		results = append(results, guidance.FieldGuidance{
			Field:  "ownerProducer",
			Status: guidance.StatusAvailable,
		})
	default:
		results = append(results, guidance.FieldGuidance{
			Field:  "ownerProducer",
			Status: guidance.StatusNotApplicable,
			Reason: "ownerProducer is not required for Unclassified documents",
		})
	}

	// joint — required when multiple ownerProducers, not applicable otherwise.
	if len(ism.OwnerProducer) > 1 {
		results = append(results, guidance.FieldGuidance{
			Field:    "joint",
			Status:   guidance.StatusRequired,
			Required: true,
			Reason:   "joint must be true when multiple ownerProducers are specified",
		})
	} else {
		results = append(results, guidance.FieldGuidance{
			Field:  "joint",
			Status: guidance.StatusNotApplicable,
			Reason: "joint only applies when multiple ownerProducers are present",
		})
	}

	// version — always available.
	results = append(results, guidance.FieldGuidance{
		Field:  "version",
		Status: guidance.StatusAvailable,
	})

	return results
}

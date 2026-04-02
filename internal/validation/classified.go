package validation

import (
	"github.com/danielmontgomery/ism-api/internal/model"
	"github.com/danielmontgomery/ism-api/internal/refdata"
)

// ClassifiedRule validates classification authority block fields for classified
// (C, S) markings. Missing authority fields produce warnings since they may be
// populated at different stages of the document lifecycle.
type ClassifiedRule struct{}

func (r *ClassifiedRule) Name() string { return "classified" }

func (r *ClassifiedRule) Applies(ism *model.ISM) bool {
	return ism.Classification.AtLeast(model.ClassificationC)
}

func (r *ClassifiedRule) Validate(ism *model.ISM, _ *refdata.Registry) *ValidationResult {
	res := &ValidationResult{Valid: true}

	hasOriginal := ism.ClassifiedBy != ""
	hasDerivative := ism.DerivativelyClassifiedBy != "" || ism.DerivedFrom != ""

	// At least one authority type should be present.
	if !hasOriginal && !hasDerivative {
		res.AddWarning("classifiedBy", "classified.missing_authority",
			"classified markings should include classifiedBy or derivativelyClassifiedBy/derivedFrom")
	}

	// Original classification requires a classificationReason.
	if hasOriginal && ism.ClassificationReason == "" {
		res.AddWarning("classificationReason", "classified.missing_reason",
			"original classification should include classificationReason")
	}

	// Derivative classification requires derivedFrom.
	if ism.DerivativelyClassifiedBy != "" && ism.DerivedFrom == "" {
		res.AddWarning("derivedFrom", "classified.missing_derived_from",
			"derivativelyClassifiedBy should be accompanied by derivedFrom")
	}

	return res
}

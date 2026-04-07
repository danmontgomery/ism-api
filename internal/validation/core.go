package validation

import (
	"expr.ai/ism-api/internal/model"
	"expr.ai/ism-api/internal/refdata"
)

// CoreRule validates fundamental ISM fields: classification enum validity,
// ownerProducer requirements, and joint consistency.
type CoreRule struct{}

func (r *CoreRule) Name() string { return "core" }

func (r *CoreRule) Applies(_ *model.ISM) bool { return true }

func (r *CoreRule) Validate(ism *model.ISM, reg *refdata.Registry) *ValidationResult {
	res := &ValidationResult{Valid: true}

	// Classification must be a valid enum value.
	if !ism.Classification.Valid() {
		res.AddError("classification", "core.invalid_classification",
			"classification must be one of: U, CUI, C, S, TS")
		return res // Further checks depend on valid classification.
	}

	// ownerProducer is required for classified (C, S) markings.
	if ism.Classification.AtLeast(model.ClassificationC) {
		if len(ism.OwnerProducer) == 0 {
			res.AddError("ownerProducer", "core.owner_producer_required",
				"ownerProducer is required for classified markings")
		}
	}

	// Joint consistency: joint must be true iff len(ownerProducer) > 1.
	if ism.Joint && len(ism.OwnerProducer) < 2 {
		res.AddError("joint", "core.joint_requires_multiple_owners",
			"joint=true requires at least two ownerProducer entries")
	}
	if len(ism.OwnerProducer) > 1 && !ism.Joint {
		res.AddError("joint", "core.joint_required_for_multiple_owners",
			"joint must be true when multiple ownerProducers are specified")
	}

	// Validate ownerProducer country codes.
	for _, code := range ism.OwnerProducer {
		if !reg.ValidCountryCode(code) {
			res.AddError("ownerProducer", "core.invalid_owner_producer",
				"unknown ownerProducer country code: "+code)
		}
	}

	return res
}

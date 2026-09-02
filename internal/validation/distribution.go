package validation

import (
	"dmontgomery/ism-api/internal/model"
	"dmontgomery/ism-api/internal/refdata"
)

// DistributionRule validates distribution statement codes (A-F) and
// classification constraints (e.g., Statement A is U/CUI only).
type DistributionRule struct{}

func (r *DistributionRule) Name() string { return "distribution" }

func (r *DistributionRule) Applies(ism *model.ISM) bool {
	return ism.DistributionStatement != ""
}

func (r *DistributionRule) Validate(ism *model.ISM, reg *refdata.Registry) *ValidationResult {
	res := &ValidationResult{Valid: true}

	// Code must be a valid distribution statement (A-F).
	if !reg.ValidDistributionStatement(ism.DistributionStatement) {
		res.AddError("distributionStatement", "distribution.invalid_code",
			"unknown distribution statement code: "+ism.DistributionStatement)
		return res
	}

	// Check classification constraints.
	for _, constraint := range refdata.DistributionClassificationConstraints() {
		if constraint.Statement != ism.DistributionStatement {
			continue
		}
		allowed := false
		for _, lvl := range constraint.AllowedLevels {
			if ism.Classification == lvl {
				allowed = true
				break
			}
		}
		if !allowed {
			res.AddError("distributionStatement", "distribution.classification_mismatch",
				"distribution statement "+ism.DistributionStatement+
					" is not permitted at classification "+ism.Classification.String())
		}
	}

	return res
}

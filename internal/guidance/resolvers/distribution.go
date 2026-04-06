package resolvers

import (
	"expr.ai/ism-api/internal/guidance"
	"expr.ai/ism-api/internal/model"
	"expr.ai/ism-api/internal/refdata"
)

// DistributionResolver provides guidance for distributionStatement and
// third-party distribution fields.
type DistributionResolver struct{}

func (r *DistributionResolver) Fields() []string {
	return []string{
		"distributionStatement",
		"3rdPartyDistributionStatement",
		"3rdPartyDistributionWarning",
		"3rdPartyDistributionContract",
		"copyright",
	}
}

func (r *DistributionResolver) Resolve(ism *model.ISM, reg *refdata.Registry) []guidance.FieldGuidance {
	var results []guidance.FieldGuidance

	// distributionStatement — always available, filtered by classification.
	constraints := refdata.DistributionClassificationConstraints()
	var allowed []guidance.AllowedValue
	for _, ds := range reg.DistributionStatements {
		if !statementAllowedForClassification(ds.Code, ism.Classification, constraints) {
			continue
		}
		allowed = append(allowed, guidance.AllowedValue{Code: ds.Code, Label: ds.Label})
	}
	results = append(results, guidance.FieldGuidance{
		Field:         "distributionStatement",
		Status:        guidance.StatusAvailable,
		AllowedValues: allowed,
	})

	// 3rdParty fields — available when distributionStatement is set.
	if ism.DistributionStatement != "" {
		results = append(results, guidance.FieldGuidance{
			Field:  "3rdPartyDistributionStatement",
			Status: guidance.StatusAvailable,
		})

		if ism.ThirdPartyDistributionStatement != "" {
			results = append(results, guidance.FieldGuidance{
				Field:    "3rdPartyDistributionWarning",
				Status:   guidance.StatusRequired,
				Required: true,
			})
			results = append(results, guidance.FieldGuidance{
				Field:    "3rdPartyDistributionContract",
				Status:   guidance.StatusRequired,
				Required: true,
			})
		} else {
			results = append(results, guidance.FieldGuidance{
				Field:      "3rdPartyDistributionWarning",
				Status:     guidance.StatusNotApplicable,
				RequiredIf: "3rdPartyDistributionStatement is set",
				Reason:     "Only applicable when 3rdPartyDistributionStatement is provided",
			})
			results = append(results, guidance.FieldGuidance{
				Field:      "3rdPartyDistributionContract",
				Status:     guidance.StatusNotApplicable,
				RequiredIf: "3rdPartyDistributionStatement is set",
				Reason:     "Only applicable when 3rdPartyDistributionStatement is provided",
			})
		}
	} else {
		results = append(results,
			guidance.FieldGuidance{
				Field:  "3rdPartyDistributionStatement",
				Status: guidance.StatusNotApplicable,
				Reason: "Only applicable when distributionStatement is set",
			},
			guidance.FieldGuidance{
				Field:  "3rdPartyDistributionWarning",
				Status: guidance.StatusNotApplicable,
				Reason: "Only applicable when distributionStatement is set",
			},
			guidance.FieldGuidance{
				Field:  "3rdPartyDistributionContract",
				Status: guidance.StatusNotApplicable,
				Reason: "Only applicable when distributionStatement is set",
			},
		)
	}

	// copyright — always available.
	results = append(results, guidance.FieldGuidance{
		Field:  "copyright",
		Status: guidance.StatusAvailable,
	})

	return results
}

func statementAllowedForClassification(code string, class model.Classification, constraints []refdata.DistributionConstraint) bool {
	for _, c := range constraints {
		if c.Statement != code {
			continue
		}
		for _, allowed := range c.AllowedLevels {
			if class == allowed {
				return true
			}
		}
		return false
	}
	return true // no constraint means allowed at all levels
}

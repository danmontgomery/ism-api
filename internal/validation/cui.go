package validation

import (
	"sort"

	"dmontgomery/ism-api/internal/model"
	"dmontgomery/ism-api/internal/refdata"
)

// CUIRule validates CUI-specific fields: category validity, alphabetization,
// and SP- prefix requirements for specified categories.
type CUIRule struct{}

func (r *CUIRule) Name() string { return "cui" }

func (r *CUIRule) Applies(ism *model.ISM) bool {
	return ism.Classification == model.ClassificationCUI
}

func (r *CUIRule) Validate(ism *model.ISM, reg *refdata.Registry) *ValidationResult {
	res := &ValidationResult{Valid: true}

	// Each category must be a known code.
	for _, cat := range ism.CategoryMarkings {
		if !reg.ValidCUICategory(cat) {
			res.AddError("categoryMarkings", "cui.invalid_category",
				"unknown CUI category: "+cat)
		}
	}

	// Categories must be in alphabetical order.
	if len(ism.CategoryMarkings) > 1 {
		if !sort.StringsAreSorted(ism.CategoryMarkings) {
			res.AddError("categoryMarkings", "cui.not_alphabetized",
				"CUI category markings must be in alphabetical order")
		}
	}

	// Specified categories must have SP- prefix (validates reference data integrity
	// from the consumer side — a category marked type=specified should start with SP-).
	for _, cat := range ism.CategoryMarkings {
		cuiCat := lookupCUICategory(reg, cat)
		if cuiCat != nil && cuiCat.Type == "specified" {
			if len(cat) < 3 || cat[:3] != "SP-" {
				res.AddError("categoryMarkings", "cui.specified_requires_sp_prefix",
					"specified CUI category must have SP- prefix: "+cat)
			}
		}
	}

	return res
}

func lookupCUICategory(reg *refdata.Registry, code string) *refdata.CUICategory {
	for i := range reg.CUICategories {
		if reg.CUICategories[i].Code == code {
			return &reg.CUICategories[i]
		}
	}
	return nil
}

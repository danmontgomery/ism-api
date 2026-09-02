package validation

import (
	"dmontgomery/ism-api/internal/model"
	"dmontgomery/ism-api/internal/refdata"
)

// Rule is the interface that all validation rules implement.
type Rule interface {
	// Name returns a short identifier for the rule (e.g. "core.classification").
	Name() string
	// Applies returns true if this rule is relevant to the given ISM.
	Applies(ism *model.ISM) bool
	// Validate checks the ISM and returns any findings.
	Validate(ism *model.ISM, reg *refdata.Registry) *ValidationResult
}

// Engine runs all registered rules against an ISM, aggregating results
// without short-circuiting.
type Engine struct {
	rules []Rule
	reg   *refdata.Registry
}

// NewEngine returns an engine loaded with all standard rules.
func NewEngine(reg *refdata.Registry) *Engine {
	return &Engine{
		reg: reg,
		rules: []Rule{
			&CoreRule{},
			&CUIRule{},
			&ClassifiedRule{},
			&DisseminationRule{},
			&DistributionRule{},
			&ThirdPartyRule{},
			&DeclassRule{},
			&FGIRule{},
			&NonUSControlsRule{},
			&ExemptFromRule{},
			&CompliesWithRule{},
			&SCIRule{},
			&AtomicEnergyRule{},
		},
	}
}

// Validate runs every applicable rule and merges results.
func (e *Engine) Validate(ism *model.ISM) *ValidationResult {
	result := &ValidationResult{Valid: true}
	for _, rule := range e.rules {
		if rule.Applies(ism) {
			result.Merge(rule.Validate(ism, e.reg))
		}
	}
	return result
}

package guidance

import (
	"dmontgomery/ism-api/internal/model"
	"dmontgomery/ism-api/internal/refdata"
)

// Resolver inspects partial ISM state and returns field guidance.
// Defined here so the engine can accept resolvers without an import cycle.
// The resolvers package implements this interface.
type Resolver interface {
	Fields() []string
	Resolve(ism *model.ISM, reg *refdata.Registry) []FieldGuidance
}

// Engine runs all registered resolvers against a partial ISM and aggregates results.
type Engine struct {
	registry  *refdata.Registry
	resolvers []Resolver
}

// NewEngine creates a guidance engine with the given registry and resolvers.
func NewEngine(reg *refdata.Registry, resolvers ...Resolver) *Engine {
	return &Engine{
		registry:  reg,
		resolvers: resolvers,
	}
}

// Evaluate runs all resolvers against the given partial ISM and returns
// the combined field guidance.
func (e *Engine) Evaluate(ism *model.ISM) []FieldGuidance {
	var results []FieldGuidance
	for _, r := range e.resolvers {
		results = append(results, r.Resolve(ism, e.registry)...)
	}
	return results
}

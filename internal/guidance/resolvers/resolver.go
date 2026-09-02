package resolvers

import (
	"dmontgomery/ism-api/internal/guidance"
	"dmontgomery/ism-api/internal/model"
	"dmontgomery/ism-api/internal/refdata"
)

// Resolver inspects partial ISM state and returns field guidance.
type Resolver interface {
	// Fields returns the ISM field names this resolver is responsible for.
	Fields() []string

	// Resolve inspects the current ISM state and returns guidance for each field.
	Resolve(ism *model.ISM, reg *refdata.Registry) []guidance.FieldGuidance
}

package guidance

import (
	"testing"

	"dmontgomery/ism-api/internal/model"
	"dmontgomery/ism-api/internal/refdata"
)

// stubResolver returns fixed guidance for testing the engine aggregation.
type stubResolver struct {
	fields  []string
	results []FieldGuidance
}

func (s *stubResolver) Fields() []string { return s.fields }
func (s *stubResolver) Resolve(_ *model.ISM, _ *refdata.Registry) []FieldGuidance {
	return s.results
}

func TestEngine_Evaluate(t *testing.T) {
	reg := refdata.NewRegistry()

	r1 := &stubResolver{
		fields: []string{"fieldA"},
		results: []FieldGuidance{
			{Field: "fieldA", Status: StatusRequired, Required: true},
		},
	}
	r2 := &stubResolver{
		fields: []string{"fieldB", "fieldC"},
		results: []FieldGuidance{
			{Field: "fieldB", Status: StatusAvailable},
			{Field: "fieldC", Status: StatusNotApplicable, Reason: "test reason"},
		},
	}

	engine := NewEngine(reg, r1, r2)
	ism := &model.ISM{Classification: model.ClassificationU}
	results := engine.Evaluate(ism)

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}

	checks := map[string]FieldStatus{
		"fieldA": StatusRequired,
		"fieldB": StatusAvailable,
		"fieldC": StatusNotApplicable,
	}
	for _, fg := range results {
		want, ok := checks[fg.Field]
		if !ok {
			t.Errorf("unexpected field %q", fg.Field)
			continue
		}
		if fg.Status != want {
			t.Errorf("field %q: status = %q, want %q", fg.Field, fg.Status, want)
		}
	}
}

func TestEngine_NoResolvers(t *testing.T) {
	reg := refdata.NewRegistry()
	engine := NewEngine(reg)
	results := engine.Evaluate(&model.ISM{})
	if len(results) != 0 {
		t.Errorf("expected 0 results with no resolvers, got %d", len(results))
	}
}

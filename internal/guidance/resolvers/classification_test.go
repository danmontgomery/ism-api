package resolvers

import (
	"testing"

	"expr.ai/ism-api/internal/guidance"
	"expr.ai/ism-api/internal/model"
	"expr.ai/ism-api/internal/refdata"
)

func TestClassificationResolver(t *testing.T) {
	reg := refdata.NewRegistry()
	r := &ClassificationResolver{}

	tests := []struct {
		name   string
		ism    model.ISM
		checks []fieldCheck
	}{
		{
			name: "empty ISM — classification required, ownerProducer not applicable",
			ism:  model.ISM{},
			checks: []fieldCheck{
				{field: "classification", status: guidance.StatusRequired, required: true, hasAllowed: true},
				{field: "ownerProducer", status: guidance.StatusNotApplicable},
				{field: "joint", status: guidance.StatusNotApplicable},
				{field: "version", status: guidance.StatusAvailable},
			},
		},
		{
			name: "U — ownerProducer not applicable",
			ism:  model.ISM{Classification: model.ClassificationU},
			checks: []fieldCheck{
				{field: "classification", status: guidance.StatusRequired, required: true, hasAllowed: true},
				{field: "ownerProducer", status: guidance.StatusNotApplicable},
				{field: "joint", status: guidance.StatusNotApplicable},
			},
		},
		{
			name: "CUI — ownerProducer available",
			ism:  model.ISM{Classification: model.ClassificationCUI},
			checks: []fieldCheck{
				{field: "classification", status: guidance.StatusRequired, required: true},
				{field: "ownerProducer", status: guidance.StatusAvailable},
			},
		},
		{
			name: "C — ownerProducer required",
			ism:  model.ISM{Classification: model.ClassificationC},
			checks: []fieldCheck{
				{field: "ownerProducer", status: guidance.StatusRequired, required: true},
			},
		},
		{
			name: "S — ownerProducer required",
			ism:  model.ISM{Classification: model.ClassificationS},
			checks: []fieldCheck{
				{field: "ownerProducer", status: guidance.StatusRequired, required: true},
			},
		},
		{
			name: "single ownerProducer — joint not applicable",
			ism: model.ISM{
				Classification: model.ClassificationC,
				OwnerProducer:  []string{"USA"},
			},
			checks: []fieldCheck{
				{field: "joint", status: guidance.StatusNotApplicable},
			},
		},
		{
			name: "multiple ownerProducers — joint required",
			ism: model.ISM{
				Classification: model.ClassificationC,
				OwnerProducer:  []string{"USA", "GBR"},
			},
			checks: []fieldCheck{
				{field: "joint", status: guidance.StatusRequired, required: true},
			},
		},
		{
			name: "no ownerProducer — joint not applicable",
			ism:  model.ISM{Classification: model.ClassificationU},
			checks: []fieldCheck{
				{field: "joint", status: guidance.StatusNotApplicable},
			},
		},
		{
			name: "classification allowed values include all levels",
			ism:  model.ISM{},
			checks: []fieldCheck{
				{field: "classification", required: true, hasAllowed: true, allowedCount: 5},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := r.Resolve(&tt.ism, reg)
			for _, check := range tt.checks {
				assertFieldGuidance(t, results, check)
			}
		})
	}
}

// fieldCheck describes expected properties for a single field in guidance results.
type fieldCheck struct {
	field        string
	status       guidance.FieldStatus
	required     bool
	hasAllowed   bool
	allowedCount int
}

// assertFieldGuidance finds a field in results and verifies expected properties.
func assertFieldGuidance(t *testing.T, results []guidance.FieldGuidance, check fieldCheck) {
	t.Helper()
	var found *guidance.FieldGuidance
	for i := range results {
		if results[i].Field == check.field {
			found = &results[i]
			break
		}
	}
	if found == nil {
		t.Errorf("field %q not found in results", check.field)
		return
	}
	if check.status != "" && found.Status != check.status {
		t.Errorf("field %q: status = %q, want %q", check.field, found.Status, check.status)
	}
	if check.required && !found.Required {
		t.Errorf("field %q: expected required=true", check.field)
	}
	if !check.required && found.Required {
		t.Errorf("field %q: expected required=false", check.field)
	}
	if check.hasAllowed && len(found.AllowedValues) == 0 {
		t.Errorf("field %q: expected allowedValues to be non-empty", check.field)
	}
	if check.allowedCount > 0 && len(found.AllowedValues) != check.allowedCount {
		t.Errorf("field %q: allowedValues count = %d, want %d", check.field, len(found.AllowedValues), check.allowedCount)
	}
}

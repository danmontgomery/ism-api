package resolvers

import (
	"testing"

	"expr.ai/ism-api/internal/guidance"
	"expr.ai/ism-api/internal/model"
	"expr.ai/ism-api/internal/refdata"
)

func TestSCIResolver(t *testing.T) {
	reg := refdata.NewRegistry()
	r := &SCIResolver{}

	tests := []struct {
		name   string
		ism    model.ISM
		checks []fieldCheck
	}{
		{
			name: "U — sciControls not applicable",
			ism:  model.ISM{Classification: model.ClassificationU},
			checks: []fieldCheck{
				{field: "sciControls", status: guidance.StatusNotApplicable},
			},
		},
		{
			name: "CUI — sciControls not applicable",
			ism:  model.ISM{Classification: model.ClassificationCUI},
			checks: []fieldCheck{
				{field: "sciControls", status: guidance.StatusNotApplicable},
			},
		},
		{
			name: "C — sciControls not applicable",
			ism:  model.ISM{Classification: model.ClassificationC},
			checks: []fieldCheck{
				{field: "sciControls", status: guidance.StatusNotApplicable},
			},
		},
		{
			name: "S — sciControls not applicable",
			ism:  model.ISM{Classification: model.ClassificationS},
			checks: []fieldCheck{
				{field: "sciControls", status: guidance.StatusNotApplicable},
			},
		},
		{
			name: "TS — sciControls available with 21 controls",
			ism:  model.ISM{Classification: model.ClassificationTS},
			checks: []fieldCheck{
				{field: "sciControls", status: guidance.StatusAvailable, hasAllowed: true, allowedCount: 21},
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

func TestSCIResolver_Fields(t *testing.T) {
	r := &SCIResolver{}
	fields := r.Fields()
	if len(fields) != 1 || fields[0] != "sciControls" {
		t.Errorf("Fields() = %v, want [sciControls]", fields)
	}
}

func TestSCIResolver_AllowedValuesHaveCategory(t *testing.T) {
	reg := refdata.NewRegistry()
	r := &SCIResolver{}

	results := r.Resolve(&model.ISM{Classification: model.ClassificationTS}, reg)
	sci := findField(results, "sciControls")
	if sci == nil {
		t.Fatal("sciControls not found")
	}

	for _, av := range sci.AllowedValues {
		if av.Code == "" || av.Label == "" || av.Type == "" {
			t.Errorf("allowed value missing fields: %+v", av)
		}
	}
}

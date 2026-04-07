package resolvers

import (
	"testing"

	"expr.ai/ism-api/internal/guidance"
	"expr.ai/ism-api/internal/model"
	"expr.ai/ism-api/internal/refdata"
)

func TestDeclassResolver(t *testing.T) {
	reg := refdata.NewRegistry()
	r := &DeclassResolver{}

	tests := []struct {
		name   string
		ism    model.ISM
		checks []fieldCheck
	}{
		{
			name: "U — all declass fields not applicable",
			ism:  model.ISM{Classification: model.ClassificationU},
			checks: []fieldCheck{
				{field: "declassDate", status: guidance.StatusNotApplicable},
				{field: "declassEvent", status: guidance.StatusNotApplicable},
				{field: "declassException", status: guidance.StatusNotApplicable},
			},
		},
		{
			name: "CUI — all declass fields not applicable",
			ism:  model.ISM{Classification: model.ClassificationCUI},
			checks: []fieldCheck{
				{field: "declassDate", status: guidance.StatusNotApplicable},
				{field: "declassEvent", status: guidance.StatusNotApplicable},
				{field: "declassException", status: guidance.StatusNotApplicable},
			},
		},
		{
			name: "C — all declass fields available",
			ism:  model.ISM{Classification: model.ClassificationC},
			checks: []fieldCheck{
				{field: "declassDate", status: guidance.StatusAvailable},
				{field: "declassEvent", status: guidance.StatusAvailable},
				{field: "declassException", status: guidance.StatusAvailable, hasAllowed: true, allowedCount: 25},
			},
		},
		{
			name: "S — all declass fields available",
			ism:  model.ISM{Classification: model.ClassificationS},
			checks: []fieldCheck{
				{field: "declassDate", status: guidance.StatusAvailable},
				{field: "declassEvent", status: guidance.StatusAvailable},
				{field: "declassException", status: guidance.StatusAvailable, hasAllowed: true},
			},
		},
		{
			name: "S with declassDate — declassEvent locked",
			ism: model.ISM{
				Classification: model.ClassificationS,
				DeclassDate:    "20360101",
			},
			checks: []fieldCheck{
				{field: "declassDate", status: guidance.StatusAvailable},
				{field: "declassEvent", status: guidance.StatusLocked},
				{field: "declassException", status: guidance.StatusAvailable, hasAllowed: true},
			},
		},
		{
			name: "S with declassEvent — declassDate locked",
			ism: model.ISM{
				Classification: model.ClassificationS,
				DeclassEvent:   "Completion of Operation EXAMPLE",
			},
			checks: []fieldCheck{
				{field: "declassDate", status: guidance.StatusLocked},
				{field: "declassEvent", status: guidance.StatusAvailable},
				{field: "declassException", status: guidance.StatusAvailable, hasAllowed: true},
			},
		},
		{
			name: "C with declassDate — same mutual exclusion",
			ism: model.ISM{
				Classification: model.ClassificationC,
				DeclassDate:    "20360601",
			},
			checks: []fieldCheck{
				{field: "declassDate", status: guidance.StatusAvailable},
				{field: "declassEvent", status: guidance.StatusLocked},
			},
		},
		{
			name: "TS — all declass fields available",
			ism:  model.ISM{Classification: model.ClassificationTS},
			checks: []fieldCheck{
				{field: "declassDate", status: guidance.StatusAvailable},
				{field: "declassEvent", status: guidance.StatusAvailable},
				{field: "declassException", status: guidance.StatusAvailable, hasAllowed: true},
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

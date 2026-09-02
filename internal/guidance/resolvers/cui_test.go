package resolvers

import (
	"testing"

	"dmontgomery/ism-api/internal/guidance"
	"dmontgomery/ism-api/internal/model"
	"dmontgomery/ism-api/internal/refdata"
)

func TestCUIResolver(t *testing.T) {
	reg := refdata.NewRegistry()
	r := &CUIResolver{}

	tests := []struct {
		name   string
		ism    model.ISM
		checks []fieldCheck
	}{
		{
			name: "U — all CUI fields not applicable",
			ism:  model.ISM{Classification: model.ClassificationU},
			checks: []fieldCheck{
				{field: "categoryMarkings", status: guidance.StatusNotApplicable},
				{field: "controlledByName", status: guidance.StatusNotApplicable},
				{field: "controlledByOffice", status: guidance.StatusNotApplicable},
				{field: "poc", status: guidance.StatusNotApplicable},
			},
		},
		{
			name: "C — all CUI fields not applicable",
			ism:  model.ISM{Classification: model.ClassificationC},
			checks: []fieldCheck{
				{field: "categoryMarkings", status: guidance.StatusNotApplicable},
				{field: "controlledByName", status: guidance.StatusNotApplicable},
			},
		},
		{
			name: "S — all CUI fields not applicable",
			ism:  model.ISM{Classification: model.ClassificationS},
			checks: []fieldCheck{
				{field: "categoryMarkings", status: guidance.StatusNotApplicable},
			},
		},
		{
			name: "CUI — categoryMarkings available with allowed values",
			ism:  model.ISM{Classification: model.ClassificationCUI},
			checks: []fieldCheck{
				{field: "categoryMarkings", status: guidance.StatusAvailable, hasAllowed: true, allowedCount: 64},
				{field: "controlledByName", status: guidance.StatusAvailable},
				{field: "controlledByOffice", status: guidance.StatusAvailable},
				{field: "poc", status: guidance.StatusAvailable},
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

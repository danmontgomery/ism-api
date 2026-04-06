package resolvers

import (
	"testing"

	"expr.ai/ism-api/internal/guidance"
	"expr.ai/ism-api/internal/model"
	"expr.ai/ism-api/internal/refdata"
)

func TestDisseminationResolver(t *testing.T) {
	reg := refdata.NewRegistry()
	r := &DisseminationResolver{}

	tests := []struct {
		name   string
		ism    model.ISM
		checks []fieldCheck
	}{
		{
			name: "U — only unclassified controls available",
			ism:  model.ISM{Classification: model.ClassificationU},
			checks: []fieldCheck{
				{field: "disseminationControls", status: guidance.StatusAvailable, hasAllowed: true},
				{field: "releasableTo", status: guidance.StatusNotApplicable},
				{field: "displayOnlyTo", status: guidance.StatusNotApplicable},
			},
		},
		{
			name: "S — all controls available including classified-only",
			ism:  model.ISM{Classification: model.ClassificationS},
			checks: []fieldCheck{
				{field: "disseminationControls", status: guidance.StatusAvailable, hasAllowed: true},
			},
		},
		{
			name: "S with REL — releasableTo required",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				DisseminationControls: []string{"REL"},
			},
			checks: []fieldCheck{
				{field: "releasableTo", status: guidance.StatusRequired, required: true},
				{field: "displayOnlyTo", status: guidance.StatusNotApplicable},
			},
		},
		{
			name: "S with EYES — releasableTo required",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				DisseminationControls: []string{"EYES"},
			},
			checks: []fieldCheck{
				{field: "releasableTo", status: guidance.StatusRequired, required: true},
			},
		},
		{
			name: "S with DISPLAY ONLY — displayOnlyTo required",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				DisseminationControls: []string{"DISPLAY ONLY"},
			},
			checks: []fieldCheck{
				{field: "displayOnlyTo", status: guidance.StatusRequired, required: true},
				{field: "releasableTo", status: guidance.StatusNotApplicable},
			},
		},
		{
			name: "U — classified-only controls filtered out",
			ism:  model.ISM{Classification: model.ClassificationU},
			checks: []fieldCheck{
				// U should not include NOFORN, OC, etc.
				{field: "disseminationControls", status: guidance.StatusAvailable, hasAllowed: true},
			},
		},
		{
			name: "S with NOFORN — releasableTo not applicable",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				DisseminationControls: []string{"NOFORN"},
			},
			checks: []fieldCheck{
				{field: "releasableTo", status: guidance.StatusNotApplicable},
				{field: "displayOnlyTo", status: guidance.StatusNotApplicable},
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

func TestDisseminationResolver_ClassificationFiltering(t *testing.T) {
	reg := refdata.NewRegistry()
	r := &DisseminationResolver{}

	// Unclassified should not include NOFORN, OC, IMCON, etc.
	uResults := r.Resolve(&model.ISM{Classification: model.ClassificationU}, reg)
	uControls := findField(uResults, "disseminationControls")
	if uControls == nil {
		t.Fatal("disseminationControls not found")
	}
	for _, av := range uControls.AllowedValues {
		if av.Code == "NOFORN" || av.Code == "OC" || av.Code == "IMCON" {
			t.Errorf("U should not include %q in allowed dissemination controls", av.Code)
		}
	}

	// Secret should include NOFORN.
	sResults := r.Resolve(&model.ISM{Classification: model.ClassificationS}, reg)
	sControls := findField(sResults, "disseminationControls")
	if sControls == nil {
		t.Fatal("disseminationControls not found")
	}
	foundNoforn := false
	for _, av := range sControls.AllowedValues {
		if av.Code == "NOFORN" {
			foundNoforn = true
			break
		}
	}
	if !foundNoforn {
		t.Error("S should include NOFORN in allowed dissemination controls")
	}
}

func findField(results []guidance.FieldGuidance, field string) *guidance.FieldGuidance {
	for i := range results {
		if results[i].Field == field {
			return &results[i]
		}
	}
	return nil
}

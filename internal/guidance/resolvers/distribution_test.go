package resolvers

import (
	"testing"

	"expr.ai/ism-api/internal/guidance"
	"expr.ai/ism-api/internal/model"
	"expr.ai/ism-api/internal/refdata"
)

func TestDistributionResolver(t *testing.T) {
	reg := refdata.NewRegistry()
	r := &DistributionResolver{}

	tests := []struct {
		name   string
		ism    model.ISM
		checks []fieldCheck
	}{
		{
			name: "U — all statements available, Statement A included",
			ism:  model.ISM{Classification: model.ClassificationU},
			checks: []fieldCheck{
				{field: "distributionStatement", status: guidance.StatusAvailable, hasAllowed: true, allowedCount: 6},
				{field: "3rdPartyDistributionStatement", status: guidance.StatusNotApplicable},
				{field: "3rdPartyDistributionWarning", status: guidance.StatusNotApplicable},
				{field: "3rdPartyDistributionContract", status: guidance.StatusNotApplicable},
				{field: "copyright", status: guidance.StatusAvailable},
			},
		},
		{
			name: "S — Statement A filtered out",
			ism:  model.ISM{Classification: model.ClassificationS},
			checks: []fieldCheck{
				{field: "distributionStatement", status: guidance.StatusAvailable, hasAllowed: true, allowedCount: 5},
			},
		},
		{
			name: "CUI — Statement A included",
			ism:  model.ISM{Classification: model.ClassificationCUI},
			checks: []fieldCheck{
				{field: "distributionStatement", status: guidance.StatusAvailable, hasAllowed: true, allowedCount: 6},
			},
		},
		{
			name: "with distributionStatement set — 3rdParty available",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				DistributionStatement: "B",
			},
			checks: []fieldCheck{
				{field: "3rdPartyDistributionStatement", status: guidance.StatusAvailable},
				{field: "3rdPartyDistributionWarning", status: guidance.StatusNotApplicable},
				{field: "3rdPartyDistributionContract", status: guidance.StatusNotApplicable},
			},
		},
		{
			name: "with 3rdPartyDistributionStatement set — warning and contract required",
			ism: model.ISM{
				Classification:                  model.ClassificationS,
				DistributionStatement:           "B",
				ThirdPartyDistributionStatement: "Category I",
			},
			checks: []fieldCheck{
				{field: "3rdPartyDistributionStatement", status: guidance.StatusAvailable},
				{field: "3rdPartyDistributionWarning", status: guidance.StatusRequired, required: true},
				{field: "3rdPartyDistributionContract", status: guidance.StatusRequired, required: true},
			},
		},
		{
			name: "copyright always available",
			ism:  model.ISM{Classification: model.ClassificationS},
			checks: []fieldCheck{
				{field: "copyright", status: guidance.StatusAvailable},
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

func TestDistributionResolver_StatementAFiltering(t *testing.T) {
	reg := refdata.NewRegistry()
	r := &DistributionResolver{}

	// Secret should NOT include Statement A.
	sResults := r.Resolve(&model.ISM{Classification: model.ClassificationS}, reg)
	dist := findField(sResults, "distributionStatement")
	if dist == nil {
		t.Fatal("distributionStatement not found")
	}
	for _, av := range dist.AllowedValues {
		if av.Code == "A" {
			t.Error("S should not include Statement A in allowed distribution statements")
		}
	}

	// Unclassified SHOULD include Statement A.
	uResults := r.Resolve(&model.ISM{Classification: model.ClassificationU}, reg)
	uDist := findField(uResults, "distributionStatement")
	if uDist == nil {
		t.Fatal("distributionStatement not found")
	}
	foundA := false
	for _, av := range uDist.AllowedValues {
		if av.Code == "A" {
			foundA = true
			break
		}
	}
	if !foundA {
		t.Error("U should include Statement A in allowed distribution statements")
	}
}

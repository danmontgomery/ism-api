package resolvers

import (
	"testing"

	"github.com/danielmontgomery/ism-api/internal/guidance"
	"github.com/danielmontgomery/ism-api/internal/model"
	"github.com/danielmontgomery/ism-api/internal/refdata"
)

func TestAuthorityResolver(t *testing.T) {
	reg := refdata.NewRegistry()
	r := &AuthorityResolver{}

	tests := []struct {
		name   string
		ism    model.ISM
		checks []fieldCheck
	}{
		{
			name: "U — all authority fields not applicable",
			ism:  model.ISM{Classification: model.ClassificationU},
			checks: []fieldCheck{
				{field: "classifiedBy", status: guidance.StatusNotApplicable},
				{field: "classificationReason", status: guidance.StatusNotApplicable},
				{field: "derivativelyClassifiedBy", status: guidance.StatusNotApplicable},
				{field: "derivedFrom", status: guidance.StatusNotApplicable},
				{field: "compilationReason", status: guidance.StatusNotApplicable},
			},
		},
		{
			name: "CUI — all authority fields not applicable",
			ism:  model.ISM{Classification: model.ClassificationCUI},
			checks: []fieldCheck{
				{field: "classifiedBy", status: guidance.StatusNotApplicable},
				{field: "classificationReason", status: guidance.StatusNotApplicable},
			},
		},
		{
			name: "C — authority fields available",
			ism:  model.ISM{Classification: model.ClassificationC},
			checks: []fieldCheck{
				{field: "classifiedBy", status: guidance.StatusAvailable},
				{field: "classificationReason", status: guidance.StatusAvailable},
				{field: "derivativelyClassifiedBy", status: guidance.StatusAvailable},
				{field: "derivedFrom", status: guidance.StatusAvailable},
				{field: "compilationReason", status: guidance.StatusAvailable},
			},
		},
		{
			name: "S — authority fields available",
			ism:  model.ISM{Classification: model.ClassificationS},
			checks: []fieldCheck{
				{field: "classifiedBy", status: guidance.StatusAvailable},
				{field: "compilationReason", status: guidance.StatusAvailable},
			},
		},
		{
			name: "S with classifiedBy — classificationReason required",
			ism: model.ISM{
				Classification: model.ClassificationS,
				ClassifiedBy:   "OCA Name",
			},
			checks: []fieldCheck{
				{field: "classifiedBy", status: guidance.StatusAvailable},
				{field: "classificationReason", status: guidance.StatusRequired, required: true},
				{field: "derivativelyClassifiedBy", status: guidance.StatusAvailable},
				{field: "derivedFrom", status: guidance.StatusAvailable},
			},
		},
		{
			name: "S with derivativelyClassifiedBy — derivedFrom required",
			ism: model.ISM{
				Classification:           model.ClassificationS,
				DerivativelyClassifiedBy: "Derivative Classifier",
			},
			checks: []fieldCheck{
				{field: "classifiedBy", status: guidance.StatusAvailable},
				{field: "classificationReason", status: guidance.StatusAvailable},
				{field: "derivativelyClassifiedBy", status: guidance.StatusAvailable},
				{field: "derivedFrom", status: guidance.StatusRequired, required: true},
			},
		},
		{
			name: "S with both classifiedBy and derivativelyClassifiedBy",
			ism: model.ISM{
				Classification:           model.ClassificationS,
				ClassifiedBy:             "OCA",
				DerivativelyClassifiedBy: "DCA",
			},
			checks: []fieldCheck{
				{field: "classificationReason", status: guidance.StatusRequired, required: true},
				{field: "derivedFrom", status: guidance.StatusRequired, required: true},
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

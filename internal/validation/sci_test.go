package validation

import (
	"testing"

	"expr.ai/ism-api/internal/model"
)

func TestSCIRule(t *testing.T) {
	tests := []struct {
		name     string
		ism      model.ISM
		wantCode string
		valid    bool
	}{
		{
			name: "SCI with TS — valid",
			ism: model.ISM{
				Classification: model.ClassificationTS,
				OwnerProducer:  []string{"USA"},
				SCIControls:    []string{"SI"},
			},
			valid: true,
		},
		{
			name: "SCI with S — requires TS",
			ism: model.ISM{
				Classification: model.ClassificationS,
				OwnerProducer:  []string{"USA"},
				SCIControls:    []string{"SI"},
			},
			wantCode: "sci.requires_ts",
			valid:    false,
		},
		{
			name: "invalid SCI control code",
			ism: model.ISM{
				Classification: model.ClassificationTS,
				OwnerProducer:  []string{"USA"},
				SCIControls:    []string{"BOGUS"},
			},
			wantCode: "sci.invalid_control",
			valid:    false,
		},
		{
			name: "multiple valid SCI controls",
			ism: model.ISM{
				Classification: model.ClassificationTS,
				OwnerProducer:  []string{"USA"},
				SCIControls:    []string{"SI", "TK"},
			},
			valid: true,
		},
	}

	rule := &SCIRule{}
	r := reg()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !rule.Applies(&tt.ism) {
				t.Fatal("rule should apply when SCIControls present")
			}
			res := rule.Validate(&tt.ism, r)
			if res.Valid != tt.valid {
				t.Errorf("valid=%v, want %v; errors: %+v", res.Valid, tt.valid, res.Errors)
			}
			if tt.wantCode != "" && !res.HasCode(tt.wantCode) {
				t.Errorf("expected code %q in %+v", tt.wantCode, res.Errors)
			}
		})
	}
}

func TestSCIRule_DoesNotApplyWhenEmpty(t *testing.T) {
	rule := &SCIRule{}
	if rule.Applies(&model.ISM{Classification: model.ClassificationTS}) {
		t.Error("SCI rule should not apply when no SCIControls")
	}
}

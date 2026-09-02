package validation

import (
	"testing"

	"dmontgomery/ism-api/internal/model"
	"dmontgomery/ism-api/internal/refdata"
)

func reg() *refdata.Registry { return refdata.NewRegistry() }

// --- Engine ---

func TestEngine_AggregatesWithoutShortCircuit(t *testing.T) {
	e := NewEngine(reg())
	// Invalid classification + declass on U = multiple errors from different rules.
	ism := &model.ISM{
		Classification: "BOGUS",
	}
	r := e.Validate(ism)
	if r.Valid {
		t.Error("expected invalid")
	}
	if len(r.Errors) == 0 {
		t.Error("expected errors")
	}
}

func TestEngine_ValidMinimalU(t *testing.T) {
	e := NewEngine(reg())
	r := e.Validate(&model.ISM{Classification: model.ClassificationU})
	if !r.Valid {
		t.Errorf("expected valid, got errors: %+v", r.Errors)
	}
}

func TestEngine_ValidSecret(t *testing.T) {
	e := NewEngine(reg())
	r := e.Validate(&model.ISM{
		Classification: model.ClassificationS,
		OwnerProducer:  []string{"USA"},
		ClassifiedBy:   "John Doe",
		ClassificationReason: "National Security",
		DeclassDate:    "20300101",
	})
	if !r.Valid {
		t.Errorf("expected valid, got errors: %+v", r.Errors)
	}
}

// --- CoreRule ---

func TestCoreRule(t *testing.T) {
	tests := []struct {
		name     string
		ism      model.ISM
		wantCode string
		valid    bool
	}{
		{
			name:  "valid U",
			ism:   model.ISM{Classification: model.ClassificationU},
			valid: true,
		},
		{
			name:     "invalid classification",
			ism:      model.ISM{Classification: "BOGUS"},
			wantCode: "core.invalid_classification",
			valid:    false,
		},
		{
			name:     "TS missing ownerProducer",
			ism:      model.ISM{Classification: model.ClassificationTS},
			wantCode: "core.owner_producer_required",
			valid:    false,
		},
		{
			name: "C missing ownerProducer",
			ism: model.ISM{
				Classification: model.ClassificationC,
			},
			wantCode: "core.owner_producer_required",
			valid:    false,
		},
		{
			name: "C with ownerProducer",
			ism: model.ISM{
				Classification: model.ClassificationC,
				OwnerProducer:  []string{"USA"},
			},
			valid: true,
		},
		{
			name: "S missing ownerProducer",
			ism: model.ISM{
				Classification: model.ClassificationS,
			},
			wantCode: "core.owner_producer_required",
			valid:    false,
		},
		{
			name: "joint with one owner",
			ism: model.ISM{
				Classification: model.ClassificationU,
				Joint:          true,
				OwnerProducer:  []string{"USA"},
			},
			wantCode: "core.joint_requires_multiple_owners",
			valid:    false,
		},
		{
			name: "joint with two owners",
			ism: model.ISM{
				Classification: model.ClassificationC,
				Joint:          true,
				OwnerProducer:  []string{"USA", "GBR"},
			},
			valid: true,
		},
		{
			name: "multiple owners without joint",
			ism: model.ISM{
				Classification: model.ClassificationC,
				OwnerProducer:  []string{"USA", "GBR"},
				Joint:          false,
			},
			wantCode: "core.joint_required_for_multiple_owners",
			valid:    false,
		},
		{
			name: "invalid ownerProducer code",
			ism: model.ISM{
				Classification: model.ClassificationC,
				OwnerProducer:  []string{"ZZZ"},
			},
			wantCode: "core.invalid_owner_producer",
			valid:    false,
		},
		{
			name: "U does not require ownerProducer",
			ism: model.ISM{
				Classification: model.ClassificationU,
			},
			valid: true,
		},
		{
			name: "CUI does not require ownerProducer",
			ism: model.ISM{
				Classification: model.ClassificationCUI,
			},
			valid: true,
		},
	}

	rule := &CoreRule{}
	r := reg()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

// --- CUIRule ---

func TestCUIRule(t *testing.T) {
	tests := []struct {
		name     string
		ism      model.ISM
		wantCode string
		valid    bool
	}{
		{
			name: "valid CUI no categories",
			ism: model.ISM{
				Classification: model.ClassificationCUI,
			},
			valid: true,
		},
		{
			name: "valid CUI with categories",
			ism: model.ISM{
				Classification:   model.ClassificationCUI,
				CategoryMarkings: []string{"CRIT", "PHYS"},
			},
			valid: true,
		},
		{
			name: "invalid category",
			ism: model.ISM{
				Classification:   model.ClassificationCUI,
				CategoryMarkings: []string{"BOGUS"},
			},
			wantCode: "cui.invalid_category",
			valid:    false,
		},
		{
			name: "categories not alphabetized",
			ism: model.ISM{
				Classification:   model.ClassificationCUI,
				CategoryMarkings: []string{"PHYS", "CRIT"},
			},
			wantCode: "cui.not_alphabetized",
			valid:    false,
		},
		{
			name: "categories alphabetized",
			ism: model.ISM{
				Classification:   model.ClassificationCUI,
				CategoryMarkings: []string{"CRIT", "EMGT", "PHYS"},
			},
			valid: true,
		},
		{
			name: "specified category with SP- prefix",
			ism: model.ISM{
				Classification:   model.ClassificationCUI,
				CategoryMarkings: []string{"SP-CEII"},
			},
			valid: true,
		},
		{
			name: "single category no order issue",
			ism: model.ISM{
				Classification:   model.ClassificationCUI,
				CategoryMarkings: []string{"WATER"},
			},
			valid: true,
		},
	}

	rule := &CUIRule{}
	r := reg()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !rule.Applies(&tt.ism) {
				t.Fatal("rule should apply to CUI")
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

func TestCUIRule_DoesNotApplyToOther(t *testing.T) {
	rule := &CUIRule{}
	if rule.Applies(&model.ISM{Classification: model.ClassificationU}) {
		t.Error("CUI rule should not apply to U")
	}
	if rule.Applies(&model.ISM{Classification: model.ClassificationS}) {
		t.Error("CUI rule should not apply to S")
	}
}

// --- ClassifiedRule ---

func TestClassifiedRule(t *testing.T) {
	tests := []struct {
		name     string
		ism      model.ISM
		wantCode string
		valid    bool // all findings are warnings, so valid is always true
	}{
		{
			name: "C with full original authority",
			ism: model.ISM{
				Classification:       model.ClassificationC,
				OwnerProducer:        []string{"USA"},
				ClassifiedBy:         "John Doe",
				ClassificationReason: "National Security",
			},
			valid: true,
		},
		{
			name: "C with derivative authority",
			ism: model.ISM{
				Classification:           model.ClassificationC,
				OwnerProducer:            []string{"USA"},
				DerivativelyClassifiedBy: "Jane Smith",
				DerivedFrom:              "SCG v3",
			},
			valid: true,
		},
		{
			name: "C missing all authority",
			ism: model.ISM{
				Classification: model.ClassificationC,
				OwnerProducer:  []string{"USA"},
			},
			wantCode: "classified.missing_authority",
			valid:    true, // warning only
		},
		{
			name: "C original without reason",
			ism: model.ISM{
				Classification: model.ClassificationC,
				OwnerProducer:  []string{"USA"},
				ClassifiedBy:   "John Doe",
			},
			wantCode: "classified.missing_reason",
			valid:    true, // warning only
		},
		{
			name: "C derivative without derivedFrom",
			ism: model.ISM{
				Classification:           model.ClassificationC,
				OwnerProducer:            []string{"USA"},
				DerivativelyClassifiedBy: "Jane Smith",
			},
			wantCode: "classified.missing_derived_from",
			valid:    true, // warning only
		},
	}

	rule := &ClassifiedRule{}
	r := reg()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !rule.Applies(&tt.ism) {
				t.Fatal("rule should apply to classified ISM")
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

func TestClassifiedRule_DoesNotApplyToU(t *testing.T) {
	rule := &ClassifiedRule{}
	if rule.Applies(&model.ISM{Classification: model.ClassificationU}) {
		t.Error("classified rule should not apply to U")
	}
	if rule.Applies(&model.ISM{Classification: model.ClassificationCUI}) {
		t.Error("classified rule should not apply to CUI")
	}
}

// --- DisseminationRule ---

func TestDisseminationRule(t *testing.T) {
	tests := []struct {
		name     string
		ism      model.ISM
		wantCode string
		valid    bool
	}{
		{
			name: "REL with releasableTo",
			ism: model.ISM{
				Classification:        model.ClassificationC,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"REL"},
				ReleasableTo:          []string{"USA", "GBR"},
			},
			valid: true,
		},
		{
			name: "REL without releasableTo",
			ism: model.ISM{
				Classification:        model.ClassificationC,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"REL"},
			},
			wantCode: "dissemination.missing_releasable_to",
			valid:    false,
		},
		{
			name: "DISPLAY ONLY with displayOnlyTo",
			ism: model.ISM{
				Classification:        model.ClassificationC,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"DISPLAY ONLY"},
				DisplayOnlyTo:         []string{"USA", "GBR"},
			},
			valid: true,
		},
		{
			name: "DISPLAY ONLY without displayOnlyTo",
			ism: model.ISM{
				Classification:        model.ClassificationC,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"DISPLAY ONLY"},
			},
			wantCode: "dissemination.missing_display_only_to",
			valid:    false,
		},
		{
			name: "NOFORN and REL conflict",
			ism: model.ISM{
				Classification:        model.ClassificationC,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"NOFORN", "REL"},
				ReleasableTo:          []string{"USA"},
			},
			wantCode: "dissemination.exclusive_conflict",
			valid:    false,
		},
		{
			name: "NOFORN and RELIDO conflict",
			ism: model.ISM{
				Classification:        model.ClassificationC,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"NOFORN", "RELIDO"},
			},
			wantCode: "dissemination.exclusive_conflict",
			valid:    false,
		},
		{
			name: "NOFORN at U — insufficient classification",
			ism: model.ISM{
				Classification:        model.ClassificationU,
				DisseminationControls: []string{"NOFORN"},
			},
			wantCode: "dissemination.insufficient_classification",
			valid:    false,
		},
		{
			name: "OC at C — meets minimum",
			ism: model.ISM{
				Classification:        model.ClassificationC,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"OC"},
			},
			valid: true,
		},
		{
			name: "OC at CUI — insufficient",
			ism: model.ISM{
				Classification:        model.ClassificationCUI,
				DisseminationControls: []string{"OC"},
			},
			wantCode: "dissemination.insufficient_classification",
			valid:    false,
		},
		{
			name: "unknown control code",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"BOGUS"},
			},
			wantCode: "dissemination.invalid_control",
			valid:    false,
		},
		{
			name: "invalid releasableTo country code",
			ism: model.ISM{
				Classification:        model.ClassificationC,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"REL"},
				ReleasableTo:          []string{"ZZZ"},
			},
			wantCode: "dissemination.invalid_country_code",
			valid:    false,
		},
		{
			name: "invalid displayOnlyTo country code",
			ism: model.ISM{
				Classification:        model.ClassificationC,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"DISPLAY ONLY"},
				DisplayOnlyTo:         []string{"ZZZ"},
			},
			wantCode: "dissemination.invalid_country_code",
			valid:    false,
		},
		{
			name: "PROPIN at U — no min classification",
			ism: model.ISM{
				Classification:        model.ClassificationU,
				DisseminationControls: []string{"PROPIN"},
			},
			valid: true,
		},
		{
			name: "EYES requires releasableTo",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"EYES"},
			},
			wantCode: "dissemination.missing_releasable_to",
			valid:    false,
		},
		{
			name: "EYES with releasableTo",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"EYES"},
				ReleasableTo:          []string{"USA", "GBR"},
			},
			valid: true,
		},
	}

	rule := &DisseminationRule{}
	r := reg()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !rule.Applies(&tt.ism) {
				t.Fatal("rule should apply when disseminationControls present")
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

func TestDisseminationRule_DoesNotApplyWhenEmpty(t *testing.T) {
	rule := &DisseminationRule{}
	if rule.Applies(&model.ISM{Classification: model.ClassificationU}) {
		t.Error("dissemination rule should not apply when no controls")
	}
}

// --- DistributionRule ---

func TestDistributionRule(t *testing.T) {
	tests := []struct {
		name     string
		ism      model.ISM
		wantCode string
		valid    bool
	}{
		{
			name: "Statement A at U",
			ism: model.ISM{
				Classification:        model.ClassificationU,
				DistributionStatement: "A",
			},
			valid: true,
		},
		{
			name: "Statement A at CUI",
			ism: model.ISM{
				Classification:        model.ClassificationCUI,
				DistributionStatement: "A",
			},
			valid: true,
		},
		{
			name: "Statement A at C — not allowed",
			ism: model.ISM{
				Classification:        model.ClassificationC,
				OwnerProducer:         []string{"USA"},
				DistributionStatement: "A",
			},
			wantCode: "distribution.classification_mismatch",
			valid:    false,
		},
		{
			name: "Statement A at S — not allowed",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA"},
				DistributionStatement: "A",
			},
			wantCode: "distribution.classification_mismatch",
			valid:    false,
		},
		{
			name: "Statement B at S — allowed",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA"},
				DistributionStatement: "B",
			},
			valid: true,
		},
		{
			name: "Statement F at C — allowed",
			ism: model.ISM{
				Classification:        model.ClassificationC,
				OwnerProducer:         []string{"USA"},
				DistributionStatement: "F",
			},
			valid: true,
		},
		{
			name: "invalid code",
			ism: model.ISM{
				Classification:        model.ClassificationU,
				DistributionStatement: "G",
			},
			wantCode: "distribution.invalid_code",
			valid:    false,
		},
	}

	rule := &DistributionRule{}
	r := reg()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !rule.Applies(&tt.ism) {
				t.Fatal("rule should apply when distributionStatement present")
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

func TestDistributionRule_DoesNotApplyWhenEmpty(t *testing.T) {
	rule := &DistributionRule{}
	if rule.Applies(&model.ISM{Classification: model.ClassificationU}) {
		t.Error("distribution rule should not apply when no statement")
	}
}

// --- ThirdPartyRule ---

func TestThirdPartyRule(t *testing.T) {
	tests := []struct {
		name     string
		ism      model.ISM
		wantCode string
		valid    bool
	}{
		{
			name: "complete contract",
			ism: model.ISM{
				Classification:                  model.ClassificationU,
				ThirdPartyDistributionStatement: "Custom statement",
				ThirdPartyDistributionContract: &model.ThirdPartyDistributionContract{
					ContractNumber:    "W15QKN-12-D-0001",
					ContractorName:    "Acme Corp",
					ContractorAddress: "123 Main St",
					ExpirationDate:    "2030-12-31",
				},
			},
			valid: true,
		},
		{
			name: "missing contract entirely",
			ism: model.ISM{
				Classification:                  model.ClassificationU,
				ThirdPartyDistributionStatement: "Custom statement",
			},
			wantCode: "third_party.contract_required",
			valid:    false,
		},
		{
			name: "contract missing contractNumber",
			ism: model.ISM{
				Classification:                  model.ClassificationU,
				ThirdPartyDistributionStatement: "Custom statement",
				ThirdPartyDistributionContract: &model.ThirdPartyDistributionContract{
					ContractorName:    "Acme Corp",
					ContractorAddress: "123 Main St",
					ExpirationDate:    "2030-12-31",
				},
			},
			wantCode: "third_party.missing_contract_number",
			valid:    false,
		},
		{
			name: "contract missing contractorName",
			ism: model.ISM{
				Classification:                  model.ClassificationU,
				ThirdPartyDistributionStatement: "Custom statement",
				ThirdPartyDistributionContract: &model.ThirdPartyDistributionContract{
					ContractNumber:    "W15QKN-12-D-0001",
					ContractorAddress: "123 Main St",
					ExpirationDate:    "2030-12-31",
				},
			},
			wantCode: "third_party.missing_contractor_name",
			valid:    false,
		},
		{
			name: "contract missing contractorAddress",
			ism: model.ISM{
				Classification:                  model.ClassificationU,
				ThirdPartyDistributionStatement: "Custom statement",
				ThirdPartyDistributionContract: &model.ThirdPartyDistributionContract{
					ContractNumber: "W15QKN-12-D-0001",
					ContractorName: "Acme Corp",
					ExpirationDate: "2030-12-31",
				},
			},
			wantCode: "third_party.missing_contractor_address",
			valid:    false,
		},
		{
			name: "contract missing expirationDate",
			ism: model.ISM{
				Classification:                  model.ClassificationU,
				ThirdPartyDistributionStatement: "Custom statement",
				ThirdPartyDistributionContract: &model.ThirdPartyDistributionContract{
					ContractNumber:    "W15QKN-12-D-0001",
					ContractorName:    "Acme Corp",
					ContractorAddress: "123 Main St",
				},
			},
			wantCode: "third_party.missing_expiration_date",
			valid:    false,
		},
	}

	rule := &ThirdPartyRule{}
	r := reg()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !rule.Applies(&tt.ism) {
				t.Fatal("rule should apply when 3rdPartyDistributionStatement present")
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

func TestThirdPartyRule_DoesNotApplyWhenEmpty(t *testing.T) {
	rule := &ThirdPartyRule{}
	if rule.Applies(&model.ISM{Classification: model.ClassificationU}) {
		t.Error("third_party rule should not apply when no statement")
	}
}

// --- DeclassRule ---

func TestDeclassRule(t *testing.T) {
	tests := []struct {
		name     string
		ism      model.ISM
		wantCode string
		valid    bool
	}{
		{
			name: "C with declassDate YYYYMMDD",
			ism: model.ISM{
				Classification: model.ClassificationC,
				OwnerProducer:  []string{"USA"},
				DeclassDate:    "20300101",
			},
			valid: true,
		},
		{
			name: "S with declassEvent",
			ism: model.ISM{
				Classification: model.ClassificationS,
				OwnerProducer:  []string{"USA"},
				DeclassEvent:   "Peace treaty signed",
			},
			valid: true,
		},
		{
			name: "C with date and event — exclusive",
			ism: model.ISM{
				Classification: model.ClassificationC,
				OwnerProducer:  []string{"USA"},
				DeclassDate:    "20300101",
				DeclassEvent:   "Peace treaty",
			},
			wantCode: "declass.date_event_exclusive",
			valid:    false,
		},
		{
			name: "U with declassDate — not applicable",
			ism: model.ISM{
				Classification: model.ClassificationU,
				DeclassDate:    "20300101",
			},
			wantCode: "declass.not_applicable",
			valid:    false,
		},
		{
			name: "C with YYYY-MM-DD format — invalid",
			ism: model.ISM{
				Classification: model.ClassificationC,
				OwnerProducer:  []string{"USA"},
				DeclassDate:    "2030-01-01",
			},
			wantCode: "declass.invalid_date_format",
			valid:    false,
		},
		{
			name: "C with invalid calendar date",
			ism: model.ISM{
				Classification: model.ClassificationC,
				OwnerProducer:  []string{"USA"},
				DeclassDate:    "20301301",
			},
			wantCode: "declass.invalid_date",
			valid:    false,
		},
		{
			name: "CUI with declassEvent — not applicable",
			ism: model.ISM{
				Classification: model.ClassificationCUI,
				DeclassEvent:   "Some event",
			},
			wantCode: "declass.not_applicable",
			valid:    false,
		},
		{
			name: "S with valid 25X exception",
			ism: model.ISM{
				Classification:   model.ClassificationS,
				OwnerProducer:    []string{"USA"},
				DeclassException: "25X1",
			},
			valid: true,
		},
		{
			name: "S with invalid exception code",
			ism: model.ISM{
				Classification:   model.ClassificationS,
				OwnerProducer:    []string{"USA"},
				DeclassException: "99X",
			},
			wantCode: "declass.invalid_exception",
			valid:    false,
		},
		{
			name: "C with 50X exception",
			ism: model.ISM{
				Classification:   model.ClassificationC,
				OwnerProducer:    []string{"USA"},
				DeclassException: "50X1-HUM",
			},
			valid: true,
		},
	}

	rule := &DeclassRule{}
	r := reg()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !rule.Applies(&tt.ism) {
				t.Fatal("rule should apply when declass fields present")
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

func TestDeclassRule_DoesNotApplyWhenEmpty(t *testing.T) {
	rule := &DeclassRule{}
	if rule.Applies(&model.ISM{Classification: model.ClassificationS}) {
		t.Error("declass rule should not apply when no declass fields")
	}
}

// --- FGIRule ---

func TestFGIRule(t *testing.T) {
	tests := []struct {
		name     string
		ism      model.ISM
		wantCode string
		valid    bool
	}{
		{
			name: "valid fgiSourceOpen",
			ism: model.ISM{
				Classification: model.ClassificationS,
				OwnerProducer:  []string{"USA"},
				FGISourceOpen:  []string{"GBR", "FRA"},
			},
			valid: true,
		},
		{
			name: "valid fgiSourceProtected",
			ism: model.ISM{
				Classification:     model.ClassificationS,
				OwnerProducer:      []string{"USA"},
				FGISourceProtected: []string{"DEU"},
			},
			valid: true,
		},
		{
			name: "invalid fgiSourceOpen code",
			ism: model.ISM{
				Classification: model.ClassificationS,
				OwnerProducer:  []string{"USA"},
				FGISourceOpen:  []string{"ZZZ"},
			},
			wantCode: "fgi.invalid_country_code",
			valid:    false,
		},
		{
			name: "invalid fgiSourceProtected code",
			ism: model.ISM{
				Classification:     model.ClassificationS,
				OwnerProducer:      []string{"USA"},
				FGISourceProtected: []string{"ZZZ"},
			},
			wantCode: "fgi.invalid_country_code",
			valid:    false,
		},
		{
			name: "coalition code in fgiSourceOpen",
			ism: model.ISM{
				Classification: model.ClassificationS,
				OwnerProducer:  []string{"USA"},
				FGISourceOpen:  []string{"NATO"},
			},
			valid: true,
		},
	}

	rule := &FGIRule{}
	r := reg()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !rule.Applies(&tt.ism) {
				t.Fatal("rule should apply when FGI fields present")
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

func TestFGIRule_DoesNotApplyWhenEmpty(t *testing.T) {
	rule := &FGIRule{}
	if rule.Applies(&model.ISM{Classification: model.ClassificationS}) {
		t.Error("fgi rule should not apply when no FGI fields")
	}
}

// --- AtomicEnergyRule ---

func TestAtomicEnergyRule(t *testing.T) {
	tests := []struct {
		name     string
		ism      model.ISM
		wantCode string
		valid    bool
	}{
		{
			name: "RD at TS — valid",
			ism: model.ISM{
				Classification:       model.ClassificationTS,
				OwnerProducer:        []string{"USA"},
				AtomicEnergyMarkings: []string{"RD"},
			},
			valid: true,
		},
		{
			name: "RD at S — valid",
			ism: model.ISM{
				Classification:       model.ClassificationS,
				OwnerProducer:        []string{"USA"},
				AtomicEnergyMarkings: []string{"RD"},
			},
			valid: true,
		},
		{
			name: "RD at C — valid",
			ism: model.ISM{
				Classification:       model.ClassificationC,
				OwnerProducer:        []string{"USA"},
				AtomicEnergyMarkings: []string{"RD"},
			},
			valid: true,
		},
		{
			name: "RD at U — insufficient",
			ism: model.ISM{
				Classification:       model.ClassificationU,
				AtomicEnergyMarkings: []string{"RD"},
			},
			wantCode: "atomic_energy.insufficient_classification",
			valid:    false,
		},
		{
			name: "RD at CUI — insufficient",
			ism: model.ISM{
				Classification:       model.ClassificationCUI,
				AtomicEnergyMarkings: []string{"RD"},
			},
			wantCode: "atomic_energy.insufficient_classification",
			valid:    false,
		},
		{
			name: "FRD at C — valid",
			ism: model.ISM{
				Classification:       model.ClassificationC,
				OwnerProducer:        []string{"USA"},
				AtomicEnergyMarkings: []string{"FRD"},
			},
			valid: true,
		},
		{
			name: "FRD at U — insufficient",
			ism: model.ISM{
				Classification:       model.ClassificationU,
				AtomicEnergyMarkings: []string{"FRD"},
			},
			wantCode: "atomic_energy.insufficient_classification",
			valid:    false,
		},
		{
			name: "CNWDI at TS — valid",
			ism: model.ISM{
				Classification:       model.ClassificationTS,
				OwnerProducer:        []string{"USA"},
				AtomicEnergyMarkings: []string{"RD-CNWDI"},
			},
			valid: true,
		},
		{
			name: "CNWDI at S — valid",
			ism: model.ISM{
				Classification:       model.ClassificationS,
				OwnerProducer:        []string{"USA"},
				AtomicEnergyMarkings: []string{"RD-CNWDI"},
			},
			valid: true,
		},
		{
			name: "CNWDI at C — insufficient",
			ism: model.ISM{
				Classification:       model.ClassificationC,
				OwnerProducer:        []string{"USA"},
				AtomicEnergyMarkings: []string{"RD-CNWDI"},
			},
			wantCode: "atomic_energy.insufficient_classification",
			valid:    false,
		},
		{
			name: "CNWDI at U — insufficient",
			ism: model.ISM{
				Classification:       model.ClassificationU,
				AtomicEnergyMarkings: []string{"RD-CNWDI"},
			},
			wantCode: "atomic_energy.insufficient_classification",
			valid:    false,
		},
		{
			name: "SIGMA (RD-SG-14) at C — valid",
			ism: model.ISM{
				Classification:       model.ClassificationC,
				OwnerProducer:        []string{"USA"},
				AtomicEnergyMarkings: []string{"RD-SG-14"},
			},
			valid: true,
		},
		{
			name: "SIGMA (RD-SG-14) at U — insufficient",
			ism: model.ISM{
				Classification:       model.ClassificationU,
				AtomicEnergyMarkings: []string{"RD-SG-14"},
			},
			wantCode: "atomic_energy.insufficient_classification",
			valid:    false,
		},
		{
			name: "FRD SIGMA (FRD-SG-18) at S — valid",
			ism: model.ISM{
				Classification:       model.ClassificationS,
				OwnerProducer:        []string{"USA"},
				AtomicEnergyMarkings: []string{"FRD-SG-18"},
			},
			valid: true,
		},
		{
			name: "DCNI at U — no classification gate",
			ism: model.ISM{
				Classification:       model.ClassificationU,
				AtomicEnergyMarkings: []string{"DCNI"},
			},
			valid: true,
		},
		{
			name: "UCNI at U — no classification gate",
			ism: model.ISM{
				Classification:       model.ClassificationU,
				AtomicEnergyMarkings: []string{"UCNI"},
			},
			valid: true,
		},
		{
			name: "TFNI at U — no classification gate",
			ism: model.ISM{
				Classification:       model.ClassificationU,
				AtomicEnergyMarkings: []string{"TFNI"},
			},
			valid: true,
		},
	}

	rule := &AtomicEnergyRule{}
	r := reg()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !rule.Applies(&tt.ism) {
				t.Fatal("rule should apply when atomicEnergyMarkings present")
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

func TestAtomicEnergyRule_DoesNotApplyWhenEmpty(t *testing.T) {
	rule := &AtomicEnergyRule{}
	if rule.Applies(&model.ISM{Classification: model.ClassificationU}) {
		t.Error("atomic_energy rule should not apply when no markings")
	}
}

// --- DisseminationRule FOUO max classification ---

func TestDisseminationRule_FOUO_MaxClassification(t *testing.T) {
	tests := []struct {
		name     string
		ism      model.ISM
		wantCode string
		valid    bool
	}{
		{
			name: "FOUO at U — valid",
			ism: model.ISM{
				Classification:        model.ClassificationU,
				DisseminationControls: []string{"FOUO"},
			},
			valid: true,
		},
		{
			name: "FOUO at S — exceeds max",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"FOUO"},
			},
			wantCode: "dissemination.exceeds_max_classification",
			valid:    false,
		},
		{
			name: "FOUO at TS — exceeds max",
			ism: model.ISM{
				Classification:        model.ClassificationTS,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"FOUO"},
			},
			wantCode: "dissemination.exceeds_max_classification",
			valid:    false,
		},
	}

	rule := &DisseminationRule{}
	r := reg()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !rule.Applies(&tt.ism) {
				t.Fatal("rule should apply when disseminationControls present")
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

// --- Result helpers ---

func TestValidationResult_AddErrorAndWarning(t *testing.T) {
	r := &ValidationResult{Valid: true}
	r.AddWarning("f", "w1", "warning")
	if !r.Valid {
		t.Error("warning should not invalidate result")
	}
	if len(r.Errors) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(r.Errors))
	}
	if r.Errors[0].Severity != SeverityWarning {
		t.Error("expected warning severity")
	}

	r.AddError("f", "e1", "error")
	if r.Valid {
		t.Error("error should invalidate result")
	}
	if len(r.Errors) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(r.Errors))
	}
}

func TestValidationResult_Merge(t *testing.T) {
	a := &ValidationResult{Valid: true}
	a.AddWarning("f", "w1", "warning")

	b := &ValidationResult{Valid: true}
	b.AddError("f", "e1", "error")

	a.Merge(b)
	if a.Valid {
		t.Error("merge with error should invalidate")
	}
	if len(a.Errors) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(a.Errors))
	}
}

func TestValidationResult_HasCode(t *testing.T) {
	r := &ValidationResult{Valid: true}
	r.AddError("f", "test.code", "msg")
	if !r.HasCode("test.code") {
		t.Error("should find code")
	}
	if r.HasCode("other.code") {
		t.Error("should not find absent code")
	}
}

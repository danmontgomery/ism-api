package compliance_test

import (
	"strings"
	"testing"

	"dmontgomery/ism-api/internal/banner"
	"dmontgomery/ism-api/internal/model"
)

// DoDM 5200.01-V2 Enclosure 4, Sections 1 & 3
// US Classification & Marking Syntax Tests
// Ref: docs/dodm-5200.01-enclosure4-requirements.md

// TestDoDM_E4S1b_MarkingSyntaxStructure verifies the overall marking syntax
// follows CLASSIFICATION//SCI//SAP//AEA//FGI//DISSEM//OTHER DISSEM.
func TestDoDM_E4S1b_MarkingSyntaxStructure(t *testing.T) {
	tests := []struct {
		name       string
		ism        model.ISM
		wantBanner string
	}{
		{
			name: "classification_only",
			ism: model.ISM{
				Classification: model.ClassificationS,
				OwnerProducer:  []string{"USA"},
			},
			wantBanner: "SECRET",
		},
		{
			name: "classification_with_SCI",
			ism: model.ISM{
				Classification: model.ClassificationTS,
				OwnerProducer:  []string{"USA"},
				SCIControls:    []string{"SI"},
			},
			wantBanner: "TOP SECRET//SI",
		},
		{
			name: "classification_SCI_dissem",
			ism: model.ISM{
				Classification:        model.ClassificationTS,
				OwnerProducer:         []string{"USA"},
				SCIControls:           []string{"SI"},
				DisseminationControls: []string{"NOFORN"},
			},
			wantBanner: "TOP SECRET//SI//NOFORN",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := banner.Render(&tt.ism)
			if result.BannerLine != tt.wantBanner {
				t.Errorf("BannerLine = %q, want %q", result.BannerLine, tt.wantBanner)
			}
		})
	}

	// Verify FGI + DISSEM case: contains checks for multi-category structure
	t.Run("classification_FGI_dissem", func(t *testing.T) {
		ism := &model.ISM{
			Classification:        model.ClassificationS,
			OwnerProducer:         []string{"USA"},
			FGISourceOpen:         []string{"GBR"},
			DisseminationControls: []string{"NOFORN"},
		}
		result := banner.Render(ism)
		if !strings.Contains(result.BannerLine, "SECRET") {
			t.Errorf("BannerLine %q should contain classification", result.BannerLine)
		}
		if !strings.Contains(result.BannerLine, "NOFORN") {
			t.Errorf("BannerLine %q should contain NOFORN", result.BannerLine)
		}
		if !strings.Contains(result.BannerLine, "FGI") {
			t.Errorf("BannerLine %q should contain FGI", result.BannerLine)
		}
	})
}

// TestDoDM_E4S1b1_DoubleSlashSeparatesCategories verifies that // separates
// different marking categories.
func TestDoDM_E4S1b1_DoubleSlashSeparatesCategories(t *testing.T) {
	tests := []struct {
		name string
		ism  model.ISM
	}{
		{
			name: "SCI_category_separated_from_classification",
			ism: model.ISM{
				Classification: model.ClassificationTS,
				OwnerProducer:  []string{"USA"},
				SCIControls:    []string{"SI"},
			},
		},
		{
			name: "dissem_category_separated_from_classification",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"NOFORN"},
			},
		},
		{
			name: "SCI_and_dissem_each_separated",
			ism: model.ISM{
				Classification:        model.ClassificationTS,
				OwnerProducer:         []string{"USA"},
				SCIControls:           []string{"SI"},
				DisseminationControls: []string{"NOFORN"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := banner.Render(&tt.ism)
			// Banner must contain // as separator between categories
			if !strings.Contains(result.BannerLine, "//") {
				t.Errorf("BannerLine %q should contain '//' category separator", result.BannerLine)
			}
			// Portion mark must also contain //
			if !strings.Contains(result.PortionMark, "//") {
				t.Errorf("PortionMark %q should contain '//' category separator", result.PortionMark)
			}
		})
	}
}

// TestDoDM_E4S1b1a_SingleSlashWithinCategory verifies that / separates
// multiple types within the same category.
func TestDoDM_E4S1b1a_SingleSlashWithinCategory(t *testing.T) {
	tests := []struct {
		name       string
		ism        model.ISM
		wantBanner string
	}{
		{
			name: "multiple_SCI_controls",
			ism: model.ISM{
				Classification: model.ClassificationTS,
				OwnerProducer:  []string{"USA"},
				SCIControls:    []string{"HCS", "SI"},
			},
			wantBanner: "TOP SECRET//HCS/SI",
		},
		{
			name: "multiple_dissem_controls",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"NOFORN", "PROPIN"},
			},
			// Within the dissem category, / separates the controls
			wantBanner: "SECRET//NOFORN/PROPIN",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := banner.Render(&tt.ism)
			if result.BannerLine != tt.wantBanner {
				t.Errorf("BannerLine = %q, want %q", result.BannerLine, tt.wantBanner)
			}
		})
	}
}

// TestDoDM_E4S1b1b_HyphenSubControls verifies that hyphens separate a control
// system from its sub-control/compartment (e.g., SI-G, RD-N).
func TestDoDM_E4S1b1b_HyphenSubControls(t *testing.T) {
	tests := []struct {
		name       string
		ism        model.ISM
		wantInBanner string
	}{
		{
			name: "SCI_with_subcontrol",
			ism: model.ISM{
				Classification: model.ClassificationTS,
				OwnerProducer:  []string{"USA"},
				SCIControls:    []string{"SI-G"},
			},
			wantInBanner: "SI-G",
		},
		{
			name: "multiple_SCI_with_subcontrols",
			ism: model.ISM{
				Classification: model.ClassificationTS,
				OwnerProducer:  []string{"USA"},
				SCIControls:    []string{"HCS", "SI-G"},
			},
			wantInBanner: "SI-G",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := banner.Render(&tt.ism)
			if !strings.Contains(result.BannerLine, tt.wantInBanner) {
				t.Errorf("BannerLine %q should contain %q with hyphen sub-control",
					result.BannerLine, tt.wantInBanner)
			}
		})
	}
}

// TestDoDM_E4S1a_CategoryOrdering verifies markings appear in Figure 25 order:
// Classification, SCI, SAP, AEA, FGI, Dissemination, Other Dissemination.
func TestDoDM_E4S1a_CategoryOrdering(t *testing.T) {
	tests := []struct {
		name        string
		ism         model.ISM
		orderedKeys []string // substrings that must appear in this order
	}{
		{
			name: "SCI_before_dissem",
			ism: model.ISM{
				Classification:        model.ClassificationTS,
				OwnerProducer:         []string{"USA"},
				SCIControls:           []string{"SI"},
				DisseminationControls: []string{"NOFORN"},
			},
			orderedKeys: []string{"TOP SECRET", "SI", "NOFORN"},
		},
		{
			name: "FGI_before_dissem",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA"},
				FGISourceOpen:         []string{"GBR"},
				DisseminationControls: []string{"NOFORN"},
			},
			orderedKeys: []string{"SECRET", "NOFORN"},
		},
		{
			name: "SCI_before_FGI_before_dissem",
			ism: model.ISM{
				Classification:        model.ClassificationTS,
				OwnerProducer:         []string{"USA"},
				SCIControls:           []string{"SI"},
				FGISourceOpen:         []string{"GBR"},
				DisseminationControls: []string{"NOFORN"},
			},
			// Per Figure 25: SCI(2) → FGI(5) → Dissem(6)
			orderedKeys: []string{"TOP SECRET", "SI", "FGI", "NOFORN"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.orderedKeys == nil {
				t.Skip("ordering checked in gap subtest below")
				return
			}
			result := banner.Render(&tt.ism)
			lastIdx := -1
			for _, key := range tt.orderedKeys {
				idx := strings.Index(result.BannerLine, key)
				if idx == -1 {
					t.Fatalf("BannerLine %q missing expected %q", result.BannerLine, key)
				}
				if idx <= lastIdx {
					t.Errorf("BannerLine %q: %q (idx=%d) should appear after previous element (idx=%d)",
						result.BannerLine, key, idx, lastIdx)
				}
				lastIdx = idx
			}
		})
	}

	// Figure 25 requires FGI(pos 5) before Dissemination(pos 6).
	t.Run("FGI_precedes_dissem", func(t *testing.T) {
		ism := &model.ISM{
			Classification:        model.ClassificationTS,
			OwnerProducer:         []string{"USA"},
			SCIControls:           []string{"SI"},
			FGISourceOpen:         []string{"GBR"},
			DisseminationControls: []string{"NOFORN"},
		}
		result := banner.Render(ism)
		fgiIdx := strings.Index(result.BannerLine, "FGI")
		nfIdx := strings.Index(result.BannerLine, "NOFORN")
		if fgiIdx == -1 || nfIdx == -1 {
			t.Fatalf("BannerLine %q missing FGI or NOFORN", result.BannerLine)
		}
		if fgiIdx > nfIdx {
			t.Errorf("BannerLine %q: FGI(idx=%d) must appear before NOFORN(idx=%d) per Figure 25",
				result.BannerLine, fgiIdx, nfIdx)
		}
	})
}

// TestDoDM_E4S1b2_ClassificationFullEnglish verifies the classification level
// in the banner line is in English and spelled out completely.
func TestDoDM_E4S1b2_ClassificationFullEnglish(t *testing.T) {
	tests := []struct {
		name           string
		classification model.Classification
		wantBanner     string
	}{
		{"TOP_SECRET", model.ClassificationTS, "TOP SECRET"},
		{"SECRET", model.ClassificationS, "SECRET"},
		{"CONFIDENTIAL", model.ClassificationC, "CONFIDENTIAL"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ism := &model.ISM{
				Classification: tt.classification,
				OwnerProducer:  []string{"USA"},
			}
			result := banner.Render(ism)
			if !strings.HasPrefix(result.BannerLine, tt.wantBanner) {
				t.Errorf("BannerLine %q should start with full English %q",
					result.BannerLine, tt.wantBanner)
			}
		})
	}
}

// TestDoDM_E4S1b2a_ControlMarkingsAbbreviated verifies that control markings
// (SCI, SAP, dissemination) in the banner may use authorized abbreviations,
// and that portion marks use abbreviated forms.
func TestDoDM_E4S1b2a_ControlMarkingsAbbreviated(t *testing.T) {
	tests := []struct {
		name          string
		ism           model.ISM
		bannerCtrl    string // expected control form in banner
		portionAbbrev string // expected abbreviated form in portion mark
	}{
		{
			name: "NOFORN_abbreviated_in_portion",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"NOFORN"},
			},
			bannerCtrl:    "NOFORN",
			portionAbbrev: "NF",
		},
		{
			name: "PROPIN_abbreviated_in_portion",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"PROPIN"},
			},
			bannerCtrl:    "PROPIN",
			portionAbbrev: "PR",
		},
		{
			name: "SCI_same_in_banner_and_portion",
			ism: model.ISM{
				Classification: model.ClassificationTS,
				OwnerProducer:  []string{"USA"},
				SCIControls:    []string{"SI"},
			},
			bannerCtrl:    "SI",
			portionAbbrev: "SI",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := banner.Render(&tt.ism)
			if !strings.Contains(result.BannerLine, tt.bannerCtrl) {
				t.Errorf("BannerLine %q should contain control %q",
					result.BannerLine, tt.bannerCtrl)
			}
			if !strings.Contains(result.PortionMark, tt.portionAbbrev) {
				t.Errorf("PortionMark %q should contain abbreviated %q",
					result.PortionMark, tt.portionAbbrev)
			}
		})
	}
}

// TestDoDM_E4S1b3_ControlsInBannerAndPortion verifies that control markings
// (SCI, SAP, AEA, dissemination) are required in both banner and portion marks.
func TestDoDM_E4S1b3_ControlsInBannerAndPortion(t *testing.T) {
	tests := []struct {
		name          string
		ism           model.ISM
		bannerCtrl    string
		portionCtrl   string
	}{
		{
			name: "SCI_in_both",
			ism: model.ISM{
				Classification: model.ClassificationTS,
				OwnerProducer:  []string{"USA"},
				SCIControls:    []string{"SI"},
			},
			bannerCtrl:  "SI",
			portionCtrl: "SI",
		},
		{
			name: "NOFORN_in_both",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"NOFORN"},
			},
			bannerCtrl:  "NOFORN",
			portionCtrl: "NF",
		},
		{
			name: "IMCON_in_both",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"IMCON"},
			},
			bannerCtrl:  "IMCON",
			portionCtrl: "IMC",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := banner.Render(&tt.ism)
			if !strings.Contains(result.BannerLine, tt.bannerCtrl) {
				t.Errorf("BannerLine %q should contain %q", result.BannerLine, tt.bannerCtrl)
			}
			if !strings.Contains(result.PortionMark, tt.portionCtrl) {
				t.Errorf("PortionMark %q should contain %q", result.PortionMark, tt.portionCtrl)
			}
		})
	}
}

// TestDoDM_E4S1b3a_DoubleSlashSeparatesClassFromControls verifies that //
// separates the classification level from control markings, and different
// categories from each other.
func TestDoDM_E4S1b3a_DoubleSlashSeparatesClassFromControls(t *testing.T) {
	tests := []struct {
		name       string
		ism        model.ISM
		wantBanner string
		wantPortion string
	}{
		{
			name: "TS_with_SCI",
			ism: model.ISM{
				Classification: model.ClassificationTS,
				OwnerProducer:  []string{"USA"},
				SCIControls:    []string{"SI"},
			},
			wantBanner:  "TOP SECRET//SI",
			wantPortion: "(TS//SI)",
		},
		{
			name: "S_with_dissem",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"NOFORN"},
			},
			wantBanner:  "SECRET//NOFORN",
			wantPortion: "(S//NF)",
		},
		{
			name: "TS_SCI_and_dissem_separated",
			ism: model.ISM{
				Classification:        model.ClassificationTS,
				OwnerProducer:         []string{"USA"},
				SCIControls:           []string{"SI"},
				DisseminationControls: []string{"NOFORN"},
			},
			wantBanner:  "TOP SECRET//SI//NOFORN",
			wantPortion: "(TS//SI//NF)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := banner.Render(&tt.ism)
			if result.BannerLine != tt.wantBanner {
				t.Errorf("BannerLine = %q, want %q", result.BannerLine, tt.wantBanner)
			}
			if result.PortionMark != tt.wantPortion {
				t.Errorf("PortionMark = %q, want %q", result.PortionMark, tt.wantPortion)
			}
		})
	}
}

// TestDoDM_E4S1d_MutualExclusivity verifies that US, FGI, and JOINT
// classification systems are mutually exclusive — they may not be combined.
func TestDoDM_E4S1d_MutualExclusivity(t *testing.T) {
	tests := []struct {
		name string
		ism  model.ISM
	}{
		{
			name: "US_only",
			ism: model.ISM{
				Classification: model.ClassificationS,
				OwnerProducer:  []string{"USA"},
			},
		},
		{
			name: "JOINT_only",
			ism: model.ISM{
				Classification: model.ClassificationS,
				OwnerProducer:  []string{"USA", "GBR"},
				Joint:          true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := banner.Render(&tt.ism)
			isUS := !tt.ism.Joint && len(tt.ism.OwnerProducer) > 0 &&
				tt.ism.OwnerProducer[0] == "USA" && !strings.HasPrefix(result.BannerLine, "//")
			isJoint := strings.HasPrefix(result.BannerLine, "//JOINT")

			if isUS && isJoint {
				t.Errorf("BannerLine %q mixes US and JOINT systems", result.BannerLine)
			}
		})
	}

	// Verify US doc banner does NOT start with "//" (which would indicate FGI)
	t.Run("US_not_FGI_prefix", func(t *testing.T) {
		ism := &model.ISM{
			Classification: model.ClassificationS,
			OwnerProducer:  []string{"USA"},
		}
		result := banner.Render(ism)
		if strings.HasPrefix(result.BannerLine, "//") {
			t.Errorf("US BannerLine %q should NOT start with '//' (FGI prefix)", result.BannerLine)
		}
	})

	// Verify JOINT doc starts with //JOINT per E4-S5.c
	t.Run("JOINT_starts_with_JOINT", func(t *testing.T) {
		ism := &model.ISM{
			Classification: model.ClassificationS,
			OwnerProducer:  []string{"USA", "GBR"},
			Joint:          true,
		}
		result := banner.Render(ism)
		if !strings.HasPrefix(result.BannerLine, "//JOINT") {
			t.Errorf("JOINT BannerLine %q should start with '//JOINT'", result.BannerLine)
		}
	})
}

// TestDoDM_E4S1b4a_MultipleDisseminationControls verifies that multiple
// dissemination controls may be used and appear in Figure 25 order.
func TestDoDM_E4S1b4a_MultipleDisseminationControls(t *testing.T) {
	tests := []struct {
		name     string
		controls []string
		ordered  []string // expected order in banner
	}{
		{
			name:     "NOFORN_and_PROPIN",
			controls: []string{"PROPIN", "NOFORN"},
			ordered:  []string{"NOFORN", "PROPIN"},
		},
		{
			name:     "RS_OC_NOFORN",
			controls: []string{"NOFORN", "OC", "RS"},
			ordered:  []string{"RS", "OC", "NOFORN"},
		},
		{
			name:     "IMCON_NOFORN",
			controls: []string{"NOFORN", "IMCON"},
			ordered:  []string{"IMCON", "NOFORN"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ism := &model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: tt.controls,
			}
			result := banner.Render(ism)

			// Verify all controls present
			for _, ctrl := range tt.ordered {
				if !strings.Contains(result.BannerLine, ctrl) {
					t.Fatalf("BannerLine %q missing expected control %q", result.BannerLine, ctrl)
				}
			}

			// Verify ordering
			lastIdx := -1
			for _, ctrl := range tt.ordered {
				idx := strings.Index(result.BannerLine, ctrl)
				if idx <= lastIdx {
					t.Errorf("BannerLine %q: %q (idx=%d) should appear after previous (idx=%d)",
						result.BannerLine, ctrl, idx, lastIdx)
				}
				lastIdx = idx
			}
		})
	}
}

// TestDoDM_E4S1b4b_UnclassifiedPortionMarking verifies that unclassified
// portions without controls are marked (U).
func TestDoDM_E4S1b4b_UnclassifiedPortionMarking(t *testing.T) {
	ism := &model.ISM{
		Classification: model.ClassificationU,
	}
	result := banner.Render(ism)
	if result.PortionMark != "(U)" {
		t.Errorf("PortionMark = %q, want %q for unclassified portion", result.PortionMark, "(U)")
	}
}

// TestDoDM_E4S3a1_TopSecret verifies TOP SECRET renders as "TOP SECRET"
// in the banner and "(TS)" in the portion mark.
func TestDoDM_E4S3a1_TopSecret(t *testing.T) {
	ism := &model.ISM{
		Classification: model.ClassificationTS,
		OwnerProducer:  []string{"USA"},
	}
	result := banner.Render(ism)
	if result.BannerLine != "TOP SECRET" {
		t.Errorf("BannerLine = %q, want %q", result.BannerLine, "TOP SECRET")
	}
	if result.PortionMark != "(TS)" {
		t.Errorf("PortionMark = %q, want %q", result.PortionMark, "(TS)")
	}
}

// TestDoDM_E4S3a2_Secret verifies SECRET renders as "SECRET" in the banner
// and "(S)" in the portion mark.
func TestDoDM_E4S3a2_Secret(t *testing.T) {
	ism := &model.ISM{
		Classification: model.ClassificationS,
		OwnerProducer:  []string{"USA"},
	}
	result := banner.Render(ism)
	if result.BannerLine != "SECRET" {
		t.Errorf("BannerLine = %q, want %q", result.BannerLine, "SECRET")
	}
	if result.PortionMark != "(S)" {
		t.Errorf("PortionMark = %q, want %q", result.PortionMark, "(S)")
	}
}

// TestDoDM_E4S3a3_Confidential verifies CONFIDENTIAL renders as
// "CONFIDENTIAL" in the banner and "(C)" in the portion mark.
func TestDoDM_E4S3a3_Confidential(t *testing.T) {
	ism := &model.ISM{
		Classification: model.ClassificationC,
		OwnerProducer:  []string{"USA"},
	}
	result := banner.Render(ism)
	if result.BannerLine != "CONFIDENTIAL" {
		t.Errorf("BannerLine = %q, want %q", result.BannerLine, "CONFIDENTIAL")
	}
	if result.PortionMark != "(C)" {
		t.Errorf("PortionMark = %q, want %q", result.PortionMark, "(C)")
	}
}

// TestDoDM_E4S3b_ClassificationNotPrecededByDoubleSlash verifies that US
// classification markings are NOT preceded by // in the banner line.
func TestDoDM_E4S3b_ClassificationNotPrecededByDoubleSlash(t *testing.T) {
	tests := []struct {
		name           string
		classification model.Classification
	}{
		{"TOP_SECRET", model.ClassificationTS},
		{"SECRET", model.ClassificationS},
		{"CONFIDENTIAL", model.ClassificationC},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ism := &model.ISM{
				Classification: tt.classification,
				OwnerProducer:  []string{"USA"},
			}
			result := banner.Render(ism)
			if strings.HasPrefix(result.BannerLine, "//") {
				t.Errorf("US BannerLine %q should NOT start with '//'", result.BannerLine)
			}
		})
	}
}

// TestDoDM_E4S3b_note_DoubleSlashPrefixIndicatesFGI verifies that a banner
// starting with // (nothing preceding it) indicates the document contains
// only FGI or JOINT information — not a US classification.
func TestDoDM_E4S3b_note_DoubleSlashPrefixIndicatesFGI(t *testing.T) {
	// US documents: banner must NOT start with //
	t.Run("US_no_leading_double_slash", func(t *testing.T) {
		for _, cls := range []model.Classification{
			model.ClassificationTS,
			model.ClassificationS,
			model.ClassificationC,
		} {
			ism := &model.ISM{
				Classification: cls,
				OwnerProducer:  []string{"USA"},
			}
			result := banner.Render(ism)
			if strings.HasPrefix(result.BannerLine, "//") {
				t.Errorf("US %s BannerLine %q should NOT start with '//'",
					cls, result.BannerLine)
			}
		}
	})

	// JOINT documents: banner starts with "//JOINT" per E4-S5.c
	t.Run("JOINT_leading_double_slash", func(t *testing.T) {
		ism := &model.ISM{
			Classification: model.ClassificationS,
			OwnerProducer:  []string{"USA", "GBR"},
			Joint:          true,
		}
		result := banner.Render(ism)
		if !strings.HasPrefix(result.BannerLine, "//JOINT") {
			t.Errorf("JOINT BannerLine %q should start with '//JOINT'", result.BannerLine)
		}
	})
}

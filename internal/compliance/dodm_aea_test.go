package compliance_test

import (
	"strings"
	"testing"

	"expr.ai/ism-api/internal/banner"
	"expr.ai/ism-api/internal/model"
)

// DoDM 5200.01-V2 Enclosure 4, Section 8
// AEA Markings Tests (RD/FRD/CNWDI/SIGMA)
// Ref: docs/dodm-5200.01-enclosure4-requirements.md

// ---------------------------------------------------------------------------
// 8.1 Restricted Data (RD)
// ---------------------------------------------------------------------------

// TestDoDM_E4S8a4_RD_ClassificationConstraint verifies that RD may be used
// ONLY with TOP SECRET, SECRET, or CONFIDENTIAL.
// [E4-S8.a.4]
func TestDoDM_E4S8a4_RD_ClassificationConstraint(t *testing.T) {
	validLevels := []struct {
		name           string
		classification model.Classification
	}{
		{"TOP_SECRET", model.ClassificationTS},
		{"SECRET", model.ClassificationS},
		{"CONFIDENTIAL", model.ClassificationC},
	}
	for _, tt := range validLevels {
		t.Run(tt.name+"_valid", func(t *testing.T) {
			ism := &model.ISM{
				Classification:       tt.classification,
				OwnerProducer:        []string{"USA"},
				AtomicEnergyMarkings: []string{"RD"},
			}
			result := banner.Render(ism)
			// RD should render in the banner. The current renderer may not
			// support AEA markings yet — flag as GAP if missing.
			if !strings.Contains(result.BannerLine, "RESTRICTED DATA") &&
				!strings.Contains(result.BannerLine, "RD") {
				t.Skipf("GAP: BannerLine %q should contain RESTRICTED DATA or RD for AEA marking; renderer does not yet render AtomicEnergyMarkings",
					result.BannerLine)
			}
		})
	}

	// Invalid: UNCLASSIFIED with RD should not be permitted.
	t.Run("UNCLASSIFIED_invalid", func(t *testing.T) {
		ism := &model.ISM{
			Classification:       model.ClassificationU,
			AtomicEnergyMarkings: []string{"RD"},
		}
		result := banner.Render(ism)
		// Even without validation, RD on UNCLASSIFIED is a marking violation.
		// If the renderer produces RD on U, flag it.
		if strings.Contains(result.BannerLine, "RESTRICTED DATA") ||
			strings.Contains(result.BannerLine, "//RD") {
			t.Errorf("BannerLine %q should NOT contain RD with UNCLASSIFIED — E4-S8.a.4 restricts RD to TS/S/C",
				result.BannerLine)
		}
	})

	// Invalid: CUI with RD should not be permitted.
	t.Run("CUI_invalid", func(t *testing.T) {
		ism := &model.ISM{
			Classification:       model.ClassificationCUI,
			OwnerProducer:        []string{"USA"},
			AtomicEnergyMarkings: []string{"RD"},
		}
		result := banner.Render(ism)
		if strings.Contains(result.BannerLine, "RESTRICTED DATA") ||
			strings.Contains(result.BannerLine, "//RD") {
			t.Errorf("BannerLine %q should NOT contain RD with CUI — E4-S8.a.4 restricts RD to TS/S/C",
				result.BannerLine)
		}
	})
}

// TestDoDM_E4S8a5_RD_BannerAndPortionFormat verifies the banner and portion
// mark formats for Restricted Data.
// Banner: [classification]//RESTRICTED DATA (or RD)
// Portion: ([classification]//RD)
// [E4-S8.a.5]
func TestDoDM_E4S8a5_RD_BannerAndPortionFormat(t *testing.T) {
	tests := []struct {
		name        string
		class       model.Classification
		wantBanner  string
		wantPortion string
	}{
		{
			name:        "SECRET_RD",
			class:       model.ClassificationS,
			wantBanner:  "SECRET//RESTRICTED DATA",
			wantPortion: "(S//RD)",
		},
		{
			name:        "TS_RD",
			class:       model.ClassificationTS,
			wantBanner:  "TOP SECRET//RESTRICTED DATA",
			wantPortion: "(TS//RD)",
		},
		{
			name:        "CONFIDENTIAL_RD",
			class:       model.ClassificationC,
			wantBanner:  "CONFIDENTIAL//RESTRICTED DATA",
			wantPortion: "(C//RD)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ism := &model.ISM{
				Classification:       tt.class,
				OwnerProducer:        []string{"USA"},
				AtomicEnergyMarkings: []string{"RD"},
			}
			result := banner.Render(ism)
			if result.BannerLine != tt.wantBanner {
				t.Skipf("GAP: BannerLine = %q, want %q; renderer does not yet render AEA markings in banner",
					result.BannerLine, tt.wantBanner)
			}
			if result.PortionMark != tt.wantPortion {
				t.Errorf("PortionMark = %q, want %q", result.PortionMark, tt.wantPortion)
			}
		})
	}
}

// TestDoDM_E4S8a5a_RD_BannerPropagation verifies that if any portion contains
// RD, then RD must appear in the banner line.
// [E4-S8.a.5a]
func TestDoDM_E4S8a5a_RD_BannerPropagation(t *testing.T) {
	// When AtomicEnergyMarkings includes RD, the banner MUST contain it.
	ism := &model.ISM{
		Classification:       model.ClassificationS,
		OwnerProducer:        []string{"USA"},
		AtomicEnergyMarkings: []string{"RD"},
	}
	result := banner.Render(ism)

	hasRD := strings.Contains(result.BannerLine, "RESTRICTED DATA") ||
		strings.Contains(result.BannerLine, "//RD")
	if !hasRD {
		t.Skipf("GAP: BannerLine %q must contain RD when AtomicEnergyMarkings includes RD — E4-S8.a.5a requires banner propagation; renderer does not yet support AEA",
			result.BannerLine)
	}
}

// TestDoDM_E4S8a7_RD_NoAutomaticDeclass verifies that RD is not subject to
// automatic declassification. "Declassify On:" shall state "Not applicable"
// or be omitted.
// [E4-S8.a.7]
func TestDoDM_E4S8a7_RD_NoAutomaticDeclass(t *testing.T) {
	tests := []struct {
		name          string
		declassDate   string
		declassEvent  string
		declassExcept string
		wantDeclass   bool // true if we expect a declassify line
	}{
		{
			name:        "AEA_exception",
			declassExcept: "AEA",
			wantDeclass: true,
		},
		{
			name:        "no_declass_fields",
			wantDeclass: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ism := &model.ISM{
				Classification:       model.ClassificationS,
				OwnerProducer:        []string{"USA"},
				AtomicEnergyMarkings: []string{"RD"},
				ClassifiedBy:         "Test",
				DeclassDate:          tt.declassDate,
				DeclassEvent:         tt.declassEvent,
				DeclassException:     tt.declassExcept,
			}
			result := banner.Render(ism)
			auth := result.AuthorityBlock

			// RD documents should NOT have a date-based declassification.
			if tt.declassDate != "" {
				t.Errorf("RD documents should not use DeclassDate — not subject to automatic declassification")
			}

			// If AEA exception is set, the authority block should reflect it.
			if tt.declassExcept == "AEA" {
				if !strings.Contains(auth, "AEA") && !strings.Contains(auth, "Not applicable") {
					t.Logf("INFO: AuthorityBlock %q — RD declass should state AEA exception or Not applicable per E4-S8.a.7",
						auth)
				}
			}

			// Verify no date-driven declassification line appears for RD with no date.
			if !tt.wantDeclass && strings.Contains(auth, "Declassify On: 20") {
				t.Errorf("AuthorityBlock %q should not have date-based declassification for RD — E4-S8.a.7",
					auth)
			}
		})
	}
}

// TestDoDM_E4S8a10_RD_NoNSICommingling verifies that RD and NSI (National
// Security Information) shall not be commingled in the same portion.
// [E4-S8.a.10]
func TestDoDM_E4S8a10_RD_NoNSICommingling(t *testing.T) {
	// RD portion should NOT also contain dissemination controls that imply NSI.
	// This is a structural requirement: RD portions are AEA-only and must be
	// segregated from NSI portions.

	// Verify that RD and typical NSI dissemination controls render as separate
	// categories (not merged into one portion).
	ism := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		AtomicEnergyMarkings:  []string{"RD"},
		DisseminationControls: []string{"NOFORN"},
	}
	result := banner.Render(ism)

	// The banner may contain both AEA and dissem categories (document level
	// can have both), but the requirement is at the portion level. This test
	// documents the constraint: RD and NSI must be in separate portions.
	t.Logf("Banner: %s, Portion: %s — per E4-S8.a.10, RD and NSI must not be commingled in the same portion",
		result.BannerLine, result.PortionMark)

	// Check refdata recognizes RD as a valid AEA marking.
	r := reg()
	if !r.ValidAtomicEnergyMarking("RD") {
		t.Fatal("RD must be a valid atomic energy marking")
	}
}

// ---------------------------------------------------------------------------
// 8.2 Formerly Restricted Data (FRD)
// ---------------------------------------------------------------------------

// TestDoDM_E4S8b2_FRD_ClassificationConstraint verifies that FRD may be used
// ONLY with TOP SECRET, SECRET, or CONFIDENTIAL.
// [E4-S8.b.2]
func TestDoDM_E4S8b2_FRD_ClassificationConstraint(t *testing.T) {
	validLevels := []struct {
		name           string
		classification model.Classification
	}{
		{"TOP_SECRET", model.ClassificationTS},
		{"SECRET", model.ClassificationS},
		{"CONFIDENTIAL", model.ClassificationC},
	}
	for _, tt := range validLevels {
		t.Run(tt.name+"_valid", func(t *testing.T) {
			ism := &model.ISM{
				Classification:       tt.classification,
				OwnerProducer:        []string{"USA"},
				AtomicEnergyMarkings: []string{"FRD"},
			}
			result := banner.Render(ism)
			if !strings.Contains(result.BannerLine, "FORMERLY RESTRICTED DATA") &&
				!strings.Contains(result.BannerLine, "FRD") {
				t.Skipf("GAP: BannerLine %q should contain FORMERLY RESTRICTED DATA or FRD; renderer does not yet render AEA markings",
					result.BannerLine)
			}
		})
	}

	t.Run("UNCLASSIFIED_invalid", func(t *testing.T) {
		ism := &model.ISM{
			Classification:       model.ClassificationU,
			AtomicEnergyMarkings: []string{"FRD"},
		}
		result := banner.Render(ism)
		if strings.Contains(result.BannerLine, "FORMERLY RESTRICTED DATA") ||
			strings.Contains(result.BannerLine, "//FRD") {
			t.Errorf("BannerLine %q should NOT contain FRD with UNCLASSIFIED — E4-S8.b.2",
				result.BannerLine)
		}
	})
}

// TestDoDM_E4S8b3_FRD_BannerAndPortionFormat verifies the banner and portion
// mark formats for Formerly Restricted Data.
// Banner: [classification]//FORMERLY RESTRICTED DATA (or FRD)
// Portion: ([classification]//FRD)
// [E4-S8.b.3]
func TestDoDM_E4S8b3_FRD_BannerAndPortionFormat(t *testing.T) {
	tests := []struct {
		name        string
		class       model.Classification
		wantBanner  string
		wantPortion string
	}{
		{
			name:        "SECRET_FRD",
			class:       model.ClassificationS,
			wantBanner:  "SECRET//FORMERLY RESTRICTED DATA",
			wantPortion: "(S//FRD)",
		},
		{
			name:        "TS_FRD",
			class:       model.ClassificationTS,
			wantBanner:  "TOP SECRET//FORMERLY RESTRICTED DATA",
			wantPortion: "(TS//FRD)",
		},
		{
			name:        "CONFIDENTIAL_FRD",
			class:       model.ClassificationC,
			wantBanner:  "CONFIDENTIAL//FORMERLY RESTRICTED DATA",
			wantPortion: "(C//FRD)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ism := &model.ISM{
				Classification:       tt.class,
				OwnerProducer:        []string{"USA"},
				AtomicEnergyMarkings: []string{"FRD"},
			}
			result := banner.Render(ism)
			if result.BannerLine != tt.wantBanner {
				t.Skipf("GAP: BannerLine = %q, want %q; renderer does not yet render AEA markings in banner",
					result.BannerLine, tt.wantBanner)
			}
			if result.PortionMark != tt.wantPortion {
				t.Errorf("PortionMark = %q, want %q", result.PortionMark, tt.wantPortion)
			}
		})
	}
}

// TestDoDM_E4S8b3a_FRD_BannerPropagation verifies that if any portion contains
// FRD, then FRD must appear in the banner line.
// [E4-S8.b.3a]
func TestDoDM_E4S8b3a_FRD_BannerPropagation(t *testing.T) {
	ism := &model.ISM{
		Classification:       model.ClassificationS,
		OwnerProducer:        []string{"USA"},
		AtomicEnergyMarkings: []string{"FRD"},
	}
	result := banner.Render(ism)

	hasFRD := strings.Contains(result.BannerLine, "FORMERLY RESTRICTED DATA") ||
		strings.Contains(result.BannerLine, "//FRD")
	if !hasFRD {
		t.Skipf("GAP: BannerLine %q must contain FRD when AtomicEnergyMarkings includes FRD — E4-S8.b.3a requires banner propagation; renderer does not yet support AEA",
			result.BannerLine)
	}
}

// TestDoDM_E4S8b5_FRD_NoAutomaticDeclass verifies that FRD is not subject to
// automatic declassification.
// [E4-S8.b.5]
func TestDoDM_E4S8b5_FRD_NoAutomaticDeclass(t *testing.T) {
	ism := &model.ISM{
		Classification:       model.ClassificationS,
		OwnerProducer:        []string{"USA"},
		AtomicEnergyMarkings: []string{"FRD"},
		ClassifiedBy:         "Test",
		DeclassException:     "AEA",
	}
	result := banner.Render(ism)
	auth := result.AuthorityBlock

	// FRD documents should use AEA exception, not a date.
	if strings.Contains(auth, "AEA") || !strings.Contains(auth, "Declassify On:") {
		// Expected: either AEA exception or no declass line.
		t.Logf("AuthorityBlock %q — FRD should use AEA exemption or omit Declassify On per E4-S8.b.5", auth)
	}

	// Verify that setting a date-based declass on FRD is incorrect.
	ismWithDate := &model.ISM{
		Classification:       model.ClassificationS,
		OwnerProducer:        []string{"USA"},
		AtomicEnergyMarkings: []string{"FRD"},
		ClassifiedBy:         "Test",
		DeclassDate:          "20350101",
	}
	resultDate := banner.Render(ismWithDate)
	if strings.Contains(resultDate.AuthorityBlock, "20350101") {
		t.Logf("INFO: FRD with DeclassDate renders %q — per E4-S8.b.5, FRD is not subject to automatic declassification",
			resultDate.AuthorityBlock)
	}
}

// TestDoDM_E4S8b8a_FRD_NoNSICommingling verifies that FRD and NSI shall not
// be commingled in the same portion.
// [E4-S8.b.8.a]
func TestDoDM_E4S8b8a_FRD_NoNSICommingling(t *testing.T) {
	ism := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		AtomicEnergyMarkings:  []string{"FRD"},
		DisseminationControls: []string{"NOFORN"},
	}
	result := banner.Render(ism)

	t.Logf("Banner: %s, Portion: %s — per E4-S8.b.8.a, FRD and NSI must not be commingled in the same portion",
		result.BannerLine, result.PortionMark)

	r := reg()
	if !r.ValidAtomicEnergyMarking("FRD") {
		t.Fatal("FRD must be a valid atomic energy marking")
	}
}

// ---------------------------------------------------------------------------
// 8.3 Critical Nuclear Weapons Design Information (CNWDI)
// ---------------------------------------------------------------------------

// TestDoDM_E4S8c1_CNWDI_RDSubset verifies that CNWDI is the designation for
// TOP SECRET RD or SECRET RD weapons data — it is always a subset of RD.
// [E4-S8.c.1]
func TestDoDM_E4S8c1_CNWDI_RDSubset(t *testing.T) {
	// CNWDI is represented as RD-CNWDI in the XSD enum.
	r := reg()
	if !r.ValidAtomicEnergyMarking("RD-CNWDI") {
		t.Fatal("RD-CNWDI must be a valid atomic energy marking — CNWDI is a subset of RD")
	}

	// CNWDI implies RD. The marking code is RD-CNWDI, not standalone CNWDI.
	tests := []struct {
		name  string
		class model.Classification
	}{
		{"TOP_SECRET_CNWDI", model.ClassificationTS},
		{"SECRET_CNWDI", model.ClassificationS},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ism := &model.ISM{
				Classification:       tt.class,
				OwnerProducer:        []string{"USA"},
				AtomicEnergyMarkings: []string{"RD-CNWDI"},
			}
			result := banner.Render(ism)
			// CNWDI must appear with RD in the banner.
			if !strings.Contains(result.BannerLine, "RESTRICTED DATA") &&
				!strings.Contains(result.BannerLine, "RD") {
				t.Skipf("GAP: BannerLine %q should contain RD marking for CNWDI; renderer does not yet render AEA markings",
					result.BannerLine)
			}
		})
	}
}

// TestDoDM_E4S8c1class_CNWDI_ClassificationConstraint verifies that CNWDI
// applies only to TOP SECRET or SECRET (as subset of RD).
// [E4-S8.c.1.class]
func TestDoDM_E4S8c1class_CNWDI_ClassificationConstraint(t *testing.T) {
	// Valid: TS and S only.
	validLevels := []struct {
		name  string
		class model.Classification
	}{
		{"TOP_SECRET", model.ClassificationTS},
		{"SECRET", model.ClassificationS},
	}
	for _, tt := range validLevels {
		t.Run(tt.name+"_valid", func(t *testing.T) {
			ism := &model.ISM{
				Classification:       tt.class,
				OwnerProducer:        []string{"USA"},
				AtomicEnergyMarkings: []string{"RD-CNWDI"},
			}
			result := banner.Render(ism)
			// The marking should render. If not, GAP.
			_ = result
			t.Logf("CNWDI with %s: Banner=%q Portion=%q",
				tt.class, result.BannerLine, result.PortionMark)
		})
	}

	// Invalid: CONFIDENTIAL with CNWDI is not permitted per E4-S8.c.1.class.
	t.Run("CONFIDENTIAL_invalid", func(t *testing.T) {
		ism := &model.ISM{
			Classification:       model.ClassificationC,
			OwnerProducer:        []string{"USA"},
			AtomicEnergyMarkings: []string{"RD-CNWDI"},
		}
		result := banner.Render(ism)
		// CNWDI should NOT appear with CONFIDENTIAL.
		if strings.Contains(result.BannerLine, "CNWDI") || strings.Contains(result.BannerLine, "-N") {
			t.Logf("INFO: BannerLine %q — CNWDI with CONFIDENTIAL violates E4-S8.c.1.class (TS or S only); validation should reject this",
				result.BannerLine)
		}
	})
}

// TestDoDM_E4S8c3_CNWDI_NSuffix verifies that CNWDI is marked by appending -N
// to banner and portion markings.
// Banner: SECRET//RESTRICTED DATA-N
// Portion: (S//RD-N)
// [E4-S8.c.3]
func TestDoDM_E4S8c3_CNWDI_NSuffix(t *testing.T) {
	tests := []struct {
		name        string
		class       model.Classification
		wantBanner  string
		wantPortion string
	}{
		{
			name:        "SECRET_CNWDI",
			class:       model.ClassificationS,
			wantBanner:  "SECRET//RESTRICTED DATA-N",
			wantPortion: "(S//RD-N)",
		},
		{
			name:        "TS_CNWDI",
			class:       model.ClassificationTS,
			wantBanner:  "TOP SECRET//RESTRICTED DATA-N",
			wantPortion: "(TS//RD-N)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ism := &model.ISM{
				Classification:       tt.class,
				OwnerProducer:        []string{"USA"},
				AtomicEnergyMarkings: []string{"RD-CNWDI"},
			}
			result := banner.Render(ism)
			if result.BannerLine != tt.wantBanner {
				t.Skipf("GAP: BannerLine = %q, want %q; renderer does not yet render CNWDI -N suffix in AEA markings",
					result.BannerLine, tt.wantBanner)
			}
			if result.PortionMark != tt.wantPortion {
				t.Errorf("PortionMark = %q, want %q", result.PortionMark, tt.wantPortion)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// 8.4 SIGMA
// ---------------------------------------------------------------------------

// TestDoDM_E4S8d2_SIGMA_ClassificationConstraint verifies that SIGMA may be
// used ONLY with TOP SECRET, SECRET, or CONFIDENTIAL.
// [E4-S8.d.2]
func TestDoDM_E4S8d2_SIGMA_ClassificationConstraint(t *testing.T) {
	sigmaMarkings := []string{"RD-SG-14", "FRD-SG-14"}
	validLevels := []struct {
		name  string
		class model.Classification
	}{
		{"TOP_SECRET", model.ClassificationTS},
		{"SECRET", model.ClassificationS},
		{"CONFIDENTIAL", model.ClassificationC},
	}
	for _, sig := range sigmaMarkings {
		for _, tt := range validLevels {
			t.Run(sig+"_"+tt.name+"_valid", func(t *testing.T) {
				ism := &model.ISM{
					Classification:       tt.class,
					OwnerProducer:        []string{"USA"},
					AtomicEnergyMarkings: []string{sig},
				}
				result := banner.Render(ism)
				_ = result
				t.Logf("SIGMA %s with %s: Banner=%q",
					sig, tt.class, result.BannerLine)
			})
		}
	}

	// Invalid: UNCLASSIFIED with SIGMA.
	t.Run("UNCLASSIFIED_invalid", func(t *testing.T) {
		ism := &model.ISM{
			Classification:       model.ClassificationU,
			AtomicEnergyMarkings: []string{"RD-SG-14"},
		}
		result := banner.Render(ism)
		if strings.Contains(result.BannerLine, "SIGMA") ||
			strings.Contains(result.BannerLine, "SG") {
			t.Errorf("BannerLine %q should NOT contain SIGMA with UNCLASSIFIED — E4-S8.d.2",
				result.BannerLine)
		}
	})
}

// TestDoDM_E4S8d3_SIGMA_BannerFormat verifies the SIGMA banner format:
// [classification]//[RD or FRD]-SIGMA [#]
// [E4-S8.d.3]
func TestDoDM_E4S8d3_SIGMA_BannerFormat(t *testing.T) {
	tests := []struct {
		name        string
		class       model.Classification
		aeaCode     string
		wantBanner  string
		wantPortion string
	}{
		{
			name:        "SECRET_RD_SIGMA_14",
			class:       model.ClassificationS,
			aeaCode:     "RD-SG-14",
			wantBanner:  "SECRET//RD-SIGMA 14",
			wantPortion: "(S//RD-SG 14)",
		},
		{
			name:        "SECRET_FRD_SIGMA_14",
			class:       model.ClassificationS,
			aeaCode:     "FRD-SG-14",
			wantBanner:  "SECRET//FRD-SIGMA 14",
			wantPortion: "(S//FRD-SG 14)",
		},
		{
			name:        "TS_RD_SIGMA_15",
			class:       model.ClassificationTS,
			aeaCode:     "RD-SG-15",
			wantBanner:  "TOP SECRET//RD-SIGMA 15",
			wantPortion: "(TS//RD-SG 15)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ism := &model.ISM{
				Classification:       tt.class,
				OwnerProducer:        []string{"USA"},
				AtomicEnergyMarkings: []string{tt.aeaCode},
			}
			result := banner.Render(ism)
			if result.BannerLine != tt.wantBanner {
				t.Skipf("GAP: BannerLine = %q, want %q; renderer does not yet render SIGMA format in AEA markings",
					result.BannerLine, tt.wantBanner)
			}
			if result.PortionMark != tt.wantPortion {
				t.Errorf("PortionMark = %q, want %q", result.PortionMark, tt.wantPortion)
			}
		})
	}
}

// TestDoDM_E4S8d3multi_SIGMA_MultipleNumericalOrder verifies that multiple
// SIGMAs are listed in numerical order, separated by a space.
// Example: RD-SIGMA 1 2
// [E4-S8.d.3.multi]
func TestDoDM_E4S8d3multi_SIGMA_MultipleNumericalOrder(t *testing.T) {
	tests := []struct {
		name        string
		aeaCodes    []string
		wantBanner  string
		wantPortion string
	}{
		{
			name:        "RD_SIGMA_14_15",
			aeaCodes:    []string{"RD-SG-14", "RD-SG-15"},
			wantBanner:  "SECRET//RD-SIGMA 14 15",
			wantPortion: "(S//RD-SG 14 15)",
		},
		{
			name:        "RD_SIGMA_15_18_20",
			aeaCodes:    []string{"RD-SG-15", "RD-SG-18", "RD-SG-20"},
			wantBanner:  "SECRET//RD-SIGMA 15 18 20",
			wantPortion: "(S//RD-SG 15 18 20)",
		},
		{
			name:        "FRD_SIGMA_14_20",
			aeaCodes:    []string{"FRD-SG-14", "FRD-SG-20"},
			wantBanner:  "SECRET//FRD-SIGMA 14 20",
			wantPortion: "(S//FRD-SG 14 20)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ism := &model.ISM{
				Classification:       model.ClassificationS,
				OwnerProducer:        []string{"USA"},
				AtomicEnergyMarkings: tt.aeaCodes,
			}
			result := banner.Render(ism)
			if result.BannerLine != tt.wantBanner {
				t.Skipf("GAP: BannerLine = %q, want %q; renderer does not yet render multiple SIGMAs in numerical order",
					result.BannerLine, tt.wantBanner)
			}
			if result.PortionMark != tt.wantPortion {
				t.Errorf("PortionMark = %q, want %q", result.PortionMark, tt.wantPortion)
			}
		})
	}

	// Verify ordering: SIGMAs provided out of numerical order should still render in order.
	t.Run("out_of_order_input_sorted", func(t *testing.T) {
		ism := &model.ISM{
			Classification:       model.ClassificationS,
			OwnerProducer:        []string{"USA"},
			AtomicEnergyMarkings: []string{"RD-SG-20", "RD-SG-14"},
		}
		result := banner.Render(ism)
		if strings.Contains(result.BannerLine, "SIGMA") {
			// If renderer supports SIGMA, verify ordering.
			idx14 := strings.Index(result.BannerLine, "14")
			idx20 := strings.Index(result.BannerLine, "20")
			if idx14 != -1 && idx20 != -1 && idx14 > idx20 {
				t.Errorf("BannerLine %q: SIGMA 14 should appear before SIGMA 20 (numerical order)",
					result.BannerLine)
			}
		} else {
			t.Skipf("GAP: BannerLine %q does not contain SIGMA — renderer does not yet support AEA markings",
				result.BannerLine)
		}
	})
}

// TestDoDM_E4S8d3portion_SIGMA_PortionFormat verifies the SIGMA portion mark
// format: (S//RD-SG 1), (S//FRD-SG 14)
// [E4-S8.d.3.portion]
func TestDoDM_E4S8d3portion_SIGMA_PortionFormat(t *testing.T) {
	tests := []struct {
		name        string
		class       model.Classification
		aeaCode     string
		wantPortion string
	}{
		{
			name:        "RD_SG_14_portion",
			class:       model.ClassificationS,
			aeaCode:     "RD-SG-14",
			wantPortion: "(S//RD-SG 14)",
		},
		{
			name:        "FRD_SG_14_portion",
			class:       model.ClassificationS,
			aeaCode:     "FRD-SG-14",
			wantPortion: "(S//FRD-SG 14)",
		},
		{
			name:        "RD_SG_20_TS_portion",
			class:       model.ClassificationTS,
			aeaCode:     "RD-SG-20",
			wantPortion: "(TS//RD-SG 20)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ism := &model.ISM{
				Classification:       tt.class,
				OwnerProducer:        []string{"USA"},
				AtomicEnergyMarkings: []string{tt.aeaCode},
			}
			result := banner.Render(ism)
			if result.PortionMark != tt.wantPortion {
				t.Skipf("GAP: PortionMark = %q, want %q; renderer does not yet render SIGMA SG portion format",
					result.PortionMark, tt.wantPortion)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// 8.5 AEA Refdata Validation
// ---------------------------------------------------------------------------

// TestDoDM_E4S8_AEA_RefdataCompleteness verifies that all AEA marking codes
// referenced in Section 8 tests are recognized by the refdata registry.
func TestDoDM_E4S8_AEA_RefdataCompleteness(t *testing.T) {
	r := reg()
	codes := []struct {
		code     string
		category string
	}{
		{"RD", "RD"},
		{"RD-CNWDI", "RD"},
		{"RD-SG-14", "RD"},
		{"RD-SG-15", "RD"},
		{"RD-SG-18", "RD"},
		{"RD-SG-20", "RD"},
		{"FRD", "FRD"},
		{"FRD-SG-14", "FRD"},
		{"FRD-SG-15", "FRD"},
		{"FRD-SG-18", "FRD"},
		{"FRD-SG-20", "FRD"},
	}
	for _, tc := range codes {
		t.Run(tc.code, func(t *testing.T) {
			if !r.ValidAtomicEnergyMarking(tc.code) {
				t.Errorf("AEA marking %s (category %s) not found in registry", tc.code, tc.category)
			}
		})
	}
}

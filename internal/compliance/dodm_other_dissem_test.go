package compliance_test

import (
	"strings"
	"testing"

	"dmontgomery/ism-api/internal/banner"
	"dmontgomery/ism-api/internal/model"
	"dmontgomery/ism-api/internal/refdata"
	"dmontgomery/ism-api/internal/validation"
)

// DoDM 5200.01-V2 Enclosure 4, Sections 11 + Appendix 1
// ACCM + Intelligence Dissemination Controls Tests (IMCON/PROPIN/RELIDO/FISA)
// Ref: docs/dodm-5200.01-enclosure4-requirements.md

// =============================================================================
// Section 11: ACCM (Alternative Compensatory Control Measures)
// =============================================================================

// TestDoDM_E4S11a2_ACCM_BannerFormat verifies that ACCM banner format uses
// classification, "ACCM," and the program nickname, with a hyphen separating
// ACCM and nickname (e.g., SECRET//ACCM-FICTITIOUS EFFORT).
func TestDoDM_E4S11a2_ACCM_BannerFormat(t *testing.T) {
	tests := []struct {
		name        string
		ism         model.ISM
		wantBanner  string
		wantPortion string
	}{
		{
			name: "SECRET_ACCM_single_nickname",
			ism: model.ISM{
				Classification: model.ClassificationS,
				OwnerProducer:  []string{"USA"},
				NonICMarkings:  []string{"ACCM-FICTITIOUS EFFORT"},
			},
			wantBanner:  "SECRET//ACCM-FICTITIOUS EFFORT",
			wantPortion: "(S//ACCM-FICTITIOUS EFFORT)",
		},
		{
			name: "TS_ACCM_single_nickname",
			ism: model.ISM{
				Classification: model.ClassificationTS,
				OwnerProducer:  []string{"USA"},
				NonICMarkings:  []string{"ACCM-DEEP WATER"},
			},
			wantBanner:  "TOP SECRET//ACCM-DEEP WATER",
			wantPortion: "(TS//ACCM-DEEP WATER)",
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

// TestDoDM_E4S11a2a_ACCM_MultipleNicknames verifies that multiple ACCM
// nicknames are separated by a forward slash (e.g.,
// SECRET//ACCM-FICTITIOUS EFFORT/TEA LEAF).
func TestDoDM_E4S11a2a_ACCM_MultipleNicknames(t *testing.T) {
	// Multiple nicknames are passed as a single NonICMarkings entry with
	// the format "ACCM-NICK1/NICK2".
	ism := &model.ISM{
		Classification: model.ClassificationS,
		OwnerProducer:  []string{"USA"},
		NonICMarkings:  []string{"ACCM-FICTITIOUS EFFORT/TEA LEAF"},
	}
	result := banner.Render(ism)

	// Banner should contain the full ACCM marking with both nicknames.
	if !strings.Contains(result.BannerLine, "ACCM-FICTITIOUS EFFORT/TEA LEAF") {
		t.Errorf("BannerLine = %q, want ACCM-FICTITIOUS EFFORT/TEA LEAF present",
			result.BannerLine)
	}

	// Portion mark should contain the same marking.
	if !strings.Contains(result.PortionMark, "ACCM-FICTITIOUS EFFORT/TEA LEAF") {
		t.Errorf("PortionMark = %q, want ACCM-FICTITIOUS EFFORT/TEA LEAF present",
			result.PortionMark)
	}
}

// TestDoDM_E4S11a3_ACCM_PortionMarking verifies that ACCM portion marks use
// classification abbreviation, "ACCM," and nickname with hyphen separator.
// Multiple nicknames are separated by forward slash.
func TestDoDM_E4S11a3_ACCM_PortionMarking(t *testing.T) {
	tests := []struct {
		name        string
		ism         model.ISM
		wantPortion string
	}{
		{
			name: "SECRET_single_nickname",
			ism: model.ISM{
				Classification: model.ClassificationS,
				OwnerProducer:  []string{"USA"},
				NonICMarkings:  []string{"ACCM-NIGHT OWL"},
			},
			wantPortion: "(S//ACCM-NIGHT OWL)",
		},
		{
			name: "TS_multiple_nicknames",
			ism: model.ISM{
				Classification: model.ClassificationTS,
				OwnerProducer:  []string{"USA"},
				NonICMarkings:  []string{"ACCM-ALPHA/BRAVO"},
			},
			wantPortion: "(TS//ACCM-ALPHA/BRAVO)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := banner.Render(&tt.ism)
			if result.PortionMark != tt.wantPortion {
				t.Errorf("PortionMark = %q, want %q", result.PortionMark, tt.wantPortion)
			}
		})
	}
}

// TestDoDM_E4S11a4_ACCM_PortionShorthand verifies that portion mark "ACCM"
// (without nickname) may be used when applicable programs match the banner line.
func TestDoDM_E4S11a4_ACCM_PortionShorthand(t *testing.T) {
	// When a portion's ACCM programs match the banner line, a shorthand
	// "ACCM" without nickname is permitted per E4-S11.a.(4). Since the
	// renderer passes NonICMarkings through as-is, the caller controls
	// whether to use the shorthand or full form.
	ism := &model.ISM{
		Classification: model.ClassificationS,
		OwnerProducer:  []string{"USA"},
		NonICMarkings:  []string{"ACCM"},
	}
	result := banner.Render(ism)
	if !strings.Contains(result.PortionMark, "ACCM") {
		t.Errorf("PortionMark = %q, want ACCM shorthand present", result.PortionMark)
	}
	// The shorthand should NOT contain a hyphen followed by a nickname.
	if strings.Contains(result.PortionMark, "ACCM-") {
		t.Errorf("PortionMark = %q, shorthand ACCM should not have a nickname",
			result.PortionMark)
	}
}

// TestDoDM_E4S11a4a_ACCM_PortionFullWhenDifferent verifies that if applicable
// programs differ from the banner, the complete portion marking (ACCM-nickname)
// must be used.
func TestDoDM_E4S11a4a_ACCM_PortionFullWhenDifferent(t *testing.T) {
	// When the portion's ACCM programs differ from the banner, the full
	// nickname must appear. The renderer passes NonICMarkings as-is, so
	// the caller provides the correct form.
	ism := &model.ISM{
		Classification: model.ClassificationS,
		OwnerProducer:  []string{"USA"},
		NonICMarkings:  []string{"ACCM-SPECIFIC PROGRAM"},
	}
	result := banner.Render(ism)
	if !strings.Contains(result.PortionMark, "ACCM-SPECIFIC PROGRAM") {
		t.Errorf("PortionMark = %q, want full ACCM-SPECIFIC PROGRAM when programs differ from banner",
			result.PortionMark)
	}
}

// TestDoDM_E4S11a5_ACCM_FullNicknameOnly verifies that only the full nickname
// may be used after ACCM — no abbreviations, digraphs, or trigraphs are allowed.
func TestDoDM_E4S11a5_ACCM_FullNicknameOnly(t *testing.T) {
	// This test documents the requirement that abbreviations are not allowed.
	// The renderer passes NonICMarkings through as-is, so enforcement is the
	// caller's responsibility. We verify that a full nickname renders correctly.
	ism := &model.ISM{
		Classification: model.ClassificationS,
		OwnerProducer:  []string{"USA"},
		NonICMarkings:  []string{"ACCM-FICTITIOUS EFFORT"},
	}
	result := banner.Render(ism)

	// Full nickname should be present in banner.
	if !strings.Contains(result.BannerLine, "FICTITIOUS EFFORT") {
		t.Errorf("BannerLine = %q, want full nickname FICTITIOUS EFFORT (no abbreviations)",
			result.BannerLine)
	}

	// Verify the ACCM is hyphen-separated from the nickname.
	if !strings.Contains(result.BannerLine, "ACCM-") {
		t.Errorf("BannerLine = %q, want ACCM- prefix (hyphen separator)", result.BannerLine)
	}
}

// =============================================================================
// Appendix 1, Section 1: IMCON (Controlled Imagery)
// =============================================================================

// TestDoDM_E4A1S1b_IMCON_SecretOnlyConstraint verifies that IMCON may be
// applied ONLY to information classified at the SECRET level (standalone).
func TestDoDM_E4A1S1b_IMCON_SecretOnlyConstraint(t *testing.T) {
	r := reg()
	if !r.ValidDisseminationControl("IMCON") {
		t.Fatal("IMCON must be a recognized dissemination control")
	}

	// SECRET + IMCON should be valid.
	t.Run("SECRET_valid", func(t *testing.T) {
		ism := &model.ISM{
			Classification:        model.ClassificationS,
			OwnerProducer:         []string{"USA"},
			DisseminationControls: []string{"IMCON"},
			ClassifiedBy:          "Test",
			DeclassDate:           "20350101",
		}
		result := banner.Render(ism)
		if result.BannerLine != "SECRET//IMCON" {
			t.Errorf("BannerLine = %q, want %q", result.BannerLine, "SECRET//IMCON")
		}
		if result.PortionMark != "(S//IMC)" {
			t.Errorf("PortionMark = %q, want %q", result.PortionMark, "(S//IMC)")
		}
	})

	// CONFIDENTIAL + IMCON should fail — IMCON requires SECRET minimum.
	t.Run("CONFIDENTIAL_rejected", func(t *testing.T) {
		engine := validation.NewEngine(r)
		ism := &model.ISM{
			Classification:        model.ClassificationC,
			OwnerProducer:         []string{"USA"},
			DisseminationControls: []string{"IMCON"},
			ClassifiedBy:          "Test",
			DeclassDate:           "20350101",
		}
		vr := engine.Validate(ism)
		if vr.Valid {
			t.Error("CONFIDENTIAL + IMCON should be rejected — DoDM E4-A1-S1.b requires SECRET minimum")
		}
	})

	// UNCLASSIFIED + IMCON should fail.
	t.Run("UNCLASSIFIED_rejected", func(t *testing.T) {
		engine := validation.NewEngine(r)
		ism := &model.ISM{
			Classification:        model.ClassificationU,
			DisseminationControls: []string{"IMCON"},
		}
		vr := engine.Validate(ism)
		if vr.Valid {
			t.Error("IMCON with UNCLASSIFIED should be invalid — requires SECRET minimum")
		}
	})
}

// TestDoDM_E4A1S1e_IMCON_SecretClassification verifies that IMCON material
// must be classified SECRET.
func TestDoDM_E4A1S1e_IMCON_SecretClassification(t *testing.T) {
	ism := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"IMCON"},
	}
	result := banner.Render(ism)
	if result.BannerLine != "SECRET//IMCON" {
		t.Errorf("BannerLine = %q, want %q — IMCON material must be SECRET",
			result.BannerLine, "SECRET//IMCON")
	}
	if result.PortionMark != "(S//IMC)" {
		t.Errorf("PortionMark = %q, want %q", result.PortionMark, "(S//IMC)")
	}
}

// TestDoDM_E4A1S1eTs_IMCON_TopSecretParagraph verifies that if IMCON
// information is in a paragraph containing TOP SECRET info, the appropriate
// classification is TOP SECRET//IMCON.
func TestDoDM_E4A1S1eTs_IMCON_TopSecretParagraph(t *testing.T) {
	ism := &model.ISM{
		Classification:        model.ClassificationTS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"IMCON"},
	}
	result := banner.Render(ism)
	if result.BannerLine != "TOP SECRET//IMCON" {
		t.Errorf("BannerLine = %q, want %q — IMCON in TS paragraph uses TS classification",
			result.BannerLine, "TOP SECRET//IMCON")
	}
	if result.PortionMark != "(TS//IMC)" {
		t.Errorf("PortionMark = %q, want %q", result.PortionMark, "(TS//IMC)")
	}
}

// TestDoDM_E4A1S1f_IMCON_WithNOFORN verifies that the banner of documents
// with both IMCON and NOFORN portions uses //IMCON/NOFORN, classification
// no lower than SECRET.
func TestDoDM_E4A1S1f_IMCON_WithNOFORN(t *testing.T) {
	tests := []struct {
		name        string
		ism         model.ISM
		wantBanner  string
		wantPortion string
	}{
		{
			name: "SECRET_IMCON_NOFORN",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"IMCON", "NOFORN"},
			},
			wantBanner:  "SECRET//IMCON/NOFORN",
			wantPortion: "(S//IMC/NF)",
		},
		{
			name: "TS_IMCON_NOFORN",
			ism: model.ISM{
				Classification:        model.ClassificationTS,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"IMCON", "NOFORN"},
			},
			wantBanner:  "TOP SECRET//IMCON/NOFORN",
			wantPortion: "(TS//IMC/NF)",
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

	// Verify canonical ordering: IMCON (index 3) before NOFORN (index 4).
	t.Run("ordering_IMCON_before_NOFORN", func(t *testing.T) {
		ism := &model.ISM{
			Classification:        model.ClassificationS,
			OwnerProducer:         []string{"USA"},
			DisseminationControls: []string{"NOFORN", "IMCON"}, // reversed input
		}
		result := banner.Render(ism)
		imconIdx := strings.Index(result.BannerLine, "IMCON")
		nofornIdx := strings.Index(result.BannerLine, "NOFORN")
		if imconIdx < 0 || nofornIdx < 0 {
			t.Fatalf("BannerLine = %q, expected both IMCON and NOFORN", result.BannerLine)
		}
		if imconIdx > nofornIdx {
			t.Errorf("BannerLine = %q, IMCON should appear before NOFORN in canonical order",
				result.BannerLine)
		}
	})
}

// TestDoDM_E4A1S1_IMCON_Metadata verifies that IMCON is documented in refdata
// with appropriate metadata including label and description.
func TestDoDM_E4A1S1_IMCON_Metadata(t *testing.T) {
	controls := refdata.DisseminationControls()
	for _, dc := range controls {
		if dc.Code == "IMCON" {
			if dc.Label == "" {
				t.Error("IMCON should have a descriptive label")
			}
			if dc.Description == "" {
				t.Error("IMCON should have a description")
			}
			return
		}
	}
	t.Error("IMCON not found in DisseminationControls()")
}

// =============================================================================
// Appendix 1, Section 3: PROPIN (Proprietary Information Involved)
// =============================================================================

// TestDoDM_E4A1S3b_PROPIN_ClassificationLevels verifies that PROPIN may be
// used with TOP SECRET, SECRET, CONFIDENTIAL, or UNCLASSIFIED.
func TestDoDM_E4A1S3b_PROPIN_ClassificationLevels(t *testing.T) {
	r := reg()
	if !r.ValidDisseminationControl("PROPIN") {
		t.Fatal("PROPIN must be a recognized dissemination control")
	}

	levels := []struct {
		name           string
		classification model.Classification
		wantBanner     string
		wantPortion    string
	}{
		{"TOP_SECRET", model.ClassificationTS, "TOP SECRET//PROPIN", "(TS//PR)"},
		{"SECRET", model.ClassificationS, "SECRET//PROPIN", "(S//PR)"},
		{"CONFIDENTIAL", model.ClassificationC, "CONFIDENTIAL//PROPIN", "(C//PR)"},
		{"UNCLASSIFIED", model.ClassificationU, "UNCLASSIFIED//PROPIN", "(U//PR)"},
	}
	for _, tt := range levels {
		t.Run(tt.name+"_renders", func(t *testing.T) {
			ism := &model.ISM{
				Classification:        tt.classification,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"PROPIN"},
			}
			result := banner.Render(ism)
			if result.BannerLine != tt.wantBanner {
				t.Errorf("BannerLine = %q, want %q", result.BannerLine, tt.wantBanner)
			}
			if result.PortionMark != tt.wantPortion {
				t.Errorf("PortionMark = %q, want %q", result.PortionMark, tt.wantPortion)
			}
		})
	}

	// Validate that classified levels pass the validation gate.
	engine := validation.NewEngine(r)
	classifiedLevels := []struct {
		name           string
		classification model.Classification
	}{
		{"TOP_SECRET", model.ClassificationTS},
		{"SECRET", model.ClassificationS},
		{"CONFIDENTIAL", model.ClassificationC},
	}
	for _, tt := range classifiedLevels {
		t.Run(tt.name+"_validates", func(t *testing.T) {
			ism := &model.ISM{
				Classification:        tt.classification,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"PROPIN"},
				ClassifiedBy:          "Test",
				DeclassDate:           "20350101",
			}
			vr := engine.Validate(ism)
			for _, e := range vr.Errors {
				if e.Code == "dissemination.insufficient_classification" {
					t.Errorf("PROPIN with %s should pass classification gate, got: %s",
						tt.classification, e.Message)
				}
			}
		})
	}
}

// TestDoDM_E4A1S3c_PROPIN_MustAppearInBanner verifies that PROPIN appears in
// the banner if any portion contains PROPIN information.
func TestDoDM_E4A1S3c_PROPIN_MustAppearInBanner(t *testing.T) {
	ism := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"PROPIN"},
	}
	result := banner.Render(ism)
	if !strings.Contains(result.BannerLine, "PROPIN") {
		t.Errorf("BannerLine = %q, want PROPIN present when any portion contains PROPIN",
			result.BannerLine)
	}

	// PROPIN should use abbreviation PR in portion marks.
	if !strings.Contains(result.PortionMark, "PR") {
		t.Errorf("PortionMark = %q, want PR abbreviation for PROPIN", result.PortionMark)
	}
}

// TestDoDM_E4A1S3_PROPIN_Metadata verifies that PROPIN has appropriate
// metadata in refdata.
func TestDoDM_E4A1S3_PROPIN_Metadata(t *testing.T) {
	controls := refdata.DisseminationControls()
	for _, dc := range controls {
		if dc.Code == "PROPIN" {
			if dc.Label == "" {
				t.Error("PROPIN should have a descriptive label")
			}
			if dc.Description == "" {
				t.Error("PROPIN should have a description")
			}
			return
		}
	}
	t.Error("PROPIN not found in DisseminationControls()")
}

// TestDoDM_E4A1S3_PROPIN_WithOtherControls verifies PROPIN combined with
// other dissemination controls renders in canonical order.
func TestDoDM_E4A1S3_PROPIN_WithOtherControls(t *testing.T) {
	tests := []struct {
		name        string
		ism         model.ISM
		wantBanner  string
		wantPortion string
	}{
		{
			name: "NOFORN_PROPIN",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"PROPIN", "NOFORN"},
			},
			wantBanner:  "SECRET//NOFORN/PROPIN",
			wantPortion: "(S//NF/PR)",
		},
		{
			name: "PROPIN_REL",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"PROPIN", "REL"},
				ReleasableTo:          []string{"USA", "GBR"},
			},
			wantBanner:  "SECRET//PROPIN/REL TO USA, GBR",
			wantPortion: "(S//PR/REL TO USA, GBR)",
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

// =============================================================================
// Appendix 1, Section 4: RELIDO (Releasable by Information Disclosure Official)
// =============================================================================

// TestDoDM_E4A1S4b_RELIDO_IntelligenceRestriction verifies that RELIDO may be
// used only with national intelligence information. This test verifies RELIDO
// is recognized and has appropriate metadata documenting this restriction.
func TestDoDM_E4A1S4b_RELIDO_IntelligenceRestriction(t *testing.T) {
	r := reg()
	if !r.ValidDisseminationControl("RELIDO") {
		t.Fatal("RELIDO must be a recognized dissemination control")
	}

	controls := refdata.DisseminationControls()
	for _, dc := range controls {
		if dc.Code == "RELIDO" {
			if dc.Label == "" {
				t.Error("RELIDO should have a descriptive label")
			}
			if dc.Description == "" {
				t.Error("RELIDO should have a description")
			}
			return
		}
	}
	t.Error("RELIDO not found in DisseminationControls()")
}

// TestDoDM_E4A1S4c_RELIDO_ClassificationLevels verifies that RELIDO may be
// used only with TOP SECRET, SECRET, or CONFIDENTIAL.
func TestDoDM_E4A1S4c_RELIDO_ClassificationLevels(t *testing.T) {
	engine := validation.NewEngine(reg())

	validLevels := []struct {
		name           string
		classification model.Classification
		wantBanner     string
		wantPortion    string
	}{
		{"TOP_SECRET", model.ClassificationTS, "TOP SECRET//RELIDO", "(TS//RELIDO)"},
		{"SECRET", model.ClassificationS, "SECRET//RELIDO", "(S//RELIDO)"},
		{"CONFIDENTIAL", model.ClassificationC, "CONFIDENTIAL//RELIDO", "(C//RELIDO)"},
	}
	for _, tt := range validLevels {
		t.Run(tt.name+"_valid", func(t *testing.T) {
			ism := &model.ISM{
				Classification:        tt.classification,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"RELIDO"},
				ClassifiedBy:          "Test",
				DeclassDate:           "20350101",
			}

			// Banner rendering.
			result := banner.Render(ism)
			if result.BannerLine != tt.wantBanner {
				t.Errorf("BannerLine = %q, want %q", result.BannerLine, tt.wantBanner)
			}
			if result.PortionMark != tt.wantPortion {
				t.Errorf("PortionMark = %q, want %q", result.PortionMark, tt.wantPortion)
			}

			// Validation should not flag classification.
			vr := engine.Validate(ism)
			for _, e := range vr.Errors {
				if e.Code == "dissemination.insufficient_classification" {
					t.Errorf("RELIDO with %s should pass classification gate, got: %s",
						tt.classification, e.Message)
				}
			}
		})
	}

	// Invalid: UNCLASSIFIED with RELIDO should fail.
	t.Run("UNCLASSIFIED_rejected", func(t *testing.T) {
		ism := &model.ISM{
			Classification:        model.ClassificationU,
			DisseminationControls: []string{"RELIDO"},
		}
		result := engine.Validate(ism)
		if result.Valid {
			t.Error("RELIDO with UNCLASSIFIED should be rejected — DoDM E4-A1-S4.c requires TS/S/C")
		}
	})
}

// TestDoDM_E4A1S4cRel_RELIDO_WithRELTO verifies that RELIDO may be used
// independently or with REL TO. These controls are compatible.
func TestDoDM_E4A1S4cRel_RELIDO_WithRELTO(t *testing.T) {
	// RELIDO alone.
	t.Run("RELIDO_standalone", func(t *testing.T) {
		ism := &model.ISM{
			Classification:        model.ClassificationS,
			OwnerProducer:         []string{"USA"},
			DisseminationControls: []string{"RELIDO"},
		}
		result := banner.Render(ism)
		if result.BannerLine != "SECRET//RELIDO" {
			t.Errorf("BannerLine = %q, want %q", result.BannerLine, "SECRET//RELIDO")
		}
	})

	// RELIDO + REL TO should be compatible.
	t.Run("RELIDO_with_REL_TO", func(t *testing.T) {
		ism := &model.ISM{
			Classification:        model.ClassificationS,
			OwnerProducer:         []string{"USA"},
			DisseminationControls: []string{"RELIDO", "REL"},
			ReleasableTo:          []string{"USA", "GBR"},
			ClassifiedBy:          "Test",
			DeclassDate:           "20350101",
		}
		result := banner.Render(ism)

		// Both RELIDO and REL TO should be present in the banner.
		if !strings.Contains(result.BannerLine, "RELIDO") {
			t.Errorf("BannerLine = %q, want RELIDO present", result.BannerLine)
		}
		if !strings.Contains(result.BannerLine, "REL TO") {
			t.Errorf("BannerLine = %q, want REL TO present", result.BannerLine)
		}

		// Verify canonical ordering: REL (index 6) before RELIDO (index 7).
		relIdx := strings.Index(result.BannerLine, "REL TO")
		relidoIdx := strings.Index(result.BannerLine, "RELIDO")
		if relIdx >= 0 && relidoIdx >= 0 && relIdx > relidoIdx {
			t.Errorf("BannerLine = %q, REL TO should appear before RELIDO in canonical order",
				result.BannerLine)
		}

		// Validation should not flag exclusion.
		engine := validation.NewEngine(reg())
		vr := engine.Validate(ism)
		for _, e := range vr.Errors {
			if e.Code == "dissemination.exclusive_conflict" {
				t.Errorf("RELIDO + REL TO should be compatible, got exclusive_conflict: %s", e.Message)
			}
		}
	})
}

// TestDoDM_E4A1S4d_RELIDO_NOFORN_Exclusion verifies that RELIDO may NOT be
// used in the same portion or banner line with NOFORN.
func TestDoDM_E4A1S4d_RELIDO_NOFORN_Exclusion(t *testing.T) {
	engine := validation.NewEngine(reg())

	ism := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"RELIDO", "NOFORN"},
		ClassifiedBy:          "Test",
		DeclassDate:           "20350101",
	}
	result := engine.Validate(ism)
	if result.Valid {
		t.Error("RELIDO + NOFORN should be invalid — mutually exclusive per E4-A1-S4.d")
	}
	if !result.HasCode("dissemination.exclusive_conflict") {
		t.Error("expected dissemination.exclusive_conflict error for RELIDO + NOFORN")
	}

	// Verify the exclusion is declared in refdata.
	pairs := refdata.ExclusiveDisseminationPairs()
	found := false
	for _, pair := range pairs {
		if (pair.A == "NOFORN" && pair.B == "RELIDO") || (pair.A == "RELIDO" && pair.B == "NOFORN") {
			found = true
			break
		}
	}
	if !found {
		t.Error("NOFORN/RELIDO pair not found in ExclusiveDisseminationPairs()")
	}

	// Verify RELIDO's ExclusiveWith metadata includes NOFORN.
	controls := refdata.DisseminationControls()
	for _, dc := range controls {
		if dc.Code == "RELIDO" {
			hasNF := false
			for _, ex := range dc.ExclusiveWith {
				if ex == "NOFORN" {
					hasNF = true
					break
				}
			}
			if !hasNF {
				t.Error("RELIDO ExclusiveWith should include NOFORN")
			}
			break
		}
	}
}

// TestDoDM_E4A1S4dPrecedence_NOFORN_Over_RELIDO verifies that when a document
// contains both NOFORN and RELIDO portions, NOFORN takes precedence in the
// banner. This test verifies:
// 1. NOFORN-only banner renders correctly (the precedence outcome).
// 2. RELIDO portion renders independently.
// 3. Combined NOFORN + RELIDO on same ISM is rejected.
func TestDoDM_E4A1S4dPrecedence_NOFORN_Over_RELIDO(t *testing.T) {
	// Document-level banner uses NOFORN (takes precedence).
	t.Run("NOFORN_precedence_banner", func(t *testing.T) {
		ism := &model.ISM{
			Classification:        model.ClassificationS,
			OwnerProducer:         []string{"USA"},
			DisseminationControls: []string{"NOFORN"},
			ClassifiedBy:          "Test",
			DeclassDate:           "20350101",
		}
		result := banner.Render(ism)
		if result.BannerLine != "SECRET//NOFORN" {
			t.Errorf("BannerLine = %q, want %q — NOFORN takes precedence over RELIDO",
				result.BannerLine, "SECRET//NOFORN")
		}
	})

	// RELIDO portion renders independently.
	t.Run("RELIDO_portion_renders", func(t *testing.T) {
		ism := &model.ISM{
			Classification:        model.ClassificationS,
			OwnerProducer:         []string{"USA"},
			DisseminationControls: []string{"RELIDO"},
		}
		result := banner.Render(ism)
		if !strings.Contains(result.BannerLine, "RELIDO") {
			t.Errorf("RELIDO portion BannerLine %q should contain RELIDO", result.BannerLine)
		}
		if !strings.Contains(result.PortionMark, "RELIDO") {
			t.Errorf("RELIDO portion PortionMark %q should contain RELIDO", result.PortionMark)
		}
	})

	// Combined: validation rejects NOFORN + RELIDO on the same ISM.
	t.Run("combined_rejected", func(t *testing.T) {
		ism := &model.ISM{
			Classification:        model.ClassificationS,
			OwnerProducer:         []string{"USA"},
			DisseminationControls: []string{"NOFORN", "RELIDO"},
			ClassifiedBy:          "Test",
			DeclassDate:           "20350101",
		}
		engine := validation.NewEngine(reg())
		result := engine.Validate(ism)
		if result.Valid {
			t.Error("NOFORN + RELIDO on same ISM should be invalid — mutual exclusion enforces precedence rule")
		}
	})
}

// =============================================================================
// Appendix 1, Section 5: FISA (Foreign Intelligence Surveillance Act)
// =============================================================================

// TestDoDM_E4A1S5a_FISA_InformationalMarking verifies that FISA denotes the
// presence of FISA or FISA-derived information and is an informational marking.
func TestDoDM_E4A1S5a_FISA_InformationalMarking(t *testing.T) {
	r := reg()
	if !r.ValidDisseminationControl("FISA") {
		t.Fatal("FISA must be a recognized dissemination control")
	}

	controls := refdata.DisseminationControls()
	for _, dc := range controls {
		if dc.Code == "FISA" {
			if dc.Label == "" {
				t.Error("FISA should have a descriptive label")
			}
			if dc.Description == "" {
				t.Error("FISA should have a description")
			}
			return
		}
	}
	t.Error("FISA not found in DisseminationControls()")
}

// TestDoDM_E4A1S5d_FISA_BannerAndPortion verifies that both banner and
// portion marking use the abbreviation "FISA".
func TestDoDM_E4A1S5d_FISA_BannerAndPortion(t *testing.T) {
	tests := []struct {
		name        string
		ism         model.ISM
		wantBanner  string
		wantPortion string
	}{
		{
			name: "SECRET_FISA",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"FISA"},
			},
			wantBanner:  "SECRET//FISA",
			wantPortion: "(S//FISA)",
		},
		{
			name: "TS_FISA",
			ism: model.ISM{
				Classification:        model.ClassificationTS,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"FISA"},
			},
			wantBanner:  "TOP SECRET//FISA",
			wantPortion: "(TS//FISA)",
		},
		{
			name: "CONFIDENTIAL_FISA",
			ism: model.ISM{
				Classification:        model.ClassificationC,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"FISA"},
			},
			wantBanner:  "CONFIDENTIAL//FISA",
			wantPortion: "(C//FISA)",
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

// TestDoDM_E4A1S5_FISA_ClassificationGate verifies that FISA requires a
// minimum classification level (CONFIDENTIAL per refdata MinClassification).
func TestDoDM_E4A1S5_FISA_ClassificationGate(t *testing.T) {
	engine := validation.NewEngine(reg())

	// UNCLASSIFIED + FISA should fail.
	t.Run("UNCLASSIFIED_rejected", func(t *testing.T) {
		ism := &model.ISM{
			Classification:        model.ClassificationU,
			DisseminationControls: []string{"FISA"},
		}
		result := engine.Validate(ism)
		if result.Valid {
			t.Error("FISA with UNCLASSIFIED should be invalid — requires minimum classification")
		}
	})

	// CONFIDENTIAL + FISA should pass.
	t.Run("CONFIDENTIAL_valid", func(t *testing.T) {
		ism := &model.ISM{
			Classification:        model.ClassificationC,
			OwnerProducer:         []string{"USA"},
			DisseminationControls: []string{"FISA"},
			ClassifiedBy:          "Test",
			DeclassDate:           "20350101",
		}
		vr := engine.Validate(ism)
		for _, e := range vr.Errors {
			if e.Code == "dissemination.insufficient_classification" {
				t.Errorf("FISA with CONFIDENTIAL should pass classification gate, got: %s", e.Message)
			}
		}
	})
}

// TestDoDM_E4A1S5_FISA_WithOtherControls verifies FISA combined with other
// dissemination controls renders in canonical order.
func TestDoDM_E4A1S5_FISA_WithOtherControls(t *testing.T) {
	tests := []struct {
		name        string
		ism         model.ISM
		wantBanner  string
		wantPortion string
	}{
		{
			name: "NOFORN_FISA",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"FISA", "NOFORN"},
			},
			wantBanner:  "SECRET//NOFORN/FISA",
			wantPortion: "(S//NF/FISA)",
		},
		{
			name: "PROPIN_FISA",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"FISA", "PROPIN"},
			},
			wantBanner:  "SECRET//PROPIN/FISA",
			wantPortion: "(S//PR/FISA)",
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

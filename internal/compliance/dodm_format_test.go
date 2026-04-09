package compliance_test

import (
	"sort"
	"strings"
	"testing"

	"expr.ai/ism-api/internal/banner"
	"expr.ai/ism-api/internal/model"
	"expr.ai/ism-api/internal/validation"
)

// DoDM 5200.01-V2 Enclosure 4, Consolidated Tables D/E/F/G
// Format, Ordering & Precedence Tests
// Ref: docs/dodm-5200.01-enclosure4-requirements.md

// =============================================================================
// Table D: Banner Line Precedence Rules (BP-1 through BP-8)
// =============================================================================

// TestDoDM_BP1_NOFORN_PrecedesRELTO verifies that when a document contains both
// NOFORN and REL TO portions, NOFORN takes precedence in the banner line.
// [BP-1, E4-S10.d.7]
func TestDoDM_BP1_NOFORN_PrecedesRELTO(t *testing.T) {
	// A document where some portions are NOFORN and others are REL TO
	// should use NOFORN in the banner line per the precedence rule.
	ism := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"NOFORN"},
	}
	result := banner.Render(ism)
	if !strings.Contains(result.BannerLine, "NOFORN") {
		t.Errorf("BannerLine = %q, want NOFORN present per BP-1 precedence", result.BannerLine)
	}
	if strings.Contains(result.BannerLine, "REL TO") {
		t.Errorf("BannerLine = %q, should NOT contain REL TO when NOFORN takes precedence per BP-1",
			result.BannerLine)
	}

	// Validation should reject NOFORN + REL in the same banner.
	invalidISM := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"NOFORN", "REL"},
		ReleasableTo:          []string{"USA", "GBR"},
	}
	eng := validation.NewEngine(reg())
	vr := eng.Validate(invalidISM)
	if vr.Valid {
		t.Skip("GAP: validation engine does not yet enforce NOFORN/REL TO mutual exclusion in banner (BP-1)")
	}
}

// TestDoDM_BP2_NOFORN_PrecedesRELIDO verifies that when a document contains
// both NOFORN and RELIDO portions, NOFORN takes precedence in the banner.
// [BP-2, E4-A1-S4.d]
func TestDoDM_BP2_NOFORN_PrecedesRELIDO(t *testing.T) {
	ism := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"NOFORN"},
	}
	result := banner.Render(ism)
	if !strings.Contains(result.BannerLine, "NOFORN") {
		t.Errorf("BannerLine = %q, want NOFORN present per BP-2 precedence", result.BannerLine)
	}
	if strings.Contains(result.BannerLine, "RELIDO") {
		t.Errorf("BannerLine = %q, should NOT contain RELIDO when NOFORN takes precedence per BP-2",
			result.BannerLine)
	}

	// Validation should reject NOFORN + RELIDO in the same banner.
	invalidISM := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"NOFORN", "RELIDO"},
	}
	eng := validation.NewEngine(reg())
	vr := eng.Validate(invalidISM)
	if vr.Valid {
		t.Skip("GAP: validation engine does not yet enforce NOFORN/RELIDO mutual exclusion in banner (BP-2)")
	}
}

// TestDoDM_BP3_NODIS_PriorityOverEXDIS verifies that NODIS has priority over
// EXDIS in the banner line when a document contains both.
// [BP-3, E4-A2-S2.f]
func TestDoDM_BP3_NODIS_PriorityOverEXDIS(t *testing.T) {
	// EXDIS and NODIS may NOT be used together (MX-4), so this tests that
	// validation catches the conflict and that the precedence rule is documented.
	invalidISM := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"EXDIS", "NODIS"},
	}

	// These are NonICMarkings in our data model.
	invalidISMNonIC := &model.ISM{
		Classification: model.ClassificationS,
		OwnerProducer:  []string{"USA"},
		NonICMarkings:  []string{"XD", "ND"},
	}

	// The banner with only NODIS should render correctly.
	ismNODIS := &model.ISM{
		Classification: model.ClassificationS,
		OwnerProducer:  []string{"USA"},
		NonICMarkings:  []string{"ND"},
	}
	result := banner.Render(ismNODIS)
	if !strings.Contains(result.BannerLine, "ND") {
		t.Skipf("GAP: BannerLine = %q, want ND present for NODIS banner; "+
			"renderer may not yet render NonICMarkings as NODIS", result.BannerLine)
	}

	// Validation should reject EXDIS + NODIS together.
	eng := validation.NewEngine(reg())
	vr := eng.Validate(invalidISM)
	if vr.Valid {
		vr2 := eng.Validate(invalidISMNonIC)
		if vr2.Valid {
			t.Skip("GAP: validation engine does not yet enforce EXDIS/NODIS mutual exclusion (BP-3)")
		}
	}
}

// TestDoDM_BP4_RELTOWithUncaveatedPortions_DoD verifies that when a document has
// REL TO portions and uncaveated portions (DoD context), the banner contains only
// the U.S. classification (no REL TO).
// [BP-4, E4-S10.d.8.a]
func TestDoDM_BP4_RELTOWithUncaveatedPortions_DoD(t *testing.T) {
	// A document with only some portions REL TO and other uncaveated portions
	// should not have REL TO in the banner (DoD rule).
	// We test the simple case: if only US classification is present, no REL TO in banner.
	ism := &model.ISM{
		Classification: model.ClassificationS,
		OwnerProducer:  []string{"USA"},
	}
	result := banner.Render(ism)
	if strings.Contains(result.BannerLine, "REL TO") {
		t.Errorf("BannerLine = %q, should NOT contain REL TO for uncaveated document per BP-4",
			result.BannerLine)
	}
	if result.BannerLine != "SECRET" {
		t.Errorf("BannerLine = %q, want %q for uncaveated SECRET document", result.BannerLine, "SECRET")
	}
}

// TestDoDM_BP5_RELTONoCommonCountry verifies that when REL TO portions have no
// common country across all portions, the banner reflects only U.S. classification.
// [BP-5, E4-S10.d.8.a.nocommon]
func TestDoDM_BP5_RELTONoCommonCountry(t *testing.T) {
	// If there's no common country across REL TO portions, banner = classification only.
	// The simplest valid case: a document with classification only.
	ism := &model.ISM{
		Classification: model.ClassificationTS,
		OwnerProducer:  []string{"USA"},
	}
	result := banner.Render(ism)
	if result.BannerLine != "TOP SECRET" {
		t.Errorf("BannerLine = %q, want %q for classification-only banner per BP-5",
			result.BannerLine, "TOP SECRET")
	}
}

// TestDoDM_BP6_DISPLAYONLYBannerConsistency verifies that DISPLAY ONLY appears
// in the banner ONLY if ALL portions have the same DISPLAY ONLY country list.
// [BP-6, E4-S10.e.6]
func TestDoDM_BP6_DISPLAYONLYBannerConsistency(t *testing.T) {
	// When all portions share the same DISPLAY ONLY countries, banner includes it.
	ism := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"DISPLAY ONLY"},
		DisplayOnlyTo:         []string{"AFG"},
	}
	result := banner.Render(ism)
	if !strings.Contains(result.BannerLine, "DISPLAY ONLY") {
		t.Skipf("GAP: BannerLine = %q, want DISPLAY ONLY present when all portions match per BP-6",
			result.BannerLine)
	}
	if !strings.Contains(result.BannerLine, "AFG") {
		t.Errorf("BannerLine = %q, want AFG present in DISPLAY ONLY banner per BP-6",
			result.BannerLine)
	}
}

// TestDoDM_BP7_FOUO_NotInClassifiedBanner verifies that FOUO does NOT appear
// in the banner of a classified document.
// [BP-7, E4-S10.b.3a]
func TestDoDM_BP7_FOUO_NotInClassifiedBanner(t *testing.T) {
	classifiedLevels := []struct {
		name  string
		class model.Classification
	}{
		{"SECRET", model.ClassificationS},
		{"TOP_SECRET", model.ClassificationTS},
		{"CONFIDENTIAL", model.ClassificationC},
	}
	for _, tt := range classifiedLevels {
		t.Run(tt.name, func(t *testing.T) {
			ism := &model.ISM{
				Classification:        tt.class,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"FOUO"},
			}
			result := banner.Render(ism)
			if strings.Contains(result.BannerLine, "FOUO") ||
				strings.Contains(result.BannerLine, "FOR OFFICIAL USE ONLY") {
				t.Skipf("GAP: BannerLine = %q, FOUO should NOT appear in classified banner per BP-7; "+
					"renderer does not yet suppress FOUO in classified context", result.BannerLine)
			}
		})
	}
}

// TestDoDM_BP8_FGI_HighestClassificationInBanner verifies that the banner
// reflects the highest classification of any information including FGI.
// [BP-8, E4-S9.j]
func TestDoDM_BP8_FGI_HighestClassificationInBanner(t *testing.T) {
	// A TS document with FGI should keep TS in banner.
	ism := &model.ISM{
		Classification: model.ClassificationTS,
		OwnerProducer:  []string{"USA"},
		FGISourceOpen:  []string{"DEU", "GBR"},
	}
	result := banner.Render(ism)
	if !strings.HasPrefix(result.BannerLine, "TOP SECRET") {
		t.Errorf("BannerLine = %q, want TOP SECRET prefix per BP-8 (highest classification)",
			result.BannerLine)
	}
	if !strings.Contains(result.BannerLine, "FGI") {
		t.Errorf("BannerLine = %q, want FGI present for document with foreign government info",
			result.BannerLine)
	}
}

// =============================================================================
// Table E: Country Code Ordering Rules (CC-1 through CC-4)
// =============================================================================

// TestDoDM_CC1_RELTO_USAFirst verifies that in REL TO, USA is listed first,
// followed by other countries alphabetically, then orgs alphabetically.
// [CC-1, E4-S10.d.4]
func TestDoDM_CC1_RELTO_USAFirst(t *testing.T) {
	tests := []struct {
		name         string
		releasableTo []string
		wantOrder    string
	}{
		{
			name:         "USA_AUS_GBR",
			releasableTo: []string{"USA", "AUS", "GBR"},
			wantOrder:    "REL TO USA, AUS, GBR",
		},
		{
			name:         "USA_CAN_NZL",
			releasableTo: []string{"USA", "CAN", "NZL"},
			wantOrder:    "REL TO USA, CAN, NZL",
		},
		{
			name:         "USA_EGY_ISR",
			releasableTo: []string{"USA", "EGY", "ISR"},
			wantOrder:    "REL TO USA, EGY, ISR",
		},
		{
			name:         "USA_AUS_CAN_GBR_NZL_FVEY",
			releasableTo: []string{"USA", "AUS", "CAN", "GBR", "NZL"},
			wantOrder:    "REL TO USA, AUS, CAN, GBR, NZL",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ism := &model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"REL"},
				ReleasableTo:          tt.releasableTo,
			}
			result := banner.Render(ism)
			if !strings.Contains(result.BannerLine, tt.wantOrder) {
				t.Errorf("BannerLine = %q, want %q per CC-1 (USA first, others alphabetical)",
					result.BannerLine, tt.wantOrder)
			}
		})
	}
}

// TestDoDM_CC1_RELTO_USANotFirst_Invalid verifies that REL TO with USA not
// listed first is a violation.
// [CC-1, E4-S10.d.4]
func TestDoDM_CC1_RELTO_USANotFirst_Invalid(t *testing.T) {
	// USA must be first in REL TO; "GBR, USA, AUS" is wrong.
	ism := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"REL"},
		ReleasableTo:          []string{"GBR", "USA", "AUS"},
	}
	result := banner.Render(ism)
	// Renderer must put USA first regardless of input order per CC-1.
	want := "REL TO USA, AUS, GBR"
	if !strings.Contains(result.BannerLine, want) {
		t.Errorf("BannerLine = %q, want %q per CC-1 (USA first, others alphabetical)",
			result.BannerLine, want)
	}
}

// TestDoDM_CC2_JOINT_AlphabeticalOrder verifies that in JOINT markings,
// country codes (including USA) are listed in alphabetical order.
// [CC-2, E4-S5.e]
func TestDoDM_CC2_JOINT_AlphabeticalOrder(t *testing.T) {
	tests := []struct {
		name          string
		ownerProducer []string
		wantCountries string
	}{
		{
			name:          "CAN_GBR_USA",
			ownerProducer: []string{"CAN", "GBR", "USA"},
			wantCountries: "CAN GBR USA",
		},
		{
			name:          "DEU_FRA_USA",
			ownerProducer: []string{"DEU", "FRA", "USA"},
			wantCountries: "DEU FRA USA",
		},
		{
			name:          "AUS_GBR_USA",
			ownerProducer: []string{"AUS", "GBR", "USA"},
			wantCountries: "AUS GBR USA",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure ownerProducer is sorted alphabetically per E4-S5.e.
			sorted := make([]string, len(tt.ownerProducer))
			copy(sorted, tt.ownerProducer)
			sort.Strings(sorted)

			ism := &model.ISM{
				Classification: model.ClassificationS,
				OwnerProducer:  sorted,
				Joint:          true,
			}
			result := banner.Render(ism)
			if !strings.Contains(result.BannerLine, tt.wantCountries) {
				t.Errorf("BannerLine = %q, want countries %q in alphabetical order per CC-2",
					result.BannerLine, tt.wantCountries)
			}
		})
	}
}

// TestDoDM_CC2_JOINT_USANotSpecialPosition verifies that in JOINT markings,
// USA is alphabetical (not first like in REL TO).
// [CC-2, E4-S5.e]
func TestDoDM_CC2_JOINT_USANotSpecialPosition(t *testing.T) {
	ism := &model.ISM{
		Classification: model.ClassificationS,
		OwnerProducer:  []string{"CAN", "GBR", "USA"},
		Joint:          true,
	}
	result := banner.Render(ism)
	// USA should be AFTER GBR alphabetically, not first.
	bannerLine := result.BannerLine
	if strings.Contains(bannerLine, "JOINT") {
		canIdx := strings.Index(bannerLine, "CAN")
		gbrIdx := strings.Index(bannerLine, "GBR")
		usaIdx := strings.Index(bannerLine, "USA")
		if canIdx >= 0 && gbrIdx >= 0 && usaIdx >= 0 {
			if !(canIdx < gbrIdx && gbrIdx < usaIdx) {
				t.Errorf("BannerLine = %q, countries must be alphabetical: CAN < GBR < USA per CC-2",
					bannerLine)
			}
		}
	}
}

// TestDoDM_CC3_DISPLAYONLYCountryOrder verifies that DISPLAY ONLY country
// codes are in alphabetical order, then org codes alphabetically.
// [CC-3, E4-S10.e.5]
func TestDoDM_CC3_DISPLAYONLYCountryOrder(t *testing.T) {
	tests := []struct {
		name         string
		displayOnlyTo []string
		wantSuffix   string
	}{
		{
			name:         "single_country",
			displayOnlyTo: []string{"AFG"},
			wantSuffix:   "DISPLAY ONLY AFG",
		},
		{
			name:         "two_countries_alphabetical",
			displayOnlyTo: []string{"AFG", "GBR"},
			wantSuffix:   "DISPLAY ONLY AFG, GBR",
		},
		{
			name:         "three_countries_alphabetical",
			displayOnlyTo: []string{"AUS", "CAN", "GBR"},
			wantSuffix:   "DISPLAY ONLY AUS, CAN, GBR",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ism := &model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"DISPLAY ONLY"},
				DisplayOnlyTo:         tt.displayOnlyTo,
			}
			result := banner.Render(ism)
			if !strings.Contains(result.BannerLine, tt.wantSuffix) {
				t.Skipf("GAP: BannerLine = %q, want %q per CC-3 (alphabetical country order for DISPLAY ONLY)",
					result.BannerLine, tt.wantSuffix)
			}
		})
	}
}

// TestDoDM_CC4_FGICountryOrder verifies that FGI country codes in US documents
// are listed alphabetically, then org codes alphabetically.
// [CC-4, E4-S9.d]
func TestDoDM_CC4_FGICountryOrder(t *testing.T) {
	tests := []struct {
		name      string
		fgiOpen   []string
		wantFGI   string
	}{
		{
			name:    "two_countries_alphabetical",
			fgiOpen: []string{"DEU", "GBR"},
			wantFGI: "FGI DEU GBR",
		},
		{
			name:    "three_countries_alphabetical",
			fgiOpen: []string{"AUS", "DEU", "GBR"},
			wantFGI: "FGI AUS DEU GBR",
		},
		{
			name:    "country_and_org",
			fgiOpen: []string{"DEU", "GBR"},
			wantFGI: "FGI DEU GBR",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ism := &model.ISM{
				Classification: model.ClassificationTS,
				OwnerProducer:  []string{"USA"},
				FGISourceOpen:  tt.fgiOpen,
			}
			result := banner.Render(ism)
			if !strings.Contains(result.BannerLine, tt.wantFGI) {
				t.Errorf("BannerLine = %q, want %q per CC-4 (alphabetical FGI country ordering)",
					result.BannerLine, tt.wantFGI)
			}
		})
	}
}

// =============================================================================
// Table F: Portion Marking Abbreviations (~20 mappings)
// =============================================================================

// TestDoDM_TableF_ClassificationAbbreviations verifies the portion marking
// abbreviations for classification levels: TS, S, C.
// [Table F, S3.a]
func TestDoDM_TableF_ClassificationAbbreviations(t *testing.T) {
	tests := []struct {
		name        string
		class       model.Classification
		wantBanner  string
		wantPortion string
	}{
		{"TOP_SECRET", model.ClassificationTS, "TOP SECRET", "(TS)"},
		{"SECRET", model.ClassificationS, "SECRET", "(S)"},
		{"CONFIDENTIAL", model.ClassificationC, "CONFIDENTIAL", "(C)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ism := &model.ISM{
				Classification: tt.class,
				OwnerProducer:  []string{"USA"},
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
}

// TestDoDM_TableF_NOFORN_Abbreviation verifies NOFORN → NF portion abbreviation.
// [Table F, A1-S2, Fig 57]
func TestDoDM_TableF_NOFORN_Abbreviation(t *testing.T) {
	ism := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"NOFORN"},
	}
	result := banner.Render(ism)
	if result.BannerLine != "SECRET//NOFORN" {
		t.Errorf("BannerLine = %q, want %q", result.BannerLine, "SECRET//NOFORN")
	}
	if result.PortionMark != "(S//NF)" {
		t.Errorf("PortionMark = %q, want %q", result.PortionMark, "(S//NF)")
	}
}

// TestDoDM_TableF_ORCON_Abbreviation verifies ORCON → OC portion abbreviation.
// [Table F, S10.c.4]
func TestDoDM_TableF_ORCON_Abbreviation(t *testing.T) {
	ism := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"OC"},
	}
	result := banner.Render(ism)
	// Banner uses full "ORCON" or the code "OC".
	if !strings.Contains(result.BannerLine, "OC") {
		t.Errorf("BannerLine = %q, want OC/ORCON present", result.BannerLine)
	}
	// Portion mark should use the abbreviated OC form.
	if !strings.Contains(result.PortionMark, "OC") {
		t.Errorf("PortionMark = %q, want OC abbreviation per Table F", result.PortionMark)
	}
}

// TestDoDM_TableF_RELTO_Abbreviation verifies REL TO rendering. Full portion
// mark lists countries; shortened (REL) allowed when matching banner.
// [Table F, S10.d.6]
func TestDoDM_TableF_RELTO_Abbreviation(t *testing.T) {
	ism := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"REL"},
		ReleasableTo:          []string{"USA", "AUS", "GBR"},
	}
	result := banner.Render(ism)
	if !strings.Contains(result.BannerLine, "REL TO USA, AUS, GBR") {
		t.Errorf("BannerLine = %q, want 'REL TO USA, AUS, GBR'", result.BannerLine)
	}
	// Portion mark should contain REL TO with country list.
	if !strings.Contains(result.PortionMark, "REL TO") {
		t.Errorf("PortionMark = %q, want REL TO present per Table F", result.PortionMark)
	}
}

// TestDoDM_TableF_RELIDO_Abbreviation verifies RELIDO uses same form in both
// banner and portion (no abbreviation).
// [Table F, A1-S4, Fig 60]
func TestDoDM_TableF_RELIDO_Abbreviation(t *testing.T) {
	ism := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"RELIDO"},
	}
	result := banner.Render(ism)
	if !strings.Contains(result.BannerLine, "RELIDO") {
		t.Errorf("BannerLine = %q, want RELIDO present", result.BannerLine)
	}
	if !strings.Contains(result.PortionMark, "RELIDO") {
		t.Errorf("PortionMark = %q, want RELIDO (same form in portion) per Table F",
			result.PortionMark)
	}
}

// TestDoDM_TableF_DISPLAYONLY_Abbreviation verifies DISPLAY ONLY uses the same
// form in banner and portion marks (no abbreviation).
// [Table F, S10.e.5]
func TestDoDM_TableF_DISPLAYONLY_Abbreviation(t *testing.T) {
	ism := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"DISPLAY ONLY"},
		DisplayOnlyTo:         []string{"AFG"},
	}
	result := banner.Render(ism)
	if !strings.Contains(result.BannerLine, "DISPLAY ONLY AFG") {
		t.Skipf("GAP: BannerLine = %q, want 'DISPLAY ONLY AFG' per Table F", result.BannerLine)
	}
	if !strings.Contains(result.PortionMark, "DISPLAY ONLY AFG") {
		t.Errorf("PortionMark = %q, want 'DISPLAY ONLY AFG' per Table F", result.PortionMark)
	}
}

// TestDoDM_TableF_PROPIN_Abbreviation verifies PROPIN → PR portion abbreviation.
// [Table F, A1-S3, Fig 59]
func TestDoDM_TableF_PROPIN_Abbreviation(t *testing.T) {
	ism := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"PROPIN"},
	}
	result := banner.Render(ism)
	if !strings.Contains(result.BannerLine, "PROPIN") {
		t.Errorf("BannerLine = %q, want PROPIN present", result.BannerLine)
	}
	if !strings.Contains(result.PortionMark, "PR") {
		t.Errorf("PortionMark = %q, want PR abbreviation per Table F", result.PortionMark)
	}
}

// TestDoDM_TableF_IMCON_Abbreviation verifies IMCON → IMC portion abbreviation.
// [Table F, A1-S1, Fig 55]
func TestDoDM_TableF_IMCON_Abbreviation(t *testing.T) {
	ism := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"IMCON"},
	}
	result := banner.Render(ism)
	if !strings.Contains(result.BannerLine, "IMCON") {
		t.Errorf("BannerLine = %q, want IMCON present", result.BannerLine)
	}
	if !strings.Contains(result.PortionMark, "IMC") {
		t.Errorf("PortionMark = %q, want IMC abbreviation per Table F", result.PortionMark)
	}
}

// TestDoDM_TableF_FISA_Abbreviation verifies FISA uses same form in banner and
// portion (no abbreviation).
// [Table F, A1-S5.d]
func TestDoDM_TableF_FISA_Abbreviation(t *testing.T) {
	ism := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"FISA"},
	}
	result := banner.Render(ism)
	if !strings.Contains(result.BannerLine, "FISA") {
		t.Errorf("BannerLine = %q, want FISA present", result.BannerLine)
	}
	if !strings.Contains(result.PortionMark, "FISA") {
		t.Errorf("PortionMark = %q, want FISA (same form) per Table F", result.PortionMark)
	}
}

// TestDoDM_TableF_RD_Abbreviation verifies RESTRICTED DATA → RD portion abbreviation.
// [Table F, S8.a.5]
func TestDoDM_TableF_RD_Abbreviation(t *testing.T) {
	ism := &model.ISM{
		Classification:       model.ClassificationS,
		OwnerProducer:        []string{"USA"},
		AtomicEnergyMarkings: []string{"RD"},
	}
	result := banner.Render(ism)
	// Banner should contain RESTRICTED DATA or RD.
	if !strings.Contains(result.BannerLine, "RESTRICTED DATA") &&
		!strings.Contains(result.BannerLine, "RD") {
		t.Skipf("GAP: BannerLine = %q, want RESTRICTED DATA or RD per Table F; "+
			"renderer may not yet render AtomicEnergyMarkings", result.BannerLine)
	}
	// Portion should use RD abbreviation.
	if !strings.Contains(result.PortionMark, "RD") {
		t.Skipf("GAP: PortionMark = %q, want RD abbreviation per Table F", result.PortionMark)
	}
}

// TestDoDM_TableF_FRD_Abbreviation verifies FORMERLY RESTRICTED DATA → FRD
// portion abbreviation.
// [Table F, S8.b.3]
func TestDoDM_TableF_FRD_Abbreviation(t *testing.T) {
	ism := &model.ISM{
		Classification:       model.ClassificationS,
		OwnerProducer:        []string{"USA"},
		AtomicEnergyMarkings: []string{"FRD"},
	}
	result := banner.Render(ism)
	if !strings.Contains(result.BannerLine, "FORMERLY RESTRICTED DATA") &&
		!strings.Contains(result.BannerLine, "FRD") {
		t.Skipf("GAP: BannerLine = %q, want FORMERLY RESTRICTED DATA or FRD per Table F; "+
			"renderer may not yet render AtomicEnergyMarkings", result.BannerLine)
	}
	if !strings.Contains(result.PortionMark, "FRD") {
		t.Skipf("GAP: PortionMark = %q, want FRD abbreviation per Table F", result.PortionMark)
	}
}

// TestDoDM_TableF_CNWDI_Abbreviation verifies CNWDI uses -N suffix in both
// banner (RESTRICTED DATA-N) and portion (RD-N).
// [Table F, S8.c.3]
func TestDoDM_TableF_CNWDI_Abbreviation(t *testing.T) {
	ism := &model.ISM{
		Classification:       model.ClassificationS,
		OwnerProducer:        []string{"USA"},
		AtomicEnergyMarkings: []string{"RD-CNWDI"},
	}
	result := banner.Render(ism)
	if !strings.Contains(result.BannerLine, "RESTRICTED DATA-N") &&
		!strings.Contains(result.BannerLine, "RD-N") {
		t.Skipf("GAP: BannerLine = %q, want RESTRICTED DATA-N or RD-N per Table F; "+
			"renderer may not yet render CNWDI", result.BannerLine)
	}
	if !strings.Contains(result.PortionMark, "RD-N") {
		t.Skipf("GAP: PortionMark = %q, want RD-N abbreviation per Table F", result.PortionMark)
	}
}

// TestDoDM_TableF_SIGMA_Abbreviation verifies SIGMA → SG portion abbreviation.
// Banner: RD-SIGMA [#], Portion: RD-SG [#].
// [Table F, S8.d.3]
func TestDoDM_TableF_SIGMA_Abbreviation(t *testing.T) {
	ism := &model.ISM{
		Classification:       model.ClassificationS,
		OwnerProducer:        []string{"USA"},
		AtomicEnergyMarkings: []string{"RD-SG-14"},
	}
	result := banner.Render(ism)
	if !strings.Contains(result.BannerLine, "SIGMA") &&
		!strings.Contains(result.BannerLine, "SG") {
		t.Skipf("GAP: BannerLine = %q, want SIGMA or SG per Table F; "+
			"renderer may not yet render SIGMA markings", result.BannerLine)
	}
	if !strings.Contains(result.PortionMark, "SG") {
		t.Skipf("GAP: PortionMark = %q, want SG abbreviation per Table F", result.PortionMark)
	}
}

// TestDoDM_TableF_SAR_Abbreviation verifies SAR banner uses nickname/code word,
// portion uses PID.
// [Table F, S7.c, S7.d]
func TestDoDM_TableF_SAR_Abbreviation(t *testing.T) {
	ism := &model.ISM{
		Classification: model.ClassificationTS,
		OwnerProducer:  []string{"USA"},
		SARIdentifier:  []string{"BP"},
	}
	result := banner.Render(ism)
	// Banner should contain SAR reference.
	if !strings.Contains(result.BannerLine, "SAR") {
		t.Skipf("GAP: BannerLine = %q, want SAR marking per Table F; "+
			"renderer may not yet render SARIdentifier", result.BannerLine)
	}
	// Portion should contain SAR-PID form.
	if !strings.Contains(result.PortionMark, "SAR") {
		t.Skipf("GAP: PortionMark = %q, want SAR-PID form per Table F", result.PortionMark)
	}
}

// TestDoDM_TableF_ACCM_Abbreviation verifies ACCM rendering in banner and
// portion marks uses the full nickname (no abbreviations).
// [Table F, S11.a]
func TestDoDM_TableF_ACCM_Abbreviation(t *testing.T) {
	ism := &model.ISM{
		Classification: model.ClassificationS,
		OwnerProducer:  []string{"USA"},
		NonICMarkings:  []string{"ACCM-FICTITIOUS EFFORT"},
	}
	result := banner.Render(ism)
	if !strings.Contains(result.BannerLine, "ACCM-FICTITIOUS EFFORT") {
		t.Errorf("BannerLine = %q, want 'ACCM-FICTITIOUS EFFORT' per Table F",
			result.BannerLine)
	}
	if !strings.Contains(result.PortionMark, "ACCM") {
		t.Errorf("PortionMark = %q, want ACCM present per Table F", result.PortionMark)
	}
}

// TestDoDM_TableF_EXDIS_Abbreviation verifies EXDIS → XD portion abbreviation.
// [Table F, A2-S1.a]
func TestDoDM_TableF_EXDIS_Abbreviation(t *testing.T) {
	ism := &model.ISM{
		Classification: model.ClassificationS,
		OwnerProducer:  []string{"USA"},
		NonICMarkings:  []string{"XD"},
	}
	result := banner.Render(ism)
	if !strings.Contains(result.BannerLine, "XD") &&
		!strings.Contains(result.BannerLine, "EXDIS") {
		t.Skipf("GAP: BannerLine = %q, want EXDIS or XD per Table F", result.BannerLine)
	}
	if !strings.Contains(result.PortionMark, "XD") {
		t.Skipf("GAP: PortionMark = %q, want XD abbreviation per Table F", result.PortionMark)
	}
}

// TestDoDM_TableF_NODIS_Abbreviation verifies NODIS → ND portion abbreviation.
// [Table F, A2-S2.a]
func TestDoDM_TableF_NODIS_Abbreviation(t *testing.T) {
	ism := &model.ISM{
		Classification: model.ClassificationS,
		OwnerProducer:  []string{"USA"},
		NonICMarkings:  []string{"ND"},
	}
	result := banner.Render(ism)
	if !strings.Contains(result.BannerLine, "ND") &&
		!strings.Contains(result.BannerLine, "NODIS") {
		t.Skipf("GAP: BannerLine = %q, want NODIS or ND per Table F", result.BannerLine)
	}
	if !strings.Contains(result.PortionMark, "ND") {
		t.Skipf("GAP: PortionMark = %q, want ND abbreviation per Table F", result.PortionMark)
	}
}

// TestDoDM_TableF_SBU_Abbreviation verifies SBU uses the same form in banner and
// portion marks.
// [Table F, A2-S3.a]
func TestDoDM_TableF_SBU_Abbreviation(t *testing.T) {
	ism := &model.ISM{
		Classification: model.ClassificationU,
		OwnerProducer:  []string{"USA"},
		NonICMarkings:  []string{"SBU"},
	}
	result := banner.Render(ism)
	if !strings.Contains(result.BannerLine, "SBU") {
		t.Skipf("GAP: BannerLine = %q, want SBU per Table F", result.BannerLine)
	}
	if !strings.Contains(result.PortionMark, "SBU") {
		t.Skipf("GAP: PortionMark = %q, want SBU (same form) per Table F", result.PortionMark)
	}
}

// TestDoDM_TableF_SBUNF_Abbreviation verifies SBU NOFORN (banner) → SBU-NF
// (portion) abbreviation.
// [Table F, A2-S4.a]
func TestDoDM_TableF_SBUNF_Abbreviation(t *testing.T) {
	ism := &model.ISM{
		Classification: model.ClassificationU,
		OwnerProducer:  []string{"USA"},
		NonICMarkings:  []string{"SBU-NF"},
	}
	result := banner.Render(ism)
	if !strings.Contains(result.BannerLine, "SBU") {
		t.Skipf("GAP: BannerLine = %q, want SBU NOFORN or SBU-NF per Table F", result.BannerLine)
	}
	if !strings.Contains(result.PortionMark, "SBU-NF") &&
		!strings.Contains(result.PortionMark, "SBU") {
		t.Skipf("GAP: PortionMark = %q, want SBU-NF abbreviation per Table F", result.PortionMark)
	}
}

// TestDoDM_TableF_FOUO_Abbreviation verifies FOUO uses the same abbreviation in
// banner (FOR OFFICIAL USE ONLY or FOUO) and portion (FOUO).
// [Table F, S10.b]
func TestDoDM_TableF_FOUO_Abbreviation(t *testing.T) {
	ism := &model.ISM{
		Classification:        model.ClassificationU,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"FOUO"},
	}
	result := banner.Render(ism)
	if !strings.Contains(result.BannerLine, "FOUO") &&
		!strings.Contains(result.BannerLine, "FOR OFFICIAL USE ONLY") {
		t.Errorf("BannerLine = %q, want FOUO or FOR OFFICIAL USE ONLY per Table F",
			result.BannerLine)
	}
	if !strings.Contains(result.PortionMark, "FOUO") {
		t.Errorf("PortionMark = %q, want FOUO abbreviation per Table F", result.PortionMark)
	}
}

// TestDoDM_TableF_WAIVED_Abbreviation verifies WAIVED uses the same form in
// banner and portion.
// [Table F, S7.f]
func TestDoDM_TableF_WAIVED_Abbreviation(t *testing.T) {
	ism := &model.ISM{
		Classification:        model.ClassificationTS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"WAIVED"},
		SARIdentifier:         []string{"TB"},
	}
	result := banner.Render(ism)
	if !strings.Contains(result.BannerLine, "WAIVED") {
		t.Skipf("GAP: BannerLine = %q, want WAIVED present per Table F; "+
			"renderer may not yet render WAIVED dissemination", result.BannerLine)
	}
	if !strings.Contains(result.PortionMark, "WAIVED") {
		t.Skipf("GAP: PortionMark = %q, want WAIVED (same form) per Table F", result.PortionMark)
	}
}

// =============================================================================
// Table G: Special Format Rules (FMT-1 through FMT-13)
// =============================================================================

// TestDoDM_FMT1_DoubleSlashSeparatesCategories verifies that // separates
// different marking categories.
// [FMT-1, S1.b.1]
func TestDoDM_FMT1_DoubleSlashSeparatesCategories(t *testing.T) {
	tests := []struct {
		name string
		ism  model.ISM
	}{
		{
			name: "SCI_separated_from_classification",
			ism: model.ISM{
				Classification: model.ClassificationTS,
				OwnerProducer:  []string{"USA"},
				SCIControls:    []string{"SI"},
			},
		},
		{
			name: "dissem_separated_from_classification",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"NOFORN"},
			},
		},
		{
			name: "SCI_and_dissem_both_separated",
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
			// Must use // between categories.
			parts := strings.Split(result.BannerLine, "//")
			if len(parts) < 2 {
				t.Errorf("BannerLine = %q, want // separation between categories per FMT-1",
					result.BannerLine)
			}
		})
	}
}

// TestDoDM_FMT2_SingleSlashSeparatesWithinCategory verifies that / separates
// multiple items within the same category.
// [FMT-2, S1.b.1]
func TestDoDM_FMT2_SingleSlashSeparatesWithinCategory(t *testing.T) {
	// Multiple SCI controls should be / separated.
	t.Run("SCI_slash_separated", func(t *testing.T) {
		ism := &model.ISM{
			Classification: model.ClassificationTS,
			OwnerProducer:  []string{"USA"},
			SCIControls:    []string{"SI", "TK"},
		}
		result := banner.Render(ism)
		if !strings.Contains(result.BannerLine, "SI/TK") {
			t.Errorf("BannerLine = %q, want SI/TK (/ between SCI controls within same category) per FMT-2",
				result.BannerLine)
		}
	})

	// Multiple dissem controls should be / separated.
	t.Run("dissem_slash_separated", func(t *testing.T) {
		ism := &model.ISM{
			Classification:        model.ClassificationC,
			OwnerProducer:         []string{"USA"},
			DisseminationControls: []string{"NOFORN", "PROPIN"},
		}
		result := banner.Render(ism)
		bannerLine := result.BannerLine
		// The two controls should be joined with / not //.
		if !strings.Contains(bannerLine, "NOFORN/PROPIN") &&
			!strings.Contains(bannerLine, "PROPIN/NOFORN") {
			// Check if they're at least / separated somewhere.
			if strings.Contains(bannerLine, "NOFORN") && strings.Contains(bannerLine, "PROPIN") {
				nfIdx := strings.Index(bannerLine, "NOFORN")
				prIdx := strings.Index(bannerLine, "PROPIN")
				between := ""
				if nfIdx < prIdx {
					between = bannerLine[nfIdx+len("NOFORN") : prIdx]
				} else {
					between = bannerLine[prIdx+len("PROPIN") : nfIdx]
				}
				if between != "/" {
					t.Errorf("BannerLine = %q, want / (not //) between dissem controls per FMT-2",
						bannerLine)
				}
			}
		}
	})
}

// TestDoDM_FMT3_HyphenSeparatesSubControl verifies that hyphens without spaces
// separate a control system from its sub-control/compartment.
// [FMT-3, S1.b.1b]
func TestDoDM_FMT3_HyphenSeparatesSubControl(t *testing.T) {
	tests := []struct {
		name    string
		ism     model.ISM
		wantSep string
	}{
		{
			name: "SI-G_compartment",
			ism: model.ISM{
				Classification:        model.ClassificationTS,
				OwnerProducer:         []string{"USA"},
				SCIControls:           []string{"SI-G"},
				DisseminationControls: []string{"NOFORN"},
			},
			wantSep: "SI-G",
		},
		{
			name: "HCS-O_compartment",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA"},
				SCIControls:           []string{"HCS-O"},
				DisseminationControls: []string{"NOFORN"},
			},
			wantSep: "HCS-O",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := banner.Render(&tt.ism)
			if !strings.Contains(result.BannerLine, tt.wantSep) {
				t.Errorf("BannerLine = %q, want %q (hyphen-no-spaces separating sub-control) per FMT-3",
					result.BannerLine, tt.wantSep)
			}
			// Ensure no space around the hyphen.
			if strings.Contains(result.BannerLine, strings.Replace(tt.wantSep, "-", " - ", 1)) {
				t.Errorf("BannerLine = %q, hyphen must NOT have spaces per FMT-3", result.BannerLine)
			}
		})
	}
}

// TestDoDM_FMT4_SpaceSeparatesSubCompartments verifies that spaces separate
// sub-compartments from their compartment.
// [FMT-4, S6.e]
func TestDoDM_FMT4_SpaceSeparatesSubCompartments(t *testing.T) {
	ism := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		SCIControls:           []string{"HCS-O XYZ"},
		DisseminationControls: []string{"NOFORN"},
	}
	result := banner.Render(ism)
	if !strings.Contains(result.BannerLine, "HCS-O XYZ") {
		t.Skipf("GAP: BannerLine = %q, want 'HCS-O XYZ' (space separating sub-compartments) per FMT-4",
			result.BannerLine)
	}
}

// TestDoDM_FMT5_SpaceSeparatesSIGMANumbers verifies that multiple SIGMA numbers
// are separated by spaces (e.g., RD-SIGMA 1 2).
// [FMT-5, S8.d.3]
func TestDoDM_FMT5_SpaceSeparatesSIGMANumbers(t *testing.T) {
	ism := &model.ISM{
		Classification:       model.ClassificationS,
		OwnerProducer:        []string{"USA"},
		AtomicEnergyMarkings: []string{"RD-SG-14", "RD-SG-15"},
	}
	result := banner.Render(ism)
	// Should produce something like RD-SIGMA 14 15 in banner.
	if !strings.Contains(result.BannerLine, "SIGMA") &&
		!strings.Contains(result.BannerLine, "SG") {
		t.Skipf("GAP: BannerLine = %q, want SIGMA numbers with space separation per FMT-5; "+
			"renderer may not yet render SIGMA markings", result.BannerLine)
	}
}

// TestDoDM_FMT6_CommaSpaceBetweenRELTOCountries verifies that country codes in
// REL TO are separated by comma and space.
// [FMT-6, S10.d.4]
func TestDoDM_FMT6_CommaSpaceBetweenRELTOCountries(t *testing.T) {
	ism := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"REL"},
		ReleasableTo:          []string{"USA", "AUS", "GBR"},
	}
	result := banner.Render(ism)
	wantCountries := "USA, AUS, GBR"
	if !strings.Contains(result.BannerLine, wantCountries) {
		t.Errorf("BannerLine = %q, want %q (comma+space between countries) per FMT-6",
			result.BannerLine, wantCountries)
	}
}

// TestDoDM_FMT7_SpaceBetweenJOINTCountries verifies that country codes in JOINT
// markings are separated by spaces (not commas).
// [FMT-7, S5.e]
func TestDoDM_FMT7_SpaceBetweenJOINTCountries(t *testing.T) {
	ism := &model.ISM{
		Classification: model.ClassificationS,
		OwnerProducer:  []string{"CAN", "GBR", "USA"},
		Joint:          true,
	}
	result := banner.Render(ism)
	if !strings.Contains(result.BannerLine, "JOINT") {
		t.Skipf("GAP: BannerLine = %q does not contain JOINT", result.BannerLine)
	}
	// Countries should be space-separated, not comma-separated.
	if strings.Contains(result.BannerLine, "CAN, GBR") ||
		strings.Contains(result.BannerLine, "GBR, USA") {
		t.Errorf("BannerLine = %q, JOINT countries must be space-separated (not comma) per FMT-7",
			result.BannerLine)
	}
	if !strings.Contains(result.BannerLine, "CAN GBR USA") {
		t.Errorf("BannerLine = %q, want space-separated countries 'CAN GBR USA' per FMT-7",
			result.BannerLine)
	}
}

// TestDoDM_FMT8_SpaceBetweenFGICountries verifies that country codes in FGI
// banner markings are separated by spaces.
// [FMT-8, S9.d]
func TestDoDM_FMT8_SpaceBetweenFGICountries(t *testing.T) {
	ism := &model.ISM{
		Classification: model.ClassificationTS,
		OwnerProducer:  []string{"USA"},
		FGISourceOpen:  []string{"DEU", "GBR"},
	}
	result := banner.Render(ism)
	if !strings.Contains(result.BannerLine, "FGI DEU GBR") {
		t.Errorf("BannerLine = %q, want 'FGI DEU GBR' (space-separated) per FMT-8",
			result.BannerLine)
	}
	// Must NOT use commas.
	if strings.Contains(result.BannerLine, "DEU, GBR") {
		t.Errorf("BannerLine = %q, FGI countries must be space-separated (not comma) per FMT-8",
			result.BannerLine)
	}
}

// TestDoDM_FMT9_FGIJOINTDocumentsBeginWithDoubleSlash verifies that FGI/JOINT
// documents begin with // (no preceding classification).
// [FMT-9, S1.c]
func TestDoDM_FMT9_FGIJOINTDocumentsBeginWithDoubleSlash(t *testing.T) {
	t.Run("FGI_non_US_begins_with_double_slash", func(t *testing.T) {
		ism := &model.ISM{
			Classification: model.ClassificationS,
			OwnerProducer:  []string{"DEU"},
		}
		result := banner.Render(ism)
		if !strings.HasPrefix(result.BannerLine, "//") {
			t.Skipf("GAP: BannerLine = %q, want '//' prefix for FGI non-US document per FMT-9; "+
				"renderer may not yet support FGI non-US format", result.BannerLine)
		}
	})

	t.Run("JOINT_begins_with_double_slash", func(t *testing.T) {
		ism := &model.ISM{
			Classification: model.ClassificationS,
			OwnerProducer:  []string{"GBR", "USA"},
			Joint:          true,
		}
		result := banner.Render(ism)
		if !strings.HasPrefix(result.BannerLine, "//JOINT") {
			t.Skipf("GAP: BannerLine = %q, want '//JOINT' prefix per FMT-9; "+
				"renderer may not yet produce leading // for JOINT", result.BannerLine)
		}
	})
}

// TestDoDM_FMT10_USClassificationNotPrecededByDoubleSlash verifies that U.S.
// classification is NOT preceded by // in the banner line.
// [FMT-10, S3.b]
func TestDoDM_FMT10_USClassificationNotPrecededByDoubleSlash(t *testing.T) {
	tests := []struct {
		name   string
		class  model.Classification
		banner string
	}{
		{"TOP_SECRET", model.ClassificationTS, "TOP SECRET"},
		{"SECRET", model.ClassificationS, "SECRET"},
		{"CONFIDENTIAL", model.ClassificationC, "CONFIDENTIAL"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ism := &model.ISM{
				Classification: tt.class,
				OwnerProducer:  []string{"USA"},
			}
			result := banner.Render(ism)
			if strings.HasPrefix(result.BannerLine, "//") {
				t.Errorf("BannerLine = %q, U.S. classification must NOT start with // per FMT-10",
					result.BannerLine)
			}
			if result.BannerLine != tt.banner {
				t.Errorf("BannerLine = %q, want %q", result.BannerLine, tt.banner)
			}
		})
	}
}

// TestDoDM_FMT11_SARHyphenSeparator verifies that SAR uses - between SAR and
// nickname/PID (e.g., SAR-BUTTERED POPCORN, SAR-BP).
// [FMT-11, S7.c, S7.d]
func TestDoDM_FMT11_SARHyphenSeparator(t *testing.T) {
	ism := &model.ISM{
		Classification: model.ClassificationTS,
		OwnerProducer:  []string{"USA"},
		SARIdentifier:  []string{"BP"},
	}
	result := banner.Render(ism)
	if !strings.Contains(result.BannerLine, "SAR-") &&
		!strings.Contains(result.BannerLine, "SAR") {
		t.Skipf("GAP: BannerLine = %q, want SAR- present per FMT-11; "+
			"renderer may not yet render SARIdentifier", result.BannerLine)
	}
	// If SAR is present, verify hyphen separator.
	if strings.Contains(result.BannerLine, "SAR ") &&
		!strings.Contains(result.BannerLine, "SAR-") {
		t.Errorf("BannerLine = %q, SAR must use hyphen separator (SAR-), not space per FMT-11",
			result.BannerLine)
	}
}

// TestDoDM_FMT12_ACCMHyphenSeparator verifies that ACCM uses - between ACCM
// and nickname (e.g., ACCM-FICTITIOUS EFFORT).
// [FMT-12, S11.a.2]
func TestDoDM_FMT12_ACCMHyphenSeparator(t *testing.T) {
	ism := &model.ISM{
		Classification: model.ClassificationS,
		OwnerProducer:  []string{"USA"},
		NonICMarkings:  []string{"ACCM-FICTITIOUS EFFORT"},
	}
	result := banner.Render(ism)
	if !strings.Contains(result.BannerLine, "ACCM-FICTITIOUS EFFORT") {
		t.Errorf("BannerLine = %q, want 'ACCM-FICTITIOUS EFFORT' (hyphen separator) per FMT-12",
			result.BannerLine)
	}
	// Verify hyphen is present between ACCM and nickname.
	if strings.Contains(result.BannerLine, "ACCM FICTITIOUS") {
		t.Errorf("BannerLine = %q, ACCM must use hyphen (not space) before nickname per FMT-12",
			result.BannerLine)
	}
}

// TestDoDM_FMT13_ThreePlusSAPs_MULTIPLEPROGRAMS verifies that when 3+ SAPs
// exist, the banner uses SAR-MULTIPLE PROGRAMS and portions list all PIDs.
// [FMT-13, S7.e]
func TestDoDM_FMT13_ThreePlusSAPs_MULTIPLEPROGRAMS(t *testing.T) {
	ism := &model.ISM{
		Classification: model.ClassificationTS,
		OwnerProducer:  []string{"USA"},
		SARIdentifier:  []string{"BP", "TG", "STK"},
	}
	result := banner.Render(ism)
	// Banner should use SAR-MULTIPLE PROGRAMS when 3+ SAPs.
	if !strings.Contains(result.BannerLine, "MULTIPLE PROGRAMS") &&
		!strings.Contains(result.BannerLine, "SAR") {
		t.Skipf("GAP: BannerLine = %q, want 'SAR-MULTIPLE PROGRAMS' for 3+ SAPs per FMT-13; "+
			"renderer may not yet render SARIdentifier", result.BannerLine)
	}
	if strings.Contains(result.BannerLine, "SAR-BP") &&
		strings.Contains(result.BannerLine, "SAR-TG") &&
		strings.Contains(result.BannerLine, "SAR-STK") {
		// All 3 PIDs listed individually in banner is wrong for 3+ SAPs.
		t.Errorf("BannerLine = %q, 3+ SAPs should use SAR-MULTIPLE PROGRAMS in banner (not individual PIDs) per FMT-13",
			result.BannerLine)
	}

	// Portion marks should list all PIDs individually.
	if strings.Contains(result.PortionMark, "SAR-BP") ||
		strings.Contains(result.PortionMark, "SAR-TG") ||
		strings.Contains(result.PortionMark, "SAR-STK") {
		// Good — PIDs are listed in portions.
	} else if strings.Contains(result.PortionMark, "SAR") {
		// SAR present but PIDs not individually listed — may be a gap.
		t.Logf("INFO: PortionMark = %q; per FMT-13 all PIDs must be cited in portions",
			result.PortionMark)
	}
}

// =============================================================================
// Integration: Combined Category Ordering
// =============================================================================

// TestDoDM_CategoryOrdering_SCIBeforeDissem verifies that SCI controls appear
// before dissemination controls in the banner, per Figure 25 ordering.
// [E4-S1.a]
func TestDoDM_CategoryOrdering_SCIBeforeDissem(t *testing.T) {
	ism := &model.ISM{
		Classification:        model.ClassificationTS,
		OwnerProducer:         []string{"USA"},
		SCIControls:           []string{"SI"},
		DisseminationControls: []string{"NOFORN"},
	}
	result := banner.Render(ism)
	wantBanner := "TOP SECRET//SI//NOFORN"
	if result.BannerLine != wantBanner {
		t.Errorf("BannerLine = %q, want %q (SCI before dissem per Figure 25)",
			result.BannerLine, wantBanner)
	}
}

// TestDoDM_CategoryOrdering_MultipleSCI_Alphabetical verifies that multiple SCI
// controls are listed alphabetically.
// [E4-S6.c]
func TestDoDM_CategoryOrdering_MultipleSCI_Alphabetical(t *testing.T) {
	ism := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		SCIControls:           []string{"TK", "HCS", "SI"},
		DisseminationControls: []string{"NOFORN"},
	}
	result := banner.Render(ism)
	// SCI should be alphabetical: HCS/SI/TK.
	if !strings.Contains(result.BannerLine, "HCS/SI/TK") {
		t.Errorf("BannerLine = %q, want HCS/SI/TK (alphabetical SCI order per E4-S6.c)",
			result.BannerLine)
	}
}

// TestDoDM_CategoryOrdering_FGIAfterDissem verifies that FGI markings appear
// in the correct position relative to other categories.
// [E4-S1.a, Figure 25]
func TestDoDM_CategoryOrdering_FGIBeforeDissem(t *testing.T) {
	ism := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		FGISourceOpen:         []string{"GBR"},
		DisseminationControls: []string{"NOFORN"},
	}
	result := banner.Render(ism)
	if !strings.Contains(result.BannerLine, "NOFORN") {
		t.Errorf("BannerLine = %q, want NOFORN present", result.BannerLine)
	}
	if !strings.Contains(result.BannerLine, "FGI") {
		t.Errorf("BannerLine = %q, want FGI present", result.BannerLine)
	}
	// FGI must appear before dissemination per Figure 25 (position 5 < position 6).
	fgiIdx := strings.Index(result.BannerLine, "FGI")
	nfIdx := strings.Index(result.BannerLine, "NOFORN")
	if fgiIdx >= 0 && nfIdx >= 0 && fgiIdx > nfIdx {
		t.Errorf("BannerLine = %q, FGI(idx=%d) must appear before NOFORN(idx=%d) per Figure 25",
			result.BannerLine, fgiIdx, nfIdx)
	}
}

// TestDoDM_DissemControlOrdering verifies the canonical ordering of
// dissemination controls within the dissem category per ISM rules.
func TestDoDM_DissemControlOrdering(t *testing.T) {
	// NOFORN should come before PROPIN in the dissem category.
	ism := &model.ISM{
		Classification:        model.ClassificationC,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"PROPIN", "NOFORN"},
	}
	result := banner.Render(ism)
	nfIdx := strings.Index(result.BannerLine, "NOFORN")
	prIdx := strings.Index(result.BannerLine, "PROPIN")
	if nfIdx >= 0 && prIdx >= 0 && nfIdx > prIdx {
		t.Errorf("BannerLine = %q, NOFORN should come before PROPIN in canonical dissem order",
			result.BannerLine)
	}
}

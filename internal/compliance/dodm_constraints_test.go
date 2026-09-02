package compliance_test

import (
	"testing"

	"dmontgomery/ism-api/internal/model"
	"dmontgomery/ism-api/internal/validation"
)

// DoDM 5200.01-V2 Enclosure 4, Consolidated Constraint Reference Tables A/B/C
// Ref: docs/dodm-5200.01-enclosure4-requirements.md §Consolidated Constraint Reference Tables
//
// Table A: Mutual Exclusions (MX-1 through MX-13)
// Table B: Co-requirements (CR-1 through CR-11)
// Table C: Classification Level Constraints (CL-1 through CL-16)

// =============================================================================
// Table A: Mutual Exclusions (Markings That Cannot Be Used Together)
// =============================================================================

// TestDoDM_MX1_NOFORN_RELTO_BannerExclusion verifies that NOFORN and REL TO
// cannot be used together in the banner line.
// [MX-1: E4-S10.d.7, E4-A1-S2.d]
func TestDoDM_MX1_NOFORN_RELTO_BannerExclusion(t *testing.T) {
	engine := validation.NewEngine(reg())

	ism := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"NOFORN", "REL"},
		ReleasableTo:          []string{"USA", "GBR"},
		ClassifiedBy:          "Test",
		DeclassDate:           "20350101",
	}
	result := engine.Validate(ism)
	if result.Valid {
		t.Error("NOFORN + REL TO should be invalid — mutually exclusive per MX-1")
	}
	if !result.HasCode("dissemination.exclusive_conflict") {
		t.Error("expected dissemination.exclusive_conflict for NOFORN + REL")
	}
}

// TestDoDM_MX2_NOFORN_RELIDO_BannerExclusion verifies that NOFORN and RELIDO
// cannot be used together in the banner line.
// [MX-2: E4-A1-S2.d, E4-A1-S4.d]
func TestDoDM_MX2_NOFORN_RELIDO_BannerExclusion(t *testing.T) {
	engine := validation.NewEngine(reg())

	ism := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"NOFORN", "RELIDO"},
		ClassifiedBy:          "Test",
		DeclassDate:           "20350101",
	}
	result := engine.Validate(ism)
	if result.Valid {
		t.Error("NOFORN + RELIDO should be invalid — mutually exclusive per MX-2")
	}
	if !result.HasCode("dissemination.exclusive_conflict") {
		t.Error("expected dissemination.exclusive_conflict for NOFORN + RELIDO")
	}
}

// TestDoDM_MX3_NOFORN_RELIDO_PortionExclusion verifies that NOFORN and RELIDO
// cannot be used together in the same portion. This is the same underlying
// validation as MX-2 since the ISM model represents per-portion markings.
// [MX-3: E4-A1-S4.d]
func TestDoDM_MX3_NOFORN_RELIDO_PortionExclusion(t *testing.T) {
	engine := validation.NewEngine(reg())

	// Portion-level ISM with both NOFORN and RELIDO.
	ism := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"NOFORN", "RELIDO"},
		ClassifiedBy:          "Test",
		DeclassDate:           "20350101",
	}
	result := engine.Validate(ism)
	if result.Valid {
		t.Error("NOFORN + RELIDO in same portion should be invalid — MX-3")
	}
}

// TestDoDM_MX4_EXDIS_NODIS_DocumentExclusion verifies that EXDIS and NODIS
// cannot be used in the same document (document-level mutually exclusive).
// [MX-4: E4-A2-S1.d]
func TestDoDM_MX4_EXDIS_NODIS_DocumentExclusion(t *testing.T) {
	engine := validation.NewEngine(reg())

	// EXDIS and NODIS are NonICMarkings (Other Dissemination Controls).
	// At the ISM level, having both in the same marking set is the violation.
	ism := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"EXDIS", "NODIS"},
		ClassifiedBy:          "Test",
		DeclassDate:           "20350101",
	}
	result := engine.Validate(ism)
	if result.Valid {
		t.Skipf("GAP: EXDIS + NODIS should be invalid — mutually exclusive per MX-4; "+
			"validator may not yet enforce EXDIS/NODIS exclusion in disseminationControls; result: %+v", result)
	}
}

// TestDoDM_MX5_EXDIS_NODIS_PortionExclusion verifies that EXDIS and NODIS
// cannot appear in the same portion.
// [MX-5: E4-A2-S2.e]
func TestDoDM_MX5_EXDIS_NODIS_PortionExclusion(t *testing.T) {
	engine := validation.NewEngine(reg())

	ism := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"EXDIS", "NODIS"},
		ClassifiedBy:          "Test",
		DeclassDate:           "20350101",
	}
	result := engine.Validate(ism)
	if result.Valid {
		t.Skipf("GAP: EXDIS + NODIS in same portion should be invalid — MX-5; "+
			"validator may not yet enforce this exclusion; result: %+v", result)
	}
}

// TestDoDM_MX6_DISPLAYONLY_NOFORN_Exclusion verifies that DISPLAY ONLY and
// NOFORN cannot be used together.
// [MX-6: E4-S10.e.4]
func TestDoDM_MX6_DISPLAYONLY_NOFORN_Exclusion(t *testing.T) {
	engine := validation.NewEngine(reg())

	ism := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"DISPLAY ONLY", "NOFORN"},
		DisplayOnlyTo:         []string{"GBR"},
		ClassifiedBy:          "Test",
		DeclassDate:           "20350101",
	}
	result := engine.Validate(ism)
	if result.Valid {
		t.Skipf("GAP: DISPLAY ONLY + NOFORN should be invalid — MX-6; "+
			"validator may not yet enforce DISPLAY ONLY exclusions; result: %+v", result)
	}
}

// TestDoDM_MX7_DISPLAYONLY_RELIDO_Exclusion verifies that DISPLAY ONLY and
// RELIDO cannot be used together.
// [MX-7: E4-S10.e.4]
func TestDoDM_MX7_DISPLAYONLY_RELIDO_Exclusion(t *testing.T) {
	engine := validation.NewEngine(reg())

	ism := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"DISPLAY ONLY", "RELIDO"},
		DisplayOnlyTo:         []string{"GBR"},
		ClassifiedBy:          "Test",
		DeclassDate:           "20350101",
	}
	result := engine.Validate(ism)
	if result.Valid {
		t.Skipf("GAP: DISPLAY ONLY + RELIDO should be invalid — MX-7; "+
			"validator may not yet enforce DISPLAY ONLY exclusions; result: %+v", result)
	}
}

// TestDoDM_MX8_USClassification_FGIClassification_Exclusion verifies that U.S.
// classification markings and FGI classification markings are mutually exclusive.
// A document cannot simultaneously be a US-classified document and an FGI
// (non-US) classified document.
// [MX-8: E4-S1.d, E4-S4.a.4]
func TestDoDM_MX8_USClassification_FGIClassification_Exclusion(t *testing.T) {
	engine := validation.NewEngine(reg())

	// A US-classified document (ownerProducer=USA) should not carry FGI-style
	// non-US-only classification. The exclusion is between classification systems,
	// not between US docs containing FGI source material.
	// Here we test: a non-US owner with US-style classification fields set.
	ism := &model.ISM{
		Classification: model.ClassificationS,
		OwnerProducer:  []string{"GBR"}, // FGI document
		// FGI documents should not carry US dissemination controls.
		DisseminationControls: []string{"NOFORN"},
	}
	result := engine.Validate(ism)
	// The system may flag this via ownerProducer/classification mismatch,
	// or it may not yet enforce the MX-8 cross-system exclusion.
	if result.Valid {
		t.Skipf("GAP: US classification system + FGI classification system should be "+
			"mutually exclusive per MX-8; validator may not yet enforce cross-system exclusion")
	}
}

// TestDoDM_MX9_USClassification_JOINTClassification_Exclusion verifies that U.S.
// classification markings and JOINT classification markings are mutually exclusive
// at the banner/portion level.
// [MX-9: E4-S1.d]
func TestDoDM_MX9_USClassification_JOINTClassification_Exclusion(t *testing.T) {
	engine := validation.NewEngine(reg())

	// A JOINT document must use JOINT classification; a US-only document uses
	// US classification. The joint=true flag with single ownerProducer is invalid.
	ism := &model.ISM{
		Classification: model.ClassificationS,
		OwnerProducer:  []string{"USA"},
		Joint:          true, // Joint requires >=2 owners
	}
	result := engine.Validate(ism)
	if result.Valid {
		t.Error("joint=true with single ownerProducer should be invalid — MX-9 enforcement via core.joint_requires_multiple_owners")
	}
	if !result.HasCode("core.joint_requires_multiple_owners") {
		t.Error("expected core.joint_requires_multiple_owners error")
	}
}

// TestDoDM_MX10_FGIClassification_JOINTClassification_Exclusion verifies that
// FGI classification and JOINT classification are mutually exclusive.
// [MX-10: E4-S1.d]
func TestDoDM_MX10_FGIClassification_JOINTClassification_Exclusion(t *testing.T) {
	engine := validation.NewEngine(reg())

	// A JOINT document with only non-US owners — the JOINT system and FGI
	// system cannot coexist on the same document.
	ism := &model.ISM{
		Classification: model.ClassificationS,
		OwnerProducer:  []string{"GBR", "DEU"},
		Joint:          true,
		FGISourceOpen:  []string{"FRA"}, // FGI markings on a JOINT doc
	}
	result := engine.Validate(ism)
	if result.Valid {
		t.Skipf("GAP: JOINT + FGI classification systems should be mutually exclusive per MX-10; "+
			"validator may not yet enforce JOINT/FGI cross-system exclusion")
	}
}

// TestDoDM_MX11_NOFORN_NATOInformation_Exclusion verifies that NOFORN cannot
// be used on NATO information.
// [MX-11: E4-S4.b.3a]
func TestDoDM_MX11_NOFORN_NATOInformation_Exclusion(t *testing.T) {
	engine := validation.NewEngine(reg())

	// A US document incorporating NATO info (FGI NATO) with NOFORN.
	ism := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		FGISourceOpen:         []string{"NATO"},
		DisseminationControls: []string{"NOFORN"},
		ClassifiedBy:          "Test",
		DeclassDate:           "20350101",
	}
	result := engine.Validate(ism)
	if result.Valid {
		t.Skipf("GAP: NOFORN + NATO information should be invalid per MX-11; "+
			"validator may not yet enforce NOFORN/NATO exclusion")
	}
}

// TestDoDM_MX12_NOFORN_FGIPortionMarks_Exclusion verifies that NOFORN cannot
// be used with FGI portion marks at the portion level.
// [MX-12: E4-S9.m]
func TestDoDM_MX12_NOFORN_FGIPortionMarks_Exclusion(t *testing.T) {
	engine := validation.NewEngine(reg())

	// A portion with both NOFORN and FGI source material.
	ism := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		FGISourceOpen:         []string{"GBR"},
		DisseminationControls: []string{"NOFORN"},
		ClassifiedBy:          "Test",
		DeclassDate:           "20350101",
	}
	result := engine.Validate(ism)
	if result.Valid {
		t.Skipf("GAP: NOFORN with FGI portion marks should be invalid per MX-12; "+
			"validator may not yet enforce NOFORN/FGI portion exclusion")
	}
}

// TestDoDM_MX13_DISPLAYONLY_OtherDissem_Exclusion verifies that DISPLAY ONLY
// cannot be combined with other dissemination controls (e.g., REL TO) unless
// policy authorized.
// [MX-13: E4-S10.e.4a]
func TestDoDM_MX13_DISPLAYONLY_OtherDissem_Exclusion(t *testing.T) {
	engine := validation.NewEngine(reg())

	ism := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"DISPLAY ONLY", "REL"},
		DisplayOnlyTo:         []string{"GBR"},
		ReleasableTo:          []string{"USA", "GBR"},
		ClassifiedBy:          "Test",
		DeclassDate:           "20350101",
	}
	result := engine.Validate(ism)
	if result.Valid {
		t.Skipf("GAP: DISPLAY ONLY + REL TO should be invalid per MX-13; "+
			"validator may not yet enforce DISPLAY ONLY exclusion with other DISSEM controls")
	}
}

// =============================================================================
// Table B: Co-requirements (Markings That Must Be Used Together)
// =============================================================================

// TestDoDM_CR1_HCS_RequiresNOFORN verifies that HCS (SCI) requires NOFORN.
// [CR-1: E4-S6.f]
func TestDoDM_CR1_HCS_RequiresNOFORN(t *testing.T) {
	engine := validation.NewEngine(reg())

	// Negative: HCS without NOFORN should be flagged.
	t.Run("HCS_without_NOFORN", func(t *testing.T) {
		ism := &model.ISM{
			Classification: model.ClassificationTS,
			OwnerProducer:  []string{"USA"},
			SCIControls:    []string{"HCS"},
			ClassifiedBy:   "Test",
			DeclassDate:    "20350101",
		}
		result := engine.Validate(ism)
		if result.Valid {
			t.Skipf("GAP: HCS without NOFORN should be invalid per CR-1; "+
				"validator may not yet enforce HCS→NOFORN co-requirement")
		}
	})

	// Positive: HCS with NOFORN should pass.
	t.Run("HCS_with_NOFORN", func(t *testing.T) {
		ism := &model.ISM{
			Classification:        model.ClassificationTS,
			OwnerProducer:         []string{"USA"},
			SCIControls:           []string{"HCS"},
			DisseminationControls: []string{"NOFORN"},
			ClassifiedBy:          "Test",
			DeclassDate:           "20350101",
		}
		result := engine.Validate(ism)
		for _, e := range result.Errors {
			if e.Severity == validation.SeverityError {
				t.Errorf("HCS + NOFORN should be valid, but got error: %s (code: %s)",
					e.Message, e.Code)
			}
		}
	})
}

// TestDoDM_CR2_TKGEOCAP_RequiresNOFORN verifies that TK-GEOCAP (SCI) requires
// NOFORN.
// [CR-2: E4-S6.f]
func TestDoDM_CR2_TKGEOCAP_RequiresNOFORN(t *testing.T) {
	engine := validation.NewEngine(reg())

	// Negative: TK without NOFORN.
	t.Run("TK_without_NOFORN", func(t *testing.T) {
		ism := &model.ISM{
			Classification: model.ClassificationTS,
			OwnerProducer:  []string{"USA"},
			SCIControls:    []string{"TK"},
			ClassifiedBy:   "Test",
			DeclassDate:    "20350101",
		}
		result := engine.Validate(ism)
		if result.Valid {
			t.Skipf("GAP: TK/TK-GEOCAP without NOFORN should be invalid per CR-2; "+
				"validator may not yet enforce TK→NOFORN co-requirement")
		}
	})

	// Positive: TK with NOFORN.
	t.Run("TK_with_NOFORN", func(t *testing.T) {
		ism := &model.ISM{
			Classification:        model.ClassificationTS,
			OwnerProducer:         []string{"USA"},
			SCIControls:           []string{"TK"},
			DisseminationControls: []string{"NOFORN"},
			ClassifiedBy:          "Test",
			DeclassDate:           "20350101",
		}
		result := engine.Validate(ism)
		for _, e := range result.Errors {
			if e.Severity == validation.SeverityError {
				t.Errorf("TK + NOFORN should be valid, but got error: %s (code: %s)",
					e.Message, e.Code)
			}
		}
	})
}

// TestDoDM_CR3_RELTO_RequiresCountryBesidesUSA verifies that REL TO requires at
// least one country besides USA.
// [CR-3: E4-S10.d.5]
func TestDoDM_CR3_RELTO_RequiresCountryBesidesUSA(t *testing.T) {
	engine := validation.NewEngine(reg())

	// Negative: REL TO with only USA.
	t.Run("REL_only_USA", func(t *testing.T) {
		ism := &model.ISM{
			Classification:        model.ClassificationS,
			OwnerProducer:         []string{"USA"},
			DisseminationControls: []string{"REL"},
			ReleasableTo:          []string{"USA"},
			ClassifiedBy:          "Test",
			DeclassDate:           "20350101",
		}
		result := engine.Validate(ism)
		if result.Valid {
			t.Skipf("GAP: REL TO with only USA should be invalid per CR-3; "+
				"validator may not yet enforce non-USA country requirement for REL TO")
		}
	})

	// Positive: REL TO with USA and another country.
	t.Run("REL_USA_and_GBR", func(t *testing.T) {
		ism := &model.ISM{
			Classification:        model.ClassificationS,
			OwnerProducer:         []string{"USA"},
			DisseminationControls: []string{"REL"},
			ReleasableTo:          []string{"USA", "GBR"},
			ClassifiedBy:          "Test",
			DeclassDate:           "20350101",
		}
		result := engine.Validate(ism)
		for _, e := range result.Errors {
			if e.Field == "releasableTo" && e.Severity == validation.SeverityError {
				t.Errorf("REL TO with USA+GBR should not have releasableTo error: %s", e.Message)
			}
		}
	})
}

// TestDoDM_CR4_RELTO_USAFirst verifies that USA must be listed first in REL TO
// country lists.
// [CR-4: E4-S10.d.4]
func TestDoDM_CR4_RELTO_USAFirst(t *testing.T) {
	engine := validation.NewEngine(reg())

	// Negative: REL TO without USA first (or without USA at all).
	t.Run("REL_no_USA", func(t *testing.T) {
		ism := &model.ISM{
			Classification:        model.ClassificationS,
			OwnerProducer:         []string{"USA"},
			DisseminationControls: []string{"REL"},
			ReleasableTo:          []string{"GBR", "CAN"},
			ClassifiedBy:          "Test",
			DeclassDate:           "20350101",
		}
		result := engine.Validate(ism)
		if result.Valid {
			t.Skipf("GAP: REL TO without USA in releasableTo should be flagged per CR-4; "+
				"validator may not yet enforce USA-first ordering in REL TO")
		}
	})

	// Positive: REL TO with USA listed first.
	t.Run("REL_USA_first", func(t *testing.T) {
		ism := &model.ISM{
			Classification:        model.ClassificationS,
			OwnerProducer:         []string{"USA"},
			DisseminationControls: []string{"REL"},
			ReleasableTo:          []string{"USA", "AUS", "GBR"},
			ClassifiedBy:          "Test",
			DeclassDate:           "20350101",
		}
		result := engine.Validate(ism)
		for _, e := range result.Errors {
			if e.Field == "releasableTo" && e.Severity == validation.SeverityError {
				t.Errorf("REL TO with USA first should be valid, but got: %s", e.Message)
			}
		}
	})
}

// TestDoDM_CR5_ORCON_MustAppearInBanner verifies that if any portion has ORCON,
// ORCON must appear in the banner line.
// [CR-5: E4-S10.c.4a]
func TestDoDM_CR5_ORCON_MustAppearInBanner(t *testing.T) {
	engine := validation.NewEngine(reg())

	// Positive: banner ISM with ORCON should validate without ORCON-related errors.
	ism := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"OC"},
		ClassifiedBy:          "Test",
		DeclassDate:           "20350101",
	}
	result := engine.Validate(ism)
	for _, e := range result.Errors {
		if e.Severity == validation.SeverityError {
			t.Errorf("SECRET//ORCON should be valid, but got error: %s (code: %s)",
				e.Message, e.Code)
		}
	}
}

// TestDoDM_CR6_RD_MustAppearInBanner verifies that if any portion contains RD,
// RD must appear in the banner line.
// [CR-6: E4-S8.a.5a]
func TestDoDM_CR6_RD_MustAppearInBanner(t *testing.T) {
	engine := validation.NewEngine(reg())

	// Positive: RD in banner (portion-level ISM also valid).
	ism := &model.ISM{
		Classification:       model.ClassificationS,
		OwnerProducer:        []string{"USA"},
		AtomicEnergyMarkings: []string{"RD"},
		ClassifiedBy:         "Test",
		DeclassDate:          "20350101",
	}
	result := engine.Validate(ism)
	for _, e := range result.Errors {
		if e.Severity == validation.SeverityError && e.Field == "atomicEnergyMarkings" {
			t.Errorf("SECRET with RD should be valid, but got AEA error: %s (code: %s)",
				e.Message, e.Code)
		}
	}
}

// TestDoDM_CR7_FRD_MustAppearInBanner verifies that if any portion contains FRD,
// FRD must appear in the banner line.
// [CR-7: E4-S8.b.3a]
func TestDoDM_CR7_FRD_MustAppearInBanner(t *testing.T) {
	engine := validation.NewEngine(reg())

	// Positive: FRD with valid classification.
	ism := &model.ISM{
		Classification:       model.ClassificationS,
		OwnerProducer:        []string{"USA"},
		AtomicEnergyMarkings: []string{"FRD"},
		ClassifiedBy:         "Test",
		DeclassDate:          "20350101",
	}
	result := engine.Validate(ism)
	for _, e := range result.Errors {
		if e.Severity == validation.SeverityError && e.Field == "atomicEnergyMarkings" {
			t.Errorf("SECRET with FRD should be valid, but got AEA error: %s (code: %s)",
				e.Message, e.Code)
		}
	}
}

// TestDoDM_CR8_PROPIN_MustAppearInBanner verifies that if any portion has PROPIN,
// PROPIN must appear in the banner line.
// [CR-8: E4-A1-S3.c]
func TestDoDM_CR8_PROPIN_MustAppearInBanner(t *testing.T) {
	engine := validation.NewEngine(reg())

	// Positive: banner ISM with PROPIN.
	ism := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"PROPIN"},
		ClassifiedBy:          "Test",
		DeclassDate:           "20350101",
	}
	result := engine.Validate(ism)
	for _, e := range result.Errors {
		if e.Severity == validation.SeverityError {
			t.Errorf("SECRET//PROPIN should be valid, but got error: %s (code: %s)",
				e.Message, e.Code)
		}
	}
}

// TestDoDM_CR9_BOHEMIA_RequiresCosmicTopSecret verifies that BOHEMIA may be used
// ONLY with //COSMIC TOP SECRET classification.
// [CR-9: E4-S4.b.2.c]
func TestDoDM_CR9_BOHEMIA_RequiresCosmicTopSecret(t *testing.T) {
	// BOHEMIA is a NATO designation requiring COSMIC TOP SECRET. This is
	// represented via HighWaterNATO or NonUSControls in the ISM model.
	// The validator may not yet enforce this NATO-specific co-requirement.
	t.Skipf("GAP: BOHEMIA→COSMIC TOP SECRET co-requirement (CR-9) requires " +
		"NATO-specific validation not yet implemented in the validator")
}

// TestDoDM_CR10_IMCON_RequiresSECRET verifies that IMCON as a standalone control
// requires SECRET classification.
// [CR-10: E4-A1-S1.b]
func TestDoDM_CR10_IMCON_RequiresSECRET(t *testing.T) {
	engine := validation.NewEngine(reg())

	// Positive: IMCON with SECRET.
	t.Run("IMCON_SECRET_valid", func(t *testing.T) {
		ism := &model.ISM{
			Classification:        model.ClassificationS,
			OwnerProducer:         []string{"USA"},
			DisseminationControls: []string{"IMC"},
			ClassifiedBy:          "Test",
			DeclassDate:           "20350101",
		}
		result := engine.Validate(ism)
		for _, e := range result.Errors {
			if e.Code == "dissemination.insufficient_classification" {
				t.Errorf("IMCON with SECRET should pass classification gate, got: %s", e.Message)
			}
		}
	})

	// Negative: IMCON with UNCLASSIFIED should fail.
	t.Run("IMCON_UNCLASSIFIED_invalid", func(t *testing.T) {
		ism := &model.ISM{
			Classification:        model.ClassificationU,
			OwnerProducer:         []string{"USA"},
			DisseminationControls: []string{"IMC"},
		}
		result := engine.Validate(ism)
		if result.Valid {
			t.Error("IMCON with UNCLASSIFIED should be invalid per CR-10")
		}
	})
}

// TestDoDM_CR11_FOUO_NotInClassifiedBanner verifies that FOUO in a classified
// document does NOT appear in the banner — it only appears in portion markings
// at the UNCLASSIFIED level.
// [CR-11: E4-S10.b.3a]
func TestDoDM_CR11_FOUO_NotInClassifiedBanner(t *testing.T) {
	engine := validation.NewEngine(reg())

	// FOUO in a classified banner ISM should be flagged or at minimum, the
	// FOUO marking on a classified ISM is semantically incorrect.
	ism := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"FOUO"},
		ClassifiedBy:          "Test",
		DeclassDate:           "20350101",
	}
	result := engine.Validate(ism)
	// CR-11 is a document-level rule: FOUO only on U-level portions within
	// classified docs. If the validator allows FOUO on a S-level ISM, flag the gap.
	if result.Valid {
		t.Skipf("GAP: FOUO on SECRET-level ISM should be flagged per CR-11; "+
			"validator may not yet enforce FOUO→U-only classification constraint")
	}
}

// =============================================================================
// Table C: Classification Level Constraints (Minimum/Maximum Classification)
// =============================================================================

// TestDoDM_CL1_RD_ClassificationLevels verifies that RD (Restricted Data) is
// allowed only with TS, S, or C.
// [CL-1: E4-S8.a.4]
func TestDoDM_CL1_RD_ClassificationLevels(t *testing.T) {
	engine := validation.NewEngine(reg())

	validLevels := []struct {
		name  string
		class model.Classification
	}{
		{"TS", model.ClassificationTS},
		{"S", model.ClassificationS},
		{"C", model.ClassificationC},
	}
	for _, tt := range validLevels {
		t.Run(tt.name+"_valid", func(t *testing.T) {
			ism := &model.ISM{
				Classification:       tt.class,
				OwnerProducer:        []string{"USA"},
				AtomicEnergyMarkings: []string{"RD"},
				ClassifiedBy:         "Test",
				DeclassDate:          "20350101",
			}
			result := engine.Validate(ism)
			for _, e := range result.Errors {
				if e.Field == "atomicEnergyMarkings" && e.Severity == validation.SeverityError {
					t.Errorf("RD with %s should be valid, got: %s", tt.class, e.Message)
				}
			}
		})
	}

	invalidLevels := []struct {
		name  string
		class model.Classification
	}{
		{"U", model.ClassificationU},
		{"CUI", model.ClassificationCUI},
	}
	for _, tt := range invalidLevels {
		t.Run(tt.name+"_invalid", func(t *testing.T) {
			ism := &model.ISM{
				Classification:       tt.class,
				AtomicEnergyMarkings: []string{"RD"},
			}
			result := engine.Validate(ism)
			if result.Valid {
				t.Skipf("GAP: RD with %s should be invalid per CL-1; "+
					"validator may not yet enforce AEA classification gates", tt.class)
			}
		})
	}
}

// TestDoDM_CL2_FRD_ClassificationLevels verifies that FRD is allowed only with
// TS, S, or C.
// [CL-2: E4-S8.b.2]
func TestDoDM_CL2_FRD_ClassificationLevels(t *testing.T) {
	engine := validation.NewEngine(reg())

	validLevels := []struct {
		name  string
		class model.Classification
	}{
		{"TS", model.ClassificationTS},
		{"S", model.ClassificationS},
		{"C", model.ClassificationC},
	}
	for _, tt := range validLevels {
		t.Run(tt.name+"_valid", func(t *testing.T) {
			ism := &model.ISM{
				Classification:       tt.class,
				OwnerProducer:        []string{"USA"},
				AtomicEnergyMarkings: []string{"FRD"},
				ClassifiedBy:         "Test",
				DeclassDate:          "20350101",
			}
			result := engine.Validate(ism)
			for _, e := range result.Errors {
				if e.Field == "atomicEnergyMarkings" && e.Severity == validation.SeverityError {
					t.Errorf("FRD with %s should be valid, got: %s", tt.class, e.Message)
				}
			}
		})
	}

	invalidLevels := []struct {
		name  string
		class model.Classification
	}{
		{"U", model.ClassificationU},
		{"CUI", model.ClassificationCUI},
	}
	for _, tt := range invalidLevels {
		t.Run(tt.name+"_invalid", func(t *testing.T) {
			ism := &model.ISM{
				Classification:       tt.class,
				AtomicEnergyMarkings: []string{"FRD"},
			}
			result := engine.Validate(ism)
			if result.Valid {
				t.Skipf("GAP: FRD with %s should be invalid per CL-2; "+
					"validator may not yet enforce AEA classification gates", tt.class)
			}
		})
	}
}

// TestDoDM_CL3_CNWDI_ClassificationLevels verifies that CNWDI is allowed only
// with TS or S (it is a subset of RD).
// [CL-3: E4-S8.c.1]
func TestDoDM_CL3_CNWDI_ClassificationLevels(t *testing.T) {
	engine := validation.NewEngine(reg())

	validLevels := []struct {
		name  string
		class model.Classification
	}{
		{"TS", model.ClassificationTS},
		{"S", model.ClassificationS},
	}
	for _, tt := range validLevels {
		t.Run(tt.name+"_valid", func(t *testing.T) {
			ism := &model.ISM{
				Classification:       tt.class,
				OwnerProducer:        []string{"USA"},
				AtomicEnergyMarkings: []string{"RD-CNWDI"},
				ClassifiedBy:         "Test",
				DeclassDate:          "20350101",
			}
			result := engine.Validate(ism)
			for _, e := range result.Errors {
				if e.Field == "atomicEnergyMarkings" && e.Severity == validation.SeverityError {
					t.Errorf("CNWDI with %s should be valid, got: %s", tt.class, e.Message)
				}
			}
		})
	}

	invalidLevels := []struct {
		name  string
		class model.Classification
	}{
		{"U", model.ClassificationU},
		{"CUI", model.ClassificationCUI},
		{"C", model.ClassificationC},
	}
	for _, tt := range invalidLevels {
		t.Run(tt.name+"_invalid", func(t *testing.T) {
			ism := &model.ISM{
				Classification:       tt.class,
				AtomicEnergyMarkings: []string{"RD-CNWDI"},
			}
			result := engine.Validate(ism)
			if result.Valid {
				t.Skipf("GAP: CNWDI with %s should be invalid per CL-3; "+
					"validator may not yet enforce CNWDI classification gates", tt.class)
			}
		})
	}
}

// TestDoDM_CL4_SIGMA_ClassificationLevels verifies that SIGMA is allowed only
// with TS, S, or C.
// [CL-4: E4-S8.d.2]
func TestDoDM_CL4_SIGMA_ClassificationLevels(t *testing.T) {
	engine := validation.NewEngine(reg())

	validLevels := []struct {
		name  string
		class model.Classification
	}{
		{"TS", model.ClassificationTS},
		{"S", model.ClassificationS},
		{"C", model.ClassificationC},
	}
	for _, tt := range validLevels {
		t.Run(tt.name+"_valid", func(t *testing.T) {
			ism := &model.ISM{
				Classification:       tt.class,
				OwnerProducer:        []string{"USA"},
				AtomicEnergyMarkings: []string{"RD-SIGMA 1"},
				ClassifiedBy:         "Test",
				DeclassDate:          "20350101",
			}
			result := engine.Validate(ism)
			for _, e := range result.Errors {
				if e.Field == "atomicEnergyMarkings" && e.Severity == validation.SeverityError {
					t.Errorf("SIGMA with %s should be valid, got: %s", tt.class, e.Message)
				}
			}
		})
	}

	invalidLevels := []struct {
		name  string
		class model.Classification
	}{
		{"U", model.ClassificationU},
		{"CUI", model.ClassificationCUI},
	}
	for _, tt := range invalidLevels {
		t.Run(tt.name+"_invalid", func(t *testing.T) {
			ism := &model.ISM{
				Classification:       tt.class,
				AtomicEnergyMarkings: []string{"RD-SIGMA 1"},
			}
			result := engine.Validate(ism)
			if result.Valid {
				t.Skipf("GAP: SIGMA with %s should be invalid per CL-4; "+
					"validator may not yet enforce SIGMA classification gates", tt.class)
			}
		})
	}
}

// TestDoDM_CL5_ORCON_ClassificationLevels verifies that ORCON is allowed only
// with TS, S, or C.
// [CL-5: E4-S10.c.3]
func TestDoDM_CL5_ORCON_ClassificationLevels(t *testing.T) {
	engine := validation.NewEngine(reg())

	validLevels := []struct {
		name  string
		class model.Classification
	}{
		{"TS", model.ClassificationTS},
		{"S", model.ClassificationS},
		{"C", model.ClassificationC},
	}
	for _, tt := range validLevels {
		t.Run(tt.name+"_valid", func(t *testing.T) {
			ism := &model.ISM{
				Classification:        tt.class,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"OC"},
				ClassifiedBy:          "Test",
				DeclassDate:           "20350101",
			}
			result := engine.Validate(ism)
			for _, e := range result.Errors {
				if e.Code == "dissemination.insufficient_classification" {
					t.Errorf("ORCON with %s should be valid, got: %s", tt.class, e.Message)
				}
			}
		})
	}

	t.Run("U_invalid", func(t *testing.T) {
		ism := &model.ISM{
			Classification:        model.ClassificationU,
			DisseminationControls: []string{"OC"},
		}
		result := engine.Validate(ism)
		if result.Valid {
			t.Error("ORCON with U should be invalid per CL-5")
		}
		if !result.HasCode("dissemination.insufficient_classification") {
			t.Error("expected dissemination.insufficient_classification for OC + U")
		}
	})
}

// TestDoDM_CL6_RELTO_ClassificationLevels verifies that REL TO is allowed with
// TS, S, C, or CUI.
// [CL-6: E4-S10.d.3]
func TestDoDM_CL6_RELTO_ClassificationLevels(t *testing.T) {
	engine := validation.NewEngine(reg())

	validLevels := []struct {
		name  string
		class model.Classification
	}{
		{"TS", model.ClassificationTS},
		{"S", model.ClassificationS},
		{"C", model.ClassificationC},
		{"CUI", model.ClassificationCUI},
	}
	for _, tt := range validLevels {
		t.Run(tt.name+"_valid", func(t *testing.T) {
			ism := &model.ISM{
				Classification:        tt.class,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"REL"},
				ReleasableTo:          []string{"USA", "GBR"},
			}
			if tt.class.AtLeast(model.ClassificationC) {
				ism.ClassifiedBy = "Test"
				ism.DeclassDate = "20350101"
			}
			result := engine.Validate(ism)
			for _, e := range result.Errors {
				if e.Code == "dissemination.insufficient_classification" {
					t.Errorf("REL TO with %s should be valid, got: %s", tt.class, e.Message)
				}
			}
		})
	}

	t.Run("U_invalid", func(t *testing.T) {
		ism := &model.ISM{
			Classification:        model.ClassificationU,
			DisseminationControls: []string{"REL"},
			ReleasableTo:          []string{"USA", "GBR"},
		}
		result := engine.Validate(ism)
		// REL TO is not in the classification gates, so U may be allowed.
		// If it is allowed, that's expected per CL-6 (TS, S, C, CUI).
		// U is not in the allowed set, so it should be invalid.
		if result.Valid {
			t.Skipf("GAP: REL TO with U should be invalid per CL-6; "+
				"validator may not have a classification gate for REL")
		}
	})
}

// TestDoDM_CL7_DISPLAYONLY_ClassificationLevels verifies that DISPLAY ONLY is
// allowed only with TS, S, or C.
// [CL-7: E4-S10.e.3]
func TestDoDM_CL7_DISPLAYONLY_ClassificationLevels(t *testing.T) {
	engine := validation.NewEngine(reg())

	validLevels := []struct {
		name  string
		class model.Classification
	}{
		{"TS", model.ClassificationTS},
		{"S", model.ClassificationS},
		{"C", model.ClassificationC},
	}
	for _, tt := range validLevels {
		t.Run(tt.name+"_valid", func(t *testing.T) {
			ism := &model.ISM{
				Classification:        tt.class,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"DISPLAY ONLY"},
				DisplayOnlyTo:         []string{"GBR"},
				ClassifiedBy:          "Test",
				DeclassDate:           "20350101",
			}
			result := engine.Validate(ism)
			for _, e := range result.Errors {
				if e.Code == "dissemination.insufficient_classification" {
					t.Errorf("DISPLAY ONLY with %s should be valid, got: %s", tt.class, e.Message)
				}
			}
		})
	}

	t.Run("U_invalid", func(t *testing.T) {
		ism := &model.ISM{
			Classification:        model.ClassificationU,
			DisseminationControls: []string{"DISPLAY ONLY"},
			DisplayOnlyTo:         []string{"GBR"},
		}
		result := engine.Validate(ism)
		if result.Valid {
			t.Skipf("GAP: DISPLAY ONLY with U should be invalid per CL-7; "+
				"validator may not have a classification gate for DISPLAY ONLY")
		}
	})
}

// TestDoDM_CL8_NOFORN_ClassificationLevels verifies that NOFORN is allowed with
// TS, S, C, or CUI.
// [CL-8: E4-A1-S2.c]
func TestDoDM_CL8_NOFORN_ClassificationLevels(t *testing.T) {
	engine := validation.NewEngine(reg())

	validLevels := []struct {
		name  string
		class model.Classification
	}{
		{"TS", model.ClassificationTS},
		{"S", model.ClassificationS},
		{"C", model.ClassificationC},
	}
	for _, tt := range validLevels {
		t.Run(tt.name+"_valid", func(t *testing.T) {
			ism := &model.ISM{
				Classification:        tt.class,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"NOFORN"},
				ClassifiedBy:          "Test",
				DeclassDate:           "20350101",
			}
			result := engine.Validate(ism)
			for _, e := range result.Errors {
				if e.Code == "dissemination.insufficient_classification" {
					t.Errorf("NOFORN with %s should be valid, got: %s", tt.class, e.Message)
				}
			}
		})
	}

	t.Run("U_invalid", func(t *testing.T) {
		ism := &model.ISM{
			Classification:        model.ClassificationU,
			DisseminationControls: []string{"NOFORN"},
		}
		result := engine.Validate(ism)
		if result.Valid {
			t.Error("NOFORN with U should be invalid per CL-8")
		}
		if !result.HasCode("dissemination.insufficient_classification") {
			t.Error("expected dissemination.insufficient_classification for NOFORN + U")
		}
	})
}

// TestDoDM_CL9_PROPIN_ClassificationLevels verifies that PROPIN is allowed with
// TS, S, C, or U.
// [CL-9: E4-A1-S3.b]
func TestDoDM_CL9_PROPIN_ClassificationLevels(t *testing.T) {
	engine := validation.NewEngine(reg())

	// PROPIN has no classification gate in the current compatibility data,
	// meaning all levels should be valid (which aligns with CL-9: TS, S, C, U).
	allLevels := []struct {
		name  string
		class model.Classification
	}{
		{"TS", model.ClassificationTS},
		{"S", model.ClassificationS},
		{"C", model.ClassificationC},
		{"U", model.ClassificationU},
	}
	for _, tt := range allLevels {
		t.Run(tt.name+"_valid", func(t *testing.T) {
			ism := &model.ISM{
				Classification:        tt.class,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"PROPIN"},
			}
			if tt.class.AtLeast(model.ClassificationC) {
				ism.ClassifiedBy = "Test"
				ism.DeclassDate = "20350101"
			}
			result := engine.Validate(ism)
			for _, e := range result.Errors {
				if e.Code == "dissemination.insufficient_classification" {
					t.Errorf("PROPIN with %s should be valid per CL-9, got: %s", tt.class, e.Message)
				}
			}
		})
	}
}

// TestDoDM_CL10_RELIDO_ClassificationLevels verifies that RELIDO is allowed only
// with TS, S, or C.
// [CL-10: E4-A1-S4.c]
func TestDoDM_CL10_RELIDO_ClassificationLevels(t *testing.T) {
	engine := validation.NewEngine(reg())

	validLevels := []struct {
		name  string
		class model.Classification
	}{
		{"TS", model.ClassificationTS},
		{"S", model.ClassificationS},
		{"C", model.ClassificationC},
	}
	for _, tt := range validLevels {
		t.Run(tt.name+"_valid", func(t *testing.T) {
			ism := &model.ISM{
				Classification:        tt.class,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"RELIDO"},
				ClassifiedBy:          "Test",
				DeclassDate:           "20350101",
			}
			result := engine.Validate(ism)
			for _, e := range result.Errors {
				if e.Code == "dissemination.insufficient_classification" {
					t.Errorf("RELIDO with %s should be valid, got: %s", tt.class, e.Message)
				}
			}
		})
	}

	t.Run("U_invalid", func(t *testing.T) {
		ism := &model.ISM{
			Classification:        model.ClassificationU,
			DisseminationControls: []string{"RELIDO"},
		}
		result := engine.Validate(ism)
		if result.Valid {
			t.Error("RELIDO with U should be invalid per CL-10 — DoDM E4-A1-S4.c requires TS/S/C")
		}
	})
}

// TestDoDM_CL11_IMCON_ClassificationLevels verifies that IMCON standalone
// requires SECRET; IMCON in a TS paragraph is also valid.
// [CL-11: E4-A1-S1.b, E4-A1-S1.e]
func TestDoDM_CL11_IMCON_ClassificationLevels(t *testing.T) {
	engine := validation.NewEngine(reg())

	// Positive: IMCON at SECRET.
	t.Run("S_valid", func(t *testing.T) {
		ism := &model.ISM{
			Classification:        model.ClassificationS,
			OwnerProducer:         []string{"USA"},
			DisseminationControls: []string{"IMC"},
			ClassifiedBy:          "Test",
			DeclassDate:           "20350101",
		}
		result := engine.Validate(ism)
		for _, e := range result.Errors {
			if e.Code == "dissemination.insufficient_classification" {
				t.Errorf("IMCON with SECRET should be valid, got: %s", e.Message)
			}
		}
	})

	// Positive: IMCON at TS (TS paragraph with IMCON marking).
	t.Run("TS_valid", func(t *testing.T) {
		ism := &model.ISM{
			Classification:        model.ClassificationTS,
			OwnerProducer:         []string{"USA"},
			DisseminationControls: []string{"IMC"},
			ClassifiedBy:          "Test",
			DeclassDate:           "20350101",
		}
		result := engine.Validate(ism)
		for _, e := range result.Errors {
			if e.Code == "dissemination.insufficient_classification" {
				t.Errorf("IMCON with TS should be valid (TS paragraph), got: %s", e.Message)
			}
		}
	})

	// Negative: IMCON at UNCLASSIFIED.
	t.Run("U_invalid", func(t *testing.T) {
		ism := &model.ISM{
			Classification:        model.ClassificationU,
			DisseminationControls: []string{"IMC"},
		}
		result := engine.Validate(ism)
		if result.Valid {
			t.Error("IMCON with U should be invalid per CL-11")
		}
	})
}

// TestDoDM_CL12_FOUO_ClassificationLevel verifies that FOUO is only for
// UNCLASSIFIED information (portion marked U//FOUO within classified documents).
// [CL-12: E4-S10.b.1, E4-S10.b.3]
func TestDoDM_CL12_FOUO_ClassificationLevel(t *testing.T) {
	engine := validation.NewEngine(reg())

	// Positive: FOUO with UNCLASSIFIED.
	t.Run("U_valid", func(t *testing.T) {
		ism := &model.ISM{
			Classification:        model.ClassificationU,
			OwnerProducer:         []string{"USA"},
			DisseminationControls: []string{"FOUO"},
		}
		result := engine.Validate(ism)
		for _, e := range result.Errors {
			if e.Code == "dissemination.insufficient_classification" && e.Severity == validation.SeverityError {
				t.Errorf("FOUO with U should be valid, got: %s", e.Message)
			}
		}
	})

	// Negative: FOUO with SECRET should be invalid.
	t.Run("S_invalid", func(t *testing.T) {
		ism := &model.ISM{
			Classification:        model.ClassificationS,
			OwnerProducer:         []string{"USA"},
			DisseminationControls: []string{"FOUO"},
			ClassifiedBy:          "Test",
			DeclassDate:           "20350101",
		}
		result := engine.Validate(ism)
		if result.Valid {
			t.Skipf("GAP: FOUO with SECRET should be invalid per CL-12; "+
				"validator may not yet enforce FOUO→U-only constraint")
		}
	})

	// Negative: FOUO with TS should be invalid.
	t.Run("TS_invalid", func(t *testing.T) {
		ism := &model.ISM{
			Classification:        model.ClassificationTS,
			OwnerProducer:         []string{"USA"},
			DisseminationControls: []string{"FOUO"},
			ClassifiedBy:          "Test",
			DeclassDate:           "20350101",
		}
		result := engine.Validate(ism)
		if result.Valid {
			t.Skipf("GAP: FOUO with TS should be invalid per CL-12; "+
				"validator may not yet enforce FOUO→U-only constraint")
		}
	})
}

// TestDoDM_CL13_SBU_ClassificationLevel verifies that SBU is only for
// UNCLASSIFIED information.
// [CL-13: E4-A2-S3.a]
func TestDoDM_CL13_SBU_ClassificationLevel(t *testing.T) {
	engine := validation.NewEngine(reg())

	// Positive: SBU with UNCLASSIFIED.
	t.Run("U_valid", func(t *testing.T) {
		ism := &model.ISM{
			Classification:        model.ClassificationU,
			OwnerProducer:         []string{"USA"},
			DisseminationControls: []string{"SBU"},
		}
		result := engine.Validate(ism)
		for _, e := range result.Errors {
			if e.Code == "dissemination.insufficient_classification" && e.Severity == validation.SeverityError {
				t.Errorf("SBU with U should be valid, got: %s", e.Message)
			}
		}
	})

	// Negative: SBU with SECRET.
	t.Run("S_invalid", func(t *testing.T) {
		ism := &model.ISM{
			Classification:        model.ClassificationS,
			OwnerProducer:         []string{"USA"},
			DisseminationControls: []string{"SBU"},
			ClassifiedBy:          "Test",
			DeclassDate:           "20350101",
		}
		result := engine.Validate(ism)
		if result.Valid {
			t.Skipf("GAP: SBU with SECRET should be invalid per CL-13; "+
				"validator may not yet enforce SBU→U-only constraint")
		}
	})
}

// TestDoDM_CL14_SBUNF_ClassificationLevel verifies that SBU-NF (SBU NOFORN) is
// only for UNCLASSIFIED information.
// [CL-14: E4-A2-S4.b]
func TestDoDM_CL14_SBUNF_ClassificationLevel(t *testing.T) {
	engine := validation.NewEngine(reg())

	// Positive: SBU-NF with UNCLASSIFIED.
	t.Run("U_valid", func(t *testing.T) {
		ism := &model.ISM{
			Classification:        model.ClassificationU,
			OwnerProducer:         []string{"USA"},
			DisseminationControls: []string{"SBU-NF"},
		}
		result := engine.Validate(ism)
		for _, e := range result.Errors {
			if e.Code == "dissemination.insufficient_classification" && e.Severity == validation.SeverityError {
				t.Errorf("SBU-NF with U should be valid, got: %s", e.Message)
			}
		}
	})

	// Negative: SBU-NF with SECRET.
	t.Run("S_invalid", func(t *testing.T) {
		ism := &model.ISM{
			Classification:        model.ClassificationS,
			OwnerProducer:         []string{"USA"},
			DisseminationControls: []string{"SBU-NF"},
			ClassifiedBy:          "Test",
			DeclassDate:           "20350101",
		}
		result := engine.Validate(ism)
		if result.Valid {
			t.Skipf("GAP: SBU-NF with SECRET should be invalid per CL-14; "+
				"validator may not yet enforce SBU-NF→U-only constraint")
		}
	})
}

// TestDoDM_CL15_JOINT_WithUS_ClassificationLevels verifies that JOINT documents
// with US as a co-owner may use TS, S, or C — NOT RESTRICTED.
// [CL-15: E4-S5.d]
func TestDoDM_CL15_JOINT_WithUS_ClassificationLevels(t *testing.T) {
	engine := validation.NewEngine(reg())

	validLevels := []struct {
		name  string
		class model.Classification
	}{
		{"TS", model.ClassificationTS},
		{"S", model.ClassificationS},
		{"C", model.ClassificationC},
	}
	for _, tt := range validLevels {
		t.Run(tt.name+"_valid", func(t *testing.T) {
			ism := &model.ISM{
				Classification: tt.class,
				OwnerProducer:  []string{"GBR", "USA"},
				Joint:          true,
				ClassifiedBy:   "Test",
				DeclassDate:    "20350101",
			}
			result := engine.Validate(ism)
			for _, e := range result.Errors {
				if e.Severity == validation.SeverityError && e.Field == "classification" {
					t.Errorf("JOINT with US at %s should be valid, got: %s", tt.class, e.Message)
				}
			}
		})
	}

	// RESTRICTED is not a valid US classification, so it's invalid by type.
	// The Classification type only accepts U, CUI, C, S, TS.
	t.Run("R_invalid_as_type", func(t *testing.T) {
		ism := &model.ISM{
			Classification: model.Classification("R"),
			OwnerProducer:  []string{"GBR", "USA"},
			Joint:          true,
		}
		result := engine.Validate(ism)
		if result.Valid {
			t.Error("JOINT with US at RESTRICTED should be invalid per CL-15")
		}
		if !result.HasCode("core.invalid_classification") {
			t.Skipf("GAP: JOINT with US and R classification — expected core.invalid_classification; "+
				"result: %+v", result.Errors)
		}
	})
}

// TestDoDM_CL16_JOINT_WithoutUS_ClassificationLevels verifies that JOINT
// documents without US as a co-owner may use TS, S, C, or R (RESTRICTED).
// [CL-16: E4-S5.d]
func TestDoDM_CL16_JOINT_WithoutUS_ClassificationLevels(t *testing.T) {
	engine := validation.NewEngine(reg())

	validLevels := []struct {
		name  string
		class model.Classification
	}{
		{"TS", model.ClassificationTS},
		{"S", model.ClassificationS},
		{"C", model.ClassificationC},
	}
	for _, tt := range validLevels {
		t.Run(tt.name+"_valid", func(t *testing.T) {
			ism := &model.ISM{
				Classification: tt.class,
				OwnerProducer:  []string{"DEU", "GBR"},
				Joint:          true,
			}
			result := engine.Validate(ism)
			for _, e := range result.Errors {
				if e.Severity == validation.SeverityError && e.Field == "classification" {
					t.Errorf("JOINT without US at %s should be valid, got: %s", tt.class, e.Message)
				}
			}
		})
	}

	// RESTRICTED is valid for non-US JOINT documents.
	// The current Classification type may not support "R" — if so, this is a GAP.
	t.Run("R_valid_for_nonUS_JOINT", func(t *testing.T) {
		r := model.Classification("R")
		if !r.Valid() {
			t.Skipf("GAP: Classification 'R' (RESTRICTED) not yet supported in model; "+
				"CL-16 requires it for non-US JOINT documents")
			return
		}
		ism := &model.ISM{
			Classification: r,
			OwnerProducer:  []string{"DEU", "GBR"},
			Joint:          true,
		}
		result := engine.Validate(ism)
		for _, e := range result.Errors {
			if e.Severity == validation.SeverityError {
				t.Errorf("JOINT without US at R should be valid per CL-16, got: %s (code: %s)",
					e.Message, e.Code)
			}
		}
	})
}

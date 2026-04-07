package refdata

import (
	"testing"

	"expr.ai/ism-api/internal/model"
)

func TestNewRegistry(t *testing.T) {
	r := NewRegistry()
	if r == nil {
		t.Fatal("NewRegistry returned nil")
	}
	if len(r.Classifications) == 0 {
		t.Error("expected classifications")
	}
	if len(r.CUICategories) == 0 {
		t.Error("expected CUI categories")
	}
	if len(r.DisseminationControls) == 0 {
		t.Error("expected dissemination controls")
	}
	if len(r.DistributionStatements) == 0 {
		t.Error("expected distribution statements")
	}
	if len(r.CountryCodes) == 0 {
		t.Error("expected country codes")
	}
	if len(r.DeclassExceptions) == 0 {
		t.Error("expected declass exceptions")
	}
	if len(r.NonICMarkings) == 0 {
		t.Error("expected non-IC markings")
	}
}

func TestClassifications(t *testing.T) {
	cs := Classifications()
	if len(cs) != 4 {
		t.Fatalf("expected 4 classifications, got %d", len(cs))
	}
	// Verify ordering
	for i := 0; i < len(cs)-1; i++ {
		if cs[i].Level >= cs[i+1].Level {
			t.Errorf("classifications not in ascending order: %s(%d) >= %s(%d)",
				cs[i].Code, cs[i].Level, cs[i+1].Code, cs[i+1].Level)
		}
	}
}

func TestCUICategories(t *testing.T) {
	cats := CUICategories()
	var specified, basic int
	for _, c := range cats {
		switch c.Type {
		case "specified":
			specified++
			if len(c.Code) < 3 || c.Code[:3] != "SP-" {
				t.Errorf("specified category %q missing SP- prefix", c.Code)
			}
		case "basic":
			basic++
		default:
			t.Errorf("unknown CUI category type: %q", c.Type)
		}
		if c.Label == "" {
			t.Errorf("CUI category %q has empty label", c.Code)
		}
		if c.Description == "" {
			t.Errorf("CUI category %q has empty description", c.Code)
		}
	}
	if specified != 58 {
		t.Errorf("expected 58 specified categories, got %d", specified)
	}
	if basic != 6 {
		t.Errorf("expected 6 basic categories, got %d", basic)
	}
}

func TestDisseminationControls(t *testing.T) {
	controls := DisseminationControls()
	if len(controls) != 18 {
		t.Fatalf("expected 18 dissemination controls, got %d", len(controls))
	}
	codes := make(map[string]bool)
	for _, c := range controls {
		if codes[c.Code] {
			t.Errorf("duplicate dissemination control code: %q", c.Code)
		}
		codes[c.Code] = true
		if c.Label == "" {
			t.Errorf("dissemination control %q has empty label", c.Code)
		}
	}
	// Verify specific metadata
	for _, c := range controls {
		if c.Code == "REL" && c.RequiresField != "releasableTo" {
			t.Error("REL should require releasableTo field")
		}
		if c.Code == "DISPLAY ONLY" && c.RequiresField != "displayOnlyTo" {
			t.Error("DISPLAY ONLY should require displayOnlyTo field")
		}
		if c.Code == "NOFORN" && len(c.ExclusiveWith) == 0 {
			t.Error("NOFORN should have exclusiveWith entries")
		}
	}
}

func TestDistributionStatements(t *testing.T) {
	stmts := DistributionStatements()
	if len(stmts) != 6 {
		t.Fatalf("expected 6 distribution statements, got %d", len(stmts))
	}
	expected := []string{"A", "B", "C", "D", "E", "F"}
	for i, s := range stmts {
		if s.Code != expected[i] {
			t.Errorf("expected statement %q at index %d, got %q", expected[i], i, s.Code)
		}
		if s.Text == "" {
			t.Errorf("distribution statement %q has empty text", s.Code)
		}
	}
	// Statement A should be constrained to U/CUI
	if stmts[0].ClassificationConstraint != "U,CUI" {
		t.Errorf("Statement A should have classification constraint U,CUI, got %q",
			stmts[0].ClassificationConstraint)
	}
}

func TestCountryCodes(t *testing.T) {
	codes := CountryCodes()
	if len(codes) == 0 {
		t.Fatal("expected country codes")
	}
	types := map[string]int{}
	seen := map[string]bool{}
	for _, c := range codes {
		if seen[c.Code] {
			t.Errorf("duplicate country code: %q", c.Code)
		}
		seen[c.Code] = true
		types[c.Type]++
		if c.Name == "" {
			t.Errorf("country code %q has empty name", c.Code)
		}
	}
	if types["country"] == 0 {
		t.Error("expected at least one country type")
	}
	if types["coalition"] == 0 {
		t.Error("expected at least one coalition type")
	}
	if types["organization"] == 0 {
		t.Error("expected at least one organization type")
	}
	// FVEY and USA must be present
	if !seen["USA"] {
		t.Error("USA must be present")
	}
	if !seen["FVEY"] {
		t.Error("FVEY must be present")
	}
	if !seen["NATO"] {
		t.Error("NATO must be present")
	}
}

func TestDeclassExceptions(t *testing.T) {
	exceptions := DeclassExceptions()
	if len(exceptions) == 0 {
		t.Fatal("expected declass exceptions")
	}
	seen := map[string]bool{}
	for _, e := range exceptions {
		if seen[e.Code] {
			t.Errorf("duplicate declass exception: %q", e.Code)
		}
		seen[e.Code] = true
		if e.Label == "" {
			t.Errorf("declass exception %q has empty label", e.Code)
		}
		if e.Category == "" {
			t.Errorf("declass exception %q has empty category", e.Code)
		}
	}
	if !seen["25X1"] {
		t.Error("25X1 must be present")
	}
}

func TestNonICMarkings(t *testing.T) {
	markings := NonICMarkings()
	if len(markings) == 0 {
		t.Fatal("expected non-IC markings")
	}
	seen := map[string]bool{}
	for _, m := range markings {
		if seen[m.Code] {
			t.Errorf("duplicate non-IC marking: %q", m.Code)
		}
		seen[m.Code] = true
		if m.Label == "" {
			t.Errorf("non-IC marking %q has empty label", m.Code)
		}
	}
	if !seen["SBU"] {
		t.Error("SBU must be present")
	}
	if !seen["LES"] {
		t.Error("LES must be present")
	}
}

func TestRegistryValidation(t *testing.T) {
	r := NewRegistry()

	// Classification lookups
	if !r.ValidClassification(model.ClassificationS) {
		t.Error("S should be valid classification")
	}
	if r.ValidClassification("TS") {
		t.Error("TS should not be valid classification")
	}

	// CUI lookups
	if !r.ValidCUICategory("SP-CEII") {
		t.Error("SP-CEII should be valid CUI category")
	}
	if r.ValidCUICategory("BOGUS") {
		t.Error("BOGUS should not be valid CUI category")
	}

	// Dissemination control lookups
	if !r.ValidDisseminationControl("NOFORN") {
		t.Error("NOFORN should be valid")
	}
	if r.ValidDisseminationControl("BOGUS") {
		t.Error("BOGUS should not be valid")
	}

	// Distribution statement lookups
	if !r.ValidDistributionStatement("A") {
		t.Error("A should be valid")
	}
	if r.ValidDistributionStatement("G") {
		t.Error("G should not be valid")
	}

	// Country code lookups
	if !r.ValidCountryCode("USA") {
		t.Error("USA should be valid")
	}
	if r.ValidCountryCode("ZZZ") {
		t.Error("ZZZ should not be valid")
	}

	// Declass exception lookups
	if !r.ValidDeclassException("25X1") {
		t.Error("25X1 should be valid")
	}
	if r.ValidDeclassException("99X") {
		t.Error("99X should not be valid")
	}

	// Non-IC marking lookups
	if !r.ValidNonICMarking("SBU") {
		t.Error("SBU should be valid")
	}
	if r.ValidNonICMarking("BOGUS") {
		t.Error("BOGUS should not be valid")
	}
}

func TestLookupDisseminationControl(t *testing.T) {
	r := NewRegistry()

	rel := r.LookupDisseminationControl("REL")
	if rel == nil {
		t.Fatal("REL lookup returned nil")
	}
	if rel.RequiresField != "releasableTo" {
		t.Errorf("REL requiresField = %q, want releasableTo", rel.RequiresField)
	}

	missing := r.LookupDisseminationControl("BOGUS")
	if missing != nil {
		t.Error("BOGUS lookup should return nil")
	}
}

func TestLookupDistributionStatement(t *testing.T) {
	r := NewRegistry()

	a := r.LookupDistributionStatement("A")
	if a == nil {
		t.Fatal("Statement A lookup returned nil")
	}
	if a.ClassificationConstraint != "U,CUI" {
		t.Errorf("Statement A constraint = %q, want U,CUI", a.ClassificationConstraint)
	}

	missing := r.LookupDistributionStatement("G")
	if missing != nil {
		t.Error("Statement G lookup should return nil")
	}
}

func TestCompatibility(t *testing.T) {
	pairs := ExclusiveDisseminationPairs()
	if len(pairs) == 0 {
		t.Error("expected exclusive pairs")
	}

	reqs := DisseminationFieldRequirements()
	if len(reqs) == 0 {
		t.Error("expected field requirements")
	}

	gates := DisseminationClassificationGates()
	if len(gates) == 0 {
		t.Error("expected classification gates")
	}

	constraints := DistributionClassificationConstraints()
	if len(constraints) == 0 {
		t.Error("expected distribution constraints")
	}

	exclusions := DeclassFieldExclusions()
	if len(exclusions) == 0 {
		t.Error("expected declass field exclusions")
	}

	applicable := DeclassApplicableClassifications()
	if len(applicable) != 2 {
		t.Errorf("expected 2 declass-applicable classifications, got %d", len(applicable))
	}
}

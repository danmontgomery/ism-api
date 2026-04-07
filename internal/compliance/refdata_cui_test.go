package compliance_test

import (
	"testing"
)

// TestXSD_CUI_SpecifiedCategoriesPresent verifies CUI specified categories from
// CVEnumISMCUISpecified.xsd against the API registry.
func TestXSD_CUI_SpecifiedCategoriesPresent(t *testing.T) {
	r := reg()
	for _, code := range xsdCUISpecified {
		t.Run(code, func(t *testing.T) {
			// The API uses SP- prefix for specified categories
			spCode := "SP-" + code
			if !r.ValidCUICategory(code) && !r.ValidCUICategory(spCode) {
				t.Errorf("CUI specified category %s not in registry — required by CVEnumISMCUISpecified.xsd", code)
			}
		})
	}
}

// TestXSD_CUI_CoverageReport reports CUI category coverage.
func TestXSD_CUI_CoverageReport(t *testing.T) {
	r := reg()
	present := 0
	for _, code := range xsdCUISpecified {
		spCode := "SP-" + code
		if r.ValidCUICategory(code) || r.ValidCUICategory(spCode) {
			present++
		}
	}
	t.Logf("CUI specified coverage: %d/%d XSD categories in registry", present, len(xsdCUISpecified))
	t.Logf("API registry has %d CUI categories total", len(r.CUICategories))
}

// TestXSD_CUI_APICategoriesToXSD checks API CUI categories against XSD.
func TestXSD_CUI_APICategoriesToXSD(t *testing.T) {
	r := reg()
	xsdSet := make(map[string]bool)
	for _, code := range xsdCUISpecified {
		xsdSet[code] = true
	}

	for _, cat := range r.CUICategories {
		t.Run(cat.Code, func(t *testing.T) {
			// Strip SP- prefix for comparison
			code := cat.Code
			if len(code) > 3 && code[:3] == "SP-" {
				code = code[3:]
			}
			if !xsdSet[code] {
				t.Logf("NOTE: API CUI category %q (%s) has no direct match in CVEnumISMCUISpecified.xsd", cat.Code, cat.Label)
			}
		})
	}
}

// TestXSD_CUI_StructFieldPresent verifies the ISM struct has cuiSpecified and cuiBasic fields.
func TestXSD_CUI_StructFieldPresent(t *testing.T) {
	t.Run("cuiSpecified", func(t *testing.T) {
		if !requireStructField(t, "cuiSpecified") {
			t.Error("ISM struct missing cuiSpecified field — required by IC-ISM.xsd")
		}
	})
	t.Run("cuiBasic", func(t *testing.T) {
		if !requireStructField(t, "cuiBasic") {
			t.Error("ISM struct missing cuiBasic field — required by IC-ISM.xsd")
		}
	})
	t.Run("categoryMarkings_exists", func(t *testing.T) {
		// The API uses categoryMarkings which maps to CUI categories
		if !requireStructField(t, "categoryMarkings") {
			t.Fatal("ISM struct must have categoryMarkings field")
		}
	})
}

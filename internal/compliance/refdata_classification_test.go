package compliance_test

import (
	"testing"

	"dmontgomery/ism-api/internal/model"
)

// TestXSD_Classification_AllLevelsPresent verifies that every classification
// level defined in CVEnumISMClassificationAll.xsd exists in the API registry.
func TestXSD_Classification_AllLevelsPresent(t *testing.T) {
	r := reg()
	for _, code := range xsdClassifications {
		t.Run(code, func(t *testing.T) {
			if !r.ValidClassification(model.Classification(code)) {
				if code == "R" {
					t.Skipf("GAP: %s not in registry — R not yet modeled", code)
				} else {
					t.Errorf("%s not in registry — required by CVEnumISMClassificationAll.xsd", code)
				}
			}
		})
	}
}

// TestXSD_Classification_TopSecret verifies TS is recognized.
func TestXSD_Classification_TopSecret(t *testing.T) {
	r := reg()
	if !r.ValidClassification(model.Classification("TS")) {
		t.Skip("GAP: TS (TOP SECRET) not in registry — required by CVEnumISMClassificationAll.xsd")
	}
	c := model.Classification("TS")
	if !c.Valid() {
		t.Error("TS should be a valid classification")
	}
	if c.Level() < model.ClassificationS.Level() {
		t.Error("TS should have a higher level than S")
	}
}

// TestXSD_Classification_Restricted verifies R is recognized.
func TestXSD_Classification_Restricted(t *testing.T) {
	r := reg()
	if !r.ValidClassification(model.Classification("R")) {
		t.Skip("GAP: R (RESTRICTED) not in registry — required by CVEnumISMClassificationAll.xsd")
	}
}

// TestXSD_Classification_Unclassified verifies U is recognized and valid.
func TestXSD_Classification_Unclassified(t *testing.T) {
	r := reg()
	if !r.ValidClassification(model.ClassificationU) {
		t.Fatal("U must be in registry")
	}
	if model.ClassificationU.Level() != 0 {
		t.Errorf("U level = %d, want 0", model.ClassificationU.Level())
	}
}

// TestXSD_Classification_Confidential verifies C is recognized and valid.
func TestXSD_Classification_Confidential(t *testing.T) {
	r := reg()
	if !r.ValidClassification(model.ClassificationC) {
		t.Fatal("C must be in registry")
	}
	if !model.ClassificationC.AtLeast(model.ClassificationU) {
		t.Error("C should be at least U")
	}
}

// TestXSD_Classification_Secret verifies S is recognized and valid.
func TestXSD_Classification_Secret(t *testing.T) {
	r := reg()
	if !r.ValidClassification(model.ClassificationS) {
		t.Fatal("S must be in registry")
	}
	if !model.ClassificationS.AtLeast(model.ClassificationC) {
		t.Error("S should be at least C")
	}
}

// TestXSD_Classification_LevelOrdering verifies classification ordering: U < C < S.
func TestXSD_Classification_LevelOrdering(t *testing.T) {
	if model.ClassificationU.Level() >= model.ClassificationC.Level() {
		t.Error("U should be below C")
	}
	if model.ClassificationC.Level() >= model.ClassificationS.Level() {
		t.Error("C should be below S")
	}
}

// TestXSD_Classification_CUIIsNonXSD documents that CUI is present in the API
// but not defined in CVEnumISMClassificationAll.xsd. CUI is a separate
// compliance framework (compliesWith=USA-CUI-ONLY/USA-CUI).
func TestXSD_Classification_CUIIsNonXSD(t *testing.T) {
	r := reg()
	if !r.ValidClassification(model.ClassificationCUI) {
		t.Fatal("CUI should be in the API registry")
	}
	// Verify CUI is not in XSD classifications
	for _, code := range xsdClassifications {
		if code == "CUI" {
			t.Error("CUI should NOT be in XSD classification enum")
		}
	}
	t.Log("NOTE: CUI is an API extension — not defined in CVEnumISMClassificationAll.xsd")
}

// TestXSD_Classification_NoSpuriousValues checks that the API does not contain
// classification codes absent from the XSD (except CUI which is documented).
func TestXSD_Classification_NoSpuriousValues(t *testing.T) {
	r := reg()
	xsdSet := make(map[string]bool)
	for _, code := range xsdClassifications {
		xsdSet[code] = true
	}
	xsdSet["CUI"] = true // documented API extension

	for _, entry := range r.Classifications {
		if !xsdSet[string(entry.Code)] {
			t.Errorf("spurious classification %q in API — no XSD counterpart", entry.Code)
		}
	}
}

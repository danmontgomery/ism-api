package compliance_test

import (
	"testing"
)

// TestXSD_SCI_NotModeled verifies that SCI controls are modeled in the API.
// The ISM struct must have an sciControls field and the registry must have SCI data.
func TestXSD_SCI_NotModeled(t *testing.T) {
	if !requireStructField(t, "sciControls") {
		t.Fatal("ISM struct missing sciControls field — 20 SCI controls from CVEnumISMSCIControls.xsd not modeled")
	}
	r := reg()
	if len(r.SCIControls) == 0 {
		t.Fatal("Registry has no SCI controls data")
	}
}

// TestXSD_SCI_AllControlsPresent checks each individual SCI control.
func TestXSD_SCI_AllControlsPresent(t *testing.T) {
	if !requireStructField(t, "sciControls") {
		t.Fatal("ISM struct missing sciControls field — cannot test SCI controls")
	}
	r := reg()
	missing := assertRegistryContains(t, r.ValidSCIControl, xsdSCIControls, "CVEnumISMSCIControls.xsd")
	if missing > 0 {
		t.Errorf("%d/%d SCI controls missing from registry", missing, len(xsdSCIControls))
	}
}

// TestXSD_SCI_ControlCategories verifies SCI controls by category.
func TestXSD_SCI_ControlCategories(t *testing.T) {
	categories := map[string][]string{
		"BUR (BYELEMAN)":            {"BUR", "BUR-BLG", "BUR-DTP", "BUR-WRG"},
		"HCS (HUMINT)":              {"HCS", "HCS-O", "HCS-P", "HCS-X"},
		"KLM (KLAMATH)":             {"KLM", "KLM-R"},
		"SI (SPECIAL INTELLIGENCE)": {"SI", "SI-EU", "SI-G", "SI-NK"},
		"TK (TALENT KEYHOLE)":       {"TK", "TK-BLFH", "TK-IDIT", "TK-KAND"},
		"MVL (MARVEL)":              {"MVL"},
		"RSV (RESERVE)":             {"RSV"},
	}
	r := reg()
	for cat, codes := range categories {
		t.Run(cat, func(t *testing.T) {
			for _, code := range codes {
				ctrl := r.LookupSCIControl(code)
				if ctrl == nil {
					t.Errorf("SCI control %s not found in registry", code)
					continue
				}
				if ctrl.Category != cat {
					t.Errorf("SCI control %s: category = %q, want %q", code, ctrl.Category, cat)
				}
			}
		})
	}
}

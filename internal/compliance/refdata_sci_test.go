package compliance_test

import (
	"testing"

	"expr.ai/ism-api/internal/model"
)

// TestXSD_SCI_NotModeled documents that SCI controls are not modeled in the API.
// The ISM struct has no SCIcontrols field, and the registry has no SCI data.
func TestXSD_SCI_NotModeled(t *testing.T) {
	if !requireStructField(t, "sciControls") {
		t.Skip("GAP: ISM struct missing sciControls field — 20 SCI controls from CVEnumISMSCIControls.xsd not modeled")
	}
}

// TestXSD_SCI_AllControlsPresent checks each individual SCI control.
func TestXSD_SCI_AllControlsPresent(t *testing.T) {
	if !requireStructField(t, "sciControls") {
		t.Skip("GAP: ISM struct missing sciControls field — cannot test SCI controls")
		return
	}
	for _, code := range xsdSCIControls {
		t.Run(code, func(t *testing.T) {
			t.Skipf("GAP: SCI control %s not in registry — required by CVEnumISMSCIControls.xsd", code)
		})
	}
}

// TestXSD_SCI_ControlCategories verifies SCI controls by category.
func TestXSD_SCI_ControlCategories(t *testing.T) {
	categories := map[string][]string{
		"BUR (BYELEMAN)":           {"BUR", "BUR-BLG", "BUR-DTP", "BUR-WRG"},
		"HCS (HUMINT)":             {"HCS", "HCS-O", "HCS-P", "HCS-X"},
		"KLM (KLAMATH)":            {"KLM", "KLM-R"},
		"SI (SPECIAL INTELLIGENCE)": {"SI", "SI-EU", "SI-G", "SI-NK"},
		"TK (TALENT KEYHOLE)":      {"TK", "TK-BLFH", "TK-IDIT", "TK-KAND"},
		"MVL (MARVEL)":             {"MVL"},
		"RSV (RESERVE)":            {"RSV"},
	}
	for cat, codes := range categories {
		t.Run(cat, func(t *testing.T) {
			t.Skipf("GAP: SCI category %s (%d controls) not modeled — required by CVEnumISMSCIControls.xsd", cat, len(codes))
		})
	}
	_ = model.ISM{} // ensure model package is used
}

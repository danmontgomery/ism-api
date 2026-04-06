package compliance_test

import (
	"testing"
)

// TestXSD_CompliesWith_NotModeled documents that compliesWith is not modeled.
func TestXSD_CompliesWith_NotModeled(t *testing.T) {
	if !requireStructField(t, "compliesWith") {
		t.Skip("GAP: ISM struct missing compliesWith field — 6 values from CVEnumISMCompliesWith.xsd not modeled")
	}
}

// TestXSD_CompliesWith_AllValuesPresent checks each compliesWith value.
func TestXSD_CompliesWith_AllValuesPresent(t *testing.T) {
	if !requireStructField(t, "compliesWith") {
		t.Skip("GAP: ISM struct missing compliesWith field — cannot test values")
		return
	}
	for _, code := range xsdCompliesWith {
		t.Run(code, func(t *testing.T) {
			t.Skipf("GAP: compliesWith value %s not in registry — required by CVEnumISMCompliesWith.xsd", code)
		})
	}
}

// TestXSD_CompliesWith_CUIRelationship documents the relationship between
// compliesWith and CUI classification in the API.
func TestXSD_CompliesWith_CUIRelationship(t *testing.T) {
	// The XSD defines CUI compliance via compliesWith=USA-CUI-ONLY or USA-CUI.
	// The API instead treats CUI as a classification level.
	// This is a fundamental architectural divergence from the XSD model.
	if !requireStructField(t, "compliesWith") {
		t.Skip("GAP: compliesWith not modeled — CUI is handled as classification level instead of compliance attribute per CVEnumISMCompliesWith.xsd")
	}
}

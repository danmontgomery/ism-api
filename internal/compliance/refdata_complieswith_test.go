package compliance_test

import (
	"testing"
)

// TestXSD_CompliesWith_NotModeled verifies the ISM struct has a compliesWith field.
func TestXSD_CompliesWith_NotModeled(t *testing.T) {
	if !requireStructField(t, "compliesWith") {
		t.Fatal("ISM struct must have compliesWith field")
	}
}

// TestXSD_CompliesWith_AllValuesPresent checks each XSD compliesWith value against the registry.
func TestXSD_CompliesWith_AllValuesPresent(t *testing.T) {
	if !requireStructField(t, "compliesWith") {
		t.Fatal("ISM struct must have compliesWith field")
	}
	r := reg()
	missing := assertRegistryContains(t, r.ValidCompliesWith, xsdCompliesWith,
		"CVEnumISMCompliesWith.xsd")
	if missing > 0 {
		t.Errorf("%d compliesWith values from XSD missing in registry", missing)
	}
}

// TestXSD_CompliesWith_CUIRelationship documents the relationship between
// compliesWith and CUI classification in the API.
func TestXSD_CompliesWith_CUIRelationship(t *testing.T) {
	// The XSD defines CUI compliance via compliesWith=USA-CUI-ONLY or USA-CUI.
	// The API models compliesWith as a first-class attribute per CVEnumISMCompliesWith.xsd.
	if !requireStructField(t, "compliesWith") {
		t.Fatal("ISM struct must have compliesWith field")
	}
	r := reg()
	if !r.ValidCompliesWith("USA-CUI-ONLY") {
		t.Error("USA-CUI-ONLY must be a valid compliesWith value")
	}
	if !r.ValidCompliesWith("USA-CUI") {
		t.Error("USA-CUI must be a valid compliesWith value")
	}
}

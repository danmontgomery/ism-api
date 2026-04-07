package compliance_test

import (
	"testing"
)

// TestXSD_NonUS_StructFieldPresent verifies the ISM struct has a nonUSControls field.
func TestXSD_NonUS_StructFieldPresent(t *testing.T) {
	if !requireStructField(t, "nonUSControls") {
		t.Fatal("ISM struct must have nonUSControls field")
	}
}

// TestXSD_NonUS_AllControlsPresent checks each XSD non-US control against the registry.
func TestXSD_NonUS_AllControlsPresent(t *testing.T) {
	r := reg()
	missing := assertRegistryContains(t, r.ValidNonUSControl, xsdNonUSControls,
		"CVEnumISMNonUSControls.xsd")
	if missing > 0 {
		t.Errorf("%d non-US controls from XSD missing in registry", missing)
	}
}

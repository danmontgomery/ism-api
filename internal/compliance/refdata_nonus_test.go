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

// TestXSD_NonUS_AllControlsPresent checks each XSD non-US control against the API.
// The API currently has the field but no registry validation for non-US controls.
func TestXSD_NonUS_AllControlsPresent(t *testing.T) {
	for _, code := range xsdNonUSControls {
		t.Run(code, func(t *testing.T) {
			// The ISM struct has the field, but there's no ValidNonUSControl method
			// on the registry. The field accepts arbitrary strings.
			t.Skipf("GAP: no registry validation for non-US control %s — required by CVEnumISMNonUSControls.xsd", code)
		})
	}
}

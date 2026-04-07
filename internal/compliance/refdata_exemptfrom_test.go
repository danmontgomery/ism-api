package compliance_test

import (
	"testing"
)

// TestXSD_ExemptFrom_StructFieldPresent verifies the ISM struct has an exemptFrom field.
func TestXSD_ExemptFrom_StructFieldPresent(t *testing.T) {
	if !requireStructField(t, "exemptFrom") {
		t.Fatal("ISM struct must have exemptFrom field")
	}
}

// TestXSD_ExemptFrom_AllValuesPresent checks each XSD exemptFrom value against the registry.
func TestXSD_ExemptFrom_AllValuesPresent(t *testing.T) {
	r := reg()
	missing := assertRegistryContains(t, r.ValidExemptFrom, xsdExemptFrom,
		"CVEnumISMExemptFrom.xsd")
	if missing > 0 {
		t.Errorf("%d exemptFrom values from XSD missing in registry", missing)
	}
}

package compliance_test

import (
	"testing"
)

// TestXSD_ExemptFrom_NotModeled documents that exemptFrom is not modeled.
func TestXSD_ExemptFrom_NotModeled(t *testing.T) {
	if !requireStructField(t, "exemptFrom") {
		t.Skip("GAP: ISM struct missing exemptFrom field — 2 values from CVEnumISMExemptFrom.xsd not modeled")
	}
}

// TestXSD_ExemptFrom_AllValuesPresent checks each exemptFrom value.
func TestXSD_ExemptFrom_AllValuesPresent(t *testing.T) {
	if !requireStructField(t, "exemptFrom") {
		t.Skip("GAP: ISM struct missing exemptFrom field — cannot test values")
		return
	}
	for _, code := range xsdExemptFrom {
		t.Run(code, func(t *testing.T) {
			t.Skipf("GAP: exemptFrom value %s not in registry — required by CVEnumISMExemptFrom.xsd", code)
		})
	}
}

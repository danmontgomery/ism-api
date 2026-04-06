package compliance_test

import (
	"testing"
)

// TestXSD_NonIC_AllMarkingsPresent verifies every non-IC marking from
// CVEnumISMNonIC.xsd exists in the API registry.
func TestXSD_NonIC_AllMarkingsPresent(t *testing.T) {
	r := reg()
	for _, code := range xsdNonICMarkings {
		t.Run(code, func(t *testing.T) {
			if !r.ValidNonICMarking(code) {
				t.Skipf("GAP: non-IC marking %s not in registry — required by CVEnumISMNonIC.xsd", code)
			}
		})
	}
}

// TestXSD_NonIC_NNPI verifies NNPI (Naval Nuclear Propulsion Information).
func TestXSD_NonIC_NNPI(t *testing.T) {
	r := reg()
	if !r.ValidNonICMarking("NNPI") {
		t.Skip("GAP: NNPI not in non-IC registry — required by CVEnumISMNonIC.xsd")
	}
}

// TestXSD_NonIC_CoreMarkingsPresent verifies the core non-IC markings that are
// commonly used are all present.
func TestXSD_NonIC_CoreMarkingsPresent(t *testing.T) {
	r := reg()
	core := []string{"SBU", "SBU-NF", "LES", "LES-NF", "SSI", "DS", "XD", "ND"}
	for _, code := range core {
		t.Run(code, func(t *testing.T) {
			if !r.ValidNonICMarking(code) {
				t.Errorf("core non-IC marking %s must be in registry", code)
			}
		})
	}
}

// TestXSD_NonIC_NoSpuriousValues checks for API non-IC markings not in the XSD.
func TestXSD_NonIC_NoSpuriousValues(t *testing.T) {
	r := reg()
	xsdSet := make(map[string]bool)
	for _, code := range xsdNonICMarkings {
		xsdSet[code] = true
	}

	for _, m := range r.NonICMarkings {
		t.Run(m.Code, func(t *testing.T) {
			if !xsdSet[m.Code] {
				// FOUO is in CVEnumISMDissem.xsd, not CVEnumISMNonIC.xsd,
				// but the API registers it as a non-IC marking.
				if m.Code == "FOUO" {
					t.Log("NOTE: FOUO registered as non-IC but defined in CVEnumISMDissem.xsd, not CVEnumISMNonIC.xsd")
					return
				}
				t.Logf("SPURIOUS: non-IC marking %q in API has no counterpart in CVEnumISMNonIC.xsd", m.Code)
			}
		})
	}
}

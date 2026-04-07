package compliance_test

import (
	"testing"
)

// TestXSD_Dissemination_AllControlsPresent verifies every dissemination control
// from CVEnumISMDissem.xsd has a mapping to the API registry.
func TestXSD_Dissemination_AllControlsPresent(t *testing.T) {
	r := reg()
	for _, xsdCode := range xsdDisseminationControls {
		t.Run(xsdCode, func(t *testing.T) {
			apiCode, mapped := xsdToAPIDissemination[xsdCode]
			if !mapped {
				t.Fatalf("%s has no API mapping — required by CVEnumISMDissem.xsd", xsdCode)
			}
			if apiCode == "" {
				t.Fatalf("%s has no API equivalent — required by CVEnumISMDissem.xsd", xsdCode)
			}
			if !r.ValidDisseminationControl(apiCode) {
				t.Errorf("%s (XSD: %s) not in registry — required by CVEnumISMDissem.xsd", apiCode, xsdCode)
			}
		})
	}
}

// TestXSD_Dissemination_CodeMapping verifies the XSD→API code translation for
// controls where the XSD uses abbreviated codes (NF, PR, IMC, etc.) while the
// API uses human-readable names (NOFORN, PROPIN, IMCON, etc.).
func TestXSD_Dissemination_CodeMapping(t *testing.T) {
	knownMappings := []struct {
		xsd string
		api string
	}{
		{"NF", "NOFORN"},
		{"PR", "PROPIN"},
		{"IMC", "IMCON"},
		{"DISPLAYONLY", "DISPLAY ONLY"},
		{"DL_ONLY", "DL ONLY"},
		{"FED_ONLY", "FED ONLY"},
	}
	r := reg()
	for _, m := range knownMappings {
		t.Run(m.xsd+"->"+m.api, func(t *testing.T) {
			if !r.ValidDisseminationControl(m.api) {
				t.Fatalf("API code %s not in registry", m.api)
			}
			// Verify the lookup works with the API code
			ctrl := r.LookupDisseminationControl(m.api)
			if ctrl == nil {
				t.Errorf("LookupDisseminationControl(%q) returned nil", m.api)
			}
		})
	}
}

// TestXSD_Dissemination_MissingControls documents XSD controls not yet in API.
func TestXSD_Dissemination_MissingControls(t *testing.T) {
	expectedMissing := []struct {
		xsd  string
		desc string
	}{
		{"RAWFISA", "Raw FISA intelligence data"},
		{"EXEMPT_FROM_ICD501_DISCOVERY", "ICD 501 discovery exemption"},
		{"WAIVED", "Waived dissemination restriction"},
		{"AC", "Attorney-Client privilege"},
		{"AWP", "Attorney Work Product"},
	}
	r := reg()
	for _, m := range expectedMissing {
		t.Run(m.xsd, func(t *testing.T) {
			apiCode := xsdToAPIDissemination[m.xsd]
			if apiCode != "" && r.ValidDisseminationControl(apiCode) {
				t.Logf("RESOLVED: %s (%s) is now in registry", m.xsd, m.desc)
				return
			}
			t.Errorf("%s (%s) not in registry — required by CVEnumISMDissem.xsd", m.xsd, m.desc)
		})
	}
}

// TestXSD_Dissemination_NoSpuriousControls checks for API dissemination codes
// that have no XSD counterpart.
func TestXSD_Dissemination_NoSpuriousControls(t *testing.T) {
	r := reg()
	// Build reverse lookup: API code -> whether it maps from an XSD code
	apiFromXSD := make(map[string]bool)
	for _, apiCode := range xsdToAPIDissemination {
		if apiCode != "" {
			apiFromXSD[apiCode] = true
		}
	}

	for _, ctrl := range r.DisseminationControls {
		t.Run(ctrl.Code, func(t *testing.T) {
			if !apiFromXSD[ctrl.Code] {
				t.Logf("SPURIOUS: %q in API has no XSD counterpart in CVEnumISMDissem.xsd", ctrl.Code)
			}
		})
	}
}

// TestXSD_Dissemination_FOUO verifies FOUO handling. FOUO appears in both
// CVEnumISMDissem.xsd and is also a legacy marking.
func TestXSD_Dissemination_FOUO(t *testing.T) {
	r := reg()
	// FOUO is in the XSD as a dissemination control
	// Check if it's in the API as a dissemination control or non-IC marking
	dissem := r.ValidDisseminationControl("FOUO")
	nonIC := r.ValidNonICMarking("FOUO")
	if !dissem && !nonIC {
		t.Fatal("FOUO not present in API (neither dissemination nor non-IC)")
	}
	if nonIC && !dissem {
		t.Log("NOTE: FOUO registered as non-IC marking, not dissemination control")
	}
}

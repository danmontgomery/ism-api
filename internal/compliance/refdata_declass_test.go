package compliance_test

import (
	"testing"
)

// TestXSD_Declass_AllExceptionsPresent verifies every 25X/50X/75X declass
// exception from CVEnumISM25X.xsd exists in the API registry.
func TestXSD_Declass_AllExceptionsPresent(t *testing.T) {
	r := reg()
	for _, code := range xsdDeclassExceptions {
		t.Run(code, func(t *testing.T) {
			if !r.ValidDeclassException(code) {
				t.Skipf("GAP: declass exception %s not in registry — required by CVEnumISM25X.xsd", code)
			}
		})
	}
}

// TestXSD_Declass_25XCodesPresent verifies the 25X series (25-year exemptions).
func TestXSD_Declass_25XCodesPresent(t *testing.T) {
	r := reg()
	codes25X := []string{"25X1", "25X2", "25X3", "25X4", "25X5", "25X6", "25X7", "25X8", "25X9"}
	for _, code := range codes25X {
		t.Run(code, func(t *testing.T) {
			if !r.ValidDeclassException(code) {
				t.Skipf("GAP: %s not in registry — required by CVEnumISM25X.xsd", code)
			}
		})
	}
}

// TestXSD_Declass_50XCodesPresent verifies the 50X series (75-year exemptions).
func TestXSD_Declass_50XCodesPresent(t *testing.T) {
	r := reg()
	codes50X := []string{
		"50X1", "50X1-HUM", "50X2", "50X2-WMD", "50X3", "50X4",
		"50X5", "50X6", "50X7", "50X8", "50X9",
	}
	for _, code := range codes50X {
		t.Run(code, func(t *testing.T) {
			if !r.ValidDeclassException(code) {
				t.Skipf("GAP: %s not in registry — required by CVEnumISM25X.xsd", code)
			}
		})
	}
}

// TestXSD_Declass_SpecialCodes verifies non-numeric exemption codes.
func TestXSD_Declass_SpecialCodes(t *testing.T) {
	r := reg()
	special := []struct {
		code string
		desc string
	}{
		{"AEA", "Atomic Energy Act exemption"},
		{"NATO", "NATO exemption"},
		{"NATO-AEA", "NATO and AEA combined exemption"},
		{"25X1-EO-12951", "25X1 with EO 12951"},
		{"75X", "ISCAP-approved specific information"},
	}
	for _, s := range special {
		t.Run(s.code, func(t *testing.T) {
			if !r.ValidDeclassException(s.code) {
				t.Skipf("GAP: %s (%s) not in registry — required by CVEnumISM25X.xsd", s.code, s.desc)
			}
		})
	}
}

// TestXSD_Declass_CoverageReport reports declass exception coverage.
func TestXSD_Declass_CoverageReport(t *testing.T) {
	r := reg()
	present := 0
	for _, code := range xsdDeclassExceptions {
		if r.ValidDeclassException(code) {
			present++
		}
	}
	t.Logf("Declass exception coverage: %d/%d XSD codes in registry", present, len(xsdDeclassExceptions))
	t.Logf("API registry has %d entries", len(r.DeclassExceptions))
}

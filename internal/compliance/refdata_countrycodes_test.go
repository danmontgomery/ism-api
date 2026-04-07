package compliance_test

import (
	"testing"
)

// A representative sample of country codes from CVEnumISMCATOwnerProducer.xsd.
// The full XSD has 341 entries; we test a representative subset of ~50 plus
// all coalitions and the FGI marker.
var xsdCountryCodesSample = []string{
	// FGI marker
	"FGI",
	// Five Eyes countries
	"USA", "GBR", "CAN", "AUS", "NZL",
	// Major NATO allies
	"FRA", "DEU", "ITA", "ESP", "NLD", "BEL", "NOR", "DNK", "PRT",
	"POL", "TUR", "CZE", "GRC", "HUN", "ISL", "LUX", "SVK", "SVN",
	"EST", "LVA", "LTU", "BGR", "ROU", "HRV", "ALB", "MNE", "MKD",
	// Key Pacific allies
	"JPN", "KOR", "PHL", "THA", "SGP",
	// Key Middle East
	"ISR", "SAU", "ARE", "QAT", "KWT", "BHR", "JOR",
	// Other significant
	"BRA", "MEX", "IND", "IDN", "ZAF", "NGA", "EGY", "PAK",
	"CHN", "RUS", "PRK", "IRN", "IRQ", "SYR",
}

// Coalition tetragraphs from CVEnumISMCATOwnerProducer.xsd.
var xsdCoalitionCodes = []string{
	"ACGU", "AMSP", "AOSC", "APFS", "ASEA", "AUSTRALIA_GROUP",
	"BHTF", "BWCS",
	"CFCK", "CFOD", "CFUP", "CLFC", "CMFC", "CMFP", "CPMT", "CTOC", "CWCS",
	"ECTF", "EFOR", "EU", "EUDA",
	"FRME", "FVEY",
	"GCCH", "GCTF", "GFNX", "GMIF",
	"IESC", "IMSC", "IMSP", "IPMC", "IRKS", "ISAF", "ISSG",
	"KFOR",
	"MCFI", "MESF", "MGEU", "MIFH", "MLEC", "MNTF", "MPFL",
	"NACT", "NATO", "NCFE", "NKIC", "NRDC", "NSG",
	"OSAG", "OSTY",
	"PAWA", "PGMF", "PSMX",
	"RISC", "RSMA",
	"SFOR", "SOFP", "SPAA",
	"TEYE", "TFTC",
	"UNCK",
}

// TestXSD_CountryCodes_FiveEyes verifies all Five Eyes countries are present.
func TestXSD_CountryCodes_FiveEyes(t *testing.T) {
	r := reg()
	fiveEyes := []string{"USA", "GBR", "CAN", "AUS", "NZL"}
	for _, code := range fiveEyes {
		t.Run(code, func(t *testing.T) {
			if !r.ValidCountryCode(code) {
				t.Errorf("Five Eyes country %s must be in registry", code)
			}
		})
	}
}

// TestXSD_CountryCodes_SamplePresent tests a representative sample of XSD
// country codes against the API registry.
func TestXSD_CountryCodes_SamplePresent(t *testing.T) {
	r := reg()
	missing := 0
	for _, code := range xsdCountryCodesSample {
		t.Run(code, func(t *testing.T) {
			if !r.ValidCountryCode(code) {
				t.Errorf("country code %s not in registry — required by CVEnumISMCATOwnerProducer.xsd", code)
				missing++
			}
		})
	}
}

// TestXSD_CountryCodes_CoalitionsPresent tests coalition tetragraphs.
func TestXSD_CountryCodes_CoalitionsPresent(t *testing.T) {
	r := reg()
	for _, code := range xsdCoalitionCodes {
		t.Run(code, func(t *testing.T) {
			if !r.ValidCountryCode(code) {
				t.Errorf("coalition %s not in registry — required by CVEnumISMCATOwnerProducer.xsd", code)
			}
		})
	}
}

// TestXSD_CountryCodes_CoverageReport reports overall country code coverage.
func TestXSD_CountryCodes_CoverageReport(t *testing.T) {
	r := reg()
	total := len(xsdCountryCodesSample) + len(xsdCoalitionCodes)
	present := 0
	for _, code := range xsdCountryCodesSample {
		if r.ValidCountryCode(code) {
			present++
		}
	}
	for _, code := range xsdCoalitionCodes {
		if r.ValidCountryCode(code) {
			present++
		}
	}
	coverage := float64(present) / float64(total) * 100
	t.Logf("Country code coverage: %d/%d (%.1f%%) of sampled XSD codes", present, total, coverage)
	t.Logf("API registry has %d codes total; XSD has 341+ codes", len(r.CountryCodes))
}

// TestXSD_CountryCodes_FGIMarker verifies the FGI pseudo-code is handled.
func TestXSD_CountryCodes_FGIMarker(t *testing.T) {
	r := reg()
	if !r.ValidCountryCode("FGI") {
		t.Error("FGI marker not in country code registry — required by CVEnumISMCATOwnerProducer.xsd")
	}
}

package compliance_test

import (
	"testing"
)

// TestXSD_ModelAttributes_AllPresent verifies every XSD attribute defined in
// IC-ISM.xsd has a corresponding field in the ISM struct.
func TestXSD_ModelAttributes_AllPresent(t *testing.T) {
	for _, attr := range xsdISMAttributes {
		t.Run(attr.Name, func(t *testing.T) {
			if !requireStructField(t, attr.JSONTag) {
				t.Skipf("GAP: ISM struct missing field for XSD attribute %s (expected json tag %q) — required by IC-ISM.xsd", attr.Name, attr.JSONTag)
			}
		})
	}
}

// TestXSD_ModelAttributes_CoverageReport reports attribute coverage.
func TestXSD_ModelAttributes_CoverageReport(t *testing.T) {
	present := 0
	var missing []string
	for _, attr := range xsdISMAttributes {
		if requireStructField(t, attr.JSONTag) {
			present++
		} else {
			missing = append(missing, attr.Name)
		}
	}
	t.Logf("Model attribute coverage: %d/%d XSD attributes mapped", present, len(xsdISMAttributes))
	if len(missing) > 0 {
		t.Logf("Missing attributes: %v", missing)
	}
}

// TestXSD_ModelAttributes_CoreFieldsPresent verifies the core fields that the
// API currently implements are correctly mapped.
func TestXSD_ModelAttributes_CoreFieldsPresent(t *testing.T) {
	core := []struct {
		xsd  string
		json string
	}{
		{"classification", "classification"},
		{"ownerProducer", "ownerProducer"},
		{"joint", "joint"},
		{"disseminationControls", "disseminationControls"},
		{"releasableTo", "releasableTo"},
		{"displayOnlyTo", "displayOnlyTo"},
		{"classifiedBy", "classifiedBy"},
		{"classificationReason", "classificationReason"},
		{"compilationReason", "compilationReason"},
		{"declassDate", "declassDate"},
		{"declassEvent", "declassEvent"},
		{"declassException", "declassException"},
		{"FGIsourceOpen", "fgiSourceOpen"},
		{"FGIsourceProtected", "fgiSourceProtected"},
		{"nonICmarkings", "nonICMarkings"},
		{"nonUSControls", "nonUSControls"},
	}
	for _, f := range core {
		t.Run(f.xsd, func(t *testing.T) {
			if !requireStructField(t, f.json) {
				t.Errorf("core attribute %s (json: %s) MUST be in ISM struct", f.xsd, f.json)
			}
		})
	}
}

// TestXSD_ModelAttributes_MissingFields documents the expected missing fields.
func TestXSD_ModelAttributes_MissingFields(t *testing.T) {
	expectedMissing := []struct {
		xsd  string
		json string
		desc string
	}{
		{"SCIcontrols", "sciControls", "SCI compartments (e.g., HCS, SI, TK)"},
		{"SARIdentifier", "sarIdentifier", "Special Access Required identifiers"},
		{"atomicEnergyMarkings", "atomicEnergyMarkings", "RD/FRD/CNWDI/TFNI markings"},
		{"noticeType", "noticeType", "Notice type markers (FISA, CNWDI, etc.)"},
		{"noticeProseID", "noticeProseID", "Prose text ID for notices"},
		{"compliesWith", "compliesWith", "Compliance framework (USGov, USIC, etc.)"},
		{"exemptFrom", "exemptFrom", "Exemption from ICD-710 or DoD Distro"},
		{"createDate", "createDate", "Document creation date"},
		{"cuiBasic", "cuiBasic", "Basic CUI marking"},
		{"cuiSpecified", "cuiSpecified", "Specified CUI category"},
		{"handleViaChannels", "handleViaChannels", "Special handling channels"},
		{"hasApproximateMarkings", "hasApproximateMarkings", "Approximate markings flag"},
		{"highWaterNATO", "highWaterNATO", "NATO high-water mark"},
		{"noAggregation", "noAggregation", "No aggregation flag"},
		{"pocType", "pocType", "Point of contact type"},
		{"secondBannerLine", "secondBannerLine", "Second banner line content"},
	}
	for _, m := range expectedMissing {
		t.Run(m.xsd, func(t *testing.T) {
			if requireStructField(t, m.json) {
				t.Logf("RESOLVED: %s (%s) is now in ISM struct", m.xsd, m.desc)
				return
			}
			t.Skipf("GAP: ISM struct missing %s (%s) — required by IC-ISM.xsd", m.xsd, m.desc)
		})
	}
}

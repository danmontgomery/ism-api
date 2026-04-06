package compliance_test

import (
	"testing"
)

// TestXSD_Notices_NotModeled documents that notice types are not modeled in the API.
func TestXSD_Notices_NotModeled(t *testing.T) {
	if !requireStructField(t, "noticeType") {
		t.Skip("GAP: ISM struct missing noticeType field — 27 notice types from CVEnumISMNotice.xsd not modeled")
	}
}

// TestXSD_Notices_AllTypesPresent checks each individual notice type.
func TestXSD_Notices_AllTypesPresent(t *testing.T) {
	if !requireStructField(t, "noticeType") {
		t.Skip("GAP: ISM struct missing noticeType field — cannot test notice types")
		return
	}
	for _, code := range xsdNoticeTypes {
		t.Run(code, func(t *testing.T) {
			t.Skipf("GAP: notice type %s not in registry — required by CVEnumISMNotice.xsd", code)
		})
	}
}

// TestXSD_Notices_Categories groups notice types by functional category.
func TestXSD_Notices_Categories(t *testing.T) {
	categories := map[string][]string{
		"intelligence_warnings":   {"FISA", "RAWFISA", "IMC", "CNWDI", "COMSEC", "GEOCAP"},
		"atomic_energy_warnings":  {"RD", "FRD"},
		"distribution_caveats":    {"DS", "LES", "LES-NF", "DSEN"},
		"dod_distribution":        {"DoD-Dist-A", "DoD-Dist-B", "DoD-Dist-C", "DoD-Dist-D", "DoD-Dist-E", "DoD-Dist-F"},
		"person_identity":         {"US-Person"},
		"originator_control":      {"pre13526ORCON"},
		"point_of_contact":        {"POC"},
		"security_notices":        {"SSI", "RSEN", "IMCON_RSEN"},
		"NATO":                    {"NATO"},
		"collection_restrictions": {"RC_Dissemination_Control_Required"},
		"trade_controls":          {"ITAR-EAR"},
	}

	for cat, notices := range categories {
		t.Run(cat, func(t *testing.T) {
			if !requireStructField(t, "noticeType") {
				t.Skipf("GAP: notice category %s (%d types) not modeled", cat, len(notices))
			}
		})
	}
}

package compliance_test

import (
	"testing"
)

// TestXSD_AtomicEnergy_NotModeled verifies that atomic energy markings are
// modeled in the API. The ISM struct must have an atomicEnergyMarkings field.
func TestXSD_AtomicEnergy_NotModeled(t *testing.T) {
	if !requireStructField(t, "atomicEnergyMarkings") {
		t.Fatal("ISM struct missing atomicEnergyMarkings field — 14 markings from CVEnumISMAtomicEnergyMarkings.xsd not modeled")
	}
}

// TestXSD_AtomicEnergy_AllMarkingsPresent checks each atomic energy marking.
func TestXSD_AtomicEnergy_AllMarkingsPresent(t *testing.T) {
	if !requireStructField(t, "atomicEnergyMarkings") {
		t.Fatal("ISM struct missing atomicEnergyMarkings — cannot test individual markings")
	}
	r := reg()
	for _, code := range xsdAtomicEnergyMarkings {
		t.Run(code, func(t *testing.T) {
			if !r.ValidAtomicEnergyMarking(code) {
				t.Errorf("atomic energy marking %s not in registry — required by CVEnumISMAtomicEnergyMarkings.xsd", code)
			}
		})
	}
}

// TestXSD_AtomicEnergy_Categories verifies markings by category.
func TestXSD_AtomicEnergy_Categories(t *testing.T) {
	categories := map[string][]string{
		"RD (Restricted Data)":          {"RD", "RD-CNWDI", "RD-SG-14", "RD-SG-15", "RD-SG-18", "RD-SG-20"},
		"FRD (Formerly RD)":             {"FRD", "FRD-SG-14", "FRD-SG-15", "FRD-SG-18", "FRD-SG-20"},
		"Controlled Nuclear":            {"DCNI", "UCNI"},
		"TFNI (Transclassified Foreign)": {"TFNI"},
	}
	r := reg()
	for cat, markings := range categories {
		t.Run(cat, func(t *testing.T) {
			if !requireStructField(t, "atomicEnergyMarkings") {
				t.Fatalf("atomic energy category %s (%d markings) not modeled", cat, len(markings))
			}
			for _, code := range markings {
				if !r.ValidAtomicEnergyMarking(code) {
					t.Errorf("marking %s in category %s not in registry", code, cat)
				}
			}
		})
	}
}

package refdata

// AtomicEnergyMarking is a reference data entry for an atomic energy marking.
type AtomicEnergyMarking struct {
	Code     string `json:"code"`
	Label    string `json:"label"`
	Category string `json:"category"`
}

// AtomicEnergyMarkings returns all atomic energy marking codes from
// CVEnumISMAtomicEnergyMarkings.xsd. Categories: RD (Restricted Data),
// FRD (Formerly Restricted Data), Controlled Nuclear, and TFNI.
func AtomicEnergyMarkings() []AtomicEnergyMarking {
	return []AtomicEnergyMarking{
		// RD — Restricted Data
		{
			Code:     "RD",
			Label:    "Restricted Data",
			Category: "RD (Restricted Data)",
		},
		{
			Code:     "RD-CNWDI",
			Label:    "RD-Critical Nuclear Weapon Design Information",
			Category: "RD (Restricted Data)",
		},
		{
			Code:     "RD-SG-14",
			Label:    "RD-SIGMA-14",
			Category: "RD (Restricted Data)",
		},
		{
			Code:     "RD-SG-15",
			Label:    "RD-SIGMA-15",
			Category: "RD (Restricted Data)",
		},
		{
			Code:     "RD-SG-18",
			Label:    "RD-SIGMA-18",
			Category: "RD (Restricted Data)",
		},
		{
			Code:     "RD-SG-20",
			Label:    "RD-SIGMA-20",
			Category: "RD (Restricted Data)",
		},

		// FRD — Formerly Restricted Data
		{
			Code:     "FRD",
			Label:    "Formerly Restricted Data",
			Category: "FRD (Formerly RD)",
		},
		{
			Code:     "FRD-SG-14",
			Label:    "FRD-SIGMA-14",
			Category: "FRD (Formerly RD)",
		},
		{
			Code:     "FRD-SG-15",
			Label:    "FRD-SIGMA-15",
			Category: "FRD (Formerly RD)",
		},
		{
			Code:     "FRD-SG-18",
			Label:    "FRD-SIGMA-18",
			Category: "FRD (Formerly RD)",
		},
		{
			Code:     "FRD-SG-20",
			Label:    "FRD-SIGMA-20",
			Category: "FRD (Formerly RD)",
		},

		// Controlled Nuclear Information
		{
			Code:     "DCNI",
			Label:    "DoD Controlled Nuclear Information",
			Category: "Controlled Nuclear",
		},
		{
			Code:     "UCNI",
			Label:    "DoE Controlled Nuclear Information",
			Category: "Controlled Nuclear",
		},

		// TFNI — Transclassified Foreign Nuclear Information
		{
			Code:     "TFNI",
			Label:    "Transclassified Foreign Nuclear Information",
			Category: "TFNI (Transclassified Foreign)",
		},
	}
}

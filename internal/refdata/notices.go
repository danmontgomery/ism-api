package refdata

// NoticeTypes returns all notice type codes from CVEnumISMNotice.xsd.
// Categories: intelligence warnings, atomic energy warnings, distribution
// caveats, DoD distribution, person identity, originator control, point of
// contact, security notices, NATO, collection restrictions, and trade controls.
func NoticeTypes() []NoticeType {
	return []NoticeType{
		// Intelligence Warnings
		{
			Code:     "FISA",
			Label:    "FISA Warning Statement",
			Category: "intelligence_warnings",
		},
		{
			Code:     "RAWFISA",
			Label:    "RAWFISA Warning Statement",
			Category: "intelligence_warnings",
		},
		{
			Code:     "IMC",
			Label:    "IMCON Warning Statement",
			Category: "intelligence_warnings",
		},
		{
			Code:     "CNWDI",
			Label:    "Critical Nuclear Weapon Design Information Warning Statement",
			Category: "intelligence_warnings",
		},
		{
			Code:     "COMSEC",
			Label:    "COMSEC Notice",
			Category: "intelligence_warnings",
		},
		{
			Code:     "GEOCAP",
			Label:    "GEOCAP Warning Statement",
			Category: "intelligence_warnings",
		},

		// Atomic Energy Warnings
		{
			Code:     "RD",
			Label:    "RD Warning Statement",
			Category: "atomic_energy_warnings",
		},
		{
			Code:     "FRD",
			Label:    "FRD Warning Statement",
			Category: "atomic_energy_warnings",
		},

		// Distribution Caveats
		{
			Code:     "DS",
			Label:    "LIMDIS Caveat",
			Category: "distribution_caveats",
		},
		{
			Code:     "LES",
			Label:    "LES Notice",
			Category: "distribution_caveats",
		},
		{
			Code:     "LES-NF",
			Label:    "LES-NF Notice",
			Category: "distribution_caveats",
		},
		{
			Code:     "DSEN",
			Label:    "DSEN Notice",
			Category: "distribution_caveats",
		},

		// DoD Distribution Statements (DoDI 5230.24)
		{
			Code:     "DoD-Dist-A",
			Label:    "DoD Distribution Statement A",
			Category: "dod_distribution",
		},
		{
			Code:     "DoD-Dist-B",
			Label:    "DoD Distribution Statement B",
			Category: "dod_distribution",
		},
		{
			Code:     "DoD-Dist-C",
			Label:    "DoD Distribution Statement C",
			Category: "dod_distribution",
		},
		{
			Code:     "DoD-Dist-D",
			Label:    "DoD Distribution Statement D",
			Category: "dod_distribution",
		},
		{
			Code:     "DoD-Dist-E",
			Label:    "DoD Distribution Statement E",
			Category: "dod_distribution",
		},
		{
			Code:     "DoD-Dist-F",
			Label:    "DoD Distribution Statement F",
			Category: "dod_distribution",
		},

		// Person Identity
		{
			Code:     "US-Person",
			Label:    "US Person Information Notice",
			Category: "person_identity",
		},

		// Originator Control
		{
			Code:     "pre13526ORCON",
			Label:    "Pre-Executive Order 13526 ORCON",
			Category: "originator_control",
		},

		// Point of Contact
		{
			Code:     "POC",
			Label:    "Point of Contact Notice",
			Category: "point_of_contact",
		},

		// Security Notices
		{
			Code:     "SSI",
			Label:    "Sensitive Security Information",
			Category: "security_notices",
		},
		{
			Code:     "RSEN",
			Label:    "Risk Sensitive Notice",
			Category: "security_notices",
		},
		{
			Code:     "IMCON_RSEN",
			Label:    "IMCON and RSEN Warning Statement",
			Category: "security_notices",
		},

		// NATO
		{
			Code:     "NATO",
			Label:    "NATO Warning Statement",
			Category: "NATO",
		},

		// Collection Restrictions
		{
			Code:     "RC_Dissemination_Control_Required",
			Label:    "Restricted Collection Dissemination Control Required",
			Category: "collection_restrictions",
		},

		// Trade Controls
		{
			Code:     "ITAR-EAR",
			Label:    "ITAR/EAR Warning Statement",
			Category: "trade_controls",
		},
	}
}

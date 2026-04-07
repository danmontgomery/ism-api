package refdata

// ExemptFrom returns all valid exemptFrom values per CVEnumISMExemptFrom.xsd.
func ExemptFrom() []ExemptFromEntry {
	return []ExemptFromEntry{
		{
			Code:        "IC_710_MANDATORY_FDR",
			Label:       "ICD-710 Mandatory FD&R",
			Description: "Exemption from ICD-710 rules mandating Foreign Disclosure and Release markings",
		},
		{
			Code:        "DOD_DISTRO_STATEMENT",
			Label:       "DoD Distribution Statement",
			Description: "Exemption from DoD 5230.24 distribution statement requirements",
		},
	}
}

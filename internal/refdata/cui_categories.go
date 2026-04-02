package refdata

// CUICategories returns all supported CUI category markings.
// Specified categories (SP- prefix) require specific handling procedures defined
// by the authorizing law, regulation, or government-wide policy.
// Basic categories follow default CUI handling procedures.
func CUICategories() []CUICategory {
	return []CUICategory{
		// Specified categories (SP- prefix) — ordered alphabetically
		{
			Code:        "SP-CEII",
			Label:       "Critical Energy Infrastructure Information",
			Type:        "specified",
			Description: "Information about critical energy infrastructure that could be useful to a person planning an attack on critical infrastructure (18 CFR 388.113)",
		},
		{
			Code:        "SP-CRITAN",
			Label:       "Critical Analysis",
			Type:        "specified",
			Description: "Analysis of critical infrastructure vulnerabilities and risks produced by DHS critical infrastructure analysis programs",
		},
		{
			Code:        "SP-CVI",
			Label:       "Chemical-terrorism Vulnerability Information",
			Type:        "specified",
			Description: "Information developed under the Chemical Facility Anti-Terrorism Standards (6 CFR Part 27)",
		},
		{
			Code:        "SP-PCII",
			Label:       "Protected Critical Infrastructure Information",
			Type:        "specified",
			Description: "Voluntarily submitted critical infrastructure information validated under the PCII Program (6 CFR Part 29)",
		},
		{
			Code:        "SP-PHYS",
			Label:       "Physical Security (Specified)",
			Type:        "specified",
			Description: "Specified physical security information requiring enhanced handling under authorizing directive",
		},
		{
			Code:        "SP-TSCA",
			Label:       "Toxic Substances Control Act",
			Type:        "specified",
			Description: "Chemical data submitted under TSCA with restrictions on disclosure (15 U.S.C. §2613)",
		},

		// Basic categories — ordered alphabetically
		{
			Code:        "CRIT",
			Label:       "Critical Infrastructure",
			Type:        "basic",
			Description: "General critical infrastructure information not covered by a specified CUI category",
		},
		{
			Code:        "EMGT",
			Label:       "Emergency Management",
			Type:        "basic",
			Description: "Information related to emergency management planning, response, or recovery operations",
		},
		{
			Code:        "ISVI",
			Label:       "Infrastructure Security Vulnerability Information",
			Type:        "basic",
			Description: "Information about vulnerabilities in critical infrastructure systems",
		},
		{
			Code:        "PHYS",
			Label:       "Physical Security",
			Type:        "basic",
			Description: "Information related to physical security measures, plans, or assessments",
		},
		{
			Code:        "SAFE",
			Label:       "Safety",
			Type:        "basic",
			Description: "Information related to safety conditions, incidents, or protective measures",
		},
		{
			Code:        "WATER",
			Label:       "Water Assessments",
			Type:        "basic",
			Description: "Information about water system vulnerability assessments and security measures",
		},
	}
}

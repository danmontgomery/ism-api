package refdata

// DeclassExceptions returns all supported declassification exemption codes
// from CVEnumISM25X.xsd. These codes extend the standard 10-year or 25-year
// declassification timeframes when specific conditions apply.
func DeclassExceptions() []DeclassException {
	return []DeclassException{
		// Special exemption codes
		{
			Code:     "AEA",
			Label:    "Atomic Energy Act",
			Category: "nuclear",
		},
		{
			Code:     "NATO",
			Label:    "NATO Treaty",
			Category: "foreign",
		},
		{
			Code:     "NATO-AEA",
			Label:    "NATO and Atomic Energy Act",
			Category: "foreign",
		},

		// 25X codes — information exempt from automatic declassification at 25 years
		{
			Code:     "25X1",
			Label:    "Human Intelligence Sources",
			Category: "intelligence",
		},
		{
			Code:     "25X1-EO-12951",
			Label:    "25X1 with Executive Order 12951",
			Category: "intelligence",
		},
		{
			Code:     "25X2",
			Label:    "Weapons of Mass Destruction",
			Category: "defense",
		},
		{
			Code:     "25X3",
			Label:    "Intelligence Activities, Sources, or Methods",
			Category: "intelligence",
		},
		{
			Code:     "25X4",
			Label:    "Intelligence Sources or Methods Under DNI Protection",
			Category: "intelligence",
		},
		{
			Code:     "25X5",
			Label:    "Foreign Government Information",
			Category: "foreign",
		},
		{
			Code:     "25X6",
			Label:    "Activities Related to Foreign Relations or Foreign Activities",
			Category: "foreign",
		},
		{
			Code:     "25X7",
			Label:    "Scientific, Technological, or Economic Matters Related to National Security",
			Category: "defense",
		},
		{
			Code:     "25X8",
			Label:    "Vulnerabilities or Capabilities of Systems, Installations, or Projects",
			Category: "defense",
		},
		{
			Code:     "25X9",
			Label:    "Violations of Federal Law or Executive Privilege",
			Category: "legal",
		},

		// 50X codes — information exempt from automatic declassification at 50 years
		{
			Code:     "50X1",
			Label:    "Human Intelligence Sources (50 Year Base)",
			Category: "intelligence",
		},
		{
			Code:     "50X1-HUM",
			Label:    "Human Intelligence Sources (50 Year)",
			Category: "intelligence",
		},
		{
			Code:     "50X2",
			Label:    "Weapons of Mass Destruction (50 Year Base)",
			Category: "defense",
		},
		{
			Code:     "50X2-WMD",
			Label:    "Weapons of Mass Destruction (50 Year)",
			Category: "defense",
		},
		{
			Code:     "50X3",
			Label:    "Intelligence Activities, Sources, or Methods (50 Year)",
			Category: "intelligence",
		},
		{
			Code:     "50X4",
			Label:    "Intelligence Sources or Methods Under DNI Protection (50 Year)",
			Category: "intelligence",
		},
		{
			Code:     "50X5",
			Label:    "Foreign Government Information (50 Year)",
			Category: "foreign",
		},
		{
			Code:     "50X6",
			Label:    "Foreign Relations or Foreign Activities (50 Year)",
			Category: "foreign",
		},
		{
			Code:     "50X7",
			Label:    "Scientific, Technological, or Economic Matters (50 Year)",
			Category: "defense",
		},
		{
			Code:     "50X8",
			Label:    "Vulnerabilities or Capabilities of Systems (50 Year)",
			Category: "defense",
		},
		{
			Code:     "50X9",
			Label:    "Violations of Federal Law or Executive Privilege (50 Year)",
			Category: "legal",
		},

		// 75X — ISCAP-approved exemption beyond 75 years
		{
			Code:     "75X",
			Label:    "ISCAP-Approved Specific Information",
			Category: "special",
		},
	}
}

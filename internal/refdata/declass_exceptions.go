package refdata

// DeclassExceptions returns all supported 25X declassification exemption codes
// per EO 13526 Section 3.3(b). These codes extend the standard 10-year or 25-year
// declassification timeframes when specific conditions apply.
func DeclassExceptions() []DeclassException {
	return []DeclassException{
		// 25X codes — information exempt from automatic declassification at 25 years
		{
			Code:     "25X1",
			Label:    "Human Intelligence Sources",
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
			Code:     "50X1-HUM",
			Label:    "Human Intelligence Sources (50 Year)",
			Category: "intelligence",
		},
		{
			Code:     "50X2-WMD",
			Label:    "Weapons of Mass Destruction (50 Year)",
			Category: "defense",
		},
	}
}

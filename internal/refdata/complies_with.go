package refdata

// CompliesWith returns all valid compliesWith values per CVEnumISMCompliesWith.xsd.
func CompliesWith() []CompliesWithEntry {
	return []CompliesWithEntry{
		{Code: "USGov", Label: "US Government"},
		{Code: "USIC", Label: "US Intelligence Community"},
		{Code: "USDOD", Label: "US Department of Defense"},
		{Code: "OtherAuthority", Label: "Other Authority"},
		{Code: "USA-CUI-ONLY", Label: "CUI Only"},
		{Code: "USA-CUI", Label: "CUI"},
	}
}

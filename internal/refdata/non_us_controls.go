package refdata

// NonUSControls returns all valid non-US control markings per CVEnumISMNonUSControls.xsd.
func NonUSControls() []NonUSControl {
	return []NonUSControl{
		{Code: "NATO-ATOMAL", Label: "NATO ATOMAL"},
		{Code: "NATO-BOHEMIA", Label: "NATO BOHEMIA"},
		{Code: "NATO-BALK", Label: "NATO BALK"},
	}
}

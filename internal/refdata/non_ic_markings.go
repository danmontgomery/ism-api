package refdata

// NonICMarkings returns all supported non-Intelligence Community markings.
// These markings originate outside the IC but may appear on IC documents.
func NonICMarkings() []NonICMarking {
	return []NonICMarking{
		{Code: "SBU", Label: "Sensitive But Unclassified"},
		{Code: "SBU-NF", Label: "Sensitive But Unclassified — NOFORN"},
		{Code: "LES", Label: "Law Enforcement Sensitive"},
		{Code: "LES-NF", Label: "Law Enforcement Sensitive — NOFORN"},
		{Code: "SSI", Label: "Sensitive Security Information"},
		{Code: "FOUO", Label: "For Official Use Only"},
		{Code: "DS", Label: "Limited Distribution"},
		{Code: "XD", Label: "Exclusively for Administrative or Operational Purposes"},
		{Code: "ND", Label: "No Distribution"},
		{Code: "NNPI", Label: "Naval Nuclear Propulsion Information"},
	}
}

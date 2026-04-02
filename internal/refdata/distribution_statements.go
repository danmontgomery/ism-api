package refdata

// DistributionStatements returns all supported distribution statements per DoDI 5230.24.
func DistributionStatements() []DistributionStatement {
	return []DistributionStatement{
		{
			Code:                     "A",
			Label:                    "Statement A — Public Release",
			Text:                     "Approved for public release; distribution is unlimited.",
			ClassificationConstraint: "U,CUI",
		},
		{
			Code:  "B",
			Label: "Statement B — U.S. Government Only",
			Text:  "Distribution authorized to U.S. Government agencies only.",
		},
		{
			Code:  "C",
			Label: "Statement C — Government Agencies and Contractors",
			Text:  "Distribution authorized to U.S. Government agencies and their contractors.",
		},
		{
			Code:  "D",
			Label: "Statement D — DoD and U.S. DoD Contractors",
			Text:  "Distribution authorized to the Department of Defense and U.S. DoD contractors only.",
		},
		{
			Code:  "E",
			Label: "Statement E — DoD Components",
			Text:  "Distribution authorized to DoD components only.",
		},
		{
			Code:  "F",
			Label: "Statement F — Further Dissemination Only as Directed",
			Text:  "Further dissemination only as directed by the controlling DoD office.",
		},
	}
}

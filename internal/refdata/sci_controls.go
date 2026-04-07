package refdata

// SCIControl is a reference data entry for an SCI (Sensitive Compartmented
// Information) control marking.
type SCIControl struct {
	Code     string `json:"code"`
	Label    string `json:"label"`
	Category string `json:"category"`
}

// SCIControls returns all SCI control codes from CVEnumISMSCIControls.xsd.
// Categories: BUR (BYELEMAN), HCS (HUMINT), KLM (KLAMATH),
// SI (SPECIAL INTELLIGENCE), TK (TALENT KEYHOLE), MVL (MARVEL), RSV (RESERVE).
func SCIControls() []SCIControl {
	return []SCIControl{
		// BUR — BYELEMAN
		{
			Code:     "BUR",
			Label:    "BUR",
			Category: "BUR (BYELEMAN)",
		},
		{
			Code:     "BUR-BLG",
			Label:    "BUR-BLG",
			Category: "BUR (BYELEMAN)",
		},
		{
			Code:     "BUR-DTP",
			Label:    "BUR-DTP",
			Category: "BUR (BYELEMAN)",
		},
		{
			Code:     "BUR-WRG",
			Label:    "BUR-WRG",
			Category: "BUR (BYELEMAN)",
		},

		// HCS — HUMINT Control System
		{
			Code:     "HCS",
			Label:    "HCS",
			Category: "HCS (HUMINT)",
		},
		{
			Code:     "HCS-O",
			Label:    "HCS-O",
			Category: "HCS (HUMINT)",
		},
		{
			Code:     "HCS-P",
			Label:    "HCS-P",
			Category: "HCS (HUMINT)",
		},
		{
			Code:     "HCS-X",
			Label:    "HCS-X",
			Category: "HCS (HUMINT)",
		},

		// KLM — KLAMATH
		{
			Code:     "KLM",
			Label:    "KLAMATH",
			Category: "KLM (KLAMATH)",
		},
		{
			Code:     "KLM-R",
			Label:    "KLAMATH-R",
			Category: "KLM (KLAMATH)",
		},

		// MVL — MARVEL
		{
			Code:     "MVL",
			Label:    "MARVEL",
			Category: "MVL (MARVEL)",
		},

		// RSV — RESERVE
		{
			Code:     "RSV",
			Label:    "RESERVE",
			Category: "RSV (RESERVE)",
		},

		// SI — SPECIAL INTELLIGENCE
		{
			Code:     "SI",
			Label:    "SPECIAL INTELLIGENCE",
			Category: "SI (SPECIAL INTELLIGENCE)",
		},
		{
			Code:     "SI-EU",
			Label:    "ECRU",
			Category: "SI (SPECIAL INTELLIGENCE)",
		},
		{
			Code:     "SI-G",
			Label:    "SI-GAMMA",
			Category: "SI (SPECIAL INTELLIGENCE)",
		},
		{
			Code:     "SI-NK",
			Label:    "NONBOOK",
			Category: "SI (SPECIAL INTELLIGENCE)",
		},

		// TK — TALENT KEYHOLE
		{
			Code:     "TK",
			Label:    "TALENT KEYHOLE",
			Category: "TK (TALENT KEYHOLE)",
		},
		{
			Code:     "TK-BLFH",
			Label:    "BLUEFISH",
			Category: "TK (TALENT KEYHOLE)",
		},
		{
			Code:     "TK-IDIT",
			Label:    "IDITAROD",
			Category: "TK (TALENT KEYHOLE)",
		},
		{
			Code:     "TK-KAND",
			Label:    "KANDIK",
			Category: "TK (TALENT KEYHOLE)",
		},
	}
}

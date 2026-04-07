package refdata

import "expr.ai/ism-api/internal/model"

// DisseminationControls returns all supported dissemination control markings
// with their compatibility metadata for validation and UI guidance.
func DisseminationControls() []DisseminationControl {
	return []DisseminationControl{
		{
			Code:        "RS",
			Label:       "Risk Sensitive",
			Description: "Information requiring risk-sensitive handling",
		},
		{
			Code:          "OC",
			Label:         "ORCON — Originator Controlled",
			Description:   "Dissemination and extraction of information controlled by originator",
			MinClassification: model.ClassificationC,
		},
		{
			Code:          "OC-USGOV",
			Label:         "ORCON — US Government Only",
			Description:   "Originator controlled, restricted to US Government personnel",
			MinClassification: model.ClassificationC,
		},
		{
			Code:          "IMCON",
			Label:         "Controlled Imagery",
			Description:   "Imagery dissemination controlled by the imagery originator",
			MinClassification: model.ClassificationC,
		},
		{
			Code:          "NOFORN",
			Label:         "Not Releasable to Foreign Nationals",
			Description:   "Information may not be released to foreign nationals",
			MinClassification: model.ClassificationC,
			ExclusiveWith: []string{"REL", "RELIDO"},
		},
		{
			Code:        "PROPIN",
			Label:       "Proprietary Information Involved",
			Description: "Contains proprietary information",
		},
		{
			Code:          "REL",
			Label:         "Authorized for Release To",
			Description:   "Information approved for release to specified countries or organizations",
			RequiresField: "releasableTo",
			ExclusiveWith: []string{"NOFORN"},
		},
		{
			Code:        "RELIDO",
			Label:       "Releasable by Information Disclosure Official",
			Description: "Release determination delegated to an Information Disclosure Official",
			ExclusiveWith: []string{"NOFORN"},
		},
		{
			Code:          "EYES",
			Label:         "Eyes Only",
			Description:   "Information restricted to named countries (FVEY, etc.)",
			RequiresField: "releasableTo",
			MinClassification: model.ClassificationC,
		},
		{
			Code:          "DSEN",
			Label:         "DEA Sensitive",
			Description:   "Drug Enforcement Administration sensitive information",
			MinClassification: model.ClassificationC,
		},
		{
			Code:          "FISA",
			Label:         "Foreign Intelligence Surveillance Act",
			Description:   "Information obtained or derived under FISA authority",
			MinClassification: model.ClassificationC,
		},
		{
			Code:          "DISPLAY ONLY",
			Label:         "Display Only",
			Description:   "Information may only be displayed to specified countries (no physical transfer)",
			RequiresField: "displayOnlyTo",
		},
		{
			Code:        "FED ONLY",
			Label:       "Federal Only",
			Description: "Distribution limited to federal government entities",
		},
		{
			Code:        "FEDCON",
			Label:       "Federal and Contractor",
			Description: "Distribution limited to federal government and authorized contractors",
		},
		{
			Code:          "NOCON",
			Label:         "Not Releasable to Contractors/Consultants",
			Description:   "Information may not be released to contractors or consultants",
			MinClassification: model.ClassificationC,
		},
		{
			Code:        "DL ONLY",
			Label:       "Distribution List Only",
			Description: "Distribution limited to specified distribution list",
		},
		{
			Code:        "REL TO USA",
			Label:       "Releasable to USA",
			Description: "FGI information authorized for release to USA",
		},
		{
			Code:        "LIST",
			Label:       "Distribution List",
			Description: "Distribution limited to a specific named list",
		},
		{
			Code:              "RAWFISA",
			Label:             "Raw FISA",
			Description:       "Raw data obtained under Foreign Intelligence Surveillance Act authority",
			MinClassification: model.ClassificationC,
		},
		{
			Code:        "FOUO",
			Label:       "For Official Use Only",
			Description: "Information not warranting classification but requiring controlled dissemination",
		},
		{
			Code:        "WAIVED",
			Label:       "Waived",
			Description: "Dissemination restriction has been waived by the originator",
		},
		{
			Code:        "AC",
			Label:       "Attorney-Client Privilege",
			Description: "Information protected by attorney-client privilege",
		},
		{
			Code:        "AWP",
			Label:       "Attorney Work Product",
			Description: "Information protected as attorney work product",
		},
		{
			Code:        "EXEMPT FROM ICD501 DISCOVERY",
			Label:       "Exempt from ICD 501 Discovery",
			Description: "Information exempt from discovery provisions of Intelligence Community Directive 501",
		},
		{
			Code:              "SCI",
			Label:             "Sensitive Compartmented Information",
			Description:       "Information requiring SCI access controls; selecting this enables the SCI controls picker",
			RequiresField:     "sciControls",
			MinClassification: model.ClassificationTS,
		},
	}
}

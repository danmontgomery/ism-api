package compliance_test

import (
	"reflect"
	"testing"

	"expr.ai/ism-api/internal/model"
	"expr.ai/ism-api/internal/refdata"
)

// reg returns a fresh registry for each test.
func reg() *refdata.Registry {
	return refdata.NewRegistry()
}

// assertRegistryContains checks that every code in xsdCodes passes validFunc.
// Codes that fail are logged with the given xsdSource.
// Returns the number of missing codes.
func assertRegistryContains(t *testing.T, validFunc func(string) bool, xsdCodes []string, xsdSource string) int {
	t.Helper()
	missing := 0
	for _, code := range xsdCodes {
		if !validFunc(code) {
			t.Logf("GAP: %s not in registry — required by %s", code, xsdSource)
			missing++
		}
	}
	return missing
}

// assertClassificationRegistryContains is like assertRegistryContains but for
// Classification-typed codes.
func assertClassificationRegistryContains(t *testing.T, reg *refdata.Registry, xsdCodes []string, xsdSource string) int {
	t.Helper()
	missing := 0
	for _, code := range xsdCodes {
		if !reg.ValidClassification(model.Classification(code)) {
			t.Logf("GAP: %s not in registry — required by %s", code, xsdSource)
			missing++
		}
	}
	return missing
}

// requireStructField uses reflection to verify the ISM struct has a field with
// the given JSON tag name. Returns true if found.
func requireStructField(t *testing.T, fieldName string) bool {
	t.Helper()
	ismType := reflect.TypeOf(model.ISM{})
	for i := 0; i < ismType.NumField(); i++ {
		f := ismType.Field(i)
		tag := f.Tag.Get("json")
		// strip ",omitempty" etc.
		if idx := len(tag); idx > 0 {
			for j := 0; j < len(tag); j++ {
				if tag[j] == ',' {
					tag = tag[:j]
					break
				}
			}
		}
		if tag == fieldName {
			return true
		}
	}
	return false
}

// XSD enum constants — single source of truth extracted from XSD files.

// XSD Classification levels from CVEnumISMClassificationAll.xsd
var xsdClassifications = []string{"U", "R", "C", "S", "TS"}

// XSD Dissemination controls from CVEnumISMDissem.xsd
var xsdDisseminationControls = []string{
	"RS", "FOUO", "OC", "OC-USGOV", "IMC", "NF", "PR", "REL", "RELIDO",
	"EYES", "DSEN", "RAWFISA", "FISA", "DISPLAYONLY",
	"EXEMPT_FROM_ICD501_DISCOVERY", "WAIVED", "AC", "AWP", "DL_ONLY",
	"FED_ONLY", "FEDCON", "NOCON",
}

// Mapping from XSD dissemination codes to API codes.
// The XSD uses abbreviated codes while the API uses human-readable names.
var xsdToAPIDissemination = map[string]string{
	"RS":                            "RS",
	"FOUO":                          "FOUO",
	"OC":                            "OC",
	"OC-USGOV":                      "OC-USGOV",
	"IMC":                           "IMCON",
	"NF":                            "NOFORN",
	"PR":                            "PROPIN",
	"REL":                           "REL",
	"RELIDO":                        "RELIDO",
	"EYES":                          "EYES",
	"DSEN":                          "DSEN",
	"RAWFISA":                       "RAWFISA",
	"FISA":                          "FISA",
	"DISPLAYONLY":                   "DISPLAY ONLY",
	"EXEMPT_FROM_ICD501_DISCOVERY":  "EXEMPT FROM ICD501 DISCOVERY",
	"WAIVED":                        "WAIVED",
	"AC":                            "AC",
	"AWP":                           "AWP",
	"DL_ONLY":                       "DL ONLY",
	"FED_ONLY":                      "FED ONLY",
	"FEDCON":                        "FEDCON",
	"NOCON":                         "NOCON",
}

// XSD SCI controls from CVEnumISMSCIControls.xsd
var xsdSCIControls = []string{
	"BUR", "BUR-BLG", "BUR-DTP", "BUR-WRG",
	"HCS", "HCS-O", "HCS-P", "HCS-X",
	"KLM", "KLM-R", "MVL", "RSV",
	"SI", "SI-EU", "SI-G", "SI-NK",
	"TK", "TK-BLFH", "TK-IDIT", "TK-KAND",
}

// XSD Notice types from CVEnumISMNotice.xsd
var xsdNoticeTypes = []string{
	"FISA", "RAWFISA", "IMC", "CNWDI", "RD", "FRD", "DS", "LES", "LES-NF",
	"DSEN", "DoD-Dist-A", "DoD-Dist-B", "DoD-Dist-C", "DoD-Dist-D",
	"DoD-Dist-E", "DoD-Dist-F", "US-Person", "pre13526ORCON", "POC",
	"COMSEC", "SSI", "RSEN", "IMCON_RSEN", "GEOCAP", "NATO",
	"RC_Dissemination_Control_Required", "ITAR-EAR",
}

// XSD Atomic energy markings from CVEnumISMAtomicEnergyMarkings.xsd
var xsdAtomicEnergyMarkings = []string{
	"RD", "RD-CNWDI", "RD-SG-14", "RD-SG-15", "RD-SG-18", "RD-SG-20",
	"FRD", "FRD-SG-14", "FRD-SG-15", "FRD-SG-18", "FRD-SG-20",
	"DCNI", "UCNI", "TFNI",
}

// XSD Non-US controls from CVEnumISMNonUSControls.xsd
var xsdNonUSControls = []string{
	"NATO-ATOMAL", "NATO-BOHEMIA", "NATO-BALK",
}

// XSD Non-IC markings from CVEnumISMNonIC.xsd
var xsdNonICMarkings = []string{
	"DS", "XD", "ND", "SBU", "SBU-NF", "LES", "LES-NF", "SSI", "NNPI",
}

// XSD 25X declass exceptions from CVEnumISM25X.xsd
var xsdDeclassExceptions = []string{
	"AEA", "NATO", "NATO-AEA",
	"25X1", "25X1-EO-12951", "25X2", "25X3", "25X4", "25X5", "25X6",
	"25X7", "25X8", "25X9",
	"50X1", "50X1-HUM", "50X2", "50X2-WMD", "50X3", "50X4", "50X5",
	"50X6", "50X7", "50X8", "50X9",
	"75X",
}

// XSD compliesWith values from CVEnumISMCompliesWith.xsd
var xsdCompliesWith = []string{
	"USGov", "USIC", "USDOD", "OtherAuthority", "USA-CUI-ONLY", "USA-CUI",
}

// XSD exemptFrom values from CVEnumISMExemptFrom.xsd
var xsdExemptFrom = []string{
	"IC_710_MANDATORY_FDR", "DOD_DISTRO_STATEMENT",
}

// XSD CUI specified categories from CVEnumISMCUISpecified.xsd (56 values)
var xsdCUISpecified = []string{
	"AIV", "ADPO", "CRITAN", "ARCHR", "FSEC", "BUDG", "FUND", "CVI",
	"CHLD", "CCI", "CONTRACT", "SUB", "CTI", "CHRI", "CEII", "LDNA",
	"EXPT", "JURY", "TAX", "FISA", "FISAB", "FNC", "INTEL", "NUC",
	"PRVCY", "PROCURE", "PROPIN", "GENETIC", "GEO", "HLTH", "HISTP",
	"INF", "PRIIG", "IFNC", "ID", "INTL", "INV", "LFNC", "NPSR",
	"NNPI", "SRI", "PERS", "MFC", "PCII", "LPROT", "SGI", "SSI",
	"SSEL", "STAT", "STUD", "TSCA", "DCNI", "UCNI", "CENS", "WHSTL",
	"WIT", "WDT",
}

// XSD ISM attribute names from IC-ISM.xsd (30 attributes)
var xsdISMAttributes = []struct {
	Name    string // XSD attribute name
	JSONTag string // expected JSON tag in ISM struct
}{
	{"atomicEnergyMarkings", "atomicEnergyMarkings"},
	{"classification", "classification"},
	{"classificationReason", "classificationReason"},
	{"classifiedBy", "classifiedBy"},
	{"compilationReason", "compilationReason"},
	{"compliesWith", "compliesWith"},
	{"createDate", "createDate"},
	{"cuiBasic", "cuiBasic"},
	{"cuiSpecified", "cuiSpecified"},
	{"declassException", "declassException"},
	{"displayOnlyTo", "displayOnlyTo"},
	{"disseminationControls", "disseminationControls"},
	{"exemptFrom", "exemptFrom"},
	{"FGIsourceOpen", "fgiSourceOpen"},
	{"FGIsourceProtected", "fgiSourceProtected"},
	{"handleViaChannels", "handleViaChannels"},
	{"hasApproximateMarkings", "hasApproximateMarkings"},
	{"highWaterNATO", "highWaterNATO"},
	{"joint", "joint"},
	{"noAggregation", "noAggregation"},
	{"nonICmarkings", "nonICMarkings"},
	{"nonUSControls", "nonUSControls"},
	{"noticeType", "noticeType"},
	{"noticeProseID", "noticeProseID"},
	{"ownerProducer", "ownerProducer"},
	{"pocType", "pocType"},
	{"releasableTo", "releasableTo"},
	{"SARIdentifier", "sarIdentifier"},
	{"SCIcontrols", "sciControls"},
	{"secondBannerLine", "secondBannerLine"},
}

package banner

import (
	"sort"
	"strconv"
	"strings"

	"expr.ai/ism-api/internal/model"
)

// Result holds the rendered banner line, portion mark, and authority block for an ISM object.
type Result struct {
	BannerLine     string `json:"bannerLine"`
	PortionMark    string `json:"portionMark"`
	AuthorityBlock string `json:"authorityBlock"`
}

// Full classification labels for banner lines.
var classificationBanner = map[model.Classification]string{
	model.ClassificationU:   "UNCLASSIFIED",
	model.ClassificationCUI: "CUI",
	model.ClassificationC:   "CONFIDENTIAL",
	model.ClassificationS:   "SECRET",
	model.ClassificationTS:  "TOP SECRET",
}

// Abbreviated classification labels for portion marks.
var classificationPortion = map[model.Classification]string{
	model.ClassificationU:   "U",
	model.ClassificationCUI: "CUI",
	model.ClassificationC:   "C",
	model.ClassificationS:   "S",
	model.ClassificationTS:  "TS",
}

// disseminationOrder defines the canonical banner ordering index for each
// dissemination control code per ISM banner ordering rules.
var disseminationOrder = map[string]int{
	"RS":           0,
	"OC":           1,
	"OC-USGOV":    2,
	"IMCON":        3,
	"NOFORN":       4,
	"PROPIN":       5,
	"REL":          6,
	"RELIDO":       7,
	"EYES":         8,
	"DSEN":         9,
	"FISA":         10,
	"DISPLAY ONLY": 11,
	"FED ONLY":     12,
	"FEDCON":       13,
	"NOCON":        14,
	"DL ONLY":      15,
}

// portionAbbrev maps dissemination control codes to their abbreviated portion
// mark forms. Controls not listed here use the full code in portion marks.
var portionAbbrev = map[string]string{
	"NOFORN": "NF",
	"PROPIN": "PR",
	"IMCON":  "IMC",
	"DSEN":   "DS",
	"NOCON":  "NC",
}

// Render produces the banner line and portion mark for the given ISM object.
// Banner ordering: classification // SCI controls // dissemination controls / FGI / non-IC.
func Render(ism *model.ISM) Result {
	bannerClass := classificationBanner[ism.Classification]
	portionClass := classificationPortion[ism.Classification]

	// Joint documents: prepend "JOINT" and append ownerProducer countries.
	if ism.Joint && len(ism.OwnerProducer) > 1 {
		countries := strings.Join(ism.OwnerProducer, " ")
		bannerClass = "JOINT " + bannerClass + " " + countries
		portionClass = "J" + portionClass + " " + countries
	}

	var bannerParts []string
	var portionParts []string
	var sciParts []string

	// CUI category markings (rendered after classification for CUI).
	if ism.Classification == model.ClassificationCUI && len(ism.CategoryMarkings) > 0 {
		for _, cat := range ism.CategoryMarkings {
			bannerParts = append(bannerParts, cat)
			portionParts = append(portionParts, cat)
		}
	}

	// SCI controls — sorted alphabetically, same code in banner and portion.
	if len(ism.SCIControls) > 0 {
		sorted := make([]string, len(ism.SCIControls))
		copy(sorted, ism.SCIControls)
		sort.Strings(sorted)
		sciParts = sorted
	}

	// SAR identifiers — after SCI, before dissemination controls.
	// 1-2 SARs: SAR-[name] separated by / (FMT-11).
	// 3+ SARs: banner uses SAR-MULTIPLE PROGRAMS (FMT-13); portions always list all PIDs.
	var sarBanner, sarPortion string
	if len(ism.SARIdentifier) > 0 {
		sorted := make([]string, len(ism.SARIdentifier))
		copy(sorted, ism.SARIdentifier)
		sort.Strings(sorted)

		// Portion marks always list all PIDs individually, alphabetically.
		pids := make([]string, len(sorted))
		for i, id := range sorted {
			pids[i] = "SAR-" + id
		}
		sarPortion = strings.Join(pids, "/")

		// Banner: 3+ SAPs collapse to SAR-MULTIPLE PROGRAMS.
		if len(sorted) >= 3 {
			sarBanner = "SAR-MULTIPLE PROGRAMS"
		} else {
			sarBanner = strings.Join(pids, "/")
		}
	}

	// AEA markings — rendered between SCI and dissemination with own // separator.
	var aeaBanner, aeaPortion []string
	if len(ism.AtomicEnergyMarkings) > 0 {
		aeaBanner, aeaPortion = renderAEA(ism.AtomicEnergyMarkings, ism.Classification)
	}

	// Dissemination controls in canonical order.
	for _, ctrl := range sortControls(ism.DisseminationControls) {
		b, p := renderControl(ctrl, ism)
		bannerParts = append(bannerParts, b)
		portionParts = append(portionParts, p)
	}

	// FGI sources.
	if len(ism.FGISourceOpen) > 0 || len(ism.FGISourceProtected) > 0 {
		b, p := renderFGI(ism)
		bannerParts = append(bannerParts, b)
		portionParts = append(portionParts, p)
	}

	// Non-IC markings (passed through as-is).
	for _, m := range ism.NonICMarkings {
		bannerParts = append(bannerParts, m)
		portionParts = append(portionParts, m)
	}

	banner := bannerClass
	if len(sciParts) > 0 {
		banner += "//" + strings.Join(sciParts, "/")
	}
	if sarBanner != "" {
		banner += "//" + sarBanner
	}
	if len(aeaBanner) > 0 {
		banner += "//" + strings.Join(aeaBanner, "/")
	}
	if len(bannerParts) > 0 {
		banner += "//" + strings.Join(bannerParts, "/")
	}

	portion := portionClass
	if len(sciParts) > 0 {
		portion += "//" + strings.Join(sciParts, "/")
	}
	if sarPortion != "" {
		portion += "//" + sarPortion
	}
	if len(aeaPortion) > 0 {
		portion += "//" + strings.Join(aeaPortion, "/")
	}
	if len(portionParts) > 0 {
		portion += "//" + strings.Join(portionParts, "/")
	}

	return Result{
		BannerLine:     banner,
		PortionMark:    "(" + portion + ")",
		AuthorityBlock: renderAuthorityBlock(ism),
	}
}

// renderAuthorityBlock produces the authority block text for classified markings.
// Returns empty string for U and CUI, or if no authority fields are populated.
func renderAuthorityBlock(ism *model.ISM) string {
	if ism.Classification != model.ClassificationC && ism.Classification != model.ClassificationS && ism.Classification != model.ClassificationTS {
		return ""
	}

	var lines []string

	// Original classification authority.
	if ism.ClassifiedBy != "" {
		lines = append(lines, "Classified By: "+ism.ClassifiedBy)
		if ism.ClassificationReason != "" {
			lines = append(lines, "Reason: "+ism.ClassificationReason)
		}
	}

	// Derivative classification authority.
	if ism.DerivedFrom != "" {
		lines = append(lines, "Derived From: "+ism.DerivedFrom)
	}

	// Compilation reason.
	if ism.CompilationReason != "" {
		lines = append(lines, "Compilation: "+ism.CompilationReason)
	}

	// If no authority lines were generated, return empty.
	if len(lines) == 0 {
		return ""
	}

	// Declassification line.
	switch {
	case ism.DeclassDate != "":
		lines = append(lines, "Declassify On: "+ism.DeclassDate)
	case ism.DeclassEvent != "":
		lines = append(lines, "Declassify On: "+ism.DeclassEvent)
	case ism.DeclassException != "":
		lines = append(lines, "Declassify On: "+ism.DeclassException)
	default:
		// Only show "not specified" for derivative (which requires declass info).
		if ism.DerivedFrom != "" {
			lines = append(lines, "Declassify On: (not specified)")
		}
	}

	return strings.Join(lines, "\n")
}

// sortControls returns dissemination controls sorted by their canonical banner
// ordering. Unknown controls sort to the end, preserving their relative order.
func sortControls(controls []string) []string {
	sorted := make([]string, len(controls))
	copy(sorted, controls)
	sort.SliceStable(sorted, func(i, j int) bool {
		oi, oki := disseminationOrder[sorted[i]]
		oj, okj := disseminationOrder[sorted[j]]
		if !oki {
			oi = len(disseminationOrder)
		}
		if !okj {
			oj = len(disseminationOrder)
		}
		return oi < oj
	})
	return sorted
}

// renderControl returns the banner and portion mark representations for a
// single dissemination control. REL and DISPLAY ONLY expand their associated
// country lists.
func renderControl(ctrl string, ism *model.ISM) (banner, portion string) {
	switch ctrl {
	case "REL":
		countries := sortCountriesUSAFirst(ism.ReleasableTo)
		joined := strings.Join(countries, ", ")
		return "REL TO " + joined, "REL TO " + joined
	case "DISPLAY ONLY":
		countries := strings.Join(ism.DisplayOnlyTo, ", ")
		return "DISPLAY ONLY " + countries, "DISPLAY ONLY " + countries
	default:
		abbr, ok := portionAbbrev[ctrl]
		if !ok {
			abbr = ctrl
		}
		return ctrl, abbr
	}
}

// sortCountriesUSAFirst returns a sorted copy of the country list with USA
// first (if present), followed by the remaining countries in alphabetical
// order. This implements CC-1 (E4-S10.d.4).
func sortCountriesUSAFirst(countries []string) []string {
	sorted := make([]string, len(countries))
	copy(sorted, countries)
	sort.Strings(sorted)

	// Move USA to front if present.
	for i, c := range sorted {
		if c == "USA" {
			sorted = append(sorted[:i], sorted[i+1:]...)
			sorted = append([]string{"USA"}, sorted...)
			break
		}
	}
	return sorted
}

// renderFGI returns the banner and portion mark representations for FGI
// sources. The banner line lists all source countries; the portion mark uses
// the abbreviated "FGI" form.
func renderFGI(ism *model.ISM) (banner, portion string) {
	var sources []string
	sources = append(sources, ism.FGISourceOpen...)
	sources = append(sources, ism.FGISourceProtected...)
	return "FGI " + strings.Join(sources, " "), "FGI"
}

// renderAEA returns the banner and portion parts for Atomic Energy Act markings.
// AEA codes: RD, FRD, RD-CNWDI, RD-SG-{n}, FRD-SG-{n}, DCNI, UCNI, TFNI.
// AEA markings are suppressed for UNCLASSIFIED and CUI.
func renderAEA(markings []string, cls model.Classification) (bannerParts, portionParts []string) {
	if cls == model.ClassificationU || cls == model.ClassificationCUI {
		return nil, nil
	}

	var hasRD, hasFRD, hasCNWDI bool
	var rdSigmas, frdSigmas []int
	var otherBanner, otherPortion []string

	for _, code := range markings {
		switch {
		case code == "RD":
			hasRD = true
		case code == "FRD":
			hasFRD = true
		case code == "RD-CNWDI":
			hasCNWDI = true
		case strings.HasPrefix(code, "RD-SG-"):
			if n, err := strconv.Atoi(code[6:]); err == nil {
				rdSigmas = append(rdSigmas, n)
			}
		case strings.HasPrefix(code, "FRD-SG-"):
			if n, err := strconv.Atoi(code[7:]); err == nil {
				frdSigmas = append(frdSigmas, n)
			}
		default:
			// DCNI, UCNI, TFNI — pass through as-is.
			otherBanner = append(otherBanner, code)
			otherPortion = append(otherPortion, code)
		}
	}

	// RD-based markings: CNWDI takes precedence over plain RD.
	if hasCNWDI {
		bannerParts = append(bannerParts, "RESTRICTED DATA-N")
		portionParts = append(portionParts, "RD-N")
	} else if hasRD {
		bannerParts = append(bannerParts, "RESTRICTED DATA")
		portionParts = append(portionParts, "RD")
	}

	if len(rdSigmas) > 0 {
		sort.Ints(rdSigmas)
		nums := joinInts(rdSigmas)
		bannerParts = append(bannerParts, "RD-SIGMA "+nums)
		portionParts = append(portionParts, "RD-SG "+nums)
	}

	// FRD-based markings.
	if hasFRD {
		bannerParts = append(bannerParts, "FORMERLY RESTRICTED DATA")
		portionParts = append(portionParts, "FRD")
	}

	if len(frdSigmas) > 0 {
		sort.Ints(frdSigmas)
		nums := joinInts(frdSigmas)
		bannerParts = append(bannerParts, "FRD-SIGMA "+nums)
		portionParts = append(portionParts, "FRD-SG "+nums)
	}

	// Other AEA markings (DCNI, UCNI, TFNI).
	bannerParts = append(bannerParts, otherBanner...)
	portionParts = append(portionParts, otherPortion...)

	return bannerParts, portionParts
}

// joinInts converts a slice of ints to a space-separated string.
func joinInts(nums []int) string {
	strs := make([]string, len(nums))
	for i, n := range nums {
		strs[i] = strconv.Itoa(n)
	}
	return strings.Join(strs, " ")
}

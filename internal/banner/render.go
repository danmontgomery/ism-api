package banner

import (
	"sort"
	"strconv"
	"strings"

	"dmontgomery/ism-api/internal/model"
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
// Banner ordering per DoDM Figure 25:
// classification // SCI // AEA // SAR // FGI // dissemination // non-IC.
func Render(ism *model.ISM) Result {
	bannerClass := classificationBanner[ism.Classification]
	portionClass := classificationPortion[ism.Classification]

	// Joint documents: //JOINT [classification] [countries] per E4-S5.c/d.
	// Countries are sorted alphabetically (E4-S5.e), NOT USA-first.
	if ism.Joint && len(ism.OwnerProducer) > 1 {
		countries := make([]string, len(ism.OwnerProducer))
		copy(countries, ism.OwnerProducer)
		sort.Strings(countries)
		countryStr := strings.Join(countries, " ")
		bannerClass = "//JOINT " + bannerClass + " " + countryStr
		portionClass = "//JOINT " + portionClass + " " + countryStr
	}

	// FGI non-US documents: //[country] [classification] format (E4-S4.a.1).
	// NATO has its own rendering rules (E4-S4.b) and is excluded here.
	if !ism.Joint && len(ism.OwnerProducer) == 1 && ism.OwnerProducer[0] != "USA" && ism.OwnerProducer[0] != "NATO" {
		country := ism.OwnerProducer[0]
		bannerClass = "//" + country + " " + bannerClass
		portionClass = "//" + country + " " + portionClass
	}

	var sciParts []string
	var fgiBannerParts, fgiPortionParts []string
	var dissemBannerParts, dissemPortionParts []string
	var nonICBannerParts, nonICPortionParts []string

	// CUI category markings (rendered before dissem controls in the dissem section).
	if ism.Classification == model.ClassificationCUI && len(ism.CategoryMarkings) > 0 {
		for _, cat := range ism.CategoryMarkings {
			dissemBannerParts = append(dissemBannerParts, cat)
			dissemPortionParts = append(dissemPortionParts, cat)
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
	classified := ism.Classification.AtLeast(model.ClassificationC)
	for _, ctrl := range sortControls(ism.DisseminationControls) {
		// BP-7: FOUO is U/CUI-only; suppress from classified banners.
		if ctrl == "FOUO" && classified {
			continue
		}
		// SCI is the umbrella category; suppress when individual compartments
		// (SI, TK, etc.) are already rendered in the SCI section.
		if ctrl == "SCI" && len(ism.SCIControls) > 0 {
			continue
		}
		b, p := renderControl(ctrl, ism)
		dissemBannerParts = append(dissemBannerParts, b)
		dissemPortionParts = append(dissemPortionParts, p)
	}

	// FGI sources.
	if len(ism.FGISourceOpen) > 0 || len(ism.FGISourceProtected) > 0 {
		b, p := renderFGI(ism)
		fgiBannerParts = append(fgiBannerParts, b)
		fgiPortionParts = append(fgiPortionParts, p)
	}

	// Non-IC markings — BP-7: suppress FOUO in classified context.
	for _, m := range ism.NonICMarkings {
		if m == "FOUO" && classified {
			continue
		}
		nonICBannerParts = append(nonICBannerParts, m)
		nonICPortionParts = append(nonICPortionParts, m)
	}

	// Assemble banner: classification // SCI // AEA // SAR // FGI // dissem // nonIC
	banner := bannerClass
	if len(sciParts) > 0 {
		banner += "//" + strings.Join(sciParts, "/")
	}
	if len(aeaBanner) > 0 {
		banner += "//" + strings.Join(aeaBanner, "/")
	}
	if sarBanner != "" {
		banner += "//" + sarBanner
	}
	if len(fgiBannerParts) > 0 {
		banner += "//" + strings.Join(fgiBannerParts, "/")
	}
	if len(dissemBannerParts) > 0 {
		banner += "//" + strings.Join(dissemBannerParts, "/")
	}
	if len(nonICBannerParts) > 0 {
		banner += "//" + strings.Join(nonICBannerParts, "/")
	}

	portion := portionClass
	if len(sciParts) > 0 {
		portion += "//" + strings.Join(sciParts, "/")
	}
	if len(aeaPortion) > 0 {
		portion += "//" + strings.Join(aeaPortion, "/")
	}
	if sarPortion != "" {
		portion += "//" + sarPortion
	}
	if len(fgiPortionParts) > 0 {
		portion += "//" + strings.Join(fgiPortionParts, "/")
	}
	if len(dissemPortionParts) > 0 {
		portion += "//" + strings.Join(dissemPortionParts, "/")
	}
	if len(nonICPortionParts) > 0 {
		portion += "//" + strings.Join(nonICPortionParts, "/")
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

	// E4-S5.h: authority block is used only when US is a co-owner in JOINT documents.
	if ism.Joint {
		hasUSA := false
		for _, c := range ism.OwnerProducer {
			if c == "USA" {
				hasUSA = true
				break
			}
		}
		if !hasUSA {
			return ""
		}
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
// sources. Open sources list country codes alphabetically (E4-S9.d.alpha);
// protected sources are concealed — no country codes shown (E4-S9.e).
func renderFGI(ism *model.ISM) (banner, portion string) {
	if len(ism.FGISourceOpen) == 0 {
		// Only protected sources — conceal all country codes per E4-S9.e.
		return "FGI", "FGI"
	}
	sorted := make([]string, len(ism.FGISourceOpen))
	copy(sorted, ism.FGISourceOpen)
	sort.Strings(sorted)
	return "FGI " + strings.Join(sorted, " "), "FGI"
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

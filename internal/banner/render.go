package banner

import (
	"sort"
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
}

// Abbreviated classification labels for portion marks.
var classificationPortion = map[model.Classification]string{
	model.ClassificationU:   "U",
	model.ClassificationCUI: "CUI",
	model.ClassificationC:   "C",
	model.ClassificationS:   "S",
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
// Banner ordering: classification / SCI (future) / dissemination controls / FGI / non-IC.
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

	// CUI category markings (rendered after classification for CUI).
	if ism.Classification == model.ClassificationCUI && len(ism.CategoryMarkings) > 0 {
		for _, cat := range ism.CategoryMarkings {
			bannerParts = append(bannerParts, cat)
			portionParts = append(portionParts, cat)
		}
	}

	// SCI controls — reserved for future TS support.

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
	if len(bannerParts) > 0 {
		banner += "//" + strings.Join(bannerParts, "/")
	}

	portion := portionClass
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
	if ism.Classification != model.ClassificationC && ism.Classification != model.ClassificationS {
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
		countries := strings.Join(ism.ReleasableTo, ", ")
		return "REL TO " + countries, "REL TO " + countries
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

// renderFGI returns the banner and portion mark representations for FGI
// sources. The banner line lists all source countries; the portion mark uses
// the abbreviated "FGI" form.
func renderFGI(ism *model.ISM) (banner, portion string) {
	var sources []string
	sources = append(sources, ism.FGISourceOpen...)
	sources = append(sources, ism.FGISourceProtected...)
	return "FGI " + strings.Join(sources, " "), "FGI"
}

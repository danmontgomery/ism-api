// Package parse converts a rendered banner line or portion mark back into a
// best-effort model.ISM object. It is the inverse of internal/banner.Render:
// where that package turns structured data into text, this package turns
// text back into structured data, on a "never fail, warn honestly" basis —
// unrecognized or unrecoverable content produces warnings rather than an
// error, and a round-trip check tells the caller whether the parse was
// faithful.
package parse

import (
	"strconv"
	"strings"

	"dmontgomery/ism-api/internal/banner"
	"dmontgomery/ism-api/internal/model"
	"dmontgomery/ism-api/internal/refdata"
	"dmontgomery/ism-api/internal/validation"
)

// Warning codes emitted by the parser.
const (
	CodeUnknownToken      = "PARSE_UNKNOWN_TOKEN"
	CodeAmbiguous         = "PARSE_AMBIGUOUS"
	CodeLossy             = "PARSE_LOSSY"
	CodeRoundTripMismatch = "PARSE_ROUNDTRIP_MISMATCH"
)

// Banner ordering sections per DoDM Figure 25:
// classification // SCI // AEA // SAR // FGI // dissem // non-IC.
// Segments must claim sections in non-decreasing order.
const (
	secSCI = iota
	secAEA
	secSAR
	secFGI
	secDissem
	secNonIC
	secCount
)

// Result is the outcome of parsing a marking string.
type Result struct {
	Form      string                  `json:"form"` // "banner" | "portion"
	ISM       model.ISM               `json:"ism"`
	RoundTrip RoundTrip               `json:"roundTrip"`
	Warnings  []validation.FieldError `json:"warnings,omitempty"`
	// Inferred lists fields populated by decoding marking convention rather
	// than from a literal token in the input, e.g. "ownerProducer" when a
	// bare classification (no country prefix) was read as US-owned.
	Inferred []string `json:"inferred,omitempty"`
}

// RoundTrip reports whether re-rendering the parsed ISM reproduces the
// (normalized) input exactly.
type RoundTrip struct {
	Matches  bool   `json:"matches"`
	Rendered string `json:"rendered"`
}

// Parse converts marking (a banner line or a parenthesized portion mark)
// into a best-effort ISM object. It never fails: malformed or unrecognized
// content is reported via Warnings instead.
func Parse(marking string, reg *refdata.Registry) *Result {
	p := &parser{reg: reg, result: &Result{}}
	p.run(marking)
	return p.result
}

type parser struct {
	reg    *refdata.Registry
	result *Result
}

func (p *parser) addWarning(field, code, message string) {
	p.result.Warnings = append(p.result.Warnings, validation.FieldError{
		Field:    field,
		Code:     code,
		Message:  message,
		Severity: validation.SeverityWarning,
	})
}

func (p *parser) addInferred(field string) {
	p.result.Inferred = append(p.result.Inferred, field)
}

func (p *parser) run(marking string) {
	trimmed := strings.TrimSpace(marking)

	isPortion := len(trimmed) >= 2 && strings.HasPrefix(trimmed, "(") && strings.HasSuffix(trimmed, ")")
	var content string
	if isPortion {
		p.result.Form = "portion"
		content = strings.TrimSpace(trimmed[1 : len(trimmed)-1])
	} else {
		p.result.Form = "banner"
		content = trimmed
	}
	content = strings.ToUpper(content)

	ism := model.ISM{}

	// Strip a leading "//" before splitting, otherwise the classification
	// segment (e.g. "//JOINT SECRET GBR USA//...") yields a spurious empty
	// first element.
	body := strings.TrimPrefix(content, "//")
	segments := strings.Split(body, "//")

	clsSeg := segments[0]
	cls, joint, owners, ok := p.parseClassificationSegment(clsSeg, isPortion)
	if ok {
		ism.Classification = cls
		ism.Joint = joint
		if len(owners) > 0 {
			ism.OwnerProducer = owners
		}
		// A marking with no country prefix is US-owned: render.go emits a
		// "//<COUNTRY> <LABEL>" prefix only when the sole owner is neither
		// USA nor NATO, and DoDM E4-S4.a.1 requires non-US owners to carry
		// their country code, so a correctly-marked bare banner cannot be
		// non-US. This is decoding the convention, not fabricating data:
		// OwnerProducer:["USA"] renders byte-identically to an empty
		// OwnerProducer, so the round trip is unaffected either way.
		//
		// Known ambiguity: render.go also excludes NATO from that prefix
		// branch, because NATO banner forms (//NATO SECRET, COSMIC TOP
		// SECRET, ATOMAL) aren't implemented yet. So within this codebase a
		// bare banner is ambiguous between USA and NATO. Under real DoDM
		// rules a NATO document always carries its own prefix, so a
		// correctly-marked bare banner is never NATO — the ambiguity is an
		// artifact of the unimplemented NATO renderer, not of this
		// inference. Revisit if NATO rendering is added.
		if len(ism.OwnerProducer) == 0 {
			ism.OwnerProducer = []string{"USA"}
			p.addInferred("ownerProducer")
		}
	} else {
		p.addWarning("classification", CodeUnknownToken,
			"unable to recognize classification segment: "+strconv.Quote(clsSeg))
	}

	lastClaimed := -1
	for _, seg := range segments[1:] {
		lastClaimed = p.claimSegment(seg, lastClaimed, isPortion, &ism)
	}

	p.result.ISM = ism

	rendered := banner.Render(&ism)
	var renderedStr, normalizedInput string
	if isPortion {
		renderedStr = rendered.PortionMark
		normalizedInput = "(" + content + ")"
	} else {
		renderedStr = rendered.BannerLine
		normalizedInput = content
	}
	matches := renderedStr == normalizedInput
	p.result.RoundTrip = RoundTrip{Matches: matches, Rendered: renderedStr}
	if !matches {
		p.addWarning("marking", CodeRoundTripMismatch,
			"parsed ISM does not render back to the original marking: got "+strconv.Quote(renderedStr)+", want "+strconv.Quote(normalizedInput))
	}
}

// parseClassificationSegment parses the first "//"-delimited segment, which
// encodes classification plus an optional JOINT or single-country (FGI
// non-US) prefix. Mirrors render.go's three forms:
//
//	JOINT <LABEL> <C1> <C2> ...   -> Joint: true, OwnerProducer: [countries]
//	<COUNTRY> <LABEL>             -> OwnerProducer: [country]
//	<LABEL>                       -> classification only
func (p *parser) parseClassificationSegment(seg string, isPortion bool) (cls model.Classification, joint bool, owners []string, ok bool) {
	if seg == "" {
		return "", false, nil, false
	}
	if isPortion {
		return p.parsePortionClassSegment(seg)
	}
	return p.parseBannerClassSegment(seg)
}

func (p *parser) parsePortionClassSegment(seg string) (model.Classification, bool, []string, bool) {
	tokens := strings.Fields(seg)
	if len(tokens) == 0 {
		return "", false, nil, false
	}

	joint := false
	idx := 0
	if tokens[0] == "JOINT" {
		joint = true
		idx = 1
	}
	if idx >= len(tokens) {
		return "", joint, nil, false
	}

	// Non-joint FGI form: "<COUNTRY> <CODE>" — exactly two tokens.
	if !joint && len(tokens) == 2 && p.reg.ValidCountryCode(tokens[0]) {
		if c, ok := banner.ClassificationFromPortion(tokens[1]); ok {
			return c, false, []string{tokens[0]}, true
		}
	}

	c, ok := banner.ClassificationFromPortion(tokens[idx])
	if !ok {
		return "", joint, nil, false
	}
	return c, joint, tokens[idx+1:], true
}

func (p *parser) parseBannerClassSegment(seg string) (model.Classification, bool, []string, bool) {
	rest := seg
	joint := false
	if strings.HasPrefix(rest, "JOINT ") {
		joint = true
		rest = strings.TrimPrefix(rest, "JOINT ")
	}

	if !joint {
		// Non-joint FGI form: "<COUNTRY> <LABEL>". Labels may contain
		// spaces (e.g. "TOP SECRET"), so split on the first space only.
		if sp := strings.IndexByte(rest, ' '); sp >= 0 {
			country := rest[:sp]
			remainder := rest[sp+1:]
			if p.reg.ValidCountryCode(country) {
				if c, ok := banner.ClassificationFromBanner(remainder); ok {
					return c, false, []string{country}, true
				}
			}
		}
		if c, ok := banner.ClassificationFromBanner(rest); ok {
			return c, false, nil, true
		}
		return "", false, nil, false
	}

	// Joint form: "JOINT <LABEL> <C1> <C2> ...". Match labels longest-first
	// so a two-word label like "TOP SECRET" isn't mistaken for a one-word
	// prefix of itself.
	for _, label := range banner.BannerClassificationLabels() {
		if rest == label {
			c, _ := banner.ClassificationFromBanner(label)
			return c, true, nil, true
		}
		if strings.HasPrefix(rest, label+" ") {
			c, _ := banner.ClassificationFromBanner(label)
			remainder := strings.TrimPrefix(rest, label+" ")
			return c, true, strings.Fields(remainder), true
		}
	}
	return "", true, nil, false
}

// claimSegment assigns a "//"-delimited segment to the earliest available
// section (at or after lastClaimed) whose recognizer accepts every item in
// the segment. It returns the index of the section that was claimed (or
// lastClaimed unchanged if nothing matched at the segment level, in which
// case items are salvaged individually).
func (p *parser) claimSegment(seg string, lastClaimed int, isPortion bool, ism *model.ISM) int {
	items := strings.Split(seg, "/")

	var candidates []int
	for i := lastClaimed + 1; i < secCount; i++ {
		if p.sectionAccepts(i, seg, items, isPortion, ism) {
			candidates = append(candidates, i)
		}
	}

	if len(candidates) == 0 {
		return p.claimItemsIndividually(items, lastClaimed, isPortion, ism)
	}

	chosen := candidates[0]
	if len(candidates) > 1 {
		p.addWarning("marking", CodeAmbiguous,
			strconv.Quote(seg)+" matches more than one section ("+sectionNames(candidates)+"); using "+sectionName(chosen))
	}
	p.assignSection(chosen, seg, items, isPortion, ism)
	return chosen
}

// claimItemsIndividually is the fallback when a whole segment isn't
// recognized by any single section: each item is evaluated on its own, so a
// genuinely unrecognized item is dropped and warned about without losing
// its recognized neighbors.
func (p *parser) claimItemsIndividually(items []string, lastClaimed int, isPortion bool, ism *model.ISM) int {
	cursor := lastClaimed
	for _, it := range items {
		found := -1
		for i := cursor + 1; i < secCount; i++ {
			if p.sectionAccepts(i, it, []string{it}, isPortion, ism) {
				found = i
				break
			}
		}
		if found == -1 {
			p.addWarning("marking", CodeUnknownToken, "unrecognized token: "+strconv.Quote(it))
			continue
		}
		p.assignSection(found, it, []string{it}, isPortion, ism)
		cursor = found
	}
	return cursor
}

func sectionName(i int) string {
	switch i {
	case secSCI:
		return "sci"
	case secAEA:
		return "aea"
	case secSAR:
		return "sar"
	case secFGI:
		return "fgi"
	case secDissem:
		return "dissem"
	case secNonIC:
		return "non-ic"
	}
	return "unknown"
}

func sectionNames(idxs []int) string {
	names := make([]string, len(idxs))
	for i, idx := range idxs {
		names[i] = sectionName(idx)
	}
	return strings.Join(names, ", ")
}

// sectionAccepts reports whether every item in the segment is recognized by
// the given section's vocabulary.
func (p *parser) sectionAccepts(idx int, rawSeg string, items []string, isPortion bool, ism *model.ISM) bool {
	switch idx {
	case secSCI:
		if len(items) == 0 {
			return false
		}
		for _, it := range items {
			if !p.reg.ValidSCIControl(it) {
				return false
			}
		}
		return true
	case secAEA:
		if len(items) == 0 {
			return false
		}
		for _, it := range items {
			if _, ok := matchAEAItem(it, isPortion); !ok {
				return false
			}
		}
		return true
	case secSAR:
		if len(items) == 0 {
			return false
		}
		for _, it := range items {
			if !strings.HasPrefix(it, "SAR-") {
				return false
			}
		}
		return true
	case secFGI:
		return rawSeg == "FGI" || strings.HasPrefix(rawSeg, "FGI ")
	case secDissem:
		if len(items) == 0 {
			return false
		}
		for _, it := range items {
			if !p.isDissemItem(it, ism.Classification) {
				return false
			}
		}
		return true
	case secNonIC:
		if len(items) == 0 {
			return false
		}
		for _, it := range items {
			if !p.reg.ValidNonICMarking(it) {
				return false
			}
		}
		return true
	}
	return false
}

func (p *parser) isDissemItem(item string, cls model.Classification) bool {
	if strings.HasPrefix(item, "REL TO ") {
		return true
	}
	if strings.HasPrefix(item, "DISPLAY ONLY ") {
		return true
	}
	if _, ok := banner.ControlFromPortionAbbrev(item); ok {
		return true
	}
	if p.reg.ValidDisseminationControl(item) {
		return true
	}
	if cls == model.ClassificationCUI && p.reg.ValidCUICategory(item) {
		return true
	}
	return false
}

// assignSection mutates ism to record the given section's interpretation of
// the segment/item.
func (p *parser) assignSection(idx int, rawSeg string, items []string, isPortion bool, ism *model.ISM) {
	switch idx {
	case secSCI:
		ism.SCIControls = append(ism.SCIControls, items...)
	case secAEA:
		for _, it := range items {
			codes, _ := matchAEAItem(it, isPortion)
			ism.AtomicEnergyMarkings = append(ism.AtomicEnergyMarkings, codes...)
		}
	case secSAR:
		p.assignSAR(items, ism)
	case secFGI:
		p.assignFGI(rawSeg, ism)
	case secDissem:
		for _, it := range items {
			p.assignDissemItem(it, ism)
		}
	case secNonIC:
		ism.NonICMarkings = append(ism.NonICMarkings, items...)
	}
}

// assignSAR handles both forms: 1-2 PIDs listed individually ("SAR-A/SAR-B"),
// and the banner-only collapsed form ("SAR-MULTIPLE PROGRAMS") used for 3+
// PIDs, which is permanently lossy — the individual PIDs cannot be recovered.
func (p *parser) assignSAR(items []string, ism *model.ISM) {
	for _, it := range items {
		if it == "SAR-MULTIPLE PROGRAMS" {
			p.addWarning("sarIdentifier", CodeLossy, "SAR-MULTIPLE PROGRAMS conceals individual PIDs")
			continue
		}
		ism.SARIdentifier = append(ism.SARIdentifier, strings.TrimPrefix(it, "SAR-"))
	}
}

// assignFGI handles both forms: "FGI <codes...>" for open sources, and bare
// "FGI" which is inherently ambiguous — it's how both protected-only sources
// and (in portion marks) any FGI source at all are rendered. Per E4-S9.e /
// renderFGI, bare "FGI" is treated as a concealed protected source; this is
// lossy for the specific countries but round-trips correctly, since
// FGISourceProtected always renders back to bare "FGI".
func (p *parser) assignFGI(rawSeg string, ism *model.ISM) {
	if rawSeg == "FGI" {
		ism.FGISourceProtected = append(ism.FGISourceProtected, "FGI")
		p.addWarning("fgiSourceOpen", CodeLossy, "bare FGI marking conceals the specific source country codes")
		return
	}
	codes := strings.Fields(strings.TrimPrefix(rawSeg, "FGI "))
	ism.FGISourceOpen = append(ism.FGISourceOpen, codes...)
}

func (p *parser) assignDissemItem(item string, ism *model.ISM) {
	switch {
	case strings.HasPrefix(item, "REL TO "):
		ism.DisseminationControls = append(ism.DisseminationControls, "REL")
		for _, c := range strings.Split(strings.TrimPrefix(item, "REL TO "), ",") {
			ism.ReleasableTo = append(ism.ReleasableTo, strings.TrimSpace(c))
		}
	case strings.HasPrefix(item, "DISPLAY ONLY "):
		ism.DisseminationControls = append(ism.DisseminationControls, "DISPLAY ONLY")
		for _, c := range strings.Split(strings.TrimPrefix(item, "DISPLAY ONLY "), ",") {
			ism.DisplayOnlyTo = append(ism.DisplayOnlyTo, strings.TrimSpace(c))
		}
	default:
		if ctrl, ok := banner.ControlFromPortionAbbrev(item); ok {
			ism.DisseminationControls = append(ism.DisseminationControls, ctrl)
			return
		}
		if p.reg.ValidDisseminationControl(item) {
			ism.DisseminationControls = append(ism.DisseminationControls, item)
			return
		}
		if ism.Classification == model.ClassificationCUI && p.reg.ValidCUICategory(item) {
			ism.CategoryMarkings = append(ism.CategoryMarkings, item)
			return
		}
	}
}

// matchAEAItem recognizes a single Atomic Energy Act marking item and
// returns the canonical model codes it contributes. Mirrors renderAEA in
// internal/banner/render.go, in both banner and portion forms.
func matchAEAItem(item string, isPortion bool) ([]string, bool) {
	if isPortion {
		switch item {
		case "RD":
			return []string{"RD"}, true
		case "RD-N":
			return []string{"RD-CNWDI"}, true
		case "FRD":
			return []string{"FRD"}, true
		case "DCNI", "UCNI", "TFNI":
			return []string{item}, true
		}
		if rest, ok := strings.CutPrefix(item, "RD-SG "); ok {
			return sigmaCodes("RD-SG-", rest)
		}
		if rest, ok := strings.CutPrefix(item, "FRD-SG "); ok {
			return sigmaCodes("FRD-SG-", rest)
		}
		return nil, false
	}

	switch item {
	case "RESTRICTED DATA":
		return []string{"RD"}, true
	case "RESTRICTED DATA-N":
		return []string{"RD-CNWDI"}, true
	case "FORMERLY RESTRICTED DATA":
		return []string{"FRD"}, true
	case "DCNI", "UCNI", "TFNI":
		return []string{item}, true
	}
	if rest, ok := strings.CutPrefix(item, "RD-SIGMA "); ok {
		return sigmaCodes("RD-SG-", rest)
	}
	if rest, ok := strings.CutPrefix(item, "FRD-SIGMA "); ok {
		return sigmaCodes("FRD-SG-", rest)
	}
	return nil, false
}

// sigmaCodes turns a space-separated list of sigma numbers into canonical
// "<prefix><n>" codes, e.g. ("RD-SG-", "14 15") -> ["RD-SG-14", "RD-SG-15"].
func sigmaCodes(prefix, nums string) ([]string, bool) {
	fields := strings.Fields(nums)
	if len(fields) == 0 {
		return nil, false
	}
	out := make([]string, 0, len(fields))
	for _, f := range fields {
		if _, err := strconv.Atoi(f); err != nil {
			return nil, false
		}
		out = append(out, prefix+f)
	}
	return out, true
}

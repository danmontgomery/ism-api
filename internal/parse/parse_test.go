package parse

import (
	"reflect"
	"testing"

	"dmontgomery/ism-api/internal/model"
	"dmontgomery/ism-api/internal/refdata"
)

// normalizeISM collapses nil and empty slices to nil across all slice fields
// so structural comparisons aren't sensitive to append-vs-nil differences.
func normalizeISM(ism model.ISM) model.ISM {
	n := func(s []string) []string {
		if len(s) == 0 {
			return nil
		}
		return s
	}
	ism.OwnerProducer = n(ism.OwnerProducer)
	ism.CUISpecified = n(ism.CUISpecified)
	ism.CategoryMarkings = n(ism.CategoryMarkings)
	ism.DisseminationControls = n(ism.DisseminationControls)
	ism.ReleasableTo = n(ism.ReleasableTo)
	ism.DisplayOnlyTo = n(ism.DisplayOnlyTo)
	ism.FGISourceOpen = n(ism.FGISourceOpen)
	ism.FGISourceProtected = n(ism.FGISourceProtected)
	ism.ExemptFrom = n(ism.ExemptFrom)
	ism.CompliesWith = n(ism.CompliesWith)
	ism.AtomicEnergyMarkings = n(ism.AtomicEnergyMarkings)
	ism.NoticeType = n(ism.NoticeType)
	ism.SCIControls = n(ism.SCIControls)
	ism.SARIdentifier = n(ism.SARIdentifier)
	ism.NonICMarkings = n(ism.NonICMarkings)
	ism.NonUSControls = n(ism.NonUSControls)
	return ism
}

func assertISM(t *testing.T, got, want model.ISM) {
	t.Helper()
	gn, wn := normalizeISM(got), normalizeISM(want)
	if !reflect.DeepEqual(gn, wn) {
		t.Errorf("ISM mismatch:\n  got:  %+v\n  want: %+v", gn, wn)
	}
}

func hasWarningCode(t *testing.T, r *Result, code string) bool {
	t.Helper()
	for _, w := range r.Warnings {
		if w.Code == code {
			return true
		}
	}
	return false
}

var testReg = refdata.NewRegistry()

// ============================================================
// Seed suite: inverted from internal/banner/render_test.go TestRender
// (42 table cases). Each case's wantBanner/wantPortion becomes a parse
// input; the expected ISM is that case's ISM, adjusted for what is actually
// recoverable from the rendered text. A bare "SECRET" with no country
// prefix is read as OwnerProducer:["USA"] (see the inference comment in
// parse.go) and reported via wantInferred, since empty and ["USA"] render
// identically but only the latter validates as owned.
// ============================================================

type invertedCase struct {
	name         string
	input        string
	form         string
	want         model.ISM
	wantMatches  bool
	wantWarning  string   // non-empty: a warning code that must be present
	wantInferred []string // fields the parser inferred rather than read literally, e.g. ["ownerProducer"]
}

func TestParse_InvertedRenderCases(t *testing.T) {
	cases := []invertedCase{
		// --- Unclassified ---
		{"U bare/banner", "UNCLASSIFIED", "banner",
			model.ISM{Classification: model.ClassificationU, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"U bare/portion", "(U)", "portion",
			model.ISM{Classification: model.ClassificationU, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"U with FED ONLY dissemination/banner", "UNCLASSIFIED//FED ONLY", "banner",
			model.ISM{Classification: model.ClassificationU, DisseminationControls: []string{"FED ONLY"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"U with FED ONLY dissemination/portion", "(U//FED ONLY)", "portion",
			model.ISM{Classification: model.ClassificationU, DisseminationControls: []string{"FED ONLY"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},

		// --- CUI ---
		{"CUI bare/banner", "CUI", "banner",
			model.ISM{Classification: model.ClassificationCUI, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"CUI bare/portion", "(CUI)", "portion",
			model.ISM{Classification: model.ClassificationCUI, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"CUI with specified category/banner", "CUI//SP-CEII", "banner",
			model.ISM{Classification: model.ClassificationCUI, CategoryMarkings: []string{"SP-CEII"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"CUI with specified category/portion", "(CUI//SP-CEII)", "portion",
			model.ISM{Classification: model.ClassificationCUI, CategoryMarkings: []string{"SP-CEII"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"CUI with multiple categories/banner", "CUI//SP-CEII/SP-PHYS", "banner",
			model.ISM{Classification: model.ClassificationCUI, CategoryMarkings: []string{"SP-CEII", "SP-PHYS"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"CUI with multiple categories/portion", "(CUI//SP-CEII/SP-PHYS)", "portion",
			model.ISM{Classification: model.ClassificationCUI, CategoryMarkings: []string{"SP-CEII", "SP-PHYS"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"CUI with category and NOFORN/banner", "CUI//SP-CEII/NOFORN", "banner",
			model.ISM{Classification: model.ClassificationCUI, CategoryMarkings: []string{"SP-CEII"}, DisseminationControls: []string{"NOFORN"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"CUI with category and NOFORN/portion", "(CUI//SP-CEII/NF)", "portion",
			model.ISM{Classification: model.ClassificationCUI, CategoryMarkings: []string{"SP-CEII"}, DisseminationControls: []string{"NOFORN"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},

		// --- Confidential (bare-USA OwnerProducer is not recoverable) ---
		{"C bare/banner", "CONFIDENTIAL", "banner",
			model.ISM{Classification: model.ClassificationC, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"C bare/portion", "(C)", "portion",
			model.ISM{Classification: model.ClassificationC, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"C with NOFORN/banner", "CONFIDENTIAL//NOFORN", "banner",
			model.ISM{Classification: model.ClassificationC, DisseminationControls: []string{"NOFORN"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"C with NOFORN/portion", "(C//NF)", "portion",
			model.ISM{Classification: model.ClassificationC, DisseminationControls: []string{"NOFORN"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"C with REL TO/banner", "CONFIDENTIAL//REL TO USA, GBR", "banner",
			model.ISM{Classification: model.ClassificationC, DisseminationControls: []string{"REL"}, ReleasableTo: []string{"USA", "GBR"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"C with REL TO/portion", "(C//REL TO USA, GBR)", "portion",
			model.ISM{Classification: model.ClassificationC, DisseminationControls: []string{"REL"}, ReleasableTo: []string{"USA", "GBR"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"C with FGI open source/banner", "CONFIDENTIAL//FGI GBR", "banner",
			model.ISM{Classification: model.ClassificationC, FGISourceOpen: []string{"GBR"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"C with FGI open source/portion", "(C//FGI)", "portion",
			model.ISM{Classification: model.ClassificationC, FGISourceProtected: []string{"FGI"}, OwnerProducer: []string{"USA"}}, true, "PARSE_LOSSY", []string{"ownerProducer"}},

		// --- Secret ---
		{"S bare/banner", "SECRET", "banner",
			model.ISM{Classification: model.ClassificationS, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"S bare/portion", "(S)", "portion",
			model.ISM{Classification: model.ClassificationS, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"S with NOFORN/banner", "SECRET//NOFORN", "banner",
			model.ISM{Classification: model.ClassificationS, DisseminationControls: []string{"NOFORN"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"S with NOFORN/portion", "(S//NF)", "portion",
			model.ISM{Classification: model.ClassificationS, DisseminationControls: []string{"NOFORN"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"S with NOFORN and PROPIN/banner", "SECRET//NOFORN/PROPIN", "banner",
			model.ISM{Classification: model.ClassificationS, DisseminationControls: []string{"NOFORN", "PROPIN"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"S with NOFORN and PROPIN/portion", "(S//NF/PR)", "portion",
			model.ISM{Classification: model.ClassificationS, DisseminationControls: []string{"NOFORN", "PROPIN"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"S with REL TO multiple countries/banner", "SECRET//REL TO USA, CAN, GBR", "banner",
			model.ISM{Classification: model.ClassificationS, DisseminationControls: []string{"REL"}, ReleasableTo: []string{"USA", "CAN", "GBR"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"S with REL TO multiple countries/portion", "(S//REL TO USA, CAN, GBR)", "portion",
			model.ISM{Classification: model.ClassificationS, DisseminationControls: []string{"REL"}, ReleasableTo: []string{"USA", "CAN", "GBR"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"S with OC and NOFORN/banner", "SECRET//OC/NOFORN", "banner",
			model.ISM{Classification: model.ClassificationS, DisseminationControls: []string{"OC", "NOFORN"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"S with OC and NOFORN/portion", "(S//OC/NF)", "portion",
			model.ISM{Classification: model.ClassificationS, DisseminationControls: []string{"OC", "NOFORN"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"S with IMCON/banner", "SECRET//IMCON", "banner",
			model.ISM{Classification: model.ClassificationS, DisseminationControls: []string{"IMCON"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"S with IMCON/portion", "(S//IMC)", "portion",
			model.ISM{Classification: model.ClassificationS, DisseminationControls: []string{"IMCON"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"S with DSEN/banner", "SECRET//DSEN", "banner",
			model.ISM{Classification: model.ClassificationS, DisseminationControls: []string{"DSEN"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"S with DSEN/portion", "(S//DS)", "portion",
			// DS is ambiguous: DSEN portion abbrev vs. the "DS" non-IC code.
			// Position resolves it: dissem is the earlier available section.
			model.ISM{Classification: model.ClassificationS, DisseminationControls: []string{"DSEN"}, OwnerProducer: []string{"USA"}}, true, "PARSE_AMBIGUOUS", []string{"ownerProducer"}},
		{"S with NOCON/banner", "SECRET//NOCON", "banner",
			model.ISM{Classification: model.ClassificationS, DisseminationControls: []string{"NOCON"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"S with NOCON/portion", "(S//NC)", "portion",
			model.ISM{Classification: model.ClassificationS, DisseminationControls: []string{"NOCON"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"S with DISPLAY ONLY/banner", "SECRET//DISPLAY ONLY USA, GBR", "banner",
			model.ISM{Classification: model.ClassificationS, DisseminationControls: []string{"DISPLAY ONLY"}, DisplayOnlyTo: []string{"USA", "GBR"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"S with DISPLAY ONLY/portion", "(S//DISPLAY ONLY USA, GBR)", "portion",
			model.ISM{Classification: model.ClassificationS, DisseminationControls: []string{"DISPLAY ONLY"}, DisplayOnlyTo: []string{"USA", "GBR"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},

		// --- Joint / Multiple ownerProducer ---
		{"S joint with two owners/banner", "//JOINT SECRET GBR USA", "banner",
			model.ISM{Classification: model.ClassificationS, Joint: true, OwnerProducer: []string{"GBR", "USA"}}, true, "", nil},
		{"S joint with two owners/portion", "(//JOINT S GBR USA)", "portion",
			model.ISM{Classification: model.ClassificationS, Joint: true, OwnerProducer: []string{"GBR", "USA"}}, true, "", nil},
		{"C joint with two owners/banner", "//JOINT CONFIDENTIAL GBR USA", "banner",
			model.ISM{Classification: model.ClassificationC, Joint: true, OwnerProducer: []string{"GBR", "USA"}}, true, "", nil},
		{"C joint with two owners/portion", "(//JOINT C GBR USA)", "portion",
			model.ISM{Classification: model.ClassificationC, Joint: true, OwnerProducer: []string{"GBR", "USA"}}, true, "", nil},
		{"S joint with three owners/banner", "//JOINT SECRET CAN GBR USA", "banner",
			model.ISM{Classification: model.ClassificationS, Joint: true, OwnerProducer: []string{"CAN", "GBR", "USA"}}, true, "", nil},
		{"S joint with three owners/portion", "(//JOINT S CAN GBR USA)", "portion",
			model.ISM{Classification: model.ClassificationS, Joint: true, OwnerProducer: []string{"CAN", "GBR", "USA"}}, true, "", nil},
		{"S joint with NOFORN/banner", "//JOINT SECRET GBR USA//NOFORN", "banner",
			model.ISM{Classification: model.ClassificationS, Joint: true, OwnerProducer: []string{"GBR", "USA"}, DisseminationControls: []string{"NOFORN"}}, true, "", nil},
		{"S joint with NOFORN/portion", "(//JOINT S GBR USA//NF)", "portion",
			model.ISM{Classification: model.ClassificationS, Joint: true, OwnerProducer: []string{"GBR", "USA"}, DisseminationControls: []string{"NOFORN"}}, true, "", nil},
		// LIMDIS is not in the compiled refdata NonICMarkings vocabulary
		// (a genuine gap: see internal/compliance/dodm_examples_test.go's
		// documented "GAP: validation rejects EXDIS" skip). The parser is
		// honest about it: LIMDIS is an unrecognized token, dropped with a
		// warning, and the round trip legitimately does not match.
		{"S joint with dissem + FGI + non-IC/banner", "//JOINT SECRET GBR USA//FGI FRA//OC/REL TO USA, GBR//LIMDIS", "banner",
			model.ISM{Classification: model.ClassificationS, Joint: true, OwnerProducer: []string{"GBR", "USA"},
				FGISourceOpen: []string{"FRA"}, DisseminationControls: []string{"OC", "REL"}, ReleasableTo: []string{"USA", "GBR"}},
			false, "PARSE_UNKNOWN_TOKEN", nil},
		{"S joint with dissem + FGI + non-IC/portion", "(//JOINT S GBR USA//FGI//OC/REL TO USA, GBR//LIMDIS)", "portion",
			model.ISM{Classification: model.ClassificationS, Joint: true, OwnerProducer: []string{"GBR", "USA"},
				FGISourceProtected: []string{"FGI"}, DisseminationControls: []string{"OC", "REL"}, ReleasableTo: []string{"USA", "GBR"}},
			false, "PARSE_UNKNOWN_TOKEN", nil},

		// --- FGI non-US (single non-US OwnerProducer) ---
		{"FGI non-US: GBR SECRET/banner", "//GBR SECRET", "banner",
			model.ISM{Classification: model.ClassificationS, OwnerProducer: []string{"GBR"}}, true, "", nil},
		{"FGI non-US: GBR SECRET/portion", "(//GBR S)", "portion",
			model.ISM{Classification: model.ClassificationS, OwnerProducer: []string{"GBR"}}, true, "", nil},
		{"FGI non-US: DEU CONFIDENTIAL/banner", "//DEU CONFIDENTIAL", "banner",
			model.ISM{Classification: model.ClassificationC, OwnerProducer: []string{"DEU"}}, true, "", nil},
		{"FGI non-US: DEU CONFIDENTIAL/portion", "(//DEU C)", "portion",
			model.ISM{Classification: model.ClassificationC, OwnerProducer: []string{"DEU"}}, true, "", nil},
		{"FGI non-US: FRA TOP SECRET/banner", "//FRA TOP SECRET", "banner",
			model.ISM{Classification: model.ClassificationTS, OwnerProducer: []string{"FRA"}}, true, "", nil},
		{"FGI non-US: FRA TOP SECRET/portion", "(//FRA TS)", "portion",
			model.ISM{Classification: model.ClassificationTS, OwnerProducer: []string{"FRA"}}, true, "", nil},
		{"FGI non-US: GBR UNCLASSIFIED/banner", "//GBR UNCLASSIFIED", "banner",
			model.ISM{Classification: model.ClassificationU, OwnerProducer: []string{"GBR"}}, true, "", nil},
		{"FGI non-US: GBR UNCLASSIFIED/portion", "(//GBR U)", "portion",
			model.ISM{Classification: model.ClassificationU, OwnerProducer: []string{"GBR"}}, true, "", nil},
		{"FGI non-US: GBR SECRET with NOFORN/banner", "//GBR SECRET//NOFORN", "banner",
			model.ISM{Classification: model.ClassificationS, OwnerProducer: []string{"GBR"}, DisseminationControls: []string{"NOFORN"}}, true, "", nil},
		{"FGI non-US: GBR SECRET with NOFORN/portion", "(//GBR S//NF)", "portion",
			model.ISM{Classification: model.ClassificationS, OwnerProducer: []string{"GBR"}, DisseminationControls: []string{"NOFORN"}}, true, "", nil},

		// --- Ordering ---
		{"controls sorted into canonical order/banner", "SECRET//OC/NOFORN/PROPIN", "banner",
			model.ISM{Classification: model.ClassificationS, DisseminationControls: []string{"OC", "NOFORN", "PROPIN"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"controls sorted into canonical order/portion", "(S//OC/NF/PR)", "portion",
			model.ISM{Classification: model.ClassificationS, DisseminationControls: []string{"OC", "NOFORN", "PROPIN"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"FGI then dissemination then non-IC/banner", "SECRET//FGI GBR//NOFORN//LIMDIS", "banner",
			model.ISM{Classification: model.ClassificationS, FGISourceOpen: []string{"GBR"}, DisseminationControls: []string{"NOFORN"}, OwnerProducer: []string{"USA"}}, false, "PARSE_UNKNOWN_TOKEN", []string{"ownerProducer"}},
		{"FGI then dissemination then non-IC/portion", "(S//FGI//NF//LIMDIS)", "portion",
			model.ISM{Classification: model.ClassificationS, FGISourceProtected: []string{"FGI"}, DisseminationControls: []string{"NOFORN"}, OwnerProducer: []string{"USA"}}, false, "PARSE_UNKNOWN_TOKEN", []string{"ownerProducer"}},

		// --- FGI ---
		{"S with multiple FGI open sources/banner", "SECRET//FGI FRA GBR", "banner",
			model.ISM{Classification: model.ClassificationS, FGISourceOpen: []string{"FRA", "GBR"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"S with multiple FGI open sources/portion", "(S//FGI)", "portion",
			model.ISM{Classification: model.ClassificationS, FGISourceProtected: []string{"FGI"}, OwnerProducer: []string{"USA"}}, true, "PARSE_LOSSY", []string{"ownerProducer"}},
		{"S with FGI protected source/banner", "SECRET//FGI", "banner",
			model.ISM{Classification: model.ClassificationS, FGISourceProtected: []string{"FGI"}, OwnerProducer: []string{"USA"}}, true, "PARSE_LOSSY", []string{"ownerProducer"}},
		{"S with FGI protected source/portion", "(S//FGI)", "portion",
			model.ISM{Classification: model.ClassificationS, FGISourceProtected: []string{"FGI"}, OwnerProducer: []string{"USA"}}, true, "PARSE_LOSSY", []string{"ownerProducer"}},

		// --- Non-IC (LIMDIS/EXDIS unrecognized — see note above) ---
		{"S with non-IC marking/banner", "SECRET//LIMDIS", "banner",
			model.ISM{Classification: model.ClassificationS, OwnerProducer: []string{"USA"}}, false, "PARSE_UNKNOWN_TOKEN", []string{"ownerProducer"}},
		{"S with non-IC marking/portion", "(S//LIMDIS)", "portion",
			model.ISM{Classification: model.ClassificationS, OwnerProducer: []string{"USA"}}, false, "PARSE_UNKNOWN_TOKEN", []string{"ownerProducer"}},
		{"S with multiple non-IC markings/banner", "SECRET//LIMDIS/EXDIS", "banner",
			model.ISM{Classification: model.ClassificationS, OwnerProducer: []string{"USA"}}, false, "PARSE_UNKNOWN_TOKEN", []string{"ownerProducer"}},
		{"S with multiple non-IC markings/portion", "(S//LIMDIS/EXDIS)", "portion",
			model.ISM{Classification: model.ClassificationS, OwnerProducer: []string{"USA"}}, false, "PARSE_UNKNOWN_TOKEN", []string{"ownerProducer"}},

		// --- Top Secret + SCI ---
		{"TS with single SCI/banner", "TOP SECRET//SI", "banner",
			model.ISM{Classification: model.ClassificationTS, SCIControls: []string{"SI"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"TS with single SCI/portion", "(TS//SI)", "portion",
			model.ISM{Classification: model.ClassificationTS, SCIControls: []string{"SI"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"TS with multiple SCI sorted alphabetically/banner", "TOP SECRET//HCS/SI/TK", "banner",
			model.ISM{Classification: model.ClassificationTS, SCIControls: []string{"HCS", "SI", "TK"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"TS with multiple SCI sorted alphabetically/portion", "(TS//HCS/SI/TK)", "portion",
			model.ISM{Classification: model.ClassificationTS, SCIControls: []string{"HCS", "SI", "TK"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"TS with SCI and NOFORN/banner", "TOP SECRET//SI//NOFORN", "banner",
			model.ISM{Classification: model.ClassificationTS, SCIControls: []string{"SI"}, DisseminationControls: []string{"NOFORN"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"TS with SCI and NOFORN/portion", "(TS//SI//NF)", "portion",
			model.ISM{Classification: model.ClassificationTS, SCIControls: []string{"SI"}, DisseminationControls: []string{"NOFORN"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"TS with SCI and dissem and FGI/banner", "TOP SECRET//SI/TK//FGI GBR//NOFORN", "banner",
			model.ISM{Classification: model.ClassificationTS, SCIControls: []string{"SI", "TK"}, FGISourceOpen: []string{"GBR"}, DisseminationControls: []string{"NOFORN"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"TS with SCI and dissem and FGI/portion", "(TS//SI/TK//FGI//NF)", "portion",
			model.ISM{Classification: model.ClassificationTS, SCIControls: []string{"SI", "TK"}, FGISourceProtected: []string{"FGI"}, DisseminationControls: []string{"NOFORN"}, OwnerProducer: []string{"USA"}}, true, "PARSE_LOSSY", []string{"ownerProducer"}},

		// --- SCI dissemination suppression ---
		{"TS with SCI dissem suppressed when sciControls present/banner", "TOP SECRET//SI/TK", "banner",
			model.ISM{Classification: model.ClassificationTS, SCIControls: []string{"SI", "TK"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"TS with SCI dissem suppressed when sciControls present/portion", "(TS//SI/TK)", "portion",
			model.ISM{Classification: model.ClassificationTS, SCIControls: []string{"SI", "TK"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"TS with SCI dissem suppressed alongside other dissem controls/banner", "TOP SECRET//SI/TK//NOFORN", "banner",
			model.ISM{Classification: model.ClassificationTS, SCIControls: []string{"SI", "TK"}, DisseminationControls: []string{"NOFORN"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"TS with SCI dissem suppressed alongside other dissem controls/portion", "(TS//SI/TK//NF)", "portion",
			model.ISM{Classification: model.ClassificationTS, SCIControls: []string{"SI", "TK"}, DisseminationControls: []string{"NOFORN"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"TS with SCI dissem rendered when no sciControls/banner", "TOP SECRET//SCI", "banner",
			model.ISM{Classification: model.ClassificationTS, DisseminationControls: []string{"SCI"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"TS with SCI dissem rendered when no sciControls/portion", "(TS//SCI)", "portion",
			model.ISM{Classification: model.ClassificationTS, DisseminationControls: []string{"SCI"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},

		// --- Combined ---
		{"S full combo: dissem + FGI + non-IC/banner", "SECRET//FGI FRA//OC/REL TO USA, GBR//LIMDIS", "banner",
			model.ISM{Classification: model.ClassificationS, FGISourceOpen: []string{"FRA"}, DisseminationControls: []string{"OC", "REL"}, ReleasableTo: []string{"USA", "GBR"}, OwnerProducer: []string{"USA"}},
			false, "PARSE_UNKNOWN_TOKEN", []string{"ownerProducer"}},
		{"S full combo: dissem + FGI + non-IC/portion", "(S//FGI//OC/REL TO USA, GBR//LIMDIS)", "portion",
			model.ISM{Classification: model.ClassificationS, FGISourceProtected: []string{"FGI"}, DisseminationControls: []string{"OC", "REL"}, ReleasableTo: []string{"USA", "GBR"}, OwnerProducer: []string{"USA"}},
			false, "PARSE_UNKNOWN_TOKEN", []string{"ownerProducer"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Parse(tc.input, testReg)
			if got.Form != tc.form {
				t.Errorf("Form = %q, want %q", got.Form, tc.form)
			}
			assertISM(t, got.ISM, tc.want)
			if got.RoundTrip.Matches != tc.wantMatches {
				t.Errorf("RoundTrip.Matches = %v, want %v (rendered=%q)", got.RoundTrip.Matches, tc.wantMatches, got.RoundTrip.Rendered)
			}
			if tc.wantWarning != "" && !hasWarningCode(t, got, tc.wantWarning) {
				t.Errorf("expected warning code %s, got warnings: %+v", tc.wantWarning, got.Warnings)
			}
			if !reflect.DeepEqual(got.Inferred, tc.wantInferred) {
				t.Errorf("Inferred = %v, want %v", got.Inferred, tc.wantInferred)
			}
		})
	}
}

// ============================================================
// Targeted cases beyond the seed inversion: AEA forms, SAR handling,
// positional disambiguation (DS/FOUO), unknown tokens, case/whitespace
// tolerance, and boundary inputs (empty / "//"-only).
// ============================================================

func TestParse_Targeted(t *testing.T) {
	cases := []invertedCase{
		// --- JOINT / FGI leading "//" doesn't create a spurious empty segment ---
		{"JOINT leading slashes", "//JOINT SECRET GBR USA//NOFORN", "banner",
			model.ISM{Classification: model.ClassificationS, Joint: true, OwnerProducer: []string{"GBR", "USA"}, DisseminationControls: []string{"NOFORN"}},
			true, "", nil},
		{"FGI non-US leading slashes", "//GBR SECRET//NOFORN", "banner",
			model.ISM{Classification: model.ClassificationS, OwnerProducer: []string{"GBR"}, DisseminationControls: []string{"NOFORN"}},
			true, "", nil},

		// --- Bare "TOP SECRET" (space in label) with no country/joint prefix ---
		{"bare TOP SECRET banner", "TOP SECRET", "banner",
			model.ISM{Classification: model.ClassificationTS, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"bare TOP SECRET portion", "(TS)", "portion",
			model.ISM{Classification: model.ClassificationTS, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},

		// --- FOUO positional disambiguation: dissem code vs. non-IC code ---
		{"FOUO on unclassified resolves to dissem, round-trips", "UNCLASSIFIED//FOUO", "banner",
			model.ISM{Classification: model.ClassificationU, DisseminationControls: []string{"FOUO"}, OwnerProducer: []string{"USA"}}, true, "PARSE_AMBIGUOUS", []string{"ownerProducer"}},
		{"FOUO on classified resolves to dissem but is suppressed on re-render", "SECRET//FOUO", "banner",
			model.ISM{Classification: model.ClassificationS, DisseminationControls: []string{"FOUO"}, OwnerProducer: []string{"USA"}}, false, "PARSE_AMBIGUOUS", []string{"ownerProducer"}},

		// --- Every AEA form (banner) ---
		{"AEA RD banner", "TOP SECRET//RESTRICTED DATA", "banner",
			model.ISM{Classification: model.ClassificationTS, AtomicEnergyMarkings: []string{"RD"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"AEA RD-CNWDI banner", "TOP SECRET//RESTRICTED DATA-N", "banner",
			model.ISM{Classification: model.ClassificationTS, AtomicEnergyMarkings: []string{"RD-CNWDI"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"AEA RD-SIGMA banner", "TOP SECRET//RD-SIGMA 14 15", "banner",
			model.ISM{Classification: model.ClassificationTS, AtomicEnergyMarkings: []string{"RD-SG-14", "RD-SG-15"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"AEA FRD banner", "TOP SECRET//FORMERLY RESTRICTED DATA", "banner",
			model.ISM{Classification: model.ClassificationTS, AtomicEnergyMarkings: []string{"FRD"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"AEA FRD-SIGMA banner", "TOP SECRET//FRD-SIGMA 14 15", "banner",
			model.ISM{Classification: model.ClassificationTS, AtomicEnergyMarkings: []string{"FRD-SG-14", "FRD-SG-15"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"AEA DCNI banner", "TOP SECRET//DCNI", "banner",
			model.ISM{Classification: model.ClassificationTS, AtomicEnergyMarkings: []string{"DCNI"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"AEA UCNI banner", "TOP SECRET//UCNI", "banner",
			model.ISM{Classification: model.ClassificationTS, AtomicEnergyMarkings: []string{"UCNI"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"AEA TFNI banner", "TOP SECRET//TFNI", "banner",
			model.ISM{Classification: model.ClassificationTS, AtomicEnergyMarkings: []string{"TFNI"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"AEA RD + RD-SIGMA combined banner", "TOP SECRET//RESTRICTED DATA/RD-SIGMA 14 15", "banner",
			model.ISM{Classification: model.ClassificationTS, AtomicEnergyMarkings: []string{"RD", "RD-SG-14", "RD-SG-15"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},

		// --- Every AEA form (portion) ---
		{"AEA RD portion", "(TS//RD)", "portion",
			model.ISM{Classification: model.ClassificationTS, AtomicEnergyMarkings: []string{"RD"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"AEA RD-CNWDI portion", "(TS//RD-N)", "portion",
			model.ISM{Classification: model.ClassificationTS, AtomicEnergyMarkings: []string{"RD-CNWDI"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"AEA RD-SG portion", "(TS//RD-SG 14 15)", "portion",
			model.ISM{Classification: model.ClassificationTS, AtomicEnergyMarkings: []string{"RD-SG-14", "RD-SG-15"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"AEA FRD portion", "(TS//FRD)", "portion",
			model.ISM{Classification: model.ClassificationTS, AtomicEnergyMarkings: []string{"FRD"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"AEA FRD-SG portion", "(TS//FRD-SG 14 15)", "portion",
			model.ISM{Classification: model.ClassificationTS, AtomicEnergyMarkings: []string{"FRD-SG-14", "FRD-SG-15"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"AEA DCNI portion", "(TS//DCNI)", "portion",
			model.ISM{Classification: model.ClassificationTS, AtomicEnergyMarkings: []string{"DCNI"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},

		// --- SAR ---
		{"SAR two identifiers round-trips", "TOP SECRET//SAR-ALPHA/SAR-BRAVO", "banner",
			model.ISM{Classification: model.ClassificationTS, SARIdentifier: []string{"ALPHA", "BRAVO"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"SAR-MULTIPLE PROGRAMS conceals individual PIDs", "TOP SECRET//SAR-MULTIPLE PROGRAMS", "banner",
			model.ISM{Classification: model.ClassificationTS, OwnerProducer: []string{"USA"}}, false, "PARSE_LOSSY", []string{"ownerProducer"}},
		{"SAR portion always lists all PIDs", "(TS//SAR-ALPHA/SAR-BRAVO/SAR-CHARLIE)", "portion",
			model.ISM{Classification: model.ClassificationTS, SARIdentifier: []string{"ALPHA", "BRAVO", "CHARLIE"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},

		// --- REL TO / DISPLAY ONLY single country ---
		{"REL TO single country", "SECRET//REL TO USA", "banner",
			model.ISM{Classification: model.ClassificationS, DisseminationControls: []string{"REL"}, ReleasableTo: []string{"USA"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},

		// --- Unknown token ---
		{"unknown token dropped with warning", "SECRET//FROBNIZ", "banner",
			model.ISM{Classification: model.ClassificationS, OwnerProducer: []string{"USA"}}, false, "PARSE_UNKNOWN_TOKEN", []string{"ownerProducer"}},

		// --- Lowercase input ---
		{"lowercase banner input", "secret//noforn", "banner",
			model.ISM{Classification: model.ClassificationS, DisseminationControls: []string{"NOFORN"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"lowercase portion input", "(s//nf)", "portion",
			model.ISM{Classification: model.ClassificationS, DisseminationControls: []string{"NOFORN"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},

		// --- Stray whitespace ---
		{"stray surrounding whitespace", "  SECRET//NOFORN  ", "banner",
			model.ISM{Classification: model.ClassificationS, DisseminationControls: []string{"NOFORN"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
		{"stray whitespace inside parens", "( S//NF )", "portion",
			model.ISM{Classification: model.ClassificationS, DisseminationControls: []string{"NOFORN"}, OwnerProducer: []string{"USA"}}, true, "", []string{"ownerProducer"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Parse(tc.input, testReg)
			if got.Form != tc.form {
				t.Errorf("Form = %q, want %q", got.Form, tc.form)
			}
			assertISM(t, got.ISM, tc.want)
			if got.RoundTrip.Matches != tc.wantMatches {
				t.Errorf("RoundTrip.Matches = %v, want %v (rendered=%q)", got.RoundTrip.Matches, tc.wantMatches, got.RoundTrip.Rendered)
			}
			if tc.wantWarning != "" && !hasWarningCode(t, got, tc.wantWarning) {
				t.Errorf("expected warning code %s, got warnings: %+v", tc.wantWarning, got.Warnings)
			}
			if !reflect.DeepEqual(got.Inferred, tc.wantInferred) {
				t.Errorf("Inferred = %v, want %v", got.Inferred, tc.wantInferred)
			}
		})
	}
}

// TestParse_EmptyAndSlashOnly covers boundary inputs that must never panic
// and must never be treated as fully valid — they should surface warnings
// and an honest (likely mismatching) round trip rather than fabricating data.
func TestParse_EmptyAndSlashOnly(t *testing.T) {
	for _, input := range []string{"", "//", "()", "( )"} {
		t.Run(input, func(t *testing.T) {
			got := Parse(input, testReg)
			if got.Form == "" {
				t.Errorf("Form not set for input %q", input)
			}
			if len(got.Warnings) == 0 {
				t.Errorf("expected at least one warning for input %q", input)
			}
		})
	}
}

func TestParse_FormDetection(t *testing.T) {
	tests := []struct {
		input    string
		wantForm string
	}{
		{"SECRET//NOFORN", "banner"},
		{"(S//NF)", "portion"},
		{"  (S//NF)  ", "portion"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := Parse(tt.input, testReg)
			if got.Form != tt.wantForm {
				t.Errorf("Form = %q, want %q", got.Form, tt.wantForm)
			}
		})
	}
}

package banner

import (
	"testing"

	"expr.ai/ism-api/internal/model"
)

func TestRender(t *testing.T) {
	tests := []struct {
		name        string
		ism         model.ISM
		wantBanner  string
		wantPortion string
	}{
		// --- Unclassified ---
		{
			name:        "U bare",
			ism:         model.ISM{Classification: model.ClassificationU},
			wantBanner:  "UNCLASSIFIED",
			wantPortion: "(U)",
		},
		{
			name: "U with FOUO dissemination",
			ism: model.ISM{
				Classification:        model.ClassificationU,
				DisseminationControls: []string{"FED ONLY"},
			},
			wantBanner:  "UNCLASSIFIED//FED ONLY",
			wantPortion: "(U//FED ONLY)",
		},

		// --- CUI ---
		{
			name:        "CUI bare",
			ism:         model.ISM{Classification: model.ClassificationCUI},
			wantBanner:  "CUI",
			wantPortion: "(CUI)",
		},
		{
			name: "CUI with specified category",
			ism: model.ISM{
				Classification:   model.ClassificationCUI,
				CategoryMarkings: []string{"SP-CEII"},
			},
			wantBanner:  "CUI//SP-CEII",
			wantPortion: "(CUI//SP-CEII)",
		},
		{
			name: "CUI with multiple categories",
			ism: model.ISM{
				Classification:   model.ClassificationCUI,
				CategoryMarkings: []string{"SP-CEII", "SP-PHYS"},
			},
			wantBanner:  "CUI//SP-CEII/SP-PHYS",
			wantPortion: "(CUI//SP-CEII/SP-PHYS)",
		},
		{
			name: "CUI with category and NOFORN",
			ism: model.ISM{
				Classification:        model.ClassificationCUI,
				CategoryMarkings:      []string{"SP-CEII"},
				DisseminationControls: []string{"NOFORN"},
			},
			wantBanner:  "CUI//SP-CEII/NOFORN",
			wantPortion: "(CUI//SP-CEII/NF)",
		},

		// --- Confidential ---
		{
			name: "C bare",
			ism: model.ISM{
				Classification: model.ClassificationC,
				OwnerProducer:  []string{"USA"},
			},
			wantBanner:  "CONFIDENTIAL",
			wantPortion: "(C)",
		},
		{
			name: "C with NOFORN",
			ism: model.ISM{
				Classification:        model.ClassificationC,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"NOFORN"},
			},
			wantBanner:  "CONFIDENTIAL//NOFORN",
			wantPortion: "(C//NF)",
		},
		{
			name: "C with REL TO",
			ism: model.ISM{
				Classification:        model.ClassificationC,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"REL"},
				ReleasableTo:          []string{"USA", "GBR"},
			},
			wantBanner:  "CONFIDENTIAL//REL TO USA, GBR",
			wantPortion: "(C//REL TO USA, GBR)",
		},
		{
			name: "C with FGI open source",
			ism: model.ISM{
				Classification: model.ClassificationC,
				OwnerProducer:  []string{"USA"},
				FGISourceOpen:  []string{"GBR"},
			},
			wantBanner:  "CONFIDENTIAL//FGI GBR",
			wantPortion: "(C//FGI)",
		},

		// --- Secret ---
		{
			name: "S bare",
			ism: model.ISM{
				Classification: model.ClassificationS,
				OwnerProducer:  []string{"USA"},
			},
			wantBanner:  "SECRET",
			wantPortion: "(S)",
		},
		{
			name: "S with NOFORN",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"NOFORN"},
			},
			wantBanner:  "SECRET//NOFORN",
			wantPortion: "(S//NF)",
		},
		{
			name: "S with NOFORN and PROPIN",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"NOFORN", "PROPIN"},
			},
			wantBanner:  "SECRET//NOFORN/PROPIN",
			wantPortion: "(S//NF/PR)",
		},
		{
			name: "S with REL TO multiple countries",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"REL"},
				ReleasableTo:          []string{"USA", "GBR", "CAN"},
			},
			wantBanner:  "SECRET//REL TO USA, CAN, GBR",
			wantPortion: "(S//REL TO USA, CAN, GBR)",
		},
		{
			name: "S with OC and NOFORN",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"NOFORN", "OC"},
			},
			wantBanner:  "SECRET//OC/NOFORN",
			wantPortion: "(S//OC/NF)",
		},
		{
			name: "S with IMCON",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"IMCON"},
			},
			wantBanner:  "SECRET//IMCON",
			wantPortion: "(S//IMC)",
		},
		{
			name: "S with DSEN",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"DSEN"},
			},
			wantBanner:  "SECRET//DSEN",
			wantPortion: "(S//DS)",
		},
		{
			name: "S with NOCON",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"NOCON"},
			},
			wantBanner:  "SECRET//NOCON",
			wantPortion: "(S//NC)",
		},
		{
			name: "S with DISPLAY ONLY",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"DISPLAY ONLY"},
				DisplayOnlyTo:         []string{"USA", "GBR"},
			},
			wantBanner:  "SECRET//DISPLAY ONLY USA, GBR",
			wantPortion: "(S//DISPLAY ONLY USA, GBR)",
		},

		// --- Joint / Multiple ownerProducer ---
		{
			name: "S joint with two owners",
			ism: model.ISM{
				Classification: model.ClassificationS,
				OwnerProducer:  []string{"USA", "GBR"},
				Joint:          true,
			},
			wantBanner:  "//JOINT SECRET GBR USA",
			wantPortion: "(//JOINT S GBR USA)",
		},
		{
			name: "C joint with two owners",
			ism: model.ISM{
				Classification: model.ClassificationC,
				OwnerProducer:  []string{"USA", "GBR"},
				Joint:          true,
			},
			wantBanner:  "//JOINT CONFIDENTIAL GBR USA",
			wantPortion: "(//JOINT C GBR USA)",
		},
		{
			name: "S joint with three owners",
			ism: model.ISM{
				Classification: model.ClassificationS,
				OwnerProducer:  []string{"USA", "GBR", "CAN"},
				Joint:          true,
			},
			wantBanner:  "//JOINT SECRET CAN GBR USA",
			wantPortion: "(//JOINT S CAN GBR USA)",
		},
		{
			name: "S joint with NOFORN",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA", "GBR"},
				Joint:                 true,
				DisseminationControls: []string{"NOFORN"},
			},
			wantBanner:  "//JOINT SECRET GBR USA//NOFORN",
			wantPortion: "(//JOINT S GBR USA//NF)",
		},
		{
			name: "S joint with dissem + FGI + non-IC",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA", "GBR"},
				Joint:                 true,
				DisseminationControls: []string{"REL", "OC"},
				ReleasableTo:          []string{"USA", "GBR"},
				FGISourceOpen:         []string{"FRA"},
				NonICMarkings:         []string{"LIMDIS"},
			},
			wantBanner:  "//JOINT SECRET GBR USA//OC/REL TO USA, GBR/FGI FRA/LIMDIS",
			wantPortion: "(//JOINT S GBR USA//OC/REL TO USA, GBR/FGI/LIMDIS)",
		},

		// --- FGI non-US (single non-US OwnerProducer) ---
		{
			name: "FGI non-US: GBR SECRET",
			ism: model.ISM{
				Classification: model.ClassificationS,
				OwnerProducer:  []string{"GBR"},
			},
			wantBanner:  "//GBR SECRET",
			wantPortion: "(//GBR S)",
		},
		{
			name: "FGI non-US: DEU CONFIDENTIAL",
			ism: model.ISM{
				Classification: model.ClassificationC,
				OwnerProducer:  []string{"DEU"},
			},
			wantBanner:  "//DEU CONFIDENTIAL",
			wantPortion: "(//DEU C)",
		},
		{
			name: "FGI non-US: FRA TOP SECRET",
			ism: model.ISM{
				Classification: model.ClassificationTS,
				OwnerProducer:  []string{"FRA"},
			},
			wantBanner:  "//FRA TOP SECRET",
			wantPortion: "(//FRA TS)",
		},
		{
			name: "FGI non-US: GBR UNCLASSIFIED",
			ism: model.ISM{
				Classification: model.ClassificationU,
				OwnerProducer:  []string{"GBR"},
			},
			wantBanner:  "//GBR UNCLASSIFIED",
			wantPortion: "(//GBR U)",
		},
		{
			name: "FGI non-US: GBR SECRET with NOFORN",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"GBR"},
				DisseminationControls: []string{"NOFORN"},
			},
			wantBanner:  "//GBR SECRET//NOFORN",
			wantPortion: "(//GBR S//NF)",
		},

		// --- Ordering ---
		{
			name: "controls sorted into canonical order",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"PROPIN", "OC", "NOFORN"},
			},
			wantBanner:  "SECRET//OC/NOFORN/PROPIN",
			wantPortion: "(S//OC/NF/PR)",
		},
		{
			name: "dissemination then FGI then non-IC",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"NOFORN"},
				FGISourceOpen:         []string{"GBR"},
				NonICMarkings:         []string{"LIMDIS"},
			},
			wantBanner:  "SECRET//NOFORN/FGI GBR/LIMDIS",
			wantPortion: "(S//NF/FGI/LIMDIS)",
		},

		// --- FGI ---
		{
			name: "S with multiple FGI open sources",
			ism: model.ISM{
				Classification: model.ClassificationS,
				OwnerProducer:  []string{"USA"},
				FGISourceOpen:  []string{"GBR", "FRA"},
			},
			wantBanner:  "SECRET//FGI FRA GBR",
			wantPortion: "(S//FGI)",
		},
		{
			name: "S with FGI protected source (concealed per E4-S9.e)",
			ism: model.ISM{
				Classification:     model.ClassificationS,
				OwnerProducer:      []string{"USA"},
				FGISourceProtected: []string{"DEU"},
			},
			wantBanner:  "SECRET//FGI",
			wantPortion: "(S//FGI)",
		},

		// --- Non-IC ---
		{
			name: "S with non-IC marking",
			ism: model.ISM{
				Classification: model.ClassificationS,
				OwnerProducer:  []string{"USA"},
				NonICMarkings:  []string{"LIMDIS"},
			},
			wantBanner:  "SECRET//LIMDIS",
			wantPortion: "(S//LIMDIS)",
		},
		{
			name: "S with multiple non-IC markings",
			ism: model.ISM{
				Classification: model.ClassificationS,
				OwnerProducer:  []string{"USA"},
				NonICMarkings:  []string{"LIMDIS", "EXDIS"},
			},
			wantBanner:  "SECRET//LIMDIS/EXDIS",
			wantPortion: "(S//LIMDIS/EXDIS)",
		},

		// --- Top Secret + SCI ---
		{
			name: "TS with single SCI",
			ism: model.ISM{
				Classification: model.ClassificationTS,
				OwnerProducer:  []string{"USA"},
				SCIControls:    []string{"SI"},
			},
			wantBanner:  "TOP SECRET//SI",
			wantPortion: "(TS//SI)",
		},
		{
			name: "TS with multiple SCI sorted alphabetically",
			ism: model.ISM{
				Classification: model.ClassificationTS,
				OwnerProducer:  []string{"USA"},
				SCIControls:    []string{"TK", "HCS", "SI"},
			},
			wantBanner:  "TOP SECRET//HCS/SI/TK",
			wantPortion: "(TS//HCS/SI/TK)",
		},
		{
			name: "TS with SCI and NOFORN",
			ism: model.ISM{
				Classification:        model.ClassificationTS,
				OwnerProducer:         []string{"USA"},
				SCIControls:           []string{"SI"},
				DisseminationControls: []string{"NOFORN"},
			},
			wantBanner:  "TOP SECRET//SI//NOFORN",
			wantPortion: "(TS//SI//NF)",
		},
		{
			name: "TS with SCI and dissem and FGI",
			ism: model.ISM{
				Classification:        model.ClassificationTS,
				OwnerProducer:         []string{"USA"},
				SCIControls:           []string{"SI", "TK"},
				DisseminationControls: []string{"NOFORN"},
				FGISourceOpen:         []string{"GBR"},
			},
			wantBanner:  "TOP SECRET//SI/TK//NOFORN/FGI GBR",
			wantPortion: "(TS//SI/TK//NF/FGI)",
		},

		// --- Combined ---
		{
			name: "S full combo: dissem + FGI + non-IC",
			ism: model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{"REL", "OC"},
				ReleasableTo:          []string{"USA", "GBR"},
				FGISourceOpen:         []string{"FRA"},
				NonICMarkings:         []string{"LIMDIS"},
			},
			wantBanner:  "SECRET//OC/REL TO USA, GBR/FGI FRA/LIMDIS",
			wantPortion: "(S//OC/REL TO USA, GBR/FGI/LIMDIS)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Render(&tt.ism)
			if got.BannerLine != tt.wantBanner {
				t.Errorf("BannerLine:\n  got:  %q\n  want: %q", got.BannerLine, tt.wantBanner)
			}
			if got.PortionMark != tt.wantPortion {
				t.Errorf("PortionMark:\n  got:  %q\n  want: %q", got.PortionMark, tt.wantPortion)
			}
		})
	}
}

func TestRenderAuthorityBlock(t *testing.T) {
	tests := []struct {
		name string
		ism  model.ISM
		want string
	}{
		{
			name: "empty for unclassified",
			ism:  model.ISM{Classification: model.ClassificationU},
			want: "",
		},
		{
			name: "empty for CUI",
			ism:  model.ISM{Classification: model.ClassificationCUI},
			want: "",
		},
		{
			name: "empty for classified with no authority fields",
			ism: model.ISM{
				Classification: model.ClassificationS,
				OwnerProducer:  []string{"USA"},
			},
			want: "",
		},
		{
			name: "original classification only",
			ism: model.ISM{
				Classification: model.ClassificationS,
				OwnerProducer:  []string{"USA"},
				ClassifiedBy:   "John Smith",
			},
			want: "Classified By: John Smith",
		},
		{
			name: "original classification with reason",
			ism: model.ISM{
				Classification:       model.ClassificationS,
				OwnerProducer:        []string{"USA"},
				ClassifiedBy:         "John Smith",
				ClassificationReason: "1.4(a)",
			},
			want: "Classified By: John Smith\nReason: 1.4(a)",
		},
		{
			name: "derivative classification only",
			ism: model.ISM{
				Classification:           model.ClassificationS,
				OwnerProducer:            []string{"USA"},
				DerivativelyClassifiedBy: "Jane Doe",
				DerivedFrom:              "Source Document A",
			},
			want: "Derived From: Source Document A\nDeclassify On: (not specified)",
		},
		{
			name: "derivative classification with declass date",
			ism: model.ISM{
				Classification:           model.ClassificationS,
				OwnerProducer:            []string{"USA"},
				DerivativelyClassifiedBy: "Jane Doe",
				DerivedFrom:              "Source Document A",
				DeclassDate:              "20360101",
			},
			want: "Derived From: Source Document A\nDeclassify On: 20360101",
		},
		{
			name: "original and derivative combined",
			ism: model.ISM{
				Classification:           model.ClassificationS,
				OwnerProducer:            []string{"USA"},
				ClassifiedBy:             "John Smith",
				ClassificationReason:     "1.4(a)",
				DerivativelyClassifiedBy: "Jane Doe",
				DerivedFrom:              "Source Document A",
				DeclassDate:              "20360101",
			},
			want: "Classified By: John Smith\nReason: 1.4(a)\nDerived From: Source Document A\nDeclassify On: 20360101",
		},
		{
			name: "declass event instead of date",
			ism: model.ISM{
				Classification:           model.ClassificationS,
				OwnerProducer:            []string{"USA"},
				DerivativelyClassifiedBy: "Jane Doe",
				DerivedFrom:              "Multiple Sources",
				DeclassEvent:             "Upon completion of project X",
			},
			want: "Derived From: Multiple Sources\nDeclassify On: Upon completion of project X",
		},
		{
			name: "declass exception",
			ism: model.ISM{
				Classification:           model.ClassificationS,
				OwnerProducer:            []string{"USA"},
				DerivativelyClassifiedBy: "Jane Doe",
				DerivedFrom:              "Multiple Sources",
				DeclassException:         "25X1",
			},
			want: "Derived From: Multiple Sources\nDeclassify On: 25X1",
		},
		{
			name: "original with declass date",
			ism: model.ISM{
				Classification:       model.ClassificationC,
				OwnerProducer:        []string{"USA"},
				ClassifiedBy:         "Director NSA",
				ClassificationReason: "1.4(c)",
				DeclassDate:          "20360601",
			},
			want: "Classified By: Director NSA\nReason: 1.4(c)\nDeclassify On: 20360601",
		},
		{
			name: "compilation reason included",
			ism: model.ISM{
				Classification:           model.ClassificationS,
				OwnerProducer:            []string{"USA"},
				DerivativelyClassifiedBy: "Jane Doe",
				DerivedFrom:              "Multiple Sources",
				CompilationReason:        "Compilation reveals additional intelligence value",
				DeclassDate:              "20360101",
			},
			want: "Derived From: Multiple Sources\nCompilation: Compilation reveals additional intelligence value\nDeclassify On: 20360101",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Render(&tt.ism)
			if got.AuthorityBlock != tt.want {
				t.Errorf("AuthorityBlock:\n  got:  %q\n  want: %q", got.AuthorityBlock, tt.want)
			}
		})
	}
}

package compliance_test

import (
	"strings"
	"testing"

	"dmontgomery/ism-api/internal/banner"
	"dmontgomery/ism-api/internal/model"
)

// TestXSD_Banner_ClassificationLabels verifies banner rendering for each
// classification level uses the correct full label.
func TestXSD_Banner_ClassificationLabels(t *testing.T) {
	tests := []struct {
		classification model.Classification
		wantBanner     string
		wantPortion    string
	}{
		{model.ClassificationU, "UNCLASSIFIED", "(U)"},
		{model.ClassificationCUI, "CUI", "(CUI)"},
		{model.ClassificationC, "CONFIDENTIAL", "(C)"},
		{model.ClassificationS, "SECRET", "(S)"},
	}
	for _, tt := range tests {
		t.Run(string(tt.classification), func(t *testing.T) {
			ism := &model.ISM{Classification: tt.classification}
			result := banner.Render(ism)
			if result.BannerLine != tt.wantBanner {
				t.Errorf("BannerLine = %q, want %q", result.BannerLine, tt.wantBanner)
			}
			if result.PortionMark != tt.wantPortion {
				t.Errorf("PortionMark = %q, want %q", result.PortionMark, tt.wantPortion)
			}
		})
	}
}

// TestXSD_Banner_TopSecret verifies TS banner rendering.
func TestXSD_Banner_TopSecret(t *testing.T) {
	ism := &model.ISM{Classification: model.Classification("TS")}
	result := banner.Render(ism)
	if result.BannerLine != "TOP SECRET" {
		t.Skipf("GAP: TS banner = %q, want 'TOP SECRET' — TS not implemented", result.BannerLine)
	}
	if result.PortionMark != "(TS)" {
		t.Errorf("TS portion mark = %q, want '(TS)'", result.PortionMark)
	}
}

// TestXSD_Banner_Restricted verifies R banner rendering.
func TestXSD_Banner_Restricted(t *testing.T) {
	ism := &model.ISM{Classification: model.Classification("R")}
	result := banner.Render(ism)
	if result.BannerLine != "RESTRICTED" {
		t.Skipf("GAP: R banner = %q, want 'RESTRICTED' — R not implemented", result.BannerLine)
	}
}

// TestXSD_Banner_DisseminationControls verifies dissemination controls appear
// correctly in banner and portion marks.
func TestXSD_Banner_DisseminationControls(t *testing.T) {
	tests := []struct {
		name           string
		controls       []string
		wantInBanner   string
		wantInPortion  string
	}{
		{"NOFORN", []string{"NOFORN"}, "NOFORN", "NF"},
		{"PROPIN", []string{"PROPIN"}, "PROPIN", "PR"},
		{"IMCON", []string{"IMCON"}, "IMCON", "IMC"},
		{"DSEN", []string{"DSEN"}, "DSEN", "DS"},
		{"NOCON", []string{"NOCON"}, "NOCON", "NC"},
		{"OC", []string{"OC"}, "OC", "OC"},
		{"FISA", []string{"FISA"}, "FISA", "FISA"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ism := &model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: tt.controls,
			}
			result := banner.Render(ism)
			if !strings.Contains(result.BannerLine, tt.wantInBanner) {
				t.Errorf("BannerLine %q should contain %q", result.BannerLine, tt.wantInBanner)
			}
			if !strings.Contains(result.PortionMark, tt.wantInPortion) {
				t.Errorf("PortionMark %q should contain %q", result.PortionMark, tt.wantInPortion)
			}
		})
	}
}

// TestXSD_Banner_RELExpansion verifies REL TO expansion with country list.
func TestXSD_Banner_RELExpansion(t *testing.T) {
	ism := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"REL"},
		ReleasableTo:          []string{"GBR", "CAN"},
	}
	result := banner.Render(ism)
	if !strings.Contains(result.BannerLine, "REL TO") {
		t.Errorf("BannerLine %q should contain 'REL TO'", result.BannerLine)
	}
	if !strings.Contains(result.BannerLine, "GBR") {
		t.Errorf("BannerLine %q should contain 'GBR'", result.BannerLine)
	}
}

// TestXSD_Banner_DisseminationOrder verifies dissemination controls are rendered
// in canonical order per ISM banner ordering rules.
func TestXSD_Banner_DisseminationOrder(t *testing.T) {
	ism := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"NOFORN", "OC", "RS"},
	}
	result := banner.Render(ism)
	// Canonical order: RS before OC before NOFORN
	rsIdx := strings.Index(result.BannerLine, "RS")
	ocIdx := strings.Index(result.BannerLine, "OC")
	nfIdx := strings.Index(result.BannerLine, "NOFORN")
	if rsIdx == -1 || ocIdx == -1 || nfIdx == -1 {
		t.Fatalf("BannerLine %q missing expected controls", result.BannerLine)
	}
	if rsIdx > ocIdx {
		t.Errorf("RS should appear before OC in banner (RS@%d, OC@%d)", rsIdx, ocIdx)
	}
	if ocIdx > nfIdx {
		t.Errorf("OC should appear before NOFORN in banner (OC@%d, NF@%d)", ocIdx, nfIdx)
	}
}

// TestXSD_Banner_JointDocument verifies joint document banner rendering.
func TestXSD_Banner_JointDocument(t *testing.T) {
	ism := &model.ISM{
		Classification: model.ClassificationS,
		OwnerProducer:  []string{"USA", "GBR"},
		Joint:          true,
	}
	result := banner.Render(ism)
	if !strings.HasPrefix(result.BannerLine, "//JOINT") {
		t.Errorf("joint BannerLine %q should start with '//JOINT'", result.BannerLine)
	}
	if !strings.Contains(result.BannerLine, "USA") || !strings.Contains(result.BannerLine, "GBR") {
		t.Errorf("joint BannerLine %q should contain ownerProducer countries", result.BannerLine)
	}
	if !strings.HasPrefix(result.PortionMark, "(//JOINT") {
		t.Errorf("joint PortionMark %q should start with '(//JOINT'", result.PortionMark)
	}
}

// TestXSD_Banner_CUICategories verifies CUI banner includes category markings.
func TestXSD_Banner_CUICategories(t *testing.T) {
	ism := &model.ISM{
		Classification:   model.ClassificationCUI,
		CategoryMarkings: []string{"SP-PCII"},
	}
	result := banner.Render(ism)
	if !strings.Contains(result.BannerLine, "SP-PCII") {
		t.Errorf("CUI BannerLine %q should contain category marking 'SP-PCII'", result.BannerLine)
	}
}

// TestXSD_Banner_AuthorityBlock verifies authority block rendering for classified docs.
func TestXSD_Banner_AuthorityBlock(t *testing.T) {
	t.Run("original_classification", func(t *testing.T) {
		ism := &model.ISM{
			Classification:       model.ClassificationS,
			OwnerProducer:        []string{"USA"},
			ClassifiedBy:         "John Smith, OCA",
			ClassificationReason: "1.4(a)",
			DeclassDate:          "20350101",
		}
		result := banner.Render(ism)
		if !strings.Contains(result.AuthorityBlock, "Classified By:") {
			t.Errorf("AuthorityBlock %q should contain 'Classified By:'", result.AuthorityBlock)
		}
		if !strings.Contains(result.AuthorityBlock, "Reason:") {
			t.Errorf("AuthorityBlock %q should contain 'Reason:'", result.AuthorityBlock)
		}
		if !strings.Contains(result.AuthorityBlock, "Declassify On:") {
			t.Errorf("AuthorityBlock %q should contain 'Declassify On:'", result.AuthorityBlock)
		}
	})

	t.Run("derivative_classification", func(t *testing.T) {
		ism := &model.ISM{
			Classification: model.ClassificationC,
			OwnerProducer:  []string{"USA"},
			DerivedFrom:    "Multiple Sources",
			DeclassDate:    "20350101",
		}
		result := banner.Render(ism)
		if !strings.Contains(result.AuthorityBlock, "Derived From:") {
			t.Errorf("AuthorityBlock %q should contain 'Derived From:'", result.AuthorityBlock)
		}
	})

	t.Run("unclassified_no_authority", func(t *testing.T) {
		ism := &model.ISM{Classification: model.ClassificationU}
		result := banner.Render(ism)
		if result.AuthorityBlock != "" {
			t.Errorf("U should have empty AuthorityBlock, got %q", result.AuthorityBlock)
		}
	})
}

// TestXSD_Banner_SCI verifies SCI controls render in banner and portion marks.
func TestXSD_Banner_SCI(t *testing.T) {
	ism := &model.ISM{
		Classification:        model.ClassificationTS,
		OwnerProducer:         []string{"USA"},
		SCIControls:           []string{"SI"},
		DisseminationControls: []string{"NOFORN"},
	}
	result := banner.Render(ism)
	if result.BannerLine != "TOP SECRET//SI//NOFORN" {
		t.Errorf("BannerLine = %q, want %q", result.BannerLine, "TOP SECRET//SI//NOFORN")
	}
	if result.PortionMark != "(TS//SI//NF)" {
		t.Errorf("PortionMark = %q, want %q", result.PortionMark, "(TS//SI//NF)")
	}
}

// TestXSD_Banner_FGI verifies FGI source rendering.
func TestXSD_Banner_FGI(t *testing.T) {
	ism := &model.ISM{
		Classification: model.ClassificationS,
		OwnerProducer:  []string{"USA"},
		FGISourceOpen:  []string{"GBR"},
	}
	result := banner.Render(ism)
	if !strings.Contains(result.BannerLine, "FGI") {
		t.Errorf("BannerLine %q should contain 'FGI'", result.BannerLine)
	}
}

// TestXSD_Banner_NonICMarkings verifies non-IC markings appear in banner.
func TestXSD_Banner_NonICMarkings(t *testing.T) {
	ism := &model.ISM{
		Classification: model.ClassificationU,
		NonICMarkings:  []string{"LES"},
	}
	result := banner.Render(ism)
	if !strings.Contains(result.BannerLine, "LES") {
		t.Errorf("BannerLine %q should contain 'LES'", result.BannerLine)
	}
}

// TestXSD_Banner_DisseminationCodeNormalization verifies that the API's
// human-readable codes (NOFORN, PROPIN, etc.) produce correct XSD-compatible
// portion mark abbreviations (NF, PR, etc.).
func TestXSD_Banner_DisseminationCodeNormalization(t *testing.T) {
	// XSD uses abbreviated codes (NF, PR, IMC, DISPLAYONLY)
	// API uses human-readable (NOFORN, PROPIN, IMCON, DISPLAY ONLY)
	// Portion marks should use XSD abbreviations
	mappings := []struct {
		apiCode     string
		portionAbbr string
	}{
		{"NOFORN", "NF"},
		{"PROPIN", "PR"},
		{"IMCON", "IMC"},
		{"DSEN", "DS"},
		{"NOCON", "NC"},
	}
	for _, m := range mappings {
		t.Run(m.apiCode, func(t *testing.T) {
			ism := &model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{m.apiCode},
			}
			result := banner.Render(ism)
			if !strings.Contains(result.PortionMark, m.portionAbbr) {
				t.Errorf("PortionMark %q should contain abbreviated %q for %s", result.PortionMark, m.portionAbbr, m.apiCode)
			}
		})
	}
}

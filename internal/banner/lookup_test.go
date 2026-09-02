package banner

import (
	"testing"

	"dmontgomery/ism-api/internal/model"
)

func TestClassificationFromBanner(t *testing.T) {
	tests := []struct {
		label string
		want  model.Classification
		ok    bool
	}{
		{"UNCLASSIFIED", model.ClassificationU, true},
		{"CUI", model.ClassificationCUI, true},
		{"CONFIDENTIAL", model.ClassificationC, true},
		{"SECRET", model.ClassificationS, true},
		{"TOP SECRET", model.ClassificationTS, true},
		{"BOGUS", "", false},
		{"", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.label, func(t *testing.T) {
			got, ok := ClassificationFromBanner(tt.label)
			if ok != tt.ok {
				t.Fatalf("ok = %v, want %v", ok, tt.ok)
			}
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestClassificationFromPortion(t *testing.T) {
	tests := []struct {
		label string
		want  model.Classification
		ok    bool
	}{
		{"U", model.ClassificationU, true},
		{"CUI", model.ClassificationCUI, true},
		{"C", model.ClassificationC, true},
		{"S", model.ClassificationS, true},
		{"TS", model.ClassificationTS, true},
		{"BOGUS", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.label, func(t *testing.T) {
			got, ok := ClassificationFromPortion(tt.label)
			if ok != tt.ok {
				t.Fatalf("ok = %v, want %v", ok, tt.ok)
			}
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestControlFromPortionAbbrev(t *testing.T) {
	tests := []struct {
		abbr string
		want string
		ok   bool
	}{
		{"NF", "NOFORN", true},
		{"PR", "PROPIN", true},
		{"IMC", "IMCON", true},
		{"DS", "DSEN", true},
		{"NC", "NOCON", true},
		{"BOGUS", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.abbr, func(t *testing.T) {
			got, ok := ControlFromPortionAbbrev(tt.abbr)
			if ok != tt.ok {
				t.Fatalf("ok = %v, want %v", ok, tt.ok)
			}
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBannerClassificationLabels_LongestFirst(t *testing.T) {
	labels := BannerClassificationLabels()
	if len(labels) != len(classificationBanner) {
		t.Fatalf("got %d labels, want %d", len(labels), len(classificationBanner))
	}
	for i := 1; i < len(labels); i++ {
		if len(labels[i-1]) < len(labels[i]) {
			t.Errorf("labels not longest-first at index %d: %q (%d) before %q (%d)",
				i, labels[i-1], len(labels[i-1]), labels[i], len(labels[i]))
		}
	}
	// "TOP SECRET" must sort before "SECRET" since it's longer and both are
	// present, otherwise prefix matching against banner text would be ambiguous.
	topIdx, secIdx := -1, -1
	for i, l := range labels {
		if l == "TOP SECRET" {
			topIdx = i
		}
		if l == "SECRET" {
			secIdx = i
		}
	}
	if topIdx == -1 || secIdx == -1 {
		t.Fatalf("expected both TOP SECRET and SECRET in labels, got %v", labels)
	}
	if topIdx > secIdx {
		t.Errorf("TOP SECRET (idx %d) should sort before SECRET (idx %d)", topIdx, secIdx)
	}
}

// TestInversionCompleteness ensures every entry in the forward maps
// (classificationBanner, classificationPortion, portionAbbrev) is reachable
// through the reverse lookups, so the two directions cannot silently drift.
func TestInversionCompleteness(t *testing.T) {
	for cls, label := range classificationBanner {
		got, ok := ClassificationFromBanner(label)
		if !ok {
			t.Errorf("ClassificationFromBanner(%q) not found, forward map has %q -> %q", label, cls, label)
			continue
		}
		if got != cls {
			t.Errorf("ClassificationFromBanner(%q) = %q, want %q", label, got, cls)
		}
	}

	for cls, label := range classificationPortion {
		got, ok := ClassificationFromPortion(label)
		if !ok {
			t.Errorf("ClassificationFromPortion(%q) not found, forward map has %q -> %q", label, cls, label)
			continue
		}
		if got != cls {
			t.Errorf("ClassificationFromPortion(%q) = %q, want %q", label, got, cls)
		}
	}

	for ctrl, abbr := range portionAbbrev {
		got, ok := ControlFromPortionAbbrev(abbr)
		if !ok {
			t.Errorf("ControlFromPortionAbbrev(%q) not found, forward map has %q -> %q", abbr, ctrl, abbr)
			continue
		}
		if got != ctrl {
			t.Errorf("ControlFromPortionAbbrev(%q) = %q, want %q", abbr, got, ctrl)
		}
	}
}

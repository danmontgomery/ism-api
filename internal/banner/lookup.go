package banner

import (
	"sort"

	"dmontgomery/ism-api/internal/model"
)

// Reverse lookup tables, built at init time by inverting the forward tables
// (classificationBanner, classificationPortion, portionAbbrev) defined in
// render.go. Keeping these as inversions of the same source maps means the
// two directions cannot drift apart.
var (
	bannerToClassification   map[string]model.Classification
	portionToClassification  map[string]model.Classification
	abbrevToControl          map[string]string
	bannerLabelsLongestFirst []string
)

func init() {
	bannerToClassification = make(map[string]model.Classification, len(classificationBanner))
	bannerLabelsLongestFirst = make([]string, 0, len(classificationBanner))
	for cls, label := range classificationBanner {
		bannerToClassification[label] = cls
		bannerLabelsLongestFirst = append(bannerLabelsLongestFirst, label)
	}
	sort.Slice(bannerLabelsLongestFirst, func(i, j int) bool {
		return len(bannerLabelsLongestFirst[i]) > len(bannerLabelsLongestFirst[j])
	})

	portionToClassification = make(map[string]model.Classification, len(classificationPortion))
	for cls, label := range classificationPortion {
		portionToClassification[label] = cls
	}

	abbrevToControl = make(map[string]string, len(portionAbbrev))
	for ctrl, abbr := range portionAbbrev {
		abbrevToControl[abbr] = ctrl
	}
}

// ClassificationFromBanner returns the Classification for a full banner label
// (e.g. "TOP SECRET" -> TS). Inverted from classificationBanner.
func ClassificationFromBanner(label string) (model.Classification, bool) {
	c, ok := bannerToClassification[label]
	return c, ok
}

// ClassificationFromPortion returns the Classification for an abbreviated
// portion label (e.g. "TS" -> TS). Inverted from classificationPortion.
func ClassificationFromPortion(label string) (model.Classification, bool) {
	c, ok := portionToClassification[label]
	return c, ok
}

// ControlFromPortionAbbrev returns the dissemination control code for a
// portion-mark abbreviation (e.g. "NF" -> "NOFORN"). Inverted from portionAbbrev.
func ControlFromPortionAbbrev(abbr string) (string, bool) {
	c, ok := abbrevToControl[abbr]
	return c, ok
}

// BannerClassificationLabels returns all banner classification labels sorted
// longest-first. Labels contain spaces (e.g. "TOP SECRET"), so the
// classification segment of a banner line cannot be tokenized by splitting on
// whitespace; longest-first prefix matching resolves the ambiguity.
func BannerClassificationLabels() []string {
	out := make([]string, len(bannerLabelsLongestFirst))
	copy(out, bannerLabelsLongestFirst)
	return out
}

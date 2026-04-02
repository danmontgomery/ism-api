package model

import "fmt"

// Classification represents a security classification level.
type Classification string

const (
	ClassificationU   Classification = "U"
	ClassificationCUI Classification = "CUI"
	ClassificationC   Classification = "C"
	ClassificationS   Classification = "S"
)

// classificationOrder defines the ordering from least to most restrictive.
var classificationOrder = map[Classification]int{
	ClassificationU:   0,
	ClassificationCUI: 1,
	ClassificationC:   2,
	ClassificationS:   3,
}

// Valid returns true if the classification is a known value.
func (c Classification) Valid() bool {
	_, ok := classificationOrder[c]
	return ok
}

// Level returns the numeric ordering (0=U, 1=CUI, 2=C, 3=S).
// Returns -1 for unknown classifications.
func (c Classification) Level() int {
	if lvl, ok := classificationOrder[c]; ok {
		return lvl
	}
	return -1
}

// AtLeast returns true if c is at or above the given classification level.
func (c Classification) AtLeast(other Classification) bool {
	return c.Level() >= other.Level()
}

// String returns the string representation.
func (c Classification) String() string {
	return string(c)
}

// ParseClassification converts a string to a Classification, returning an error
// if the value is not recognized.
func ParseClassification(s string) (Classification, error) {
	c := Classification(s)
	if !c.Valid() {
		return "", fmt.Errorf("unknown classification: %q", s)
	}
	return c, nil
}

// AllClassifications returns all supported classification levels in order.
func AllClassifications() []Classification {
	return []Classification{
		ClassificationU,
		ClassificationCUI,
		ClassificationC,
		ClassificationS,
	}
}

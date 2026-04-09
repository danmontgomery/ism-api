package refdata

import "expr.ai/ism-api/internal/model"

// ExclusivePair represents two dissemination controls that cannot coexist.
type ExclusivePair struct {
	A string
	B string
}

// FieldRequirement maps a dissemination control to a required ISM field.
type FieldRequirement struct {
	Control string
	Field   string
}

// ClassificationGate maps a dissemination control to its minimum classification.
type ClassificationGate struct {
	Control           string
	MinClassification model.Classification
}

// DistributionConstraint maps a distribution statement to its allowed classifications.
type DistributionConstraint struct {
	Statement        string
	AllowedLevels    []model.Classification
}

// ExclusiveDisseminationPairs returns all pairs of mutually exclusive dissemination controls.
func ExclusiveDisseminationPairs() []ExclusivePair {
	return []ExclusivePair{
		{A: "NOFORN", B: "REL"},
		{A: "NOFORN", B: "RELIDO"},
	}
}

// DisseminationFieldRequirements returns dissemination controls that require
// a specific ISM field to be populated.
func DisseminationFieldRequirements() []FieldRequirement {
	return []FieldRequirement{
		{Control: "REL", Field: "releasableTo"},
		{Control: "EYES", Field: "releasableTo"},
		{Control: "DISPLAY ONLY", Field: "displayOnlyTo"},
		{Control: "SCI", Field: "sciControls"},
	}
}

// DisseminationClassificationGates returns dissemination controls that require
// a minimum classification level.
func DisseminationClassificationGates() []ClassificationGate {
	return []ClassificationGate{
		{Control: "OC", MinClassification: model.ClassificationC},
		{Control: "OC-USGOV", MinClassification: model.ClassificationC},
		{Control: "IMCON", MinClassification: model.ClassificationS},
		{Control: "NOFORN", MinClassification: model.ClassificationCUI},
		{Control: "EYES", MinClassification: model.ClassificationC},
		{Control: "DSEN", MinClassification: model.ClassificationC},
		{Control: "FISA", MinClassification: model.ClassificationC},
		{Control: "RAWFISA", MinClassification: model.ClassificationC},
		{Control: "NOCON", MinClassification: model.ClassificationC},
		{Control: "RELIDO", MinClassification: model.ClassificationC},
		{Control: "SCI", MinClassification: model.ClassificationTS},
	}
}

// DistributionClassificationConstraints returns distribution statements with
// restrictions on which classification levels they may be applied to.
func DistributionClassificationConstraints() []DistributionConstraint {
	return []DistributionConstraint{
		{
			Statement:     "A",
			AllowedLevels: []model.Classification{model.ClassificationU, model.ClassificationCUI},
		},
	}
}

// DeclassFieldExclusions returns pairs of declassification fields that are
// mutually exclusive (only one may be set).
func DeclassFieldExclusions() [][2]string {
	return [][2]string{
		{"declassDate", "declassEvent"},
	}
}

// DeclassApplicableClassifications returns the classification levels for which
// declassification fields are applicable.
func DeclassApplicableClassifications() []model.Classification {
	return []model.Classification{
		model.ClassificationC,
		model.ClassificationS,
		model.ClassificationTS,
	}
}

// CUIRequiresOwnerProducer returns false — ownerProducer is not required for CUI.
// This function exists so the compatibility matrix is the single source of truth
// for field-level requirements.
func CUIRequiresOwnerProducer() bool {
	return false
}

// ClassifiedRequiresOwnerProducer returns true — ownerProducer is required for
// Confidential and Secret classifications.
func ClassifiedRequiresOwnerProducer() bool {
	return true
}

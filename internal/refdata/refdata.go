package refdata

import (
	"cmp"
	"slices"

	"expr.ai/ism-api/internal/model"
)

// ClassificationEntry is a reference data entry for a classification level.
type ClassificationEntry struct {
	Code  model.Classification `json:"code"`
	Label string               `json:"label"`
	Level int                  `json:"level"`
}

// CUICategory is a reference data entry for a CUI category marking.
type CUICategory struct {
	Code        string `json:"code"`
	Label       string `json:"label"`
	Type        string `json:"type"` // "specified" or "basic"
	Description string `json:"description"`
}

// DisseminationControl is a reference data entry for a dissemination control.
type DisseminationControl struct {
	Code              string               `json:"code"`
	Label             string               `json:"label"`
	Description       string               `json:"description"`
	RequiresField     string               `json:"requiresField,omitempty"`
	MinClassification model.Classification  `json:"minClassification,omitempty"`
	ExclusiveWith     []string             `json:"exclusiveWith,omitempty"`
}

// DistributionStatement is a reference data entry for a distribution statement.
type DistributionStatement struct {
	Code                     string               `json:"code"`
	Label                    string               `json:"label"`
	Text                     string               `json:"text"`
	ClassificationConstraint string               `json:"classificationConstraint,omitempty"`
}

// CountryCode is a reference data entry for a country, coalition, or organization.
type CountryCode struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Type string `json:"type"` // "country", "coalition", "organization", or "marker"
}

// DeclassException is a reference data entry for a 25X declassification exemption.
type DeclassException struct {
	Code     string `json:"code"`
	Label    string `json:"label"`
	Category string `json:"category"`
}

// NonICMarking is a reference data entry for a non-IC marking.
type NonICMarking struct {
	Code  string `json:"code"`
	Label string `json:"label"`
}

// NonUSControl is a reference data entry for a non-US control marking.
type NonUSControl struct {
	Code  string `json:"code"`
	Label string `json:"label"`
}

// ExemptFromEntry is a reference data entry for an ISM exemption.
type ExemptFromEntry struct {
	Code        string `json:"code"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

// CompliesWithEntry is a reference data entry for an ISM compliance framework.
type CompliesWithEntry struct {
	Code  string `json:"code"`
	Label string `json:"label"`
}

// NoticeType is a reference data entry for an ISM notice type.
type NoticeType struct {
	Code     string `json:"code"`
	Label    string `json:"label"`
	Category string `json:"category"`
}

// Registry aggregates all compiled-in reference data for the ISM API.
type Registry struct {
	Classifications        []ClassificationEntry
	CUICategories          []CUICategory
	DisseminationControls  []DisseminationControl
	DistributionStatements []DistributionStatement
	CountryCodes           []CountryCode
	DeclassExceptions      []DeclassException
	NonICMarkings          []NonICMarking
	NonUSControls          []NonUSControl
	ExemptFrom             []ExemptFromEntry
	CompliesWith           []CompliesWithEntry
	AtomicEnergyMarkings   []AtomicEnergyMarking
	NoticeTypes            []NoticeType
	SCIControls            []SCIControl
}

// NewRegistry returns a Registry populated with all compiled-in reference data.
// All slices are sorted alphabetically by code for consistent API responses.
func NewRegistry() *Registry {
	r := &Registry{
		Classifications:        Classifications(),
		CUICategories:          CUICategories(),
		DisseminationControls:  DisseminationControls(),
		DistributionStatements: DistributionStatements(),
		CountryCodes:           CountryCodes(),
		DeclassExceptions:      DeclassExceptions(),
		NonICMarkings:          NonICMarkings(),
		NonUSControls:          NonUSControls(),
		ExemptFrom:             ExemptFrom(),
		CompliesWith:           CompliesWith(),
		AtomicEnergyMarkings:   AtomicEnergyMarkings(),
		NoticeTypes:            NoticeTypes(),
		SCIControls:            SCIControls(),
	}

	slices.SortFunc(r.Classifications, func(a, b ClassificationEntry) int {
		return cmp.Compare(string(a.Code), string(b.Code))
	})
	slices.SortFunc(r.CUICategories, func(a, b CUICategory) int {
		return cmp.Compare(a.Code, b.Code)
	})
	slices.SortFunc(r.DisseminationControls, func(a, b DisseminationControl) int {
		return cmp.Compare(a.Code, b.Code)
	})
	slices.SortFunc(r.DistributionStatements, func(a, b DistributionStatement) int {
		return cmp.Compare(a.Code, b.Code)
	})
	slices.SortFunc(r.CountryCodes, func(a, b CountryCode) int {
		return cmp.Compare(a.Code, b.Code)
	})
	slices.SortFunc(r.DeclassExceptions, func(a, b DeclassException) int {
		return cmp.Compare(a.Code, b.Code)
	})
	slices.SortFunc(r.NonICMarkings, func(a, b NonICMarking) int {
		return cmp.Compare(a.Code, b.Code)
	})
	slices.SortFunc(r.NonUSControls, func(a, b NonUSControl) int {
		return cmp.Compare(a.Code, b.Code)
	})
	slices.SortFunc(r.ExemptFrom, func(a, b ExemptFromEntry) int {
		return cmp.Compare(a.Code, b.Code)
	})
	slices.SortFunc(r.CompliesWith, func(a, b CompliesWithEntry) int {
		return cmp.Compare(a.Code, b.Code)
	})
	slices.SortFunc(r.AtomicEnergyMarkings, func(a, b AtomicEnergyMarking) int {
		return cmp.Compare(a.Code, b.Code)
	})
	slices.SortFunc(r.NoticeTypes, func(a, b NoticeType) int {
		return cmp.Compare(a.Code, b.Code)
	})
	slices.SortFunc(r.SCIControls, func(a, b SCIControl) int {
		return cmp.Compare(a.Code, b.Code)
	})

	return r
}

// ValidClassification returns true if code is a known classification.
func (r *Registry) ValidClassification(code model.Classification) bool {
	for _, c := range r.Classifications {
		if c.Code == code {
			return true
		}
	}
	return false
}

// ValidCUICategory returns true if code is a known CUI category.
func (r *Registry) ValidCUICategory(code string) bool {
	for _, c := range r.CUICategories {
		if c.Code == code {
			return true
		}
	}
	return false
}

// ValidDisseminationControl returns true if code is a known dissemination control.
func (r *Registry) ValidDisseminationControl(code string) bool {
	for _, d := range r.DisseminationControls {
		if d.Code == code {
			return true
		}
	}
	return false
}

// ValidDistributionStatement returns true if code is a known distribution statement.
func (r *Registry) ValidDistributionStatement(code string) bool {
	for _, d := range r.DistributionStatements {
		if d.Code == code {
			return true
		}
	}
	return false
}

// ValidCountryCode returns true if code is a known country/coalition/org code.
func (r *Registry) ValidCountryCode(code string) bool {
	for _, c := range r.CountryCodes {
		if c.Code == code {
			return true
		}
	}
	return false
}

// ValidDeclassException returns true if code is a known 25X declass exemption.
func (r *Registry) ValidDeclassException(code string) bool {
	for _, d := range r.DeclassExceptions {
		if d.Code == code {
			return true
		}
	}
	return false
}

// ValidNonICMarking returns true if code is a known non-IC marking.
func (r *Registry) ValidNonICMarking(code string) bool {
	for _, n := range r.NonICMarkings {
		if n.Code == code {
			return true
		}
	}
	return false
}

// ValidNonUSControl returns true if code is a known non-US control marking.
func (r *Registry) ValidNonUSControl(code string) bool {
	for _, n := range r.NonUSControls {
		if n.Code == code {
			return true
		}
	}
	return false
}

// ValidExemptFrom returns true if code is a known exemptFrom value.
func (r *Registry) ValidExemptFrom(code string) bool {
	for _, e := range r.ExemptFrom {
		if e.Code == code {
			return true
		}
	}
	return false
}

// ValidCompliesWith returns true if code is a known compliesWith value.
func (r *Registry) ValidCompliesWith(code string) bool {
	for _, c := range r.CompliesWith {
		if c.Code == code {
			return true
		}
	}
	return false
}

// ValidAtomicEnergyMarking returns true if code is a known atomic energy marking.
func (r *Registry) ValidAtomicEnergyMarking(code string) bool {
	for _, a := range r.AtomicEnergyMarkings {
		if a.Code == code {
			return true
		}
	}
	return false
}

// ValidNoticeType returns true if code is a known notice type.
func (r *Registry) ValidNoticeType(code string) bool {
	for _, n := range r.NoticeTypes {
		if n.Code == code {
			return true
		}
	}
	return false
}

// ValidSCIControl returns true if code is a known SCI control.
func (r *Registry) ValidSCIControl(code string) bool {
	for _, s := range r.SCIControls {
		if s.Code == code {
			return true
		}
	}
	return false
}

// LookupSCIControl returns the SCIControl for the given code, or nil.
func (r *Registry) LookupSCIControl(code string) *SCIControl {
	for i := range r.SCIControls {
		if r.SCIControls[i].Code == code {
			return &r.SCIControls[i]
		}
	}
	return nil
}

// LookupDisseminationControl returns the DisseminationControl for the given code, or nil.
func (r *Registry) LookupDisseminationControl(code string) *DisseminationControl {
	for i := range r.DisseminationControls {
		if r.DisseminationControls[i].Code == code {
			return &r.DisseminationControls[i]
		}
	}
	return nil
}

// LookupDistributionStatement returns the DistributionStatement for the given code, or nil.
func (r *Registry) LookupDistributionStatement(code string) *DistributionStatement {
	for i := range r.DistributionStatements {
		if r.DistributionStatements[i].Code == code {
			return &r.DistributionStatements[i]
		}
	}
	return nil
}

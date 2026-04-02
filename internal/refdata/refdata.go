package refdata

import "github.com/danielmontgomery/ism-api/internal/model"

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
	Type string `json:"type"` // "country", "coalition", or "organization"
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

// Registry aggregates all compiled-in reference data for the ISM API.
type Registry struct {
	Classifications        []ClassificationEntry
	CUICategories          []CUICategory
	DisseminationControls  []DisseminationControl
	DistributionStatements []DistributionStatement
	CountryCodes           []CountryCode
	DeclassExceptions      []DeclassException
	NonICMarkings          []NonICMarking
}

// NewRegistry returns a Registry populated with all compiled-in reference data.
func NewRegistry() *Registry {
	return &Registry{
		Classifications:        Classifications(),
		CUICategories:          CUICategories(),
		DisseminationControls:  DisseminationControls(),
		DistributionStatements: DistributionStatements(),
		CountryCodes:           CountryCodes(),
		DeclassExceptions:      DeclassExceptions(),
		NonICMarkings:          NonICMarkings(),
	}
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

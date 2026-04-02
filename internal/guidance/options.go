package guidance

// FieldStatus represents the current applicability of an ISM field given partial state.
type FieldStatus string

const (
	StatusAvailable     FieldStatus = "available"
	StatusRequired      FieldStatus = "required"
	StatusNotApplicable FieldStatus = "not_applicable"
	StatusLocked        FieldStatus = "locked"
)

// AllowedValue represents a single valid option for a guidance field.
type AllowedValue struct {
	Code  string `json:"code"`
	Label string `json:"label"`
	Type  string `json:"type,omitempty"`
}

// FieldGuidance describes the current state and valid options for a single ISM field.
type FieldGuidance struct {
	Field         string        `json:"field"`
	Status        FieldStatus   `json:"status"`
	Required      bool          `json:"required,omitempty"`
	RequiredIf    string        `json:"requiredIf,omitempty"`
	AllowedValues []AllowedValue `json:"allowedValues,omitempty"`
	Reason        string        `json:"reason,omitempty"`
}

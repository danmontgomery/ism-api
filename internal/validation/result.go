package validation

// Severity indicates how serious a validation finding is.
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
)

// FieldError represents a single validation finding tied to a specific field.
type FieldError struct {
	Field    string   `json:"field"`
	Code     string   `json:"code"`
	Message  string   `json:"message"`
	Severity Severity `json:"severity"`
}

// ValidationResult aggregates all findings from running the validation engine.
type ValidationResult struct {
	Valid  bool         `json:"valid"`
	Errors []FieldError `json:"errors,omitempty"`
}

// AddError appends an error-severity finding and marks the result invalid.
func (r *ValidationResult) AddError(field, code, message string) {
	r.Errors = append(r.Errors, FieldError{
		Field:    field,
		Code:     code,
		Message:  message,
		Severity: SeverityError,
	})
	r.Valid = false
}

// AddWarning appends a warning-severity finding without affecting validity.
func (r *ValidationResult) AddWarning(field, code, message string) {
	r.Errors = append(r.Errors, FieldError{
		Field:    field,
		Code:     code,
		Message:  message,
		Severity: SeverityWarning,
	})
}

// Merge incorporates all findings from another result.
func (r *ValidationResult) Merge(other *ValidationResult) {
	for _, e := range other.Errors {
		r.Errors = append(r.Errors, e)
		if e.Severity == SeverityError {
			r.Valid = false
		}
	}
}

// HasCode returns true if any finding matches the given code.
func (r *ValidationResult) HasCode(code string) bool {
	for _, e := range r.Errors {
		if e.Code == code {
			return true
		}
	}
	return false
}

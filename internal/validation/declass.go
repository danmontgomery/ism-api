package validation

import (
	"github.com/danielmontgomery/ism-api/internal/model"
	"github.com/danielmontgomery/ism-api/internal/refdata"
)

// DeclassRule validates declassification fields: date/event mutual exclusion,
// 25X exception code validity, and applicability only to C/S classifications.
type DeclassRule struct{}

func (r *DeclassRule) Name() string { return "declass" }

func (r *DeclassRule) Applies(ism *model.ISM) bool {
	return ism.DeclassDate != "" || ism.DeclassEvent != "" || ism.DeclassException != ""
}

func (r *DeclassRule) Validate(ism *model.ISM, reg *refdata.Registry) *ValidationResult {
	res := &ValidationResult{Valid: true}

	// Declass fields only apply to C and S.
	applicable := false
	for _, c := range refdata.DeclassApplicableClassifications() {
		if ism.Classification == c {
			applicable = true
			break
		}
	}
	if !applicable {
		res.AddError("declassDate", "declass.not_applicable",
			"declassification fields are only applicable to C and S classifications")
		return res
	}

	// DeclassDate and DeclassEvent are mutually exclusive.
	for _, pair := range refdata.DeclassFieldExclusions() {
		hasA := fieldSet(ism, pair[0])
		hasB := fieldSet(ism, pair[1])
		if hasA && hasB {
			res.AddError(pair[0], "declass.date_event_exclusive",
				pair[0]+" and "+pair[1]+" are mutually exclusive")
		}
	}

	// 25X exception codes must be valid.
	if ism.DeclassException != "" {
		if !reg.ValidDeclassException(ism.DeclassException) {
			res.AddError("declassException", "declass.invalid_exception",
				"unknown declassification exception code: "+ism.DeclassException)
		}
	}

	return res
}

// fieldSet returns true if the named declass field on the ISM is non-empty.
func fieldSet(ism *model.ISM, name string) bool {
	switch name {
	case "declassDate":
		return ism.DeclassDate != ""
	case "declassEvent":
		return ism.DeclassEvent != ""
	default:
		return false
	}
}

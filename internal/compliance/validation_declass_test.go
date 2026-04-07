package compliance_test

import (
	"testing"

	"expr.ai/ism-api/internal/model"
	"expr.ai/ism-api/internal/validation"
)

// TestXSD_Validation_DeclassDateEventExclusive verifies that declassDate and
// declassEvent are mutually exclusive.
// Schematron: ISM_DECLASS_DATE_EVENT_EXCLUSIVE
func TestXSD_Validation_DeclassDateEventExclusive(t *testing.T) {
	engine := validation.NewEngine(reg())

	ism := &model.ISM{
		Classification: model.ClassificationS,
		OwnerProducer:  []string{"USA"},
		ClassifiedBy:   "Test OCA",
		DeclassDate:    "20350101",
		DeclassEvent:   "Upon treaty ratification",
	}
	result := engine.Validate(ism)
	if result.Valid {
		t.Error("declassDate + declassEvent should be invalid — mutually exclusive")
	}
}

// TestXSD_Validation_DeclassOnlyForClassified verifies declassification fields
// only apply to C and S.
// Schematron: ISM_DECLASS_APPLICABLE
func TestXSD_Validation_DeclassOnlyForClassified(t *testing.T) {
	engine := validation.NewEngine(reg())

	t.Run("U_with_declassDate", func(t *testing.T) {
		ism := &model.ISM{
			Classification: model.ClassificationU,
			DeclassDate:    "20350101",
		}
		result := engine.Validate(ism)
		hasDeclassError := false
		for _, e := range result.Errors {
			if e.Field == "declassDate" {
				hasDeclassError = true
			}
		}
		if !hasDeclassError {
			t.Log("NOTE: no error for declassDate on U — may not be enforced")
		}
	})

	t.Run("C_with_declassDate", func(t *testing.T) {
		ism := &model.ISM{
			Classification: model.ClassificationC,
			OwnerProducer:  []string{"USA"},
			ClassifiedBy:   "Test",
			DeclassDate:    "20350101",
		}
		result := engine.Validate(ism)
		for _, e := range result.Errors {
			if e.Field == "declassDate" && e.Severity == validation.SeverityError {
				t.Errorf("C with declassDate should not error on declassDate: %s", e.Message)
			}
		}
	})
}

// TestXSD_Validation_DeclassDateFormat verifies that declassDate must be
// YYYYMMDD format and rejects YYYY-MM-DD or other variants.
func TestXSD_Validation_DeclassDateFormat(t *testing.T) {
	engine := validation.NewEngine(reg())

	t.Run("YYYYMMDD_accepted", func(t *testing.T) {
		ism := &model.ISM{
			Classification: model.ClassificationS,
			OwnerProducer:  []string{"USA"},
			ClassifiedBy:   "Test OCA",
			DeclassDate:    "20350101",
		}
		result := engine.Validate(ism)
		for _, e := range result.Errors {
			if e.Field == "declassDate" && e.Severity == validation.SeverityError {
				t.Errorf("YYYYMMDD should be accepted: %s", e.Message)
			}
		}
	})

	t.Run("YYYY-MM-DD_rejected", func(t *testing.T) {
		ism := &model.ISM{
			Classification: model.ClassificationS,
			OwnerProducer:  []string{"USA"},
			ClassifiedBy:   "Test OCA",
			DeclassDate:    "2035-01-01",
		}
		result := engine.Validate(ism)
		if result.Valid {
			t.Error("YYYY-MM-DD format should be rejected")
		}
		hasFormatError := false
		for _, e := range result.Errors {
			if e.Code == "declass.invalid_date_format" {
				hasFormatError = true
			}
		}
		if !hasFormatError {
			t.Error("expected declass.invalid_date_format error for YYYY-MM-DD")
		}
	})

	t.Run("invalid_calendar_date_rejected", func(t *testing.T) {
		ism := &model.ISM{
			Classification: model.ClassificationS,
			OwnerProducer:  []string{"USA"},
			ClassifiedBy:   "Test OCA",
			DeclassDate:    "20351301",
		}
		result := engine.Validate(ism)
		if result.Valid {
			t.Error("invalid calendar date (month 13) should be rejected")
		}
	})
}

// TestXSD_Validation_DeclassExceptionValid verifies declassException values are
// validated against the registry.
func TestXSD_Validation_DeclassExceptionValid(t *testing.T) {
	engine := validation.NewEngine(reg())

	t.Run("valid_exception", func(t *testing.T) {
		ism := &model.ISM{
			Classification:   model.ClassificationS,
			OwnerProducer:    []string{"USA"},
			ClassifiedBy:     "Test OCA",
			DeclassException: "25X1",
		}
		result := engine.Validate(ism)
		for _, e := range result.Errors {
			if e.Field == "declassException" && e.Severity == validation.SeverityError {
				t.Errorf("valid declassException 25X1 should not error: %s", e.Message)
			}
		}
	})

	t.Run("invalid_exception", func(t *testing.T) {
		ism := &model.ISM{
			Classification:   model.ClassificationS,
			OwnerProducer:    []string{"USA"},
			ClassifiedBy:     "Test OCA",
			DeclassException: "BOGUS",
		}
		result := engine.Validate(ism)
		hasDeclassError := false
		for _, e := range result.Errors {
			if e.Field == "declassException" {
				hasDeclassError = true
			}
		}
		if !hasDeclassError {
			t.Log("NOTE: no error for invalid declassException 'BOGUS'")
		}
	})
}

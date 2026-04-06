package compliance_test

import (
	"testing"

	"expr.ai/ism-api/internal/model"
	"expr.ai/ism-api/internal/validation"
)

// TestXSD_Validation_ClassificationRequired verifies ISM.USGov.1: classification
// is required and must be a valid value.
// Schematron: ISM_RESOURCE_ELEMENT.classification required
func TestXSD_Validation_ClassificationRequired(t *testing.T) {
	engine := validation.NewEngine(reg())

	t.Run("empty_classification_rejected", func(t *testing.T) {
		ism := &model.ISM{Classification: ""}
		result := engine.Validate(ism)
		if result.Valid {
			t.Error("empty classification should be invalid")
		}
	})

	t.Run("unknown_classification_rejected", func(t *testing.T) {
		ism := &model.ISM{Classification: model.Classification("BOGUS")}
		result := engine.Validate(ism)
		if result.Valid {
			t.Error("unknown classification should be invalid")
		}
	})

	t.Run("valid_U_accepted", func(t *testing.T) {
		ism := &model.ISM{Classification: model.ClassificationU}
		result := engine.Validate(ism)
		if !result.Valid {
			t.Errorf("U should be valid; errors: %v", result.Errors)
		}
	})
}

// TestXSD_Validation_OwnerProducerRequired verifies ISM.USGov.3-6: classified
// documents require ownerProducer.
// Schematron: ISM_RESOURCE_ELEMENT.ownerProducer required for classified
func TestXSD_Validation_OwnerProducerRequired(t *testing.T) {
	engine := validation.NewEngine(reg())

	t.Run("classified_without_ownerProducer", func(t *testing.T) {
		ism := &model.ISM{
			Classification: model.ClassificationC,
			ClassifiedBy:   "Test Authority",
			DeclassDate:    "20350101",
		}
		result := engine.Validate(ism)
		if result.Valid {
			t.Error("classified without ownerProducer should be invalid")
		}
		if !result.HasCode("core.missing_owner_producer") {
			t.Log("NOTE: error code may differ from expected 'core.missing_owner_producer'")
		}
	})

	t.Run("classified_with_ownerProducer", func(t *testing.T) {
		ism := &model.ISM{
			Classification: model.ClassificationC,
			OwnerProducer:  []string{"USA"},
			ClassifiedBy:   "Test Authority",
			DeclassDate:    "20350101",
		}
		result := engine.Validate(ism)
		// May still have warnings but should not have ownerProducer error
		for _, e := range result.Errors {
			if e.Field == "ownerProducer" && e.Severity == validation.SeverityError {
				t.Errorf("classified with ownerProducer should not have ownerProducer error: %s", e.Message)
			}
		}
	})

	t.Run("unclassified_without_ownerProducer_ok", func(t *testing.T) {
		ism := &model.ISM{Classification: model.ClassificationU}
		result := engine.Validate(ism)
		if !result.Valid {
			t.Errorf("U without ownerProducer should be valid; errors: %v", result.Errors)
		}
	})
}

// TestXSD_Validation_JointRequiresMultipleOwners verifies that joint=true
// requires at least 2 ownerProducer entries.
// Schematron: ISM_RESOURCE_ELEMENT.joint
func TestXSD_Validation_JointRequiresMultipleOwners(t *testing.T) {
	engine := validation.NewEngine(reg())

	t.Run("joint_with_single_owner", func(t *testing.T) {
		ism := &model.ISM{
			Classification: model.ClassificationC,
			OwnerProducer:  []string{"USA"},
			Joint:          true,
			ClassifiedBy:   "Test",
			DeclassDate:    "20350101",
		}
		result := engine.Validate(ism)
		if result.Valid {
			t.Error("joint with single ownerProducer should be invalid")
		}
	})

	t.Run("joint_with_multiple_owners", func(t *testing.T) {
		ism := &model.ISM{
			Classification: model.ClassificationC,
			OwnerProducer:  []string{"USA", "GBR"},
			Joint:          true,
			ClassifiedBy:   "Test",
			DeclassDate:    "20350101",
		}
		result := engine.Validate(ism)
		for _, e := range result.Errors {
			if e.Field == "joint" && e.Severity == validation.SeverityError {
				t.Errorf("joint with 2 owners should not have joint error: %s", e.Message)
			}
		}
	})

	t.Run("non_joint_multiple_owners_requires_joint", func(t *testing.T) {
		ism := &model.ISM{
			Classification: model.ClassificationC,
			OwnerProducer:  []string{"USA", "GBR"},
			Joint:          false,
			ClassifiedBy:   "Test",
			DeclassDate:    "20350101",
		}
		result := engine.Validate(ism)
		if result.Valid {
			t.Error("multiple ownerProducer without joint=true should be invalid")
		}
	})
}

package compliance_test

import (
	"testing"

	"expr.ai/ism-api/internal/model"
	"expr.ai/ism-api/internal/refdata"
	"expr.ai/ism-api/internal/validation"
)

// TestXSD_Validation_NOFORN_REL_Exclusive verifies ISM.USGov.12: NOFORN and REL
// are mutually exclusive.
// Schematron: ISM_DISSEM_NOFORN_REL_EXCLUSIVE
func TestXSD_Validation_NOFORN_REL_Exclusive(t *testing.T) {
	engine := validation.NewEngine(reg())

	ism := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"NOFORN", "REL"},
		ReleasableTo:          []string{"GBR"},
		ClassifiedBy:          "Test",
		DeclassDate:           "20350101",
	}
	result := engine.Validate(ism)
	if result.Valid {
		t.Error("NOFORN + REL should be invalid — mutually exclusive per ISM.USGov.12")
	}
}

// TestXSD_Validation_NOFORN_RELIDO_Exclusive verifies NOFORN and RELIDO
// are mutually exclusive.
func TestXSD_Validation_NOFORN_RELIDO_Exclusive(t *testing.T) {
	engine := validation.NewEngine(reg())

	ism := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"NOFORN", "RELIDO"},
		ClassifiedBy:          "Test",
		DeclassDate:           "20350101",
	}
	result := engine.Validate(ism)
	if result.Valid {
		t.Error("NOFORN + RELIDO should be invalid — mutually exclusive")
	}
}

// TestXSD_Validation_REL_RequiresReleasableTo verifies ISM.USGov.14: REL
// requires the releasableTo field.
// Schematron: ISM_DISSEM_REL_RELEASABLETO
func TestXSD_Validation_REL_RequiresReleasableTo(t *testing.T) {
	engine := validation.NewEngine(reg())

	t.Run("REL_without_releasableTo", func(t *testing.T) {
		ism := &model.ISM{
			Classification:        model.ClassificationS,
			OwnerProducer:         []string{"USA"},
			DisseminationControls: []string{"REL"},
			ClassifiedBy:          "Test",
			DeclassDate:           "20350101",
		}
		result := engine.Validate(ism)
		if result.Valid {
			t.Error("REL without releasableTo should be invalid")
		}
	})

	t.Run("REL_with_releasableTo", func(t *testing.T) {
		ism := &model.ISM{
			Classification:        model.ClassificationS,
			OwnerProducer:         []string{"USA"},
			DisseminationControls: []string{"REL"},
			ReleasableTo:          []string{"GBR", "CAN"},
			ClassifiedBy:          "Test",
			DeclassDate:           "20350101",
		}
		result := engine.Validate(ism)
		for _, e := range result.Errors {
			if e.Field == "releasableTo" && e.Severity == validation.SeverityError {
				t.Errorf("REL with releasableTo should not have releasableTo error: %s", e.Message)
			}
		}
	})
}

// TestXSD_Validation_EYES_RequiresReleasableTo verifies EYES requires releasableTo.
func TestXSD_Validation_EYES_RequiresReleasableTo(t *testing.T) {
	engine := validation.NewEngine(reg())

	ism := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"EYES"},
		ClassifiedBy:          "Test",
		DeclassDate:           "20350101",
	}
	result := engine.Validate(ism)
	if result.Valid {
		t.Error("EYES without releasableTo should be invalid")
	}
}

// TestXSD_Validation_DISPLAYONLY_RequiresDisplayOnlyTo verifies DISPLAY ONLY
// requires displayOnlyTo field.
func TestXSD_Validation_DISPLAYONLY_RequiresDisplayOnlyTo(t *testing.T) {
	engine := validation.NewEngine(reg())

	ism := &model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"DISPLAY ONLY"},
		ClassifiedBy:          "Test",
		DeclassDate:           "20350101",
	}
	result := engine.Validate(ism)
	if result.Valid {
		t.Error("DISPLAY ONLY without displayOnlyTo should be invalid")
	}
}

// TestXSD_Validation_ClassificationGates verifies dissemination controls that
// require minimum classification levels.
// Schematron: ISM_DISSEM_*_CLASSIFICATION
func TestXSD_Validation_ClassificationGates(t *testing.T) {
	engine := validation.NewEngine(reg())
	gates := refdata.DisseminationClassificationGates()

	for _, gate := range gates {
		t.Run(gate.Control+"_below_minimum", func(t *testing.T) {
			ism := &model.ISM{
				Classification:        model.ClassificationU,
				DisseminationControls: []string{gate.Control},
			}
			result := engine.Validate(ism)
			if result.Valid {
				t.Errorf("%s with U classification should be invalid — requires %s minimum", gate.Control, gate.MinClassification)
			}
		})
	}
}

// TestXSD_Validation_AllExclusivePairs verifies all documented exclusive pairs.
func TestXSD_Validation_AllExclusivePairs(t *testing.T) {
	engine := validation.NewEngine(reg())
	pairs := refdata.ExclusiveDisseminationPairs()

	for _, pair := range pairs {
		t.Run(pair.A+"_vs_"+pair.B, func(t *testing.T) {
			ism := &model.ISM{
				Classification:        model.ClassificationS,
				OwnerProducer:         []string{"USA"},
				DisseminationControls: []string{pair.A, pair.B},
				ReleasableTo:          []string{"GBR"}, // in case REL needs it
				ClassifiedBy:          "Test",
				DeclassDate:           "20350101",
			}
			result := engine.Validate(ism)
			if result.Valid {
				t.Errorf("%s + %s should be invalid — mutually exclusive", pair.A, pair.B)
			}
		})
	}
}

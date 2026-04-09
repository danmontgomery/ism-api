package validation

import (
	"strings"

	"expr.ai/ism-api/internal/model"
	"expr.ai/ism-api/internal/refdata"
)

// AtomicEnergyRule validates Atomic Energy Act classification level gates:
//   - RD:    CONFIDENTIAL or higher (CL-1)
//   - FRD:   CONFIDENTIAL or higher (CL-2)
//   - CNWDI: SECRET or higher       (CL-3)
//   - SIGMA: CONFIDENTIAL or higher (CL-4)
type AtomicEnergyRule struct{}

func (r *AtomicEnergyRule) Name() string { return "atomic_energy" }

func (r *AtomicEnergyRule) Applies(ism *model.ISM) bool {
	return len(ism.AtomicEnergyMarkings) > 0
}

func (r *AtomicEnergyRule) Validate(ism *model.ISM, _ *refdata.Registry) *ValidationResult {
	res := &ValidationResult{Valid: true}

	for _, code := range ism.AtomicEnergyMarkings {
		minClass := aeaMinClassification(code)
		if minClass != "" && !ism.Classification.AtLeast(minClass) {
			res.AddError("atomicEnergyMarkings", "atomic_energy.insufficient_classification",
				code+" requires at least "+minClass.String()+" classification")
		}
	}

	return res
}

// aeaMinClassification returns the minimum classification for the given AEA
// marking code based on its category prefix:
//   - CNWDI (RD-CNWDI):       SECRET or higher (CL-3)
//   - RD, RD-SG-*, RD-SIGMA:  CONFIDENTIAL or higher (CL-1, CL-4)
//   - FRD, FRD-SG-*, FRD-SIGMA: CONFIDENTIAL or higher (CL-2, CL-4)
//   - DCNI, UCNI, TFNI:       no classification gate
func aeaMinClassification(code string) model.Classification {
	// CNWDI requires SECRET or higher (CL-3). Check before general RD prefix.
	if strings.Contains(code, "CNWDI") {
		return model.ClassificationS
	}
	// RD-based markings (including SIGMA) require CONFIDENTIAL or higher (CL-1, CL-4).
	if strings.HasPrefix(code, "RD") {
		return model.ClassificationC
	}
	// FRD-based markings (including SIGMA) require CONFIDENTIAL or higher (CL-2, CL-4).
	if strings.HasPrefix(code, "FRD") {
		return model.ClassificationC
	}
	return ""
}

package refdata

import "expr.ai/ism-api/internal/model"

// Classifications returns all supported classification levels in display order.
func Classifications() []ClassificationEntry {
	return []ClassificationEntry{
		{Code: model.ClassificationU, Label: "Unclassified", Level: 0},
		{Code: model.ClassificationCUI, Label: "Controlled Unclassified Information", Level: 1},
		{Code: model.ClassificationC, Label: "Confidential", Level: 2},
		{Code: model.ClassificationS, Label: "Secret", Level: 3},
		{Code: model.ClassificationTS, Label: "Top Secret", Level: 4},
	}
}

package handler

import (
	"strings"

	"dmontgomery/ism-api/internal/parse"
	"dmontgomery/ism-api/internal/validation"
	"github.com/gin-gonic/gin"
)

type parseRequest struct {
	Marking string `json:"marking"`
}

// parseResponse embeds the parse result and adds the validation engine's
// findings against the parsed ISM.
type parseResponse struct {
	*parse.Result
	Validation *validation.ValidationResult `json:"validation"`
}

// Parse converts a banner line or portion mark into a best-effort ISM
// object, along with a round-trip check and validation findings. Unlike
// malformed JSON or a missing marking (which 400), unrecognized or lossy
// content in the marking itself never fails the request — it surfaces as
// warnings on a 200 response, matching the existing convention where a
// semantically invalid ISM still returns 200 from /validate.
func (h *Handler) Parse(c *gin.Context) {
	var req parseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, 400, "malformed request: "+err.Error())
		return
	}
	if strings.TrimSpace(req.Marking) == "" {
		respondError(c, 400, "malformed request: marking is required")
		return
	}

	result := parse.Parse(req.Marking, h.reg)
	respondData(c, parseResponse{
		Result:     result,
		Validation: h.validator.Validate(&result.ISM),
	})
}

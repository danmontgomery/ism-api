package handler

import (
	"dmontgomery/ism-api/internal/model"
	"github.com/gin-gonic/gin"
)

type validateRequest struct {
	ISM model.ISM `json:"ism"`
}

// Validate runs all validation rules against a complete ISM object.
// Returns HTTP 200 for both valid and invalid ISM; 400 for malformed requests.
func (h *Handler) Validate(c *gin.Context) {
	var req validateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, 400, "malformed request: "+err.Error())
		return
	}

	result := h.validator.Validate(&req.ISM)
	respondData(c, result)
}

// ValidatePortion runs validation against a portion-level ISM.
// Same engine — the rules self-select based on what fields are present.
func (h *Handler) ValidatePortion(c *gin.Context) {
	var req validateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, 400, "malformed request: "+err.Error())
		return
	}

	result := h.validator.Validate(&req.ISM)
	respondData(c, result)
}

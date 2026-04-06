package handler

import (
	"expr.ai/ism-api/internal/model"
	"github.com/gin-gonic/gin"
)

type guidanceRequest struct {
	ISM model.ISM `json:"ism"`
}

type guidanceResponse struct {
	Fields any `json:"fields"`
}

// Guidance returns field-level guidance for a partial ISM object.
func (h *Handler) Guidance(c *gin.Context) {
	var req guidanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, 400, "malformed request: "+err.Error())
		return
	}

	fields := h.guider.Evaluate(&req.ISM)
	respondData(c, guidanceResponse{Fields: fields})
}

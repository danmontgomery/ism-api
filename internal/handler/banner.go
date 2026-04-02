package handler

import (
	"github.com/danielmontgomery/ism-api/internal/banner"
	"github.com/danielmontgomery/ism-api/internal/model"
	"github.com/gin-gonic/gin"
)

type bannerRequest struct {
	ISM model.ISM `json:"ism"`
}

// Banner renders the banner line and portion mark for an ISM object.
func (h *Handler) Banner(c *gin.Context) {
	var req bannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, 400, "malformed request: "+err.Error())
		return
	}

	result := banner.Render(&req.ISM)
	respondData(c, result)
}

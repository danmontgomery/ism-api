package handler

import "github.com/gin-gonic/gin"

// RefClassifications returns all classification levels.
func (h *Handler) RefClassifications(c *gin.Context) {
	respondData(c, h.reg.Classifications)
}

// RefCUICategories returns all CUI category markings.
func (h *Handler) RefCUICategories(c *gin.Context) {
	respondData(c, h.reg.CUICategories)
}

// RefDisseminationControls returns all dissemination controls.
func (h *Handler) RefDisseminationControls(c *gin.Context) {
	respondData(c, h.reg.DisseminationControls)
}

// RefDistributionStatements returns all distribution statements.
func (h *Handler) RefDistributionStatements(c *gin.Context) {
	respondData(c, h.reg.DistributionStatements)
}

// RefCountryCodes returns all country, coalition, and organization codes.
func (h *Handler) RefCountryCodes(c *gin.Context) {
	respondData(c, h.reg.CountryCodes)
}

// RefDeclassExceptions returns all declassification exception codes.
func (h *Handler) RefDeclassExceptions(c *gin.Context) {
	respondData(c, h.reg.DeclassExceptions)
}

// RefNonICMarkings returns all non-IC markings.
func (h *Handler) RefNonICMarkings(c *gin.Context) {
	respondData(c, h.reg.NonICMarkings)
}

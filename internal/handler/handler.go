package handler

import (
	"io/fs"

	"dmontgomery/ism-api/internal/guidance"
	"dmontgomery/ism-api/internal/refdata"
	"dmontgomery/ism-api/internal/validation"
	"github.com/gin-gonic/gin"
)

// Handler holds the dependencies for all route handlers.
type Handler struct {
	reg       *refdata.Registry
	validator *validation.Engine
	guider    *guidance.Engine
	docsFS    fs.FS
}

// New creates a Handler with the given dependencies.
func New(reg *refdata.Registry, validator *validation.Engine, guider *guidance.Engine, docsFS fs.FS) *Handler {
	return &Handler{
		reg:       reg,
		validator: validator,
		guider:    guider,
		docsFS:    docsFS,
	}
}

// Register mounts all routes on the given Gin engine.
func (h *Handler) Register(r *gin.Engine) {
	r.Use(CORS(), RequestID(), Logger(), Recovery())

	r.GET("/healthz", h.Health)
	r.GET("/openapi.yaml", h.OpenAPISpec)
	r.GET("/docs", h.ScalarDocs)

	v1 := r.Group("/api/v1")
	{
		ref := v1.Group("/ref")
		{
			ref.GET("/classifications", h.RefClassifications)
			ref.GET("/cui-categories", h.RefCUICategories)
			ref.GET("/dissemination-controls", h.RefDisseminationControls)
			ref.GET("/distribution-statements", h.RefDistributionStatements)
			ref.GET("/country-codes", h.RefCountryCodes)
			ref.GET("/declass-exceptions", h.RefDeclassExceptions)
			ref.GET("/non-ic-markings", h.RefNonICMarkings)
			ref.GET("/exempt-from", h.RefExemptFrom)
			ref.GET("/complies-with", h.RefCompliesWith)
			ref.GET("/atomic-energy-markings", h.RefAtomicEnergyMarkings)
			ref.GET("/notice-types", h.RefNoticeTypes)
			ref.GET("/sci-controls", h.RefSCIControls)
		}

		v1.POST("/validate", h.Validate)
		v1.POST("/validate/portion", h.ValidatePortion)
		v1.POST("/guidance", h.Guidance)
		v1.POST("/banner", h.Banner)
	}
}

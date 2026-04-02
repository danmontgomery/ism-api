package handler

import (
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// OpenAPISpec serves the embedded OpenAPI specification as YAML.
func (h *Handler) OpenAPISpec(c *gin.Context) {
	data, err := fs.ReadFile(h.docsFS, "openapi.yaml")
	if err != nil {
		c.String(http.StatusInternalServerError, "openapi spec not found")
		return
	}
	c.Data(http.StatusOK, "application/yaml", data)
}

// ScalarDocs serves the Scalar API reference UI.
func (h *Handler) ScalarDocs(c *gin.Context) {
	data, err := fs.ReadFile(h.docsFS, "scalar.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "docs page not found")
		return
	}
	c.Data(http.StatusOK, "text/html", data)
}

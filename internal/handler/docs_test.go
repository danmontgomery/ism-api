package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"dmontgomery/ism-api/api"
	"dmontgomery/ism-api/internal/guidance"
	"dmontgomery/ism-api/internal/guidance/resolvers"
	"dmontgomery/ism-api/internal/handler"
	"dmontgomery/ism-api/internal/refdata"
	"dmontgomery/ism-api/internal/validation"
	"github.com/gin-gonic/gin"
)

func docsRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	reg := refdata.NewRegistry()
	validator := validation.NewEngine(reg)
	guider := guidance.NewEngine(reg,
		&resolvers.ClassificationResolver{},
		&resolvers.CUIResolver{},
		&resolvers.DisseminationResolver{},
		&resolvers.DistributionResolver{},
		&resolvers.AuthorityResolver{},
		&resolvers.DeclassResolver{},
	)
	r := gin.New()
	h := handler.New(reg, validator, guider, api.Content)
	h.Register(r)
	return r
}

func TestOpenAPISpec_ReturnsYAML(t *testing.T) {
	r := docsRouter()
	req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET /openapi.yaml status = %d, want %d", w.Code, http.StatusOK)
	}

	ct := w.Header().Get("Content-Type")
	if ct != "application/yaml" {
		t.Errorf("Content-Type = %q, want %q", ct, "application/yaml")
	}

	body := w.Body.String()
	if !strings.Contains(body, "openapi:") {
		t.Error("response body does not contain 'openapi:' key")
	}
}

func TestOpenAPISpec_ContainsVersion(t *testing.T) {
	r := docsRouter()
	req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	body := w.Body.String()
	if !strings.Contains(body, "3.0.3") {
		t.Error("response body does not contain OpenAPI version '3.0.3'")
	}
}

func TestScalarDocs_ReturnsHTML(t *testing.T) {
	r := docsRouter()
	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET /docs status = %d, want %d", w.Code, http.StatusOK)
	}

	ct := w.Header().Get("Content-Type")
	if ct != "text/html" {
		t.Errorf("Content-Type = %q, want %q", ct, "text/html")
	}

	body := w.Body.String()
	if !strings.Contains(body, "api-reference") {
		t.Error("response body does not contain Scalar 'api-reference' script tag")
	}
}

func TestScalarDocs_ContainsOpenAPIDataURL(t *testing.T) {
	r := docsRouter()
	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	body := w.Body.String()
	if !strings.Contains(body, `data-url="/openapi.yaml"`) {
		t.Error("response body does not contain data-url pointing to /openapi.yaml")
	}
}

func TestScalarDocs_ContainsScalarCDN(t *testing.T) {
	r := docsRouter()
	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	body := w.Body.String()
	if !strings.Contains(body, "cdn.jsdelivr.net/npm/@scalar/api-reference") {
		t.Error("response body does not contain Scalar CDN URL")
	}
}

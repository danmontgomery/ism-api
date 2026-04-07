package api_test

import (
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

// openAPISpec is the minimal structure we need to validate the spec.
type openAPISpec struct {
	OpenAPI    string                       `yaml:"openapi"`
	Paths      map[string]map[string]any    `yaml:"paths"`
	Components struct {
		Schemas map[string]any `yaml:"schemas"`
	} `yaml:"components"`
}

func TestOpenAPISpec(t *testing.T) {
	data, err := os.ReadFile("openapi.yaml")
	if err != nil {
		t.Fatalf("failed to read openapi.yaml: %v", err)
	}

	var spec openAPISpec
	if err := yaml.Unmarshal(data, &spec); err != nil {
		t.Fatalf("failed to parse openapi.yaml: %v", err)
	}

	// Must be OpenAPI 3.0.x
	if spec.OpenAPI == "" {
		t.Fatal("openapi version not set")
	}
	if spec.OpenAPI[:3] != "3.0" {
		t.Errorf("expected OpenAPI 3.0.x, got %s", spec.OpenAPI)
	}

	// All 12 paths must exist.
	expectedPaths := []struct {
		path   string
		method string
	}{
		{"/healthz", "get"},
		{"/api/v1/ref/classifications", "get"},
		{"/api/v1/ref/cui-categories", "get"},
		{"/api/v1/ref/dissemination-controls", "get"},
		{"/api/v1/ref/distribution-statements", "get"},
		{"/api/v1/ref/country-codes", "get"},
		{"/api/v1/ref/declass-exceptions", "get"},
		{"/api/v1/ref/non-ic-markings", "get"},
		{"/api/v1/ref/exempt-from", "get"},
		{"/api/v1/ref/complies-with", "get"},
		{"/api/v1/validate", "post"},
		{"/api/v1/validate/portion", "post"},
		{"/api/v1/guidance", "post"},
		{"/api/v1/banner", "post"},
	}

	for _, ep := range expectedPaths {
		pathItem, ok := spec.Paths[ep.path]
		if !ok {
			t.Errorf("missing path: %s", ep.path)
			continue
		}
		if _, ok := pathItem[ep.method]; !ok {
			t.Errorf("path %s missing method %s", ep.path, ep.method)
		}
	}

	if len(spec.Paths) != 14 {
		t.Errorf("expected 14 paths, got %d", len(spec.Paths))
	}

	// All expected schemas must exist.
	expectedSchemas := []string{
		"ISM",
		"ThirdPartyDistributionContract",
		"Envelope",
		"ErrorDetail",
		"ClassificationEntry",
		"CUICategory",
		"DisseminationControl",
		"DistributionStatement",
		"CountryCode",
		"DeclassException",
		"NonICMarking",
		"ExemptFromEntry",
		"CompliesWithEntry",
		"ValidationResult",
		"FieldError",
		"FieldGuidance",
		"AllowedValue",
		"BannerResult",
	}

	for _, name := range expectedSchemas {
		if _, ok := spec.Components.Schemas[name]; !ok {
			t.Errorf("missing schema: %s", name)
		}
	}
}

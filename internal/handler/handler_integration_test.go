package handler_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"expr.ai/ism-api/api"
	"expr.ai/ism-api/internal/guidance"
	"expr.ai/ism-api/internal/guidance/resolvers"
	"expr.ai/ism-api/internal/handler"
	"expr.ai/ism-api/internal/model"
	"expr.ai/ism-api/internal/refdata"
	"expr.ai/ism-api/internal/validation"
	"github.com/gin-gonic/gin"
)

// testRouter builds a fully-wired Gin engine identical to cmd/server/main.go.
func testRouter() *gin.Engine {
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

// envelope mirrors the API's standard response wrapper for test assertions.
type envelope struct {
	Data   json.RawMessage `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors,omitempty"`
}

// postJSON sends a POST with JSON body and returns the recorder.
func postJSON(r *gin.Engine, path string, body any) *httptest.ResponseRecorder {
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// getJSON sends a GET request and returns the recorder.
func getJSON(r *gin.Engine, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// parseEnvelope decodes the response body into an envelope struct.
func parseEnvelope(t *testing.T, body io.Reader) envelope {
	t.Helper()
	var env envelope
	if err := json.NewDecoder(body).Decode(&env); err != nil {
		t.Fatalf("failed to decode response envelope: %v", err)
	}
	return env
}

// assertContentType checks that the Content-Type header is application/json.
func assertContentType(t *testing.T, w *httptest.ResponseRecorder) {
	t.Helper()
	ct := w.Header().Get("Content-Type")
	if ct != "application/json; charset=utf-8" {
		t.Errorf("Content-Type = %q, want application/json; charset=utf-8", ct)
	}
}

// assertRequestID checks that X-Request-ID is present in the response.
func assertRequestID(t *testing.T, w *httptest.ResponseRecorder) {
	t.Helper()
	if w.Header().Get("X-Request-ID") == "" {
		t.Error("missing X-Request-ID header")
	}
}

// ============================================================
// Health endpoint
// ============================================================

func TestHealth(t *testing.T) {
	r := testRouter()
	w := getJSON(r, "/healthz")
	if w.Code != 200 {
		t.Fatalf("GET /healthz status = %d, want 200", w.Code)
	}
	var body map[string]string
	json.NewDecoder(w.Body).Decode(&body)
	if body["status"] != "ok" {
		t.Errorf("health status = %q, want ok", body["status"])
	}
}

// ============================================================
// Reference data endpoints — all 7
// ============================================================

func TestRefEndpoints(t *testing.T) {
	r := testRouter()

	paths := []struct {
		name  string
		path  string
		field string // field to check exists in each element
	}{
		{"classifications", "/api/v1/ref/classifications", "code"},
		{"cui-categories", "/api/v1/ref/cui-categories", "code"},
		{"dissemination-controls", "/api/v1/ref/dissemination-controls", "code"},
		{"distribution-statements", "/api/v1/ref/distribution-statements", "code"},
		{"country-codes", "/api/v1/ref/country-codes", "code"},
		{"declass-exceptions", "/api/v1/ref/declass-exceptions", "code"},
		{"non-ic-markings", "/api/v1/ref/non-ic-markings", "code"},
	}

	for _, tc := range paths {
		t.Run(tc.name, func(t *testing.T) {
			w := getJSON(r, tc.path)
			if w.Code != 200 {
				t.Fatalf("status = %d, want 200", w.Code)
			}
			assertContentType(t, w)
			assertRequestID(t, w)

			env := parseEnvelope(t, w.Body)
			if len(env.Errors) != 0 {
				t.Fatalf("unexpected errors: %v", env.Errors)
			}
			if env.Data == nil {
				t.Fatal("data is nil")
			}

			// Data should be a non-empty array of objects with the expected field.
			var items []map[string]any
			if err := json.Unmarshal(env.Data, &items); err != nil {
				t.Fatalf("data is not an array: %v", err)
			}
			if len(items) == 0 {
				t.Fatal("data array is empty")
			}
			for i, item := range items {
				if _, ok := item[tc.field]; !ok {
					t.Errorf("item[%d] missing field %q", i, tc.field)
					break
				}
			}
		})
	}
}

func TestRefClassificationsShape(t *testing.T) {
	r := testRouter()
	w := getJSON(r, "/api/v1/ref/classifications")
	env := parseEnvelope(t, w.Body)

	var items []struct {
		Code  string `json:"code"`
		Label string `json:"label"`
		Level int    `json:"level"`
	}
	json.Unmarshal(env.Data, &items)

	if len(items) != 4 {
		t.Fatalf("got %d classifications, want 4", len(items))
	}

	// Check expected fields present on each entry.
	for _, item := range items {
		if item.Code == "" || item.Label == "" {
			t.Errorf("classification entry missing code or label: %+v", item)
		}
	}
}

func TestRefDisseminationControlsShape(t *testing.T) {
	r := testRouter()
	w := getJSON(r, "/api/v1/ref/dissemination-controls")
	env := parseEnvelope(t, w.Body)

	var items []struct {
		Code        string `json:"code"`
		Label       string `json:"label"`
		Description string `json:"description"`
	}
	json.Unmarshal(env.Data, &items)

	if len(items) == 0 {
		t.Fatal("dissemination controls empty")
	}
	for _, item := range items {
		if item.Code == "" || item.Label == "" || item.Description == "" {
			t.Errorf("missing required fields: %+v", item)
		}
	}
}

func TestRefCountryCodesShape(t *testing.T) {
	r := testRouter()
	w := getJSON(r, "/api/v1/ref/country-codes")
	env := parseEnvelope(t, w.Body)

	var items []struct {
		Code string `json:"code"`
		Name string `json:"name"`
		Type string `json:"type"`
	}
	json.Unmarshal(env.Data, &items)

	if len(items) == 0 {
		t.Fatal("country codes empty")
	}
	types := map[string]bool{}
	for _, item := range items {
		types[item.Type] = true
	}
	if !types["country"] {
		t.Error("no entries with type=country")
	}
}

// ============================================================
// Validate endpoint
// ============================================================

func TestValidate_MalformedRequest(t *testing.T) {
	r := testRouter()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/validate", bytes.NewReader([]byte("not json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Fatalf("status = %d, want 400", w.Code)
	}
	assertContentType(t, w)
	env := parseEnvelope(t, w.Body)
	if len(env.Errors) == 0 {
		t.Error("expected errors in response")
	}
	if env.Data != nil && string(env.Data) != "null" {
		t.Errorf("expected null data, got %s", env.Data)
	}
}

func TestValidate_ValidUnclassified(t *testing.T) {
	r := testRouter()
	body := map[string]any{
		"ism": map[string]any{
			"classification": "U",
		},
	}
	w := postJSON(r, "/api/v1/validate", body)

	if w.Code != 200 {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	assertContentType(t, w)
	assertRequestID(t, w)

	env := parseEnvelope(t, w.Body)
	if len(env.Errors) != 0 {
		t.Fatalf("unexpected envelope errors: %v", env.Errors)
	}

	var result struct {
		Valid  bool `json:"valid"`
		Errors []struct {
			Field    string `json:"field"`
			Code     string `json:"code"`
			Message  string `json:"message"`
			Severity string `json:"severity"`
		} `json:"errors"`
	}
	json.Unmarshal(env.Data, &result)
	if !result.Valid {
		t.Errorf("expected valid=true, got errors: %+v", result.Errors)
	}
}

func TestValidate_InvalidClassification(t *testing.T) {
	r := testRouter()
	body := map[string]any{
		"ism": map[string]any{
			"classification": "TOPSECRET",
		},
	}
	w := postJSON(r, "/api/v1/validate", body)

	// Invalid ISM still returns 200 (validation result, not malformed request).
	if w.Code != 200 {
		t.Fatalf("status = %d, want 200", w.Code)
	}

	env := parseEnvelope(t, w.Body)
	var result struct {
		Valid  bool `json:"valid"`
		Errors []struct {
			Field    string `json:"field"`
			Code     string `json:"code"`
			Severity string `json:"severity"`
		} `json:"errors"`
	}
	json.Unmarshal(env.Data, &result)
	if result.Valid {
		t.Error("expected valid=false for invalid classification")
	}
	if len(result.Errors) == 0 {
		t.Error("expected validation errors")
	}
}

func TestValidate_MissingOwnerProducerForSecret(t *testing.T) {
	r := testRouter()
	body := map[string]any{
		"ism": map[string]any{
			"classification": "S",
			"classifiedBy":   "OCA Name",
			"classificationReason": "Reason",
			"declassDate": "20360101",
		},
	}
	w := postJSON(r, "/api/v1/validate", body)

	if w.Code != 200 {
		t.Fatalf("status = %d, want 200", w.Code)
	}

	env := parseEnvelope(t, w.Body)
	var result struct {
		Valid  bool `json:"valid"`
		Errors []struct {
			Code string `json:"code"`
		} `json:"errors"`
	}
	json.Unmarshal(env.Data, &result)
	if result.Valid {
		t.Error("expected valid=false when ownerProducer missing for SECRET")
	}
	found := false
	for _, e := range result.Errors {
		if e.Code == "core.owner_producer_required" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected core.owner_producer_required error, got: %+v", result.Errors)
	}
}

func TestValidatePortion_MalformedRequest(t *testing.T) {
	r := testRouter()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/validate/portion", bytes.NewReader([]byte("{bad")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

func TestValidatePortion_ValidSecret(t *testing.T) {
	r := testRouter()
	body := map[string]any{
		"ism": map[string]any{
			"classification": "S",
			"ownerProducer":  []string{"USA"},
		},
	}
	w := postJSON(r, "/api/v1/validate/portion", body)

	if w.Code != 200 {
		t.Fatalf("status = %d, want 200", w.Code)
	}

	env := parseEnvelope(t, w.Body)
	var result struct {
		Valid bool `json:"valid"`
	}
	json.Unmarshal(env.Data, &result)
	// Portion validation uses same engine; may warn about missing authority but should not error fatally.
}

// ============================================================
// Guidance endpoint
// ============================================================

func TestGuidance_MalformedRequest(t *testing.T) {
	r := testRouter()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/guidance", bytes.NewReader([]byte("{{{")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

func TestGuidance_Unclassified(t *testing.T) {
	r := testRouter()
	body := map[string]any{
		"ism": map[string]any{
			"classification": "U",
		},
	}
	w := postJSON(r, "/api/v1/guidance", body)

	if w.Code != 200 {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	assertContentType(t, w)

	env := parseEnvelope(t, w.Body)
	var result struct {
		Fields []struct {
			Field  string `json:"field"`
			Status string `json:"status"`
		} `json:"fields"`
	}
	json.Unmarshal(env.Data, &result)
	if len(result.Fields) == 0 {
		t.Fatal("expected field guidance entries")
	}

	// For U, CUI fields should be not_applicable.
	for _, f := range result.Fields {
		if f.Field == "categoryMarkings" && f.Status != "not_applicable" {
			t.Errorf("categoryMarkings status = %q, want not_applicable for U", f.Status)
		}
	}
}

func TestGuidance_Secret(t *testing.T) {
	r := testRouter()
	body := map[string]any{
		"ism": map[string]any{
			"classification": "S",
			"ownerProducer":  []string{"USA"},
		},
	}
	w := postJSON(r, "/api/v1/guidance", body)

	if w.Code != 200 {
		t.Fatalf("status = %d, want 200", w.Code)
	}

	env := parseEnvelope(t, w.Body)
	var result struct {
		Fields []struct {
			Field  string `json:"field"`
			Status string `json:"status"`
		} `json:"fields"`
	}
	json.Unmarshal(env.Data, &result)

	fieldStatus := map[string]string{}
	for _, f := range result.Fields {
		fieldStatus[f.Field] = f.Status
	}

	// For SECRET, declass fields should be available, CUI fields not_applicable.
	if fieldStatus["categoryMarkings"] != "not_applicable" {
		t.Errorf("categoryMarkings = %q, want not_applicable for S", fieldStatus["categoryMarkings"])
	}
}

// ============================================================
// Banner endpoint
// ============================================================

func TestBanner_MalformedRequest(t *testing.T) {
	r := testRouter()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/banner", bytes.NewReader([]byte("not-json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

func TestBanner_Unclassified(t *testing.T) {
	r := testRouter()
	body := map[string]any{
		"ism": map[string]any{
			"classification": "U",
		},
	}
	w := postJSON(r, "/api/v1/banner", body)

	if w.Code != 200 {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	assertContentType(t, w)

	env := parseEnvelope(t, w.Body)
	var result struct {
		BannerLine  string `json:"bannerLine"`
		PortionMark string `json:"portionMark"`
	}
	json.Unmarshal(env.Data, &result)

	if result.BannerLine != "UNCLASSIFIED" {
		t.Errorf("bannerLine = %q, want UNCLASSIFIED", result.BannerLine)
	}
	if result.PortionMark != "(U)" {
		t.Errorf("portionMark = %q, want (U)", result.PortionMark)
	}
}

func TestBanner_Secret(t *testing.T) {
	r := testRouter()
	body := map[string]any{
		"ism": map[string]any{
			"classification": "S",
			"ownerProducer":  []string{"USA"},
		},
	}
	w := postJSON(r, "/api/v1/banner", body)

	if w.Code != 200 {
		t.Fatalf("status = %d, want 200", w.Code)
	}

	env := parseEnvelope(t, w.Body)
	var result struct {
		BannerLine  string `json:"bannerLine"`
		PortionMark string `json:"portionMark"`
	}
	json.Unmarshal(env.Data, &result)

	if result.BannerLine != "SECRET" {
		t.Errorf("bannerLine = %q, want SECRET", result.BannerLine)
	}
	if result.PortionMark != "(S)" {
		t.Errorf("portionMark = %q, want (S)", result.PortionMark)
	}
}

func TestBanner_AuthorityBlock(t *testing.T) {
	r := testRouter()
	body := map[string]any{
		"ism": map[string]any{
			"classification":       "S",
			"ownerProducer":        []string{"USA"},
			"classifiedBy":         "John Smith",
			"classificationReason": "1.4(a)",
			"declassDate":          "20360101",
		},
	}
	w := postJSON(r, "/api/v1/banner", body)

	if w.Code != 200 {
		t.Fatalf("status = %d, want 200", w.Code)
	}

	env := parseEnvelope(t, w.Body)
	var result struct {
		BannerLine     string `json:"bannerLine"`
		PortionMark    string `json:"portionMark"`
		AuthorityBlock string `json:"authorityBlock"`
	}
	json.Unmarshal(env.Data, &result)

	if result.BannerLine != "SECRET" {
		t.Errorf("bannerLine = %q, want SECRET", result.BannerLine)
	}
	wantAuth := "Classified By: John Smith\nReason: 1.4(a)\nDeclassify On: 20360101"
	if result.AuthorityBlock != wantAuth {
		t.Errorf("authorityBlock = %q, want %q", result.AuthorityBlock, wantAuth)
	}
}

func TestBanner_AuthorityBlockEmpty(t *testing.T) {
	r := testRouter()
	body := map[string]any{
		"ism": map[string]any{
			"classification": "U",
		},
	}
	w := postJSON(r, "/api/v1/banner", body)

	if w.Code != 200 {
		t.Fatalf("status = %d, want 200", w.Code)
	}

	env := parseEnvelope(t, w.Body)
	var result struct {
		AuthorityBlock string `json:"authorityBlock"`
	}
	json.Unmarshal(env.Data, &result)

	if result.AuthorityBlock != "" {
		t.Errorf("authorityBlock = %q, want empty for U", result.AuthorityBlock)
	}
}

// ============================================================
// Middleware: X-Request-ID pass-through
// ============================================================

func TestMiddleware_RequestIDPassThrough(t *testing.T) {
	r := testRouter()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	req.Header.Set("X-Request-ID", "test-id-12345")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if got := w.Header().Get("X-Request-ID"); got != "test-id-12345" {
		t.Errorf("X-Request-ID = %q, want test-id-12345", got)
	}
}

func TestMiddleware_RequestIDGenerated(t *testing.T) {
	r := testRouter()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Header().Get("X-Request-ID") == "" {
		t.Error("expected generated X-Request-ID when not provided")
	}
}

// ============================================================
// Smoke tests: full ISM spec examples
// ============================================================

// smokeISM is a helper for smoke tests — validates ISM is valid AND renders expected banner.
func smokeISM(t *testing.T, name string, ism model.ISM, wantBanner, wantPortion string) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		r := testRouter()

		// Validate: should be valid (or at least return 200).
		valBody := map[string]any{"ism": ism}
		vw := postJSON(r, "/api/v1/validate", valBody)
		if vw.Code != 200 {
			t.Fatalf("validate status = %d, want 200", vw.Code)
		}
		vEnv := parseEnvelope(t, vw.Body)
		var vResult struct {
			Valid  bool `json:"valid"`
			Errors []struct {
				Code     string `json:"code"`
				Severity string `json:"severity"`
			} `json:"errors"`
		}
		json.Unmarshal(vEnv.Data, &vResult)
		// Check no error-severity findings (warnings are OK).
		for _, e := range vResult.Errors {
			if e.Severity == "error" {
				t.Errorf("unexpected validation error: %+v", e)
			}
		}

		// Banner: should render expected output.
		banBody := map[string]any{"ism": ism}
		bw := postJSON(r, "/api/v1/banner", banBody)
		if bw.Code != 200 {
			t.Fatalf("banner status = %d, want 200", bw.Code)
		}
		bEnv := parseEnvelope(t, bw.Body)
		var bResult struct {
			BannerLine  string `json:"bannerLine"`
			PortionMark string `json:"portionMark"`
		}
		json.Unmarshal(bEnv.Data, &bResult)

		if bResult.BannerLine != wantBanner {
			t.Errorf("bannerLine = %q, want %q", bResult.BannerLine, wantBanner)
		}
		if bResult.PortionMark != wantPortion {
			t.Errorf("portionMark = %q, want %q", bResult.PortionMark, wantPortion)
		}
	})
}

func TestSmoke_Unclassified(t *testing.T) {
	smokeISM(t, "basic-U", model.ISM{
		Classification: model.ClassificationU,
	}, "UNCLASSIFIED", "(U)")
}

func TestSmoke_UnclassifiedWithDistribution(t *testing.T) {
	smokeISM(t, "U-with-distA", model.ISM{
		Classification:        model.ClassificationU,
		DistributionStatement: "A",
	}, "UNCLASSIFIED", "(U)")
}

func TestSmoke_CUI_Basic(t *testing.T) {
	smokeISM(t, "CUI-basic", model.ISM{
		Classification: model.ClassificationCUI,
	}, "CUI", "(CUI)")
}

func TestSmoke_CUI_WithCategories(t *testing.T) {
	smokeISM(t, "CUI-categories", model.ISM{
		Classification:   model.ClassificationCUI,
		CategoryMarkings: []string{"SP-CEII", "SP-PCII"},
	}, "CUI//SP-CEII/SP-PCII", "(CUI//SP-CEII/SP-PCII)")
}

func TestSmoke_Confidential(t *testing.T) {
	smokeISM(t, "C-basic", model.ISM{
		Classification:     model.ClassificationC,
		OwnerProducer:      []string{"USA"},
		ClassifiedBy:       "OCA Name",
		ClassificationReason: "Reason 1.4(a)",
		DeclassDate:        "20360101",
	}, "CONFIDENTIAL", "(C)")
}

func TestSmoke_Secret_NOFORN(t *testing.T) {
	smokeISM(t, "S-NOFORN", model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"NOFORN"},
		ClassifiedBy:          "OCA Name",
		ClassificationReason:  "Reason 1.4(c)",
		DeclassDate:           "20360101",
	}, "SECRET//NOFORN", "(S//NF)")
}

func TestSmoke_Secret_REL(t *testing.T) {
	smokeISM(t, "S-REL", model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"REL"},
		ReleasableTo:          []string{"USA", "GBR"},
		ClassifiedBy:          "OCA Name",
		ClassificationReason:  "Reason 1.4(a)",
		DeclassDate:           "20360101",
	}, "SECRET//REL TO USA, GBR", "(S//REL TO USA, GBR)")
}

func TestSmoke_Secret_MultipleControls(t *testing.T) {
	smokeISM(t, "S-multi-controls", model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA"},
		DisseminationControls: []string{"IMCON", "NOFORN"},
		ClassifiedBy:          "OCA Name",
		ClassificationReason:  "Reason 1.4(a)",
		DeclassDate:           "20360101",
	}, "SECRET//IMCON/NOFORN", "(S//IMC/NF)")
}

func TestSmoke_Secret_FGI(t *testing.T) {
	smokeISM(t, "S-FGI", model.ISM{
		Classification:     model.ClassificationS,
		OwnerProducer:      []string{"USA"},
		FGISourceOpen:      []string{"GBR"},
		ClassifiedBy:       "OCA Name",
		ClassificationReason: "Reason 1.4(a)",
		DeclassDate:        "20360101",
	}, "SECRET//FGI GBR", "(S//FGI)")
}

func TestSmoke_Secret_DerivedClassification(t *testing.T) {
	smokeISM(t, "S-derived", model.ISM{
		Classification:           model.ClassificationS,
		OwnerProducer:            []string{"USA"},
		DerivativelyClassifiedBy: "Classifier Name",
		DerivedFrom:              "SCG v3.0",
		DeclassDate:              "20360101",
	}, "SECRET", "(S)")
}

func TestSmoke_Secret_Joint(t *testing.T) {
	smokeISM(t, "S-joint-USA-GBR", model.ISM{
		Classification:       model.ClassificationS,
		OwnerProducer:        []string{"USA", "GBR"},
		Joint:                true,
		ClassifiedBy:         "OCA Name",
		ClassificationReason: "Reason 1.4(a)",
		DeclassDate:          "20360101",
	}, "JOINT SECRET USA GBR", "(JS USA GBR)")
}

func TestSmoke_Secret_Joint_WithNOFORN(t *testing.T) {
	smokeISM(t, "S-joint-NOFORN", model.ISM{
		Classification:        model.ClassificationS,
		OwnerProducer:         []string{"USA", "GBR"},
		Joint:                 true,
		DisseminationControls: []string{"NOFORN"},
		ClassifiedBy:          "OCA Name",
		ClassificationReason:  "Reason 1.4(a)",
		DeclassDate:           "20360101",
	}, "JOINT SECRET USA GBR//NOFORN", "(JS USA GBR//NF)")
}

// ============================================================
// Error scenarios: 400 vs 200
// ============================================================

func TestValidate_EmptyBody(t *testing.T) {
	r := testRouter()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/validate", bytes.NewReader([]byte("")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Errorf("empty body status = %d, want 400", w.Code)
	}
}

func TestBanner_EmptyBody(t *testing.T) {
	r := testRouter()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/banner", bytes.NewReader([]byte("")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Errorf("empty body status = %d, want 400", w.Code)
	}
}

func TestGuidance_EmptyBody(t *testing.T) {
	r := testRouter()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/guidance", bytes.NewReader([]byte("")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Errorf("empty body status = %d, want 400", w.Code)
	}
}

func TestValidate_InvalidISMReturns200WithErrors(t *testing.T) {
	r := testRouter()

	cases := []struct {
		name string
		body map[string]any
	}{
		{
			"bad-classification",
			map[string]any{"ism": map[string]any{"classification": "BOGUS"}},
		},
		{
			"NOFORN-without-classification-gate",
			map[string]any{"ism": map[string]any{
				"classification":        "U",
				"disseminationControls": []string{"NOFORN"},
			}},
		},
		{
			"CUI-bad-category",
			map[string]any{"ism": map[string]any{
				"classification":   "CUI",
				"categoryMarkings": []string{"FAKE-CAT"},
			}},
		},
		{
			"exclusive-controls",
			map[string]any{"ism": map[string]any{
				"classification":        "S",
				"ownerProducer":         []string{"USA"},
				"disseminationControls": []string{"NOFORN", "REL"},
				"releasableTo":          []string{"GBR"},
			}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := postJSON(r, "/api/v1/validate", tc.body)
			if w.Code != 200 {
				t.Fatalf("status = %d, want 200 for invalid ISM", w.Code)
			}

			env := parseEnvelope(t, w.Body)
			var result struct {
				Valid  bool `json:"valid"`
				Errors []struct {
					Code string `json:"code"`
				} `json:"errors"`
			}
			json.Unmarshal(env.Data, &result)
			if result.Valid {
				t.Error("expected valid=false")
			}
			if len(result.Errors) == 0 {
				t.Error("expected at least one validation error")
			}
		})
	}
}

// ============================================================
// 404 on unknown routes
// ============================================================

func TestUnknownRoute_Returns404(t *testing.T) {
	r := testRouter()
	w := getJSON(r, "/api/v1/nonexistent")
	if w.Code != 404 {
		t.Errorf("status = %d, want 404", w.Code)
	}
}

func TestWrongMethod_Returns404Or405(t *testing.T) {
	r := testRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/validate", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Gin returns 404 by default for wrong method (405 requires HandleMethodNotAllowed=true).
	if w.Code != 404 && w.Code != 405 {
		t.Errorf("status = %d, want 404 or 405", w.Code)
	}
}

// ============================================================
// CORS middleware
// ============================================================

func TestCORS_SimpleRequest(t *testing.T) {
	r := testRouter()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("Access-Control-Allow-Origin = %q, want *", got)
	}
}

func TestCORS_PreflightRequest(t *testing.T) {
	r := testRouter()
	req := httptest.NewRequest(http.MethodOptions, "/api/v1/validate", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 204 {
		t.Fatalf("preflight status = %d, want 204", w.Code)
	}
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("Access-Control-Allow-Origin = %q, want *", got)
	}
	if got := w.Header().Get("Access-Control-Allow-Methods"); got == "" {
		t.Error("missing Access-Control-Allow-Methods header")
	}
	if got := w.Header().Get("Access-Control-Allow-Headers"); got == "" {
		t.Error("missing Access-Control-Allow-Headers header")
	}
}

func TestCORS_PostEndpointIncludesHeaders(t *testing.T) {
	r := testRouter()
	body := map[string]any{
		"ism": map[string]any{"classification": "U"},
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/validate", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "http://localhost:3000")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("Access-Control-Allow-Origin = %q, want *", got)
	}
}

func TestCORS_NoOriginStillSetsHeader(t *testing.T) {
	r := testRouter()
	w := getJSON(r, "/healthz")

	// CORS headers should still be present even without Origin header,
	// so server-to-server calls aren't affected.
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("Access-Control-Allow-Origin = %q, want *", got)
	}
}

func TestCORS_BrowserClientFlow(t *testing.T) {
	// Simulates the example client's full browser flow:
	// 1. Preflight + GET each reference-data endpoint
	// 2. Preflight + POST to /api/v1/banner
	// All responses must include CORS headers.
	r := testRouter()
	origin := "http://localhost:3000"

	refEndpoints := []string{
		"/api/v1/ref/classifications",
		"/api/v1/ref/cui-categories",
		"/api/v1/ref/dissemination-controls",
		"/api/v1/ref/distribution-statements",
		"/api/v1/ref/country-codes",
		"/api/v1/ref/declass-exceptions",
		"/api/v1/ref/non-ic-markings",
	}

	for _, ep := range refEndpoints {
		req := httptest.NewRequest(http.MethodGet, ep, nil)
		req.Header.Set("Origin", origin)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Errorf("%s: status = %d, want 200", ep, w.Code)
		}
		if got := w.Header().Get("Access-Control-Allow-Origin"); got != "*" {
			t.Errorf("%s: Access-Control-Allow-Origin = %q, want *", ep, got)
		}
	}

	// Preflight for POST /api/v1/banner (browser sends OPTIONS first)
	preflight := httptest.NewRequest(http.MethodOptions, "/api/v1/banner", nil)
	preflight.Header.Set("Origin", origin)
	preflight.Header.Set("Access-Control-Request-Method", "POST")
	preflight.Header.Set("Access-Control-Request-Headers", "Content-Type")
	pw := httptest.NewRecorder()
	r.ServeHTTP(pw, preflight)

	if pw.Code != 204 {
		t.Fatalf("preflight /api/v1/banner: status = %d, want 204", pw.Code)
	}
	if got := pw.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("preflight /api/v1/banner: Access-Control-Allow-Origin = %q, want *", got)
	}

	// Actual POST to /api/v1/banner (the request the browser sends after preflight)
	body := map[string]any{
		"classification": "U",
		"ownerProducer":  []string{"USA"},
	}
	b, _ := json.Marshal(body)
	postReq := httptest.NewRequest(http.MethodPost, "/api/v1/banner", bytes.NewReader(b))
	postReq.Header.Set("Content-Type", "application/json")
	postReq.Header.Set("Origin", origin)
	postW := httptest.NewRecorder()
	r.ServeHTTP(postW, postReq)

	if postW.Code != 200 {
		t.Fatalf("POST /api/v1/banner: status = %d, want 200", postW.Code)
	}
	if got := postW.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("POST /api/v1/banner: Access-Control-Allow-Origin = %q, want *", got)
	}
}

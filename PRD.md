# ISM API — Product Requirements Document

## Context

DoD and IC documents require precise security markings governed by multiple overlapping standards (DoDI 5200.48, DoDI 5230.24, ICD-710, ISM.XML). Getting these markings right is error-prone — fields are interdependent, certain controls are mutually exclusive, ordering rules apply, and the valid options change based on classification level. This API provides a programmatic interface for a client application to **construct valid ISM markings through guided workflows** and **validate completed markings** against the governing rules.

The data format follows the [json-ism spec](https://github.com/540co/json-ism), which maps IC-ISM XML attributes to JSON for REST API consumption.

## Scope

**In scope:** Unclassified (U), CUI, Confidential (C), Secret (S)
**Out of scope (for now):** Top Secret, SCI controls, SAR identifiers, Atomic Energy markings. The architecture will support adding these later with only additive changes.

## Technical Decisions

| Decision | Choice |
|----------|--------|
| Language | Go |
| Router | Gin |
| Persistence | Stateless — all reference data compiled into binary |
| Output format | JSON-ISM objects |
| Consumer | Internal client application |
| Auth | None (internal tool) |

---

## Data Model

The core ISM object per the json-ism spec (TS-specific fields excluded):

### Core Classification
- `version` (string) — JSON-ISM version identifier
- `classification` (string) — `U`, `CUI`, `C`, `S`
- `ownerProducer` ([]string) — national governments/orgs (ISO 3166 codes, e.g., `USA`)
- `joint` (bool) — true when multiple ownerProducers share joint ownership

### CUI Fields
- `categoryMarkings` ([]string) — CUI categories, SP-prefixed for Specified
  - **Specified:** SP-CEII, SP-CRITAN, SP-CVI, SP-PCII, SP-PHYS, SP-TSCA
  - **Basic:** CRIT, EMGT, ISVI, PHYS, SAFE, WATER
- `controlledByName` (string) — DoD component determining CUI status
- `controlledByOffice` (string) — DoD office
- `poc` (string) — point of contact (email or phone)

### Dissemination Controls
- `disseminationControls` ([]string) — RS, OC, OC-USGOV, IMCON, NOFORN, PROPIN, REL, RELIDO, EYES, DSEN, FISA, DISPLAY ONLY, FED ONLY, FEDCON, NOCON, DL ONLY, REL TO USA, LIST
- `releasableTo` ([]string) — country/org codes (required when REL is set)
- `displayOnlyTo` ([]string) — country/org codes (required when DISPLAY ONLY is set)

### Distribution Statements (DoDI 5230.24)
- `distributionStatement` (string) — Statements A through F
- `3rdPartyDistributionStatement` (string) — contractor IP rights category
- `3rdPartyDistributionWarning` (string) — required legal language
- `3rdPartyDistributionContract` (object) — `contractNumber`, `contractorName`, `contractorAddress`, `expirationDate`
- `copyright` (string)

### Classification Authority
- `classifiedBy` (string) — original classifier identity
- `classificationReason` (string) — basis for classification (cites EO 13526 sec 1.4)
- `derivativelyClassifiedBy` (string) — derivative classifier identity
- `derivedFrom` (string) — source material citation
- `compilationReason` (string) — justification when compilation classification exceeds components

### Declassification
- `declassDate` (string) — automatic declass date (YYYYMMDD or ISO 8601)
- `declassEvent` (string) — event-based declass trigger
- `declassException` (string) — 25X exemption code

### Foreign Government Information
- `fgiSourceOpen` ([]string) — non-concealed FGI sources
- `fgiSourceProtected` ([]string) — classified FGI sources

### Other
- `nonICMarkings` ([]string)
- `nonUSControls` ([]string)
- `bannerLine` (string) — rendered convenience field

---

## API Endpoints

All responses use a consistent envelope: `{ "data": ..., "errors": [...] }`

### Health
| Method | Path | Description |
|--------|------|-------------|
| GET | `/healthz` | Returns `{"status": "ok"}` |

### Reference Data (`GET /api/v1/ref/*`)

Stateless lookups — client caches these. No parameters.

| Path | Returns |
|------|---------|
| `/ref/classifications` | `[]{ code, label, level }` |
| `/ref/cui-categories` | `[]{ code, label, type, description }` — type is "specified" or "basic" |
| `/ref/dissemination-controls` | `[]{ code, label, description, requiresField, minClassification, exclusiveWith }` |
| `/ref/distribution-statements` | `[]{ code, label, text, classificationConstraint }` |
| `/ref/country-codes` | `[]{ code, name, type }` — type is "country", "coalition", or "organization" |
| `/ref/declass-exceptions` | `[]{ code, label, category }` |
| `/ref/non-ic-markings` | `[]{ code, label }` |

The dissemination-controls endpoint includes compatibility metadata (`requiresField`, `exclusiveWith`, `minClassification`) so the client can provide real-time UI hints.

### Validation (`POST /api/v1/validate`)

Accepts a complete ISM object, runs all applicable rules, returns the full set of errors and warnings (does not short-circuit).

**Request:** `{ "ism": { ...full ISM object } }`
**Response:** `{ "data": { "valid": bool, "errors": [...], "warnings": [...] } }`

Each error/warning:
```json
{
  "field": "releasableTo",
  "code": "REL_REQUIRES_RELEASABLE_TO",
  "message": "REL dissemination control requires releasableTo to specify recipient countries",
  "severity": "error"
}
```

- `errors` = hard failures (marking is invalid)
- `warnings` = best-practice guidance (marking is technically valid but may be incomplete)
- Error codes are stable string constants for client consumption
- Returns HTTP 200 even for invalid ISM — the call succeeded, the ISM didn't. HTTP 400 is for malformed requests only.

**Portion-level validation:** `POST /api/v1/validate/portion` — same shape, different rule set (no `version`, `distributionStatement` etc. required at portion level).

### Guidance (`POST /api/v1/guidance`)

Powers the wizard-style UI. Client sends the ISM in its current partial state; API returns which fields are relevant, required, and what values are valid.

**Request:** `{ "ism": { ...partial ISM } }`
**Response:**
```json
{
  "data": {
    "fields": [
      {
        "field": "categoryMarkings",
        "status": "available",
        "required": false,
        "requiredIf": "CUI Specified marking is intended",
        "allowedValues": [
          { "code": "SP-CEII", "label": "Critical Energy Infrastructure Information", "type": "specified" }
        ]
      },
      {
        "field": "classifiedBy",
        "status": "not_applicable",
        "reason": "Only applicable for Confidential or Secret classifications"
      }
    ]
  }
}
```

Field statuses:
- **available** — field can be set given current state
- **required** — field must be set given current state (e.g., `releasableTo` when REL is selected)
- **not_applicable** — field is irrelevant at this classification level
- **locked** — value determined by other selections (e.g., `joint` auto-set when >1 ownerProducer)

### Banner Rendering (`POST /api/v1/banner`)

**Request:** `{ "ism": { ...ISM object } }`
**Response:** `{ "data": { "bannerLine": "SECRET//NOFORN", "portionMark": "(S//NF)" } }`

---

## Validation Rules

### Core (always apply)
- `classification` must be a known value (U, CUI, C, S)
- `ownerProducer` required when classification is not U
- `joint` must be true iff len(ownerProducer) > 1

### CUI (when classification == CUI)
- `categoryMarkings` must be valid codes from the known set
- `categoryMarkings` must be alphabetized: SP-prefixed before Basic, alpha within each group
- CUI Specified requires at least one SP- category in categoryMarkings

### Classified (when classification in C, S)
- Should have `classifiedBy` or `derivativelyClassifiedBy` (warning)
- If `derivativelyClassifiedBy` set, should have `derivedFrom` (warning)
- If `classifiedBy` set, should have `classificationReason` (warning)

### Dissemination
- REL requires non-empty `releasableTo`
- DISPLAY ONLY requires non-empty `displayOnlyTo`
- NOFORN and REL are mutually exclusive
- All codes must be in the known set
- Certain controls require minimum classification level

### Distribution
- Statement A is valid only for U/CUI
- Statement code must be A-F
- `3rdPartyDistributionStatement` requires both `3rdPartyDistributionWarning` and `3rdPartyDistributionContract` (all four sub-fields)

### Declassification
- `declassDate` and `declassEvent` are mutually exclusive
- `declassException` must be a known 25X code
- Declass fields only apply to C/S

### FGI
- Country codes in `fgiSourceOpen`/`fgiSourceProtected` must be valid

---

## Project Structure

```
ism-api/
├── cmd/server/main.go                      # Entrypoint: wire Gin router, inject deps
├── internal/
│   ├── model/                              # Domain types
│   │   ├── ism.go                          # Core ISM struct
│   │   ├── classification.go              # Level enum + ordering
│   │   └── contract.go                    # ThirdPartyDistributionContract
│   ├── refdata/                            # Compiled-in reference data
│   │   ├── refdata.go                     # Registry aggregator
│   │   ├── classifications.go
│   │   ├── cui_categories.go
│   │   ├── dissemination_controls.go
│   │   ├── distribution_statements.go
│   │   ├── country_codes.go
│   │   ├── declass_exceptions.go
│   │   └── compatibility.go              # Cross-field exclusion/dependency matrices
│   ├── validation/
│   │   ├── engine.go                      # Runs all rules, collects results
│   │   ├── result.go                      # ValidationResult, FieldError types
│   │   └── rules/
│   │       ├── rule.go                    # Rule interface: Name(), Applies(), Validate()
│   │       ├── core.go                    # Classification + ownerProducer rules
│   │       ├── cui.go                     # CUI category rules
│   │       ├── classified.go             # Authority block rules
│   │       ├── dissemination.go          # Dissemination control rules
│   │       ├── distribution.go           # Distribution statement rules
│   │       ├── third_party.go            # 3rd party contract rules
│   │       ├── declass.go               # Declassification rules
│   │       └── fgi.go                    # FGI source rules
│   ├── guidance/
│   │   ├── engine.go                      # Runs all resolvers
│   │   ├── options.go                     # FieldGuidance, AllowedValue types
│   │   └── resolvers/
│   │       ├── resolver.go               # Resolver interface: Fields(), Resolve()
│   │       ├── classification.go
│   │       ├── cui.go
│   │       ├── dissemination.go
│   │       ├── distribution.go
│   │       ├── authority.go
│   │       └── declass.go
│   ├── banner/
│   │   └── render.go                      # ISM -> bannerLine + portionMark
│   └── handler/
│       ├── handler.go                     # Handler struct + constructor
│       ├── health.go
│       ├── refdata.go
│       ├── validate.go
│       ├── guidance.go
│       ├── banner.go
│       ├── middleware.go                  # Request-ID, logging, recovery
│       └── response.go                   # Envelope helpers
├── go.mod
├── go.sum
└── Makefile                               # build, test, lint, run targets
```

### Key Architectural Patterns

**Validation engine** — each rule implements `Rule` interface with `Applies(ism) bool` and `Validate(ism) []FieldError`. Rules self-select based on classification level. The engine runs all applicable rules and aggregates results. Adding TS support = adding new rule files + registering them.

**Guidance engine** — each resolver implements `Resolver` interface with `Fields() []string` and `Resolve(ism) []FieldGuidance`. Resolvers inspect the partial ISM state and return what's valid/required/not_applicable. Same self-selection pattern as validation rules.

**Reference data** — compiled as Go variables (no file I/O, no DB). The `compatibility.go` file contains cross-field matrices (exclusions, required fields, min levels) consumed by both validation rules and guidance resolvers as a single source of truth.

---

## Top Secret Extensibility

Adding TS later requires only additive changes:
1. Add `TS` to classification enum and refdata
2. Add `sciControls`, `sarIdentifiers`, `atomicEnergyMarkings` fields to ISM struct
3. Add new rule files (`rules/sci.go`, `rules/sar.go`, `rules/atomic.go`)
4. Add new resolver files (`resolvers/sci.go`)
5. Add new refdata files (`refdata/sci_controls.go`)
6. Register in engines

Zero changes to existing rules, resolvers, or handlers.

---

## Verification

1. **Unit tests** — table-driven tests for every validation rule and guidance resolver. Each test case: input ISM + expected errors/guidance.
2. **Integration tests** — `httptest` against the Gin router. POST JSON, assert response shape and status codes.
3. **End-to-end smoke tests** — construct ISM objects matching the json-ism spec examples (unclassified.json, cui.json, secret.json) and validate they pass.
4. **Banner rendering tests** — ISM in, expected banner/portion-mark strings out.
5. **Run locally** — `make run`, hit endpoints with curl/httpie to verify request/response format.

# ISM API

A REST API for constructing, validating, and rendering DoD/IC security classification markings per the [json-ism spec](https://github.com/540co/json-ism).

## Quick Start

```bash
# Start the API server (default port 8080)
go run ./cmd/server

# Or specify a custom port
PORT=9090 go run ./cmd/server
```

## Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `PORT` | HTTP server listen port | `8080` | No |

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/healthz` | Health check |
| GET | `/api/v1/ref/classifications` | Classification levels |
| GET | `/api/v1/ref/cui-categories` | CUI category codes |
| GET | `/api/v1/ref/dissemination-controls` | Dissemination control codes |
| GET | `/api/v1/ref/distribution-statements` | Distribution statement codes |
| GET | `/api/v1/ref/country-codes` | Country/org codes |
| GET | `/api/v1/ref/declass-exceptions` | Declassification exception codes |
| GET | `/api/v1/ref/non-ic-markings` | Non-IC marking codes |
| POST | `/api/v1/validate` | Validate a complete ISM object |
| POST | `/api/v1/validate/portion` | Validate a portion-level ISM object |
| POST | `/api/v1/guidance` | Get field-level guidance for partial ISM state |
| POST | `/api/v1/banner` | Render banner line and portion mark |

All responses use the envelope format: `{ "data": ..., "errors": [...] }`.

## Sample Client

A vanilla JavaScript/HTML demo application is included in [`examples/client/`](examples/client/) that demonstrates all API capabilities with zero dependencies.

### Running the Demo

1. **Start the API server:**

   ```bash
   go run ./cmd/server
   ```

2. **Open the client:**

   Open `examples/client/index.html` in a browser. The client defaults to `http://localhost:8080` as the API base URL.

   > **CORS note:** If serving the HTML from a different origin (e.g., a static file server on another port), the API server must allow cross-origin requests. For local development, you can use a simple proxy or serve both from the same origin. When opening via `file://` protocol, some browsers block fetch requests — use a local HTTP server instead:
   >
   > ```bash
   > # Python 3
   > cd examples/client && python3 -m http.server 3000
   > # Then open http://localhost:3000
   > ```

3. **Walk through a demo flow:**

   - Click **Connect** to load reference data from the API
   - Select a **classification level** (e.g., SECRET) — the wizard dynamically shows/hides fields based on the guidance engine
   - Choose **Owner/Producer** countries (required for Confidential/Secret)
   - Add **dissemination controls** (e.g., NOFORN, REL) — conditional fields like "Releasable To" appear as needed
   - Fill in the **authority block** (classifiedBy, derivedFrom, etc.)
   - Set **declassification** date, event, or exception
   - Add **FGI sources** or **non-IC markings** as applicable
   - Watch the **banner line** and **portion mark** update in real time
   - Click **Validate** to run full ISM validation and see errors/warnings
   - Expand the **ISM Object (JSON)** panel to see the raw payload

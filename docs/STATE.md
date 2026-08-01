# STATE
<!-- Overwritten at the end of every working session. Keep under ~20 lines. -->
_Last updated: 2026-07-28_

## Works now
- `/health` live on Railway; config loads from env (`PORT` 8080, `GOOGLE_PLACES_API_KEY`
  required + fails fast, `ENVIRONMENT` development)
- `godotenv` loads `.env` locally; Railway injects env in prod
- `http.Server` has read/write/idle timeouts (5s/10s/60s) — see concepts/http-server-timeouts.md
- `models.Place` struct defined (not yet used by any handler)
- `internal/places/client.go` — `Client` struct + `New()` constructor. Base URL is a
  field (defaulted, override for httptest); own `http.Client` at 5s. No request helper yet.
- Google API validated by hand via curl — fixtures saved in `internal/places/testdata/`
  (autocomplete, details save-tier, details enrich-tier). Gemini `generativeSummary` +
  `reviewSummary` both return populated — the whole reason this backend exists.

## In flight
- `client.go` next keystroke: the shared authenticated-request/decode helper
  (build req → attach key + field mask → execute → defer close → status check → decode
  into caller dest). Then endpoint siblings.

## Next
1. Request helper → `autocomplete.go` → `details.go` (tier→field-mask) → `photos.go`
2. Handlers for the three endpoints, validating inputs before any Google call
3. Point iOS `PlacesService` at savor-api instead of the SDK

## Landmines
- Parity target is **Autocomplete**, not Text Search — iOS never text-searches
- `models.Place.PriceLevel` still `string`; decided format is nullable int — change with `details.go`
- Attributions + Gemini disclosure text both need a home — already in the details JSON, see TODO
- `WriteTimeout` (10s) bounds the whole handler — `client.go`'s `http.Client` (5s) sits below it ✓
- Match highlighting dropped for Phase 1 (see decisions.md) — iOS `PlaceSuggestion` can drop `AttributedString`

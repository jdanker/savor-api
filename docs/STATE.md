# STATE
<!-- Overwritten at the end of every working session. Keep under ~20 lines. -->
_Last updated: 2026-07-26_

## Works now
- `/health` live on Railway; config loads from env (`PORT` 8080, `GOOGLE_PLACES_API_KEY`
  required + fails fast, `ENVIRONMENT` development)
- `godotenv` loads `.env` locally; Railway injects env in prod
- `models.Place` struct defined (not yet used by any handler)

## In flight
- All three Phase 1 design blockers resolved (see decisions.md): session tokens,
  photo delivery, price-level wire format. Client code is unblocked.
- Next keystroke: `internal/places/client.go` — `Client` struct, constructor, shared
  authenticated-request/decode helper. Endpoint files come after, as siblings.

## Next
1. `client.go` plumbing → `autocomplete.go` → `details.go` (tier→field-mask) → `photos.go`
2. Handlers for the three endpoints, validating inputs before any Google call
3. Point iOS `PlacesService` at savor-api instead of the SDK

## Landmines
- Parity target is **Autocomplete**, not Text Search — iOS never text-searches
- `models.Place.PriceLevel` still `string`; decided format is nullable int — change with `details.go`
- No attribution field on `Place`; Google's terms require displaying photo attributions
- `http.Server` has no Read/Write/Idle timeouts; `go mod tidy` still owed (godotenv `// indirect`)

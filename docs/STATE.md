# STATE
<!-- Overwritten at the end of every working session. Keep under ~20 lines. -->
_Last updated: 2026-07-20_

## Works now
- `/health` endpoint live, deployed on Railway
- Config loads from env: `PORT` (default 8080), `GOOGLE_PLACES_API_KEY` (required), `ENVIRONMENT` (default "development")
- `godotenv` loads `.env` locally; Railway injects env in prod
- `models.Place` struct defined (not yet used by any handler)

## In flight
- Nothing mid-implementation. Phase 1 starts fresh.

## Next
1. `internal/places/` — Google Places REST client (Text Search + Place Details)
2. `GET /places/search?q=` handler with input validation
3. `GET /places/details/{id}` handler
4. Point iOS `PlacesService` at this instead of the SDK

## Landmines
- `http.Server` has no Read/Write/Idle timeouts yet — fix when adding real handlers
- `models.Place.PriceLevel` is `string`; iOS uses `Int?` (1–4). Decide the wire format before building `/places/details`
- `go.mod` marks godotenv `// indirect` but it's a direct import — run `go mod tidy`

# STATE
<!-- Overwritten at the end of every working session. Keep under ~20 lines. -->
_Last updated: 2026-08-08_

## Works now
- **Full Phase 1 surface live locally, smoke-tested against real Google:**
  - `POST /places/autocomplete` — UUID `searchSessionID` required, restaurant-type
    filter, plain-string suggestions
  - `GET /places/{id}?tier=save|enrich|coordinate` — tier→field-mask server-side;
    token forwarded only on `save`; enrich returns populated Gemini summaries
  - `GET /places/{id}/photos?max=3` — CDN URIs via `skipHttpRedirect`, attributions included
- `models.Place` reshaped: JSON tags, `PriceLevel *int` (unknown→nil), pointer `Coordinates`
- Stdlib mux with Go 1.22+ method+wildcard patterns — still no chi
- Upstream failures: real error logged (slog), generic 502 to client

## In flight
- Nothing — Phase 1 server surface is code-complete, untested (no _test.go yet)

## Next
1. Tests: httptest fake Google against `testdata/` fixtures (see TODO for cases)
2. Point iOS `PlacesService` at savor-api; deploy current main to Railway first
3. Recapture empty `details_enrich.json` fixture

## Landmines
- `details_enrich.json` is 0 bytes — enrich verified live only, no test fixture
- Photos: sequential media calls could pathologically exceed 10s `WriteTimeout` (TODO)
- Gemini `disclosureText` currently dropped in translation — check terms before detail UI
- Attributions must render in iOS `PhotoCarouselView` — Google terms require display

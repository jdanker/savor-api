# TODO
<!-- Actionable backlog. Prune ruthlessly; done items get deleted, not checked off
     forever. STATE.md is for awareness; this file is for work. -->

## Phase 1 — SDK parity
<!-- Parity target is what the iOS app actually does today: autocomplete, tiered
     details, photos. NOT Text Search — that's a future Discover tab, not v1. -->
- [ ] `internal/places/client.go` — `Client` struct, constructor, shared
      authenticated-request/decode helper. Reused `http.Client` with a timeout
- [ ] `internal/places/autocomplete.go` — Autocomplete (New); Google response types
      stay unexported here
- [ ] `internal/places/details.go` — Place Details; `tier` → field mask; price-level
      string → nullable int (unknown values → nil, never error)
- [ ] `internal/places/photos.go` — photo media with `skipHttpRedirect`; return URIs
- [ ] `POST /places/autocomplete` — `searchSessionID` **required**, UUID-validated
      before calling Google; structured-log the token for future ratio analysis
- [ ] `GET /places/{id}?tier=save|enrich|coordinate` — reject unknown tiers;
      forward `searchSessionID` only on `tier=save`
- [ ] `GET /places/{id}/photos?max=3`
- [ ] Point iOS `PlacesService` at these endpoints (blocked on the above)

## Model gaps
- [ ] `models.Place.PriceLevel` — change `string` → nullable int when `details.go` lands
- [ ] Photo attributions — no field exists for them; Google's terms require display.
      Needs a home in the photos response shape (and eventually `PhotoCarouselView`)

## Chores
- [ ] Add Read/Write/Idle timeouts to `http.Server` before shipping real handlers
- [ ] `go mod tidy` — godotenv is a direct dep, marked `// indirect`
# TODO
<!-- Actionable backlog. Prune ruthlessly; done items get deleted, not checked off
     forever. STATE.md is for awareness; this file is for work. -->

## Design questions (answer before writing client code)
- [ ] **Session tokens on a stateless server.** The iOS SDK held the autocomplete
      session token on-device and retired it after a details fetch. Go is stateless —
      decide who mints the token and how it travels across keystrokes → details.
      → decisions.md once chosen
- [ ] **Photos: proxy bytes vs. return URIs.** Google's photo-media endpoint can
      return a short-lived CDN URI instead of the image. Proxying costs Railway
      bandwidth + latency; URIs keep the key server-side either way.
      → decisions.md once chosen
- [ ] **`Place.PriceLevel` wire format.** Google returns `"PRICE_LEVEL_MODERATE"`,
      Go model has `string`, iOS wants `Int?` (1–4). Decide where the mapping lives.
      → decisions.md once chosen

## Phase 1 — SDK parity
<!-- Parity target is what the iOS app actually does today: autocomplete, tiered
     details, photos. NOT Text Search — that's a future Discover tab, not v1. -->
- [ ] `internal/places/` — Google Places (New) REST client: Autocomplete +
      Place Details + Photo media. Owns auth, field masks, Google→`models.Place` mapping
- [ ] `POST /places/autocomplete` — returns `{placeID, primaryText, fullText}`;
      validate query length and session token before calling Google
- [ ] `GET /places/{id}?tier=save|enrich|coordinate` — intent-named tiers so field
      masks (and billing) stay server-side; iOS never sees a field mask
- [ ] `GET /places/{id}/photos?max=3`
- [ ] Point iOS `PlacesService` at these endpoints (blocked on the above)

## Chores
- [ ] Add Read/Write/Idle timeouts to `http.Server` before shipping real handlers
- [ ] `go mod tidy` — godotenv is a direct dep, marked `// indirect`
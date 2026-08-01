# TODO
<!-- Actionable backlog. Prune ruthlessly; done items get deleted, not checked off
     forever. STATE.md is for awareness; this file is for work. -->

## Phase 1 — SDK parity
<!-- Parity target is what the iOS app actually does today: autocomplete, tiered
     details, photos. NOT Text Search — that's a future Discover tab, not v1. -->
- [ ] `internal/places/client.go` — `Client` struct, constructor, shared
      authenticated-request/decode helper. Base URL is a **struct field, not a
      package const** — that's what lets `httptest` swap in a fake Google. Own
      `http.Client` with a timeout below the server's 10s `WriteTimeout`
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
- [ ] Photo attributions — `authorAttributions` (displayName, uri, photoUri) comes back
      per-photo in the details response; needs to survive into the photos response shape
      and eventually render in `PhotoCarouselView`. Google's terms require display.
- [ ] Gemini disclosure text — `generativeSummary` and `reviewSummary` both carry
      `disclosureText: "Summarized with Gemini"`. Check whether Google's terms require
      displaying it before building the detail UI.

## Deferred (not Phase 1)
- [ ] **Autocomplete match highlighting.** Dropped for Phase 1 — see decisions.md.
      Google returns `matches: [{startOffset, endOffset}]` per text field, in **Unicode
      code points**. Go maps cleanly (`[]rune`); Swift does *not* — `Character` is a
      grapheme cluster, so reconstruction must index `String.unicodeScalars`, not
      `Character`s. Can't be faked with a client-side substring search: Google's matches
      come from spell correction and transliteration too, so "papy" legitimately
      highlights "Pappy". Restoring it = wire-format change (breaks iOS decoding) +
      `PlaceSuggestion` back to `AttributedString` + scalar-offset reconstruction.

## Chores
- [ ] `gofmt` `cmd/server/main.go` — trailing whitespace in the `http.Server` literal
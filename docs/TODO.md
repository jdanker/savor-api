# TODO
<!-- Actionable backlog. Prune ruthlessly; done items get deleted, not checked off
     forever. STATE.md is for awareness; this file is for work. -->

## Phase 1 — SDK parity
- [ ] Point iOS `PlacesService` at the three savor-api endpoints (all live now);
      `PlaceSuggestion` drops `AttributedString` per decisions.md
- [ ] Tests: httptest fake Google using `testdata/` fixtures — the whole reason
      `baseUrl` is a struct field. Cover: autocomplete translation, tier→field-mask,
      price-level → nullable int (incl. unknown → nil), non-200 upstream → error
- [ ] Recapture `testdata/details_enrich.json` — file is empty (0 bytes); enrich
      tier verified live but has no fixture for tests

## Model gaps
- [ ] Photo attributions — now returned by `GET /places/{id}/photos`; still needs
      to render in iOS `PhotoCarouselView`. Google's terms require display.
- [ ] Gemini disclosure text — `generativeSummary` and `reviewSummary` both carry
      `disclosureText: "Summarized with Gemini"`. Check whether Google's terms require
      displaying it before building the detail UI. (Currently dropped in translation.)

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
- [ ] Photos worst case: 1 details + N sequential media calls, each on a 5s client
      timeout — pathological slowness could exceed the server's 10s `WriteTimeout`.
      Parallelize media calls or accept truncation if it ever shows up in logs.

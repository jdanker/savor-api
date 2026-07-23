# Decisions
<!-- Append-only. One dated paragraph per non-obvious tradeoff. Newest at top. -->

**2026-07 — stdlib `net/http`, no framework (yet).** Two GET endpoints don't justify
chi/gin. Go 1.22+ mux supports method+path patterns (`GET /places/details/{id}`),
which covers Phase 1. Revisit if middleware chains (auth, rate limiting in Phase 3)
get awkward.

**2026-07 — REST API directly, no Google Go SDK.** Learning goal: understand the
actual HTTP contract, field masks, and billing SKUs rather than hiding them behind
a client library. Also keeps the `internal/places` boundary honest — mapping raw
JSON forces an explicit translation layer.

**2026-07 — `godotenv` for local dev only.** `godotenv.Load()` silently no-ops when
`.env` is absent, so prod (Railway env injection) and local use the same code path.

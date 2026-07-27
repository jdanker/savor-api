# Decisions
<!-- Append-only. One dated paragraph per non-obvious tradeoff. Newest at top. -->

**2026-07 — Autocomplete session tokens: client mints, server enforces.** Google
bills autocomplete per *session* — all keystroke requests plus the closing Details
call bundle into one charge if they share a token. On a stateless Go server there's
nowhere to hold that token: no auth until Phase 3 (so no user key), IP is wrong
(carrier NAT, roaming), in-memory dies on redeploy and breaks on a second instance.
So iOS mints a UUID per search interaction and passes it through. Named
`searchSessionID` in our contract, not "session token" — a client-generated ID
grouping one search interaction is a legitimate API concept; that it maps 1:1 onto
Google's token stays an implementation detail inside `internal/places`. The client
owns the *UUID*, the server owns the *rules*: only `tier=save` forwards it to
Google; `enrich` and `coordinate` never do, mirroring what `refreshRestaurant` and
`fetchCoordinate` already do on-device. Mitigations available now: reject malformed
UUIDs before calling Google, make the field **required** (an optional field that iOS
silently drops = per-keystroke billing forever, undetected), and structured-log the
token on every request so the token:request ratio is already in Railway logs the day
a real gate is worth building. Not solvable statelessly: a client minting a fresh
valid token per keystroke. Accepted — the endpoint has no auth anyway, so rate
limiting (Phase 3) is the real abuse control.

**2026-07 — Photos return signed URIs, not proxied bytes.** `GET /places/{id}/photos`
responds with JSON (`uri`, dimensions, attributions), not image data. `internal/places`
calls Google's photo-media endpoint with `skipHttpRedirect`, which returns a
short-lived pre-signed CDN URI needing no API key; iOS fetches images straight from
Google's CDN. Rejected proxying the bytes: photos are the largest payloads in the app
and would cross Railway twice, making this the one endpoint that could dominate
hosting cost, for strictly worse latency than a CDN built for it. The usual objection
to URIs — they expire — doesn't apply here because the iOS cache stores decoded
`UIImage` bytes in both tiers (NSCache + disk JPEGs), never URLs. A URI that dies in
an hour is fine when nothing holds it for more than a second. Cost is two round trips
on a cold cache, which per place happens roughly once. Note: billing is identical
either way — both shapes call the media endpoint once per photo. Reversible:
proxying later is a change to what the handler returns.

**2026-07 — `internal/places` is a translation layer, not a passthrough.** Every
field arriving in Google's vocabulary leaves in ours. Forcing case was price level:
Google returns `"PRICE_LEVEL_MODERATE"`, iOS wants `Int?` (1–4). Mapping lives
server-side in `internal/places`, so `models.Place.PriceLevel` becomes a nullable int
and iOS deletes its conversion switch entirely. Rejected passing the string through —
it would make iOS re-own vocabulary it just shed, and a string that originated at
Google is a Google type in everything but declaration. Unknown enum values map to nil
rather than erroring, so a fifth Google tier degrades quietly instead of 500ing. This
generalizes: the next awkward Google representation is already decided the same way.

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

# http.Server timeouts

**What**: HTTP is connection-oriented, and a connection is a resource the *client*
controls the lifetime of unless the server says otherwise. A client can open a
socket, send half a request, and stop — the server keeps the connection, its
goroutine, and its buffers alive indefinitely waiting for the rest. Nothing about
this looks like an attack; it's indistinguishable from a phone on a bad tunnel
connection. Timeouts are how the server reclaims the decision: *I will wait this
long and no longer.* Go's `http.Server` zero value has **none of them set**, so
the out-of-the-box behavior is "wait forever," and that is the specific reason
this is a well-known Go footgun rather than a detail the runtime handles for you.

The four knobs, in the order a request encounters them:

| Field | Covers |
|---|---|
| `ReadHeaderTimeout` | accept → request headers fully read |
| `ReadTimeout` | accept → request **body** fully read (subsumes the above) |
| `WriteTimeout` | headers read → response fully written (**includes handler execution**) |
| `IdleTimeout` | how long an unused keep-alive connection survives between requests |

**Why here**: Two reasons, one general and one specific to this service.

The general one: savor-api is a public endpoint on Railway with no authentication
until Phase 3. Anyone with the URL can open connections to it. Without timeouts,
exhausting the process doesn't require a clever exploit — it requires patience and
a few hundred half-open sockets. Timeouts are the cheapest possible mitigation and
they cost nothing when everything is behaving.

The specific one: **this service's handlers make outbound calls to Google.** That
makes `WriteTimeout` a coupling point rather than an isolated setting, because it
spans the whole handler — including however long Google takes to answer. So the
`http.Client` timeout inside `internal/places` and the server's `WriteTimeout` are
one decision made in two files. If the client timeout is the larger of the two, the
server kills the connection first and the handler never reaches the code that would
have returned a clean `502`. The failure mode isn't a slow response; it's a *lost*
one, and it looks like a client-side network error rather than a server bug. Current
sizing — read 5s, write 10s, idle 60s — leaves room for a client timeout around 5s.

Sizing philosophy: these are **backstops, not tuning knobs**. Places responses land
in a few hundred milliseconds, so 10s is ~20x headroom. A timeout tight enough to
fire during normal operation is a bug generator; one loose enough to never fire
during normal operation but still bound the worst case is doing its job.

**Where**: `cmd/server/main.go:26-32` (the `http.Server` literal).

**Gotchas**:
- **Server-wide, not per-route.** The moment one endpoint has a genuinely different
  profile — a long upload, a streaming response — `WriteTimeout` stops being
  expressive enough. The escape hatch is `http.TimeoutHandler` or per-request
  context deadlines, not loosening the global value and hoping.
- `WriteTimeout` starts at *header read*, not handler entry, so a slow request body
  eats into the write budget.

**Deeper**: https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/

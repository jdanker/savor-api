# CLAUDE.md

Guidance for Claude Code when working in this repository.

## What this service is

`savor-api` is a Go HTTP API backing **Savor**, an iOS restaurant discovery app.
It is a places intelligence layer between the iOS app and third-party place APIs
(Google Places now; Yelp/Apple Maps later). The iOS app never calls external place
APIs directly — all place search and enrichment flows through this service, keeping
API keys off the device.

**Structure, boundaries, and current state live in `docs/` — read before non-trivial work:**
- `docs/ARCHITECTURE.md` — system diagrams, invariants, Code Map (every file, one line)
- `docs/STATE.md` — where the project stands right now
- `docs/TODO.md` — actionable backlog; pull next steps from here
- `docs/decisions.md` — why non-obvious choices were made
- `docs/concepts/` — one file per concept new to this project

## Commands

```bash
go run ./cmd/server          # start locally (reads .env via godotenv)
go build ./...               # build all packages
go test ./...                # run all tests
go vet ./...                 # static checks
```

Requires a local `.env` (gitignored) with `GOOGLE_PLACES_API_KEY` set —
`config.Load()` fails fast if missing. `PORT` (default `8080`) and `ENVIRONMENT`
(default `development`) are optional.

## Stack

Go stdlib `net/http` (chi only if middleware needs justify it), Google Places REST
API (no SDK — deliberate, see decisions.md), env vars for all secrets, deployed on
Railway (TLS at platform).

## Collaboration ground rules

The user is a platform/DevOps engineer learning Go web services by writing the code
himself. Guide and review; don't dump full implementations unless asked. Explain
trade-offs before implementing. Flag anti-patterns instead of quietly fixing them.

## Living docs maintenance

`docs/` is the project's memory. Rules:

- `STATE.md` — **overwrite** (never append) at session end: works now / in flight /
  next / landmines. Under ~20 lines.
- `ARCHITECTURE.md` Code Map — every new file gets a one-line entry **in the same
  change that creates it**. Structural changes update the diagrams.
- `concepts/` — create from `_template.md` when a concept/framework/technique first
  appears in the project. Include `file:line` pointers.
- `decisions.md` — append-only, one dated paragraph per non-obvious tradeoff.
  A trade-off explained during a session should land here, not evaporate in chat.
- `TODO.md` — actionable backlog. Bugs/chores found but not fixed in-session go
  here (flagged anti-patterns included). Delete done items; prune ruthlessly.

**Update triggers**: new dependency/framework/API; first use of a pattern; module
boundary or data-flow change; non-obvious tradeoff decided. **Not**: bugfixes,
styling, refactors within existing patterns. Test: "would future-me need this
re-explained in 3 weeks?"

**Cadence**: never interleave doc edits mid-implementation. At session end, propose
all doc updates as one reviewable batch, then rewrite STATE.md.

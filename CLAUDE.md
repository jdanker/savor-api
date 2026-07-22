# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
go run ./cmd/server          # start the server locally (reads .env via godotenv)
go build ./...                # build all packages
go test ./...                 # run all tests (no test files exist yet)
go vet ./...                  # static checks
```

Requires a local `.env` (gitignored) with at least `GOOGLE_PLACES_API_KEY` set — `config.Load()`
fails fast (returns an error, and `main.go` calls `log.Fatalf`) if it's missing. `PORT` (default
`8080`) and `ENVIRONMENT` (default `development`) are optional.

## What this service is

`savor-api` is a Go HTTP API that acts as the backend for **Savor**, an iOS restaurant discovery and wishlist app built with SwiftUI. It is a **places intelligence layer** between the Savor iOS app and third-party place APIs (Google Places, and eventually Yelp, Apple Maps, etc.). The iOS app never calls external place APIs directly — all place search and enrichment flows through this service, keeping API keys off the device.

## Current phase

**Phase 1** — scaffolding a basic Go HTTP server and proxying Google Places API search and detail endpoints.

## Stack

- **Language:** Go, stdlib `net/http` (chi router may be added soon)
- **External APIs:** Google Places API (HTTP REST, no Go SDK)
- **Config:** Environment variables for all secrets (`.env` locally, platform env vars in production)
- **Deployment:** Railway

## Project structure (current, not aspirational)

```
savor-api/
├── cmd/server/main.go         ← entry point: loads .env + config, registers routes, starts http.Server
├── internal/config/config.go  ← env var loading (Config struct)
├── internal/handlers/health.go← GET /health — only route that exists so far
├── internal/models/place.go   ← Place struct (Google-shaped fields; not yet consumed anywhere)
├── .env                       ← local secrets, gitignored
├── go.mod
└── go.sum
```

Note: `internal/places/` (a Google Places client) doesn't exist yet — `models.Place` is defined
but nothing populates it yet. Routing uses stdlib `net/http.ServeMux` directly — chi hasn't been
added despite being flagged as a likely near-term addition (see Stack below).

## Principles

- **No API key leakage** — Google Places API keys never touch the iOS device
- **Internal models only** — Google's response types never leak outside `internal/places`; always map to our own `Place` struct in `internal/models`
- **Input validation** — validate all query params before hitting external APIs
- **Security-first** — no hardcoded secrets, HTTPS enforced in production
- **Idiomatic Go** — errors as values, no magic, keep it simple, strict typing

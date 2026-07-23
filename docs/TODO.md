# TODO
<!-- Actionable backlog. Prune ruthlessly; done items get deleted, not checked off
     forever. STATE.md is for awareness; this file is for work. -->

## Phase 1
- [ ] `internal/places/` Google Places REST client (Text Search + Place Details)
- [ ] `GET /places/search?q=` with input validation
- [ ] `GET /places/details/{id}`
- [ ] Decide `Place.PriceLevel` wire format (Go has `string`, iOS wants `Int?` 1–4)
      → record in docs/decisions.md once chosen

## Chores
- [ ] Add Read/Write/Idle timeouts to `http.Server` before shipping real handlers
- [ ] `go mod tidy` — godotenv is a direct dep, marked `// indirect`
- [ ] Un-ignore `CLAUDE.md` / `.claude/` in .gitignore (keep ignoring
      `.claude/settings.local.json`)

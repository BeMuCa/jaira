# Next tasks — 2026-08-31 (fourth list)

**The board is the list now.** `jaira next --per-lane --json` says what has
work; the backlog holds everything the previous lists carried. Read
`HANDOFF.md` first. What remains here are the ground rules:

- **Every item goes through a GSD entry point** (`/gsd:quick` or the board's
  own lanes). Atomic commits, a SUMMARY, a STATE row.
- **Verify on the running binary**: `go build -o ~/.local/bin/jaira ./cmd/jaira`
  after every change. Never `go install`.
- **Gate on `gofmt -l core internal` listing exactly `internal/cli/tickets.go`.**
- **His real board is a different repo**; never probe it with live moves —
  `move --dry-run` exists.
- `go test ./... -race` with the cache cleared, before claiming anything.

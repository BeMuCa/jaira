# Quick Task 260812-wz3: Fix TUI board creation to write .gitignore + agent note — Summary

One-liner: moved board-privacy and agent-note logic into a stdlib-only `core/board` package with a shared `Prepare` entry point, wired the TUI's 'i' key to it, and widened the agent note from 4 to 9 commands covering the whole working loop.

## What changed

The TUI browse screen's 'i' key created a board with `st.Init()` alone, so boards created there were never gitignored and never announced in AGENTS.md/CLAUDE.md — unlike `jaira init`. `internal/tui` could not import the existing helpers because they lived in `internal/cli`, which already imports `internal/tui` (a cycle). Moved the helpers into a new `core/board` package (stdlib-only) with one entry point, `board.Prepare(root)`, and called it from both `jaira init` and the TUI's 'i' key handler.

### Task 1 — `core/board` package + rewire `internal/cli`
- Created `core/board/gitignore.go` (`IgnoreLine`, `Ignored`, `AddIgnore`, `RemoveIgnore` — moved verbatim from `internal/cli/share.go`, exported)
- Created `core/board/announce.go` (`jairaMarkerStart/End`, `agentNote`, `agentFiles`, `AnnounceInAgentFiles`, `announceInAgentFile` — moved verbatim from `internal/cli/announce.go`)
- Created `core/board/board.go` with `Prepared` struct and `Prepare(root) Prepared`, which calls `AddIgnore` then `AnnounceInAgentFiles` and records both results without stopping at the first error
- Created `core/board/board_test.go`: first coverage for this behavior (Prepare on empty dir → ignored + both files announced; Prepare again → no-op, no duplicate gitignore line)
- Rewrote `internal/cli/share.go` and `internal/cli/tickets.go` to call `board.*`; `internal/cli/announce.go` now holds only `announceLine`
- Verified `jaira init --json` output keys are byte-identical to before (`root`, `tickets_dir`, `created`, `private`, `gitignore_written`, `state_dir`, `agent_notes`)

### Task 2 — Call `board.Prepare` from the TUI 'i' key path
- `internal/tui/browse.go`: after `st.Init()` succeeds, calls `board.Prepare(st.Root)`. A `Prepare` failure does not undo the board creation (the board exists either way) but is surfaced as a one-line warning in `b.msg` naming the failure (privacy vs. agent note), rather than being swallowed
- Added `internal/tui/browse_test.go` (new file — none existed before): drives `b.key("i")` directly with `JAIRA_HOME` redirected to a temp dir, asserts the board, `.gitignore` entry and `CLAUDE.md` marker all exist

### Task 3 — Widen the agent note
- Replaced the `agentNote` body in `core/board/announce.go` verbatim per the plan: now documents `create` (with `--goal/--context/--dod` flags), `list --actionable`, `next`, `claim`, `show --for-lane`, `dod`, `note`, `move`, and `resume` — the whole working loop, not just four commands
- Rewrote the constant as concatenated interpreted string literals (`"...\n" + "...\n" + ...`) instead of the previous raw-string-plus-backtick-splicing trick, since backticks needed no escaping this way
- Also reformatted `core/board/board.go`'s struct field alignment per `gofmt` (whitespace-only; the file was written slightly out of `gofmt` canonical form in Task 1)

## Deviations from Plan

**1. [gofmt cleanup, not a deviation rule] `core/board/board.go` struct alignment**
- Found during: Task 3, running `gofmt -l` across touched files before committing
- Issue: `Prepared` struct fields had extra alignment spaces gofmt doesn't produce
- Fix: `gofmt -w core/board/board.go`
- Files modified: `core/board/board.go`
- Commit: `0dad0dc` (bundled with Task 3, since board.go was untouched by task 3's own file list otherwise)

No other deviations. All three tasks executed as written; every verify block in the plan passed on first attempt except this formatting nit.

## Commits

| Task | Commit | Message |
|------|--------|---------|
| 1 | `d72d075` | refactor: move board privacy and agent-note logic into core/board |
| 2 | `988c083` | fix: gitignore and announce a board created from the TUI browse screen |
| 3 | `0dad0dc` | feat: widen the agent note to the whole working loop |

## Verification

Final `go test ./...` (run after all three tasks committed):

```
?   	github.com/BeMuCa/jaira/cmd/jaira	[no test files]
ok  	github.com/BeMuCa/jaira/core/board	(cached)
ok  	github.com/BeMuCa/jaira/core/gate	(cached)
?   	github.com/BeMuCa/jaira/core/gitrepo	[no test files]
?   	github.com/BeMuCa/jaira/core/lane	[no test files]
ok  	github.com/BeMuCa/jaira/core/merge	(cached)
ok  	github.com/BeMuCa/jaira/core/project	(cached)
?   	github.com/BeMuCa/jaira/core/session	[no test files]
ok  	github.com/BeMuCa/jaira/core/ticket	(cached)
ok  	github.com/BeMuCa/jaira/core/validate	(cached)
ok  	github.com/BeMuCa/jaira/internal/cli	(cached)
ok  	github.com/BeMuCa/jaira/internal/tui	(cached)
?   	github.com/BeMuCa/jaira/scripts/iconpreview	[no test files]
?   	github.com/BeMuCa/jaira/scripts/shotgen	[no test files]
```

`go build ./...` and `go vet ./...` both clean. The no-unexported-copy grep (Task 1's second verify) and the "every command in the note exists" check (Task 3's second verify) both passed.

## Known Stubs

None.

## Threat Flags

None — this plan's threat register (T-wz3-01, T-wz3-02) covers the surface touched; no new network endpoints, auth paths, or schema changes were introduced.

## Self-Check: PASSED

- `core/board/gitignore.go`, `core/board/announce.go`, `core/board/board.go`, `core/board/board_test.go` — FOUND
- `internal/tui/browse_test.go` — FOUND
- Commit `d72d075` — FOUND
- Commit `988c083` — FOUND
- Commit `0dad0dc` — FOUND

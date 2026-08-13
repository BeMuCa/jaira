---
phase: quick/260813-1br
plan: 01
subsystem: cli, core/release
tags: [versioning, update-command, board-setup]
requires: []
provides: [release.Current, release.Stamp, release.Stamped, release.Since, StateDir, jaira-update]
affects: [internal/cli/root.go, internal/cli/tickets.go, internal/tui/browse.go]
tech-stack:
  added: []
  patterns: ["go:embed for shipping data inside the binary instead of copying it to disk"]
key-files:
  created:
    - core/release/release.go
    - core/release/NOTES.md
    - core/release/release_test.go
    - internal/cli/update.go
    - internal/cli/update_test.go
  modified:
    - core/ticket/store.go
    - internal/cli/root.go
    - internal/cli/tickets.go
    - internal/tui/browse.go
    - README.md
    - docs/COMMANDS.md
decisions: []
metrics:
  duration: "~35m"
  completed: "2026-08-13"
---

# Phase quick/260813-1br Plan 01: Version stamp and jaira update Summary

A board now records which jaira version last prepared it, in the per-working-tree
state directory outside the repo; every command nudges once on stderr when the
running binary is newer, and `jaira update` re-runs the setup, prints the
embedded change notes since that stamp, and re-stamps.

## What was built

- `core/release`: `Current` (the running binary's version, set once by
  `cli.Execute`), `Notes()`/`Since()` (parse and select from an embedded
  `NOTES.md`), `Stamp()`/`Stamped()` (write/read `<stateDir>/version`).
- `core/release/NOTES.md`: the format-rules header plus a real `## 0.1.0` entry
  with three reader-facing instructions (TUI board creation now stamped/prepared
  same as init; the agent note now covers the full working loop; the detail pane
  shows/copies the full ticket id).
- `core/ticket.Store.StateDir()`: exports the existing unexported `stateDir()` so
  callers outside the package can address the per-working-tree state directory.
- `internal/cli/update.go`: `nudgeIfStale` (one stderr line, silent on `dev`
  builds, on a matching stamp, or on an unreadable state directory) and
  `newUpdateCmd` (`jaira update`, re-runs `board.Prepare`, prints notes since the
  prior stamp, stamps only after Prepare succeeds, `--json` payload with
  `version`/`previous`/`gitignore_written`/`agent_notes`/`notes`).
- `internal/cli/root.go`: `Execute` sets `release.Current`; `openStore` calls
  `nudgeIfStale` right after `bindDriverIfShared`; `update` registered.
- `internal/cli/tickets.go` (`init`) and `internal/tui/browse.go` (the `i` key):
  both stamp the board immediately after a successful `board.Prepare`, so a
  freshly created board is never nagged about itself.
- `README.md` / `docs/COMMANDS.md`: `jaira update` documented alongside the
  other commands.

## What a reader must do differently

- After upgrading the `jaira` binary, run `jaira update` in each board you use —
  every other command now prints one line on stderr when it notices the binary
  is newer than what last set the board up, but only `update` actually re-applies
  the setup and shows what changed.
- A `dev` build (an unreleased local build) never prints that nudge — this is by
  design, not a bug, since a local build's version string says nothing about
  what it contains.
- The nudge and the stamp never touch the repository: the stamp lives under
  `$JAIRA_HOME/state/<key>/version` (or `~/.jaira/state/<key>/version`), never
  inside `.jaira/`.

## Deviations from Plan

None — the plan's design decisions were implemented as written. One
clarification worth recording: `nudgeIfStale`'s "unreadable state directory"
case can't be detected through `Stamped()` alone (it deliberately swallows every
error, by design, to make "missing" and "unreadable" indistinguishable for its
own purpose). `nudgeIfStale` therefore does its own `os.ReadDir(dir)` check
first — tolerating `IsNotExist` (a normal unstamped board, still nudge) but
bailing out silently on any other error (a genuinely unreadable directory) —
before ever consulting `Stamped`. This is additive precision inside the
function the plan already specified, not a change to any exported behavior.

Concurrent, unrelated work landed on `master` mid-session (commits `c3ffedb`
"the detail pane says when a ticket appeared and when it was last touched" and
`b2058b7` "docs: say plainly how a context should read", the latter touching
`core/board/announce.go` and `core/lane/builtin/15-pre-process.md` per the
plan's own warning about parallel changes). Neither touched any file this plan
modified; both commits already sit under this plan's two commits in the log,
and the full test suite passed after both.

## Known Stubs

None.

## Threat Flags

None — the update path calls the same `board.Prepare` used by `init`, reads and
writes only a plaintext version string under the existing per-working-tree state
directory, and introduces no new network, auth, or trust-boundary surface.

## Commits

- `0872e5c` feat: embed release notes and a version stamp for boards
- `3acedb0` feat: add jaira update and a per-command staleness nudge

## Verification

Smoke test (plan's second `<verify>` block) — real binaries, real stderr/stdout
separation:

```
== after init, expect NO nudge ==
== newer binary, expect ONE nudge line ==
jaira: this board was set up by an older version of jaira — run 'jaira update' to bring it up to date
== update, expect notes on stdout ==

Nothing has changed since the version that last set this board up.
== update --json ==
{
  "agent_notes": null,
  "gitignore_written": false,
  "notes": [
    {
      "changes": [
        "A board created from the TUI browse screen now gets the same setup `jaira init` gives it — `.jaira/` gitignored and the jaira section written into `AGENTS.md` and `CLAUDE.md`. Boards created that way before this release have neither; run `jaira update` in each of them.",
        "The jaira section in `AGENTS.md` and `CLAUDE.md` now names the full working loop rather than a handful of commands. Re-read it — it is the contract for how to drive the board.",
        "The ticket detail pane shows the full ticket id and `y` copies it. Use the full id when handing a ticket to another tool; a prefix can turn ambiguous as the board grows."
      ],
      "version": "0.1.0"
    }
  ],
  "previous": "9.9.9",
  "root": "/tmp/tmp.QSqufLjybY",
  "version": "9.9.9"
}
== after update, expect NO nudge ==
== stamp is outside the repo ==
9.9.9
STAMP OUTSIDE REPO OK
```

(The first `update` shows "nothing changed" because the v1 binary's own stamp,
`0.1.0`, is the only — and therefore topmost — entry in `NOTES.md`; the second
`update --json` call runs after the first already re-stamped the board to
`9.9.9`, a version that matches no heading, so it correctly falls back to
showing every entry. Both are the selection rule from the plan's design
decisions, not a bug.)

The constraint-mandated proofs, each run for real:

1. Nudge never writes to stdout — proved by
   `TestNudgeIfStaleOnDifferingStampPrintsOneStderrLineOnly` (asserts
   `stdout == ""` while `stderr` carries exactly one line) and by the smoke run
   above, where every nudging command redirects stdout to `/dev/null` and the
   line still appears.
2. Nudge never fires for `dev` — proved by `TestNudgeIfStaleNeverFiresOnDevBuild`.
3. Nudge never fails a command when the state directory is unreadable — proved
   by `TestNudgeIfStaleSkipsSilentlyWhenStateDirUnreadable` (chmod 000 on the
   real state directory, asserts no output and no error).

```
$ go test ./internal/cli/... -run 'Update|Nudge' -v
=== RUN   TestNudgeIfStaleNeverFiresOnDevBuild
--- PASS: TestNudgeIfStaleNeverFiresOnDevBuild (0.00s)
=== RUN   TestNudgeIfStaleSilentOnMatchingStamp
--- PASS: TestNudgeIfStaleSilentOnMatchingStamp (0.00s)
=== RUN   TestNudgeIfStaleOnDifferingStampPrintsOneStderrLineOnly
--- PASS: TestNudgeIfStaleOnDifferingStampPrintsOneStderrLineOnly (0.00s)
=== RUN   TestNudgeIfStaleSkipsSilentlyWhenStateDirUnreadable
--- PASS: TestNudgeIfStaleSkipsSilentlyWhenStateDirUnreadable (0.00s)
=== RUN   TestUpdateStampsSoASecondNudgeCheckIsSilent
--- PASS: TestUpdateStampsSoASecondNudgeCheckIsSilent (0.00s)
=== RUN   TestUpdateJSONCarriesTheDocumentedFieldsOnStdoutOnly
--- PASS: TestUpdateJSONCarriesTheDocumentedFieldsOnStdoutOnly (0.00s)
PASS
ok  	github.com/BeMuCa/jaira/internal/cli	0.016s
```

Final full suite, verbatim:

```
$ export PATH=$PATH:$HOME/.local/go/bin && go build ./... && go vet ./... && go test ./... -count=1
?   	github.com/BeMuCa/jaira/cmd/jaira	[no test files]
ok  	github.com/BeMuCa/jaira/core/board	0.015s
ok  	github.com/BeMuCa/jaira/core/gate	0.079s
?   	github.com/BeMuCa/jaira/core/gitrepo	[no test files]
?   	github.com/BeMuCa/jaira/core/lane	[no test files]
ok  	github.com/BeMuCa/jaira/core/merge	0.085s
ok  	github.com/BeMuCa/jaira/core/project	0.015s
ok  	github.com/BeMuCa/jaira/core/release	0.007s
?   	github.com/BeMuCa/jaira/core/session	[no test files]
ok  	github.com/BeMuCa/jaira/core/ticket	0.036s
ok  	github.com/BeMuCa/jaira/core/validate	0.093s
ok  	github.com/BeMuCa/jaira/internal/cli	0.109s
ok  	github.com/BeMuCa/jaira/internal/tui	1.615s
?   	github.com/BeMuCa/jaira/scripts/iconpreview	[no test files]
?   	github.com/BeMuCa/jaira/scripts/shotgen	[no test files]
```

`go build ./...` and `go vet ./...` produced no output (clean).

## Self-Check

- `core/release/release.go` — FOUND
- `core/release/NOTES.md` — FOUND
- `core/release/release_test.go` — FOUND
- `internal/cli/update.go` — FOUND
- `internal/cli/update_test.go` — FOUND
- Commit `0872e5c` — FOUND
- Commit `3acedb0` — FOUND

## Self-Check: PASSED

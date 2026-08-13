---
phase: quick/260813-px2
plan: 01
subsystem: cli, tui, core/lane
tags: [lanes, cli, tui, agent-plumbing]
requires: []
provides: [lanes-use-publish-adopt-default-cli, lane-settings-new-key]
affects: [internal/cli/lanes.go, internal/tui/lanes.go, core/lane/template.go]
tech-stack:
  added: []
  patterns: ["CLI subcommands wrap existing core/lane functions rather than reimplementing them"]
key-files:
  created:
    - internal/cli/lanes.go
    - core/lane/template.go
  modified:
    - internal/cli/lanes_test.go
    - internal/cli/tickets.go
    - internal/tui/lanes.go
    - internal/tui/lanes_test.go
    - internal/tui/model.go
    - docs/AGENTS.md
    - docs/COMMANDS.md
decisions:
  - "lanes adopt takes a path (the PATH column jaira lanes shared prints), not a lane id, since two teammates can publish under the same id"
  - "moved the lane skeleton (laneTemplate) from internal/cli into core/lane as lane.Template so the CLI's 'lanes template' and the TUI's new 'n' key share one skeleton"
  - "laneScreen.key(s) keeps its existing (done bool) signature; a new pendingCmd field carries the one tea.Cmd 'n' needs out to model.go, so the ~15 existing test call sites using `done := ls.key(s)` did not need to change"
metrics:
  duration: ~90m
  completed: 2026-08-13
---

# Phase quick/260813-px2 Plan 01: Give Claude the lane actions the TUI already has Summary

Added `jaira lanes use`, `publish`, `adopt` and `default` as CLI subcommands that
wrap the same `core/lane` functions (`Export`, `Publish`, `Adopt`,
`LoadDefaultBoard`/`SaveDefaultBoard`) the TUI's settings screen already calls,
and gave the lane settings screen an `n` key that writes a lane skeleton and
opens it in `$EDITOR`.

## What was built

**Task 1 — `use`, `publish`, `adopt` as commands** (`internal/cli/lanes.go`,
commit `f9e0ede`):
- `jaira lanes use <id>` copies a catalogue lane into `<root>/.jaira/lanes`
  via `lane.Export`; `jaira lanes publish <id>` copies it into
  `.jaira/shared/<slug>/` via `lane.Publish`, stamping `creator:` with the
  acting identity's slug (matching the TUI's `p` key, which passes the slug
  as both the directory name and the stamp — not the raw name).
- `jaira lanes adopt <path>` copies a teammate's shared lane into the
  catalogue via `lane.Adopt`. It takes the `PATH` column `jaira lanes shared`
  prints, not an id — documented in the command's `Long` text, since a bare
  id can collide across teammates and only the path names one file
  unambiguously.
- All three take `--force`; without it, an existing file is refused with the
  core function's own "already exists; refusing to overwrite" message,
  reported at `ExitValidation` (3) via a shared `writeConflictError` helper.
  `--json` on each carries the id and the written path.

**Task 2 — the default board from the CLI** (`internal/cli/lanes.go`, commit
`3e3b5b4`):
- `jaira lanes default` with no flags prints the effective lane list (the
  built-ins, named explicitly, when the file is absent or empty) and the
  pre-ticked options, plus the file's path.
- `--lanes a,b,c` and `--options x,y` validate every id/name against the
  currently installed lanes (`lane.Load`) before writing, refusing an unknown
  one with a usage error naming what is actually installed.
- `--clear` removes the file via `os.Remove`, tolerating "already absent" as
  success (idempotent, matching the project's concurrency-safe-retry
  convention).

**Task 3 — `n` writes a skeleton and opens it** (`internal/tui/lanes.go`,
`internal/tui/model.go`, `core/lane/template.go`, commit `bbe3ae5`):
- Moved the lane skeleton constant from `internal/cli/tickets.go` (private
  `laneTemplate`) to `core/lane` as exported `lane.Template`, since both
  `jaira lanes template` and the new TUI key needed to print the identical
  skeleton — one template to keep in sync with `parse()`, not two.
- `n` (only when the settings screen's own lane list has focus) writes
  `lane.Template` into the catalogue under the first free name —
  `new-lane.md`, then `new-lane-2.md`, `new-lane-3.md`, … — using
  `os.O_EXCL` so the file creation itself is the collision check, then opens
  it via `tea.ExecProcess` and the existing `editorCommand()` helper from
  `external.go`.
- On the editor exiting without error, the model does a full `m.reload()`
  and rebuilds the lane screen (`newLaneScreen`), so the new lane appears —
  or a duplicate-id warning does, if the skeleton was left unedited and its
  `id: my-lane` collides with an earlier untouched skeleton.
- The footer now reads `... p publish · n new · R refresh ...` (and the
  shared-list variant), so the key is discoverable.

**Task 4 — docs say the loop is complete** (`internal/cli/tickets.go`,
`docs/AGENTS.md`, `docs/COMMANDS.md`, commit `2c5941b`):
- `jaira lanes`'s `Long` text gained a 5th step, `'jaira lanes use <id>' to
  put it to work in this project`, so "no TUI required" is now true
  end-to-end rather than stopping short of the step that used to be
  TUI-only.
- `docs/COMMANDS.md`'s Writing table gained rows for `use`, `publish`,
  `adopt` and `default`.
- `docs/AGENTS.md` gained a short "Building and sharing a lane without the
  TUI" section naming the five commands in the order an agent would run
  them.

## Deviations from Plan

**1. [Rule 3 - blocking issue] `internal/cli/lanes.go` did not exist; the
`lanes` parent command and its existing subcommands live in
`internal/cli/tickets.go`.** The plan's `files_modified` frontmatter listed
`internal/cli/lanes.go` as an existing file to edit. It does not exist —
`newLanesCmd`, `newLanesShowCmd`, `newLanesPathCmd`, `newLanesTemplateCmd`
and `newLanesSharedCmd` are all defined in `tickets.go`. Created
`internal/cli/lanes.go` as a new file holding the four new subcommand
constructors (`newLanesUseCmd`, `newLanesPublishCmd`, `newLanesAdoptCmd`,
`newLanesDefaultCmd`) and edited `tickets.go` only to register them via
`cmd.AddCommand(...)` and to extend the `Long` help text (Task 4). This
keeps the new commands in a sensibly-named file without churning the
existing `lanes` command's home.

**2. [Rule 3 - blocking issue] The lane skeleton constant was private to
`internal/cli`, but the TUI's new `n` key needed the identical skeleton.**
Moved `laneTemplate` out of `internal/cli/tickets.go` into a new
`core/lane/template.go` as exported `lane.Template`, and updated
`tickets.go`'s `lanes template` command to reference it. This is the
single-source-of-truth fix the plan's own framing implies ("the same
skeleton `jaira lanes template` prints") rather than duplicating the
27-line constant into the `tui` package.

**3. [Rule 3 - blocking issue] `laneScreen.key(s)` had a `(done bool)`-only
signature; `n` needed to hand a `tea.Cmd` (the `$EDITOR` launch) back to
`model.go`, but changing the signature would have broken roughly 15 existing
test call sites (`done := ls.key(s)`) across `internal/tui/lanes_test.go`.**
Added a `pendingCmd tea.Cmd` field to `laneScreen`, reset at the top of
`key()` and set only by the `n` case; `model.go`'s `modeLanes` branch reads
`m.laneScreen.pendingCmd` after calling `key()`. This mirrors the existing
`defaultBoardScreen.key(s) (done bool, cmd tea.Cmd)` pattern in spirit
without touching every existing `laneScreen` test.

**4. [Rule 1 - bug caught during implementation] `publish`'s creator stamp
used the raw identity name, not the slug.** My first implementation of
`newLanesPublishCmd` called `who := identity()` (the raw name, e.g. "Alex
Doe") and only slugged it for the shared directory name, stamping the raw
name as `creator:`. The TUI's `publish()` (`internal/tui/lanes.go`) passes
the *slug* as `who` to `lane.Publish` for both the directory and the stamp.
Caught by `TestLanesPublishWritesUnderIdentitySlug` failing against a
built-in lane (further caught that a builtin's `creator: jaira` line means
it never gets stamped at all — fixed the test to publish a custom lane with
no `creator:` line instead). Fixed to slug once and reuse it for both,
matching the TUI exactly.

## Auth gates

None — no command in this plan requires authentication.

## Editor-launch caveat (per task constraints)

Task 3's `$EDITOR` launch cannot be driven from this environment (no
interactive terminal). What is verified by
`TestLaneScreenNewWritesSkeletonAndReloadsList`:
- the skeleton file is written to the catalogue before the editor command is
  even constructed (`writeLaneSkeleton` runs synchronously in `newLane()`);
- a non-nil `tea.Cmd` is queued (confirming `tea.ExecProcess` was invoked
  with the resolved `$EDITOR` argv — the same `editorCommand()` helper the
  ticket-body editor already uses);
- simulating the editor's exit (`m.Update(newLaneDoneMsg{})`) triggers the
  reload and the lane list growing by one.

The actual spawning and interaction with a real `$EDITOR` process is
**not** exercised by any test here and remains unverified in this session.

## Self-Check

```
FOUND: internal/cli/lanes.go
FOUND: core/lane/template.go
FOUND: internal/cli/lanes_test.go
FOUND: internal/cli/tickets.go
FOUND: internal/tui/lanes.go
FOUND: internal/tui/lanes_test.go
FOUND: internal/tui/model.go
FOUND: docs/AGENTS.md
FOUND: docs/COMMANDS.md
FOUND: f9e0ede (feat: give jaira lanes use, publish and adopt as commands)
FOUND: 3e3b5b4 (feat: read and write the default board from jaira lanes default)
FOUND: bbe3ae5 (feat: n in the lane settings screen writes a skeleton and opens it)
FOUND: 2c5941b (docs: say the by-hand lane loop is complete, because now it is)
```

## Self-Check: PASSED

## Verification

Per-task automated verify blocks: all passed (see commit history; each task
was built, vetted and tested with `-race` before its commit).

Final full-suite verbatim output (`go build ./... && go vet ./... && go test
./... -race -count=1`), run after all four commits:

```
?   	github.com/BeMuCa/jaira/cmd/jaira	[no test files]
ok  	github.com/BeMuCa/jaira/core/board	1.030s
ok  	github.com/BeMuCa/jaira/core/gate	1.758s
?   	github.com/BeMuCa/jaira/core/gitrepo	[no test files]
ok  	github.com/BeMuCa/jaira/core/identity	1.032s
ok  	github.com/BeMuCa/jaira/core/lane	2.861s
ok  	github.com/BeMuCa/jaira/core/merge	1.387s
ok  	github.com/BeMuCa/jaira/core/project	1.026s
ok  	github.com/BeMuCa/jaira/core/release	1.012s
?   	github.com/BeMuCa/jaira/core/session	[no test files]
ok  	github.com/BeMuCa/jaira/core/ticket	1.048s
ok  	github.com/BeMuCa/jaira/core/validate	1.362s
ok  	github.com/BeMuCa/jaira/internal/cli	2.216s
ok  	github.com/BeMuCa/jaira/internal/tui	11.608s
?   	github.com/BeMuCa/jaira/scripts/iconpreview	[no test files]
?   	github.com/BeMuCa/jaira/scripts/shotgen	[no test files]
```

`go vet ./...` produced no output (clean). `go build ./...` produced no
output (clean).

**End-to-end check against a real built binary**, in a throwaway repo with
`JAIRA_HOME`, `JAIRA_LANES_DIR` and `JAIRA_DEFAULT_BOARD` all redirected
away from the real user's `~/.jaira`:

```
--- catalogue file written: ---
total 12
-rw-r--r-- 1 berk berk 1634 Aug 13 19:00 e2e-lane.md
--- lanes list sees it: ---
e2e-lane       My Lane          42     true     strong   /tmp/tmp.6QLn2R6SIo/catalogue/e2e-lane.md
--- lanes use: ---
wrote /tmp/tmp.6QLn2R6SIo/repo/.jaira/lanes/e2e-lane.md
total 12
-rw-r--r-- 1 berk berk 1634 Aug 13 19:00 e2e-lane.md
--- lanes publish: ---
published /tmp/tmp.6QLn2R6SIo/repo/.jaira/shared/tester/e2e-lane.md
/tmp/tmp.6QLn2R6SIo/repo/.jaira/shared/tester/e2e-lane.md
--- lanes shared: ---
ID             FOLDER         CREATOR          PATH
e2e-lane       tester         you              /tmp/tmp.6QLn2R6SIo/repo/.jaira/shared/tester/e2e-lane.md
```

Template landed in the catalogue, `lanes use` copied it into the project's
own lane directory, `lanes publish` copied it into `.jaira/shared/tester/`,
and `lanes shared` confirmed it from a fresh read — each file landed exactly
where its command said it would.

## Known Stubs

None.

## Threat Flags

None. `use`, `publish` and `adopt` are thin wrappers over `core/lane`
functions whose overwrite-refusal and filename-from-parsed-id behavior
(guarding against a malicious source path, per the existing `T-5-03`/`T-5-07`
comments in `core/lane/share.go` and `defaultboard.go`) is unchanged — no new
trust boundary was introduced. `adopt` still requires the operator to have
already read the shared lane before running the command (documented in the
command's `Long` text), matching the TUI's own confirm-before-agree design.

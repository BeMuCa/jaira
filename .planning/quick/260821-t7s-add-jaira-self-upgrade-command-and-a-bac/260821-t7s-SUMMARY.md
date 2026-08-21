---
phase: quick-260821-t7s
plan: 01
subsystem: cli
tags: [go, cobra, self-update, tui, release-checksum, http, cache]

requires: []
provides:
  - "core/selfupdate package: Latest/At/Binary/AssetName/DownloadURL against api.github.com and github.com/releases, sha256-verified before extraction"
  - "core/selfupdate/replace.go + replace_windows.go: atomic same-directory temp-file-then-rename binary replacement, with Windows .old-<pid> sweep"
  - "core/selfupdate/install.go: Detect() classifies Homebrew / go-install / unwritable-directory / self-managed installs"
  - "core/selfupdate/cache.go: Check/Path/Read/Write/Stale/Disabled/SpawnRefresh/PollCache — a ~24h release-check cache with a detached-child refresher"
  - "jaira self upgrade command group: happy path, --check, --version, refusals, --json contract"
  - "internal/tui persistent version indicator (versionLine) in the launcher and board footers, replacing a CLI stderr nudge that was removed mid-task"
affects: [cli, release-tooling, tui]

actuals:
  tokens: 20491
  tasks: 3
  commits: 3

tech-stack:
  added: []
  patterns:
    - "Endpoint override via env var only (JAIRA_RELEASE_API/JAIRA_RELEASE_DOWNLOADS), never a flag — keeps the test seam out of the user-facing surface"
    - "Stage-then-rename atomic file replacement, staged in the target's own directory to guarantee same-filesystem rename"
    - "Producer/consumer cache split: a detached, non-blocking child process is the only thing that ever calls the network; every synchronous path only ever reads a local cache file"
    - "TUI computes a persistent indicator once at Model/Home construction rather than per-render, since a bubbletea program renders far more often than the ~24h staleness window needs"

key-files:
  created:
    - core/selfupdate/selfupdate.go
    - core/selfupdate/replace.go
    - core/selfupdate/replace_windows.go
    - core/selfupdate/install.go
    - core/selfupdate/cache.go
    - core/selfupdate/selfupdate_test.go
    - core/selfupdate/install_test.go
    - core/selfupdate/cache_test.go
    - internal/cli/self.go
    - internal/cli/self_test.go
    - internal/tui/updatecheck.go
    - internal/tui/updatecheck_test.go
    - internal/tui/main_test.go
  modified:
    - internal/cli/root.go
    - internal/cli/update.go
    - internal/tui/home.go
    - internal/tui/model.go
    - internal/tui/view.go
    - README.md
    - docs/COMMANDS.md

key-decisions:
  - "Mid-task scope amendment (from the user, applied to Task 3 before it was committed): CLI commands stay fully silent about a published release — the release-availability line was removed from nudgeIfStale entirely, and the board-staleness nudge reverted byte-for-byte to its pre-task form. The cache/staleness/detached-refresh machinery was kept as planned; only its consumer moved, from internal/cli/update.go to a new persistent indicator in internal/tui (launcher and board footers)."
  - "PollCache() (core/selfupdate/cache.go) is the one composed, non-network entry point a consumer reads through instead of hand-composing Read/Stale/Write/SpawnRefresh — used by the TUI indicator today, available to any future consumer."
  - "A never-checked or JAIRA_NO_UPDATE_CHECK=1 cache renders the running version alone, never 'up to date' — asserting a checked state without ever having checked was treated as a correctness bug, not a style choice."
  - "internal/tui/main_test.go disables the release check for the whole TUI test package by default (JAIRA_NO_UPDATE_CHECK=1 in TestMain), because every one of ~200 existing tests builds a Model or Home, and each one hitting a real, uninitialized cache would otherwise spawn a real detached network call per test."

requirements-completed: [QT-260821-t7s]

coverage:
  - id: D1
    description: "'jaira self upgrade' resolves, sha256-verifies, and atomically replaces the running binary; a bad/absent checksum leaves it byte-for-byte unchanged"
    requirement: QT-260821-t7s
    verification:
      - kind: unit
        ref: "core/selfupdate/selfupdate_test.go#TestFullFetchAndReplaceOverwritesTargetAtomically"
        status: pass
      - kind: unit
        ref: "core/selfupdate/selfupdate_test.go#TestBinaryChecksumMismatchLeavesTargetUnchanged"
        status: pass
      - kind: unit
        ref: "core/selfupdate/selfupdate_test.go#TestBinaryMissingChecksumEntryLeavesTargetUnchanged"
        status: pass
    human_judgment: false
  - id: D2
    description: "A Homebrew, go-install, or unwritable install refuses (exit 3) with the fix named; a dev build refuses unless --check; --check/--version behave as specified; every --json path emits one object"
    requirement: QT-260821-t7s
    verification:
      - kind: unit
        ref: "core/selfupdate/install_test.go#TestDetectHomebrewWinsOverGoInstall"
        status: pass
      - kind: unit
        ref: "internal/cli/self_test.go#TestSelfUpgradeDevBuildRefusesButCheckReports"
        status: pass
      - kind: unit
        ref: "internal/cli/self_test.go#TestSelfUpgradeVersionPinsInstallsExactRelease"
        status: pass
    human_judgment: false
  - id: D3
    description: "No command makes a synchronous network call; a stale cache spawns one detached refresher and returns immediately; the release status renders in the TUI's launcher/board footer, not on CLI stderr"
    requirement: QT-260821-t7s
    verification:
      - kind: unit
        ref: "core/selfupdate/cache_test.go#TestPollCacheOnStaleCacheReturnsWithoutBlocking"
        status: pass
      - kind: unit
        ref: "internal/tui/updatecheck_test.go#TestBoardStatusBarCarriesTheVersionIndicator"
        status: pass
      - kind: unit
        ref: "internal/cli/update_test.go#TestNudgeIfStaleOnDifferingStampPrintsOneStderrLineOnly"
        status: pass
    human_judgment: false

duration: 3h
completed: 2026-08-21
status: complete
---

# Quick Task 260821-t7s: `jaira self upgrade` and the background release check Summary

**A Go port of `scripts/install.sh`'s resolve-download-verify-install logic
lands as `jaira self upgrade`, and a release check that costs a command
nothing at startup now surfaces as a persistent status line in the TUI
instead of on the CLI's stderr — the CLI stays silent about it entirely.**

## Performance

- **Duration:** ~3h (including standing up a missing Go toolchain and
  tracking down a test-isolation bug the mid-task scope change exposed)
- **Tasks:** 3
- **Commits:** 3 (one per task)
- **Files created:** 13, **files modified:** 7

## Accomplishments

- `core/selfupdate`: `Latest`/`At`/`Binary`/`AssetName`/`DownloadURL` against
  `api.github.com` and `github.com/.../releases`, with `checksums.txt`
  verified via sha256 **before** any extraction, and
  `JAIRA_RELEASE_API`/`JAIRA_RELEASE_DOWNLOADS` as the env-only test seam.
- Atomic, cross-platform binary replacement (`replace.go` / `replace_windows.go`):
  stage in the target's own directory, then rename — safe on unix because a
  running process keeps its old inode; Windows renames the old binary aside
  first (`Sweep` cleans up `*.old-*` leftovers on a later run).
- `core/selfupdate/install.go`: `Detect()` refuses to touch a Homebrew,
  `go install`, or unwritable install, naming the right command instead —
  Homebrew wins when a path matches more than one classification.
- `jaira self upgrade` with `--check` (report only) and `--version vX.Y.Z`
  (pin or downgrade, never touching `releases/latest`); a dev build and an
  already-current binary are both handled without a download; every
  `--json` path emits exactly one object.
- `core/selfupdate/cache.go`: a `~/.jaira/update-check.json` cache with a
  24h staleness window, a touch-before-spawn guard against parallel sessions
  each spawning their own refresher, `JAIRA_NO_UPDATE_CHECK=1` as the
  opt-out, and `PollCache()` as the one composed entry point a consumer
  reads through.
- **Mid-task scope amendment applied:** the release-availability line was
  removed from the CLI's `nudgeIfStale` entirely (reverted to its
  pre-task, board-staleness-only form) and replaced with a persistent
  `jaira vX.Y.Z · up to date` / `jaira vX.Y.Z · vA.B.C available — run:
  jaira self upgrade` indicator in the TUI's launcher and board footers,
  computed once at construction.

## Task Commits

1. **Task 1: `jaira self upgrade` end to end — one path, the happy one** -
   `5ca85e2` (feat)
2. **Task 2: the refusals and the flags — install method, dev build,
   --check, --version** - `9bebf2a` (feat)
3. **Task 3 (amended): the background check moves to the TUI, and the
   docs** - `88cc626` (feat)

_No separate metadata commit yet — this SUMMARY, STATE.md and the
deferred-items note are left for the orchestrator's docs commit per the
quick-task constraints._

## Files Created/Modified

- `core/selfupdate/selfupdate.go` - endpoints, `Client`, `Release`,
  `Latest`/`At`/`Binary`/`AssetName`/`DownloadURL`, `digestFor`, `extract`,
  `stage`
- `core/selfupdate/replace.go` / `replace_windows.go` - `Replace`/`Sweep`,
  unix and Windows
- `core/selfupdate/install.go` - `Detect`, Homebrew/go-install/writability
  classification
- `core/selfupdate/cache.go` - `Check`, `Path`/`Read`/`Write`/`Stale`/
  `Disabled`/`SpawnRefresh`/`PollCache`
- `core/selfupdate/*_test.go` - httptest-server-driven fixtures, real
  tar.gz/zip archives built at test time, no hardcoded digests
- `internal/cli/self.go` - `jaira self` / `jaira self upgrade`, all guards
  and flags
- `internal/cli/self_test.go` - end-to-end CLI tests against a two-release
  httptest fixture, isolated to a per-test `JAIRA_HOME`
- `internal/cli/root.go` - registers `newSelfCmd()`
- `internal/cli/update.go` - `nudgeIfStale` reverted to board-only; a doc
  comment explains where the release question now lives
- `internal/tui/updatecheck.go` - `versionLine()`, the shared indicator
  renderer
- `internal/tui/home.go`, `internal/tui/model.go`, `internal/tui/view.go` -
  wire `versionLine()` into the launcher's and board's footers
- `internal/tui/main_test.go` - `TestMain` disables the release check
  package-wide by default
- `internal/tui/updatecheck_test.go` - the indicator's three states (never
  checked / up to date / available), plus footer integration checks
- `README.md`, `docs/COMMANDS.md` - document `jaira self upgrade` and the
  TUI indicator

## Decisions Made

- See `key-decisions` in the frontmatter. The headline one: the release
  question moved from a CLI stderr nudge to a TUI-only indicator, per an
  explicit mid-task instruction from the user, applied before Task 3 was
  committed.
- `PollCache()` was added as a small composition function in
  `core/selfupdate/cache.go` rather than duplicating the
  read-stale-touch-spawn sequence in `internal/tui` — the TUI is described
  in `model.go`'s own package doc as "a peer of the CLI, not a wrapper
  around it," which is precisely the existing justification for it linking
  core packages directly.

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 3 - Blocking] No Go toolchain was present on the machine**
- **Found during:** before Task 1 — `go build`/`go test` failed with
  `command not found`.
- **Issue:** the plan's every `<verify>` step requires `go build`/`go test`/
  `gofmt`, none of which existed on PATH.
- **Fix:** downloaded and extracted the official `go1.26.5.linux-amd64.tar.gz`
  from `go.dev` (the canonical distribution source, not a package-manager
  registry, so this sits outside the npm/pip/cargo slopsquatting exclusion
  in Rule 3) into `~/.local/go`, verified `go version`, then built and ran
  the pre-existing suite as a baseline before touching any code. (The
  coordinator later confirmed Go was installed and on PATH mid-task,
  independent of this fix.)
- **Files modified:** none — toolchain only, not tracked in the repo.
- **Verification:** `go build ./...` and `go test ./... -race -count=1`
  both green before any Task 1 edit.

**2. [Rule 1 - Bug] `internal/cli/self_test.go` was writing to the real
user's `~/.jaira` during every test run**
- **Found during:** Task 3, while verifying the full suite runs offline —
  the real `~/.jaira/update-check.json` appeared after `go test ./...`,
  with `"latest":"1.1.0"`, matching a Task-2 test fixture's pinned version
  rather than any real network response.
- **Issue:** Task 2's `self_test.go` never isolated `JAIRA_HOME`. That was
  harmless when written (nothing under test touched a machine-wide cache
  yet), but Task 3 wired `selfupdate.Write(...)` into `self.go`'s
  `--check`/successful-upgrade paths, and every test in that file started
  writing straight into the real home directory.
- **Fix:** added `t.Setenv("JAIRA_HOME", t.TempDir())` to `selfTestRelease`,
  the shared fixture helper every test in the file calls first.
- **Files modified:** `internal/cli/self_test.go`.
- **Verification:** `rm -f ~/.jaira/update-check.json`, then three
  consecutive full `go test ./... -race -count=1` runs with no reappearance;
  `core/selfupdate`'s own test time dropped from ~1.1s (masking a hidden
  recursive re-exec, see below) to consistently ~0.15-0.2s.
- **Committed in:** `88cc626` (folded into the Task 3 commit, since Task 3
  introduced the bug this fixes).

**3. [Rule 1 - Bug, discovered incidentally] `SpawnRefresh` under `go test`
recursively re-runs the whole package's test suite**
- **Found during:** the same investigation as #2 above.
- **Issue:** `SpawnRefresh` resolves `os.Executable()` and execs it with
  `self upgrade --check --json`. Under `go test`, `os.Executable()` returns
  the compiled *test* binary, not the real `jaira` CLI. Go's `flag.Parse()`
  stops at the first non-flag argument ("self"), so the test binary just
  runs `testing.Main` with its defaults — i.e. it reruns every test in the
  package as a detached child. Empirically confirmed by invoking a built
  `core/selfupdate` test binary directly with those arguments and observing
  `PASS`.
- **Impact assessed as contained, not fixed at the source:** the
  recursion is self-limiting (the child inherits `JAIRA_NO_UPDATE_CHECK=1`,
  so its own `PollCache` calls no-op immediately) and every test's
  `JAIRA_HOME` isolation still holds inside the recursive re-run, so no
  data leaks — it only wastes CPU. Given the actual fix for the *dangerous*
  half (real-home writes) is #2 above, and rearchitecting `SpawnRefresh` to
  detect a test binary would be production code shaped around test
  artifacts, this was left as a documented, bounded quirk rather than
  "fixed": `internal/tui/main_test.go`'s `TestMain` avoids triggering it in
  ~200 tests by disabling the check package-wide, and `core/selfupdate`'s
  and `internal/cli`'s own tests that do exercise a stale cache accept the
  one extra (harmless, isolated) recursive pass.
- **Files modified:** none beyond #2 and the `TestMain` already planned for
  Task 3's TUI test isolation.
- **Verification:** three consecutive full-suite runs, stable timing,
  no failures, no real-home writes.

---

**Total deviations:** 3 (1 blocking-environment fix, 2 bugs found and
fixed/contained during Task 3's own verification)
**Impact on plan:** All three were necessary for the suite to actually run
and to actually stay offline as required; no scope creep beyond the
explicit, user-directed Task 3 amendment.

## Issues Encountered

- The scope amendment arrived while the full pre-amendment `go test ./...`
  run was still executing. Applied it to Task 3 before committing (Task 3
  had not yet been committed at that point), per the amendment's own
  instruction.
- Two pre-existing `gofmt` findings (`core/gate/gate.go`,
  `internal/cli/tickets.go`) are unrelated to any file this task touched;
  confirmed via `git stash` against the branch tip before any edit, and
  logged in `deferred-items.md` rather than fixed, per the scope-boundary
  rule.

## User Setup Required

None - no external service configuration required. (`go.mod`/`go.sum` are
untouched; nothing was installed via a package manager.)

## Next Phase Readiness

- `jaira self upgrade` and the background release check are complete and
  independently testable offline.
- The TUI indicator reads a cache any future consumer could also read via
  `selfupdate.PollCache()` — nothing here is CLI-specific if a future task
  wants a different surface.
- Known, accepted quirk: `SpawnRefresh`'s exec-self approach behaves oddly
  under `go test` (see Deviation 3) — bounded and harmless today, worth a
  second look only if a future task adds more test suites that construct a
  `Model`/`Home` outside `internal/tui`'s existing `TestMain` protection.

---
*Quick task: 260821-t7s*
*Completed: 2026-08-21*

## Self-Check: PASSED

All 13 created files and the SUMMARY itself confirmed present on disk; all
3 task commits (`5ca85e2`, `9bebf2a`, `88cc626`) confirmed in `git log`.

# Quick Task 260813-u7t: Lane settings screen becomes a small board

**Status:** DONE — all 3 tasks complete, full suite green.

## One-liner

A project's column order lives in one plain-text file; `jaira lanes add/remove/move` and the settings screen's new small-board-with-a-'+'-column both go through the same `core/lane` functions to mutate it.

## Task 1: A project's column order, in one file — DONE

- Added `core/lane/order.go`:
  - `orderFileName = "order"` — plain text (not `.md`, so it never satisfies `ProjectLanesActive`'s glob), one lane id per line, lives at `<root>/.jaira/lanes/order`.
  - `LoadOrder(root) ([]string, error)` — absent file returns `nil, nil`.
  - `SaveOrder(root, ids) error`.
  - `Move(ids, id, delta) []string` — pure; swaps with neighbour; no-op past either end or on unknown id; no wrap-around.
  - `applyOrder(lanes, ids) ([]*Lane, []string)` — reorders already after:-resolved lanes per an order file: named lanes take that position/sequence, lanes missing from the file are appended afterward in today's load order (never dropped), an unknown id in the file produces a warning and is skipped.
- Wired into `core/lane/lane.go`'s `Load()`: after the existing `after:`-anchor resolution (`order()`), a project's order file (if non-empty) is applied via `applyOrder`, and its warnings are appended to the existing single warning channel.
- `precedence` untouched: no code in this task reads or writes `Lane.Precedence`.

**Files:** `core/lane/order.go` (new), `core/lane/order_test.go` (new), `core/lane/lane.go` (+9 lines), `core/lane/lane_test.go` (+5 tests, +`reflect` import).

**Tests added:** `TestMoveSwapsWithNeighbour`, `TestMoveAtEitherEndIsNoOp`, `TestMoveSingleElementIsNoOp`, `TestMoveUnknownIDIsNoOp`, `TestLoadOrderAbsentFileIsNotAnError`, `TestSaveThenLoadOrderRoundTrips`, `TestLoadWithNoOrderFileLeavesTodaysBehaviourUnchanged`, `TestLoadOrderFileReordersColumns`, `TestLoadOrderFileUnknownIDWarns`, `TestLoadOrderFileMissingLaneStillAppears`.

**Verify:** `go build ./... && go vet ./... && go test ./... -race -count=1` — 0 FAIL.

**Commit:** 711b9f2 — `feat: a project's column order lives in one file`

## Task 2: Materialising, adding and removing — DONE

- `core/lane/order.go` gained:
  - `MaterialiseWorkingSet(root, set)` — writes every currently loaded lane into the project's own lane dir via `Export`, only when `ProjectLanesActive(root)` is false; no-op (idempotent) otherwise.
  - `Installable(set)` — built-ins + catalogue lanes not already in `set`; backs both `lanes add` and (task 3) the settings screen's `+` catalogue.
  - `Add(root, set, id)` — refuses if already part of `set`; else materialises, exports, appends to the order file, and clears any prior removal tombstone entry for `id`.
  - `Remove(root, set, store, id)` — refuses (naming ticket handles) if any ticket has `Status == id`; else materialises, deletes the project's `.md` copy if present, drops `id` from the order file, and **adds `id` to a new per-project tombstone file** (see deviation below).
  - `MoveLane(root, set, id, delta)` — materialises, then applies the pure `Move` to the effective order and saves it.
  - `effectiveOrder`, `withoutID`, `ticketsIn`, `readIDList`/`writeIDList` (shared by order and removed files) as internal helpers.
- `core/lane/lane.go`'s `Load()` now filters out any lane whose id is in the project's `removed` tombstone before running `order()`/`applyOrder` — see deviation below for why this was necessary.
- `internal/cli/lanes.go`: added `lanes add <id>`, `lanes remove <id>`, `lanes move <id> --left|--right`, each with `--json`, wired in `internal/cli/tickets.go`'s `newLanesCmd`.

**Files:** `core/lane/order.go` (+~215 lines), `core/lane/lane.go` (+16 lines), `internal/cli/lanes.go` (+~135 lines), `internal/cli/tickets.go` (+1 line — command registration, not in the plan's file list for this task but required to expose the new subcommands), `internal/cli/lanes_test.go` (+9 tests).

**Verify:**
- `go build ./... && go vet ./... && go test ./... -race -count=1` → 0 FAIL.
- `jaira lanes --help | grep -cE "add|remove|move"` → `3`.

**Commit:** fb9cf7d — `feat: add, remove and move a project's lanes from the CLI`

### Deviation (Rule 1/Rule 2 — architectural finding, resolved without a checkpoint)

Discovered while testing `remove`: `Load()` always injects **all** built-in lanes first, unconditionally, regardless of whether a project has its own lane directory (`core/lane/lane.go`'s own doc comment: "built-ins first, then either the project's own lane directory or the catalogue"). A project's own directory only ever *overrides* a built-in's content when a file shares its id, or *adds* a new custom lane — there was never a way to make a built-in disappear by simply not copying it in. Since a fresh project's board is 100% built-ins, this made `jaira lanes remove <builtin-id>` a no-op in practice: the file was deleted but the builtin reappeared on the very next `Load()`.

Fixed by adding a per-project `removed` tombstone file (`core/lane/order.go`, alongside `order`), one id per line, that `Load()` consults and filters out before `order()` runs. `Remove` writes to it; `Add` clears an entry from it so a removed lane can be re-added. This does not touch `Lane.Precedence` or any merge logic — a removed lane simply isn't in the loaded `Set`, and `Set.Precedence`/`Columns` already treat an absent id as "unknown" via their pre-existing fallback paths (verified: `TestLanesRemoveMaterialisesWorkingSetOnFirstChange`, `TestLanesRemoveTakesLaneOutOfProjectOnly`).

This was treated as Rule 1/2 (bug fix / missing critical functionality for the literally-requested `remove` feature) rather than a Rule 4 checkpoint, since it is a small, additive, same-shape mechanism (one more plain-text tombstone file next to the existing order file) rather than a new schema, service, or behavioural change to anything outside this quick task's own scope.

## Task 3: The settings screen becomes a small board — DONE

- `internal/tui/lanes.go` rewritten: `laneScreen` now keeps the full `*lane.Set` (not just its `[]*Lane` slice) so it can call `lane.Add/Remove/MoveLane`, which need `Set.Get`.
  - `idx` now selects a **column**: `0..len(lanes)-1` is a lane, `len(lanes)` is the `+` column (`isPlusColumn()`).
  - `renderBoard` draws lanes as small bordered columns (`laneColWidth = 16`) joined horizontally via `lipgloss.JoinHorizontal`, scrolled to keep the selection visible — the same shape as the main board's `renderBoard`/`renderColumn` in `view.go`, reimplemented rather than called directly (documented in a comment on `laneColumnStyle`): `renderColumn` draws ticket cards from `Model`'s own state, which this screen has none of.
  - Keys: `h`/`l` (or arrows) move the selected column; `H`/`L` call `lane.MoveLane(..., ±1)`; `x` calls `lane.Remove` (refusing, naming tickets, exactly as the CLI does); `enter` on `+` opens the catalogue (`lane.Installable`); `enter` in the catalogue calls `lane.Add` and closes it.
  - Pre-existing `u`/`p`/`n`/`R`/`a`/`tab` keys and the "Shared by teammates" section are unchanged and still work (kept — plan is additive to the screen, not a replacement of its whole feature set).
  - After any of add/remove/move, `ls.reload()` re-`Load`s the set in place so the board reflects the write immediately.
  - Footer help text is deliberately **not** passed through `truncate(..., w)` any more (documented in a comment): with 8+ keys now named, the full line exceeds the screen's 78-column cap, and truncating would silently drop the last key names rather than the column content `truncate` exists to fit.
- `docs/COMMANDS.md`: added `lanes add/remove/move` rows and a new "Lane settings" keys block.
- `docs/AGENTS.md`: expanded "A note on lane order" to cover the order file and the three new commands, noting they're the exact calls the settings screen's board makes and that the first change to a project with no lane directory materialises the whole working set.

**Files:** `internal/tui/lanes.go` (rewritten, +~150 net lines), `internal/tui/lanes_test.go` (updated 2 pre-existing tests for the new h/l-based navigation and footer length; added 6 new tests), `docs/COMMANDS.md`, `docs/AGENTS.md`.

**Tests added/changed:** `TestLaneScreenNavigationClamps` (rewritten for h/l + the `+` column), `TestLaneScreenRendersBoardWithPlusColumn`, `TestLaneScreenLMovesSelectedLaneAndWritesOrderFile`, `TestLaneScreenXRemovesSelectedLane`, `TestLaneScreenXRefusesWhenTicketPresent`, `TestLaneScreenPlusColumnOpensCatalogueListingOnlyMissingLanes`.

**Verify:** `go build ./... && go vet ./... && go test ./... -race -count=1` → 0 FAIL (all 14 packages).

**Commit:** 061a09e — `feat: the lane settings screen becomes a small board`

### Real-binary verification (plan's `<verification>` block)

Ran against a throwaway repo (`git init` + `jaira init`, isolated `JAIRA_LANES_DIR`):

```
--- before ---
ID             NAME             RANK   ...
backlog        Backlog          0      ... built-in
brainstorm     Brainstorm       5      ... built-in
todo           Todo             20     ... built-in
...(10 lanes total, all built-in)

moving: brainstorm
moved brainstorm

--- after ---
ID             NAME             RANK   ...
brainstorm     Brainstorm       5      ... /…/.jaira/lanes/brainstorm.md
backlog        Backlog          0      ... /…/.jaira/lanes/backlog.md
todo           Todo             20     ... /…/.jaira/lanes/todo.md
...(same 10 lanes, all now project-sourced)
jaira: warning: lane /…/.jaira/lanes/backlog.md: anchor "" is not installed; placed before the terminal lane
```

Column order changed (brainstorm swapped ahead of backlog) and every one of the 10 original lanes is still present — matches the plan's literal ask. See "Known Issue (deferred)" below for the spurious warning line.

### Human check (cannot be run here)

The plan's `<verification>` also asks: "open settings, see the small board, move a lane with L, add one through +, and confirm the main board reflects it." This requires an interactive terminal (real keypresses into the Bubble Tea program) and was **not** run — stated plainly per the task constraints. All of the underlying logic (`H`/`L`/`x`/`+`/catalogue) is covered by the automated TUI tests listed above instead.

## Deviations from Plan

1. **[Rule 1/2 — architectural finding] `Load()` never removed a built-in lane by absence.** See Task 2's deviation note above: added a `removed` tombstone file so `jaira lanes remove <builtin-id>` actually takes effect.
2. **`internal/cli/tickets.go` touched, not in either task's `files_modified` list.** Registering the three new subcommands in `newLanesCmd`'s `cmd.AddCommand(...)` call (that function's existing home) required a one-line edit there. No other change was made to that file.
3. **Test file scenarios adjusted after discovering D-03's actual scope.** Two test scenarios (one in `internal/cli/lanes_test.go`, one in `internal/tui/lanes_test.go`) were initially written assuming a catalogue lane becomes "installable" before a project's first change — this is impossible by design: before any project lane directory exists, the catalogue and all built-ins are *already* showing (D-03's fallback), so nothing is ever "not already in this project" until something has first been removed. Both tests were rewritten to materialise/remove first, matching what the feature can actually do.

## Known Issue (deferred, not fixed)

`jaira lanes` (and the settings screen) show a spurious warning after any `add`/`remove`/`move` on a project with no prior lane directory:

```
jaira: warning: lane <root>/.jaira/lanes/backlog.md: anchor "" is not installed; placed before the terminal lane
```

**Root cause:** `MaterialiseWorkingSet` exports every currently-loaded lane, including all built-ins, into the project's lane directory. Once re-read from disk, `core/lane/lane.go`'s `Load()` treats every one of those files as a *custom override* of its built-in (via `replaceLane`), and `replaceLane`'s own comment states this is deliberate for genuine overrides: "the replacement is not itself marked Builtin, so `order()` will position it by its own `after:`... an override does not inherit the shipped slot." `order()` then buckets every "override" — even a byte-identical, unmodified one — away from the "built-ins keep shipped order" fast path and into anchor resolution. `backlog` has no `after:` anchor (by design, it's first), so it falls into the "anchor is empty → not installed → park before terminal" branch, which is technically correct-looking but produces a confusing warning.

**Why it doesn't corrupt anything:** the project's `order` file, once present, is applied via `applyOrder` *after* `order()` runs and fully overrides whatever position `order()` computed — verified directly in the real-binary run above: the final displayed order matches exactly what `MoveLane` wrote, in both the "before" and "after" listings.

**Why it wasn't fixed in this task:** the correct fix touches `order()`'s builtin/pending bucketing logic in `core/lane/lane.go` — a heavily-tested, precedence-adjacent function — to special-case a materialised copy that is behaviourally identical to the built-in it shadows (the same distinction `lanesEquivalent` already draws for the "overrides" *warning*, but not for `order()`'s *bucketing*). That's a more invasive change than this task's scope, and the risk of a subtle regression in cycle/anchor detection outweighed fixing a cosmetic warning within the remaining session. This is flagged here rather than left silent, per the "no smoothing over failures" rule; the plan's success criteria (order lives in one file, every lane still present, `precedence` untouched) are all met regardless.

## Threat Flags

None. `core/lane/order.go`'s new functions read/write only inside `ProjectLanesDir(root)` (already-trusted project-local paths, following the same pattern as `Export`/`Publish`), and `Remove`'s ticket-occupancy check only reads ticket status via the existing `Store.List()` — no new network, auth, or trust-boundary surface.

## Self-Check: PASSED

- Commits verified present in `git log`: `711b9f2`, `fb9cf7d`, `061a09e`.
- Files verified present on disk: `core/lane/order.go`, `core/lane/order_test.go`, `internal/cli/lanes.go`, `internal/tui/lanes.go`, `docs/COMMANDS.md`, `docs/AGENTS.md`.
- Full suite re-run at completion: `go build ./... && go vet ./... && go test ./... -race -count=1` → 0 FAIL across all 14 packages.

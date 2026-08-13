# Task 2: Root-aware lane loading, the three-source layout, and overridable built-ins

**Commit:** `dee99a7` — feat: give lane loading a project root, so a custom lane can override a built-in

## What changed

- `lane.Load()` became `lane.Load(root string)`. `root == ""` means no project
  is in hand (only the launcher hits this case) and behaves exactly as
  before: built-ins plus the catalogue.
- Added `lane.ProjectLanesDir(root string) string` — `<root>/.jaira/lanes`,
  built from `ticket.DirName` rather than a second literal.
- Implemented D-03 (authoritative-when-present): if `ProjectLanesDir(root)`
  exists and holds at least one `.md` file, it replaces the catalogue as the
  second lane source for that root. If it exists but is empty, the loader
  falls back to the catalogue (same outcome as if the directory did not
  exist) but emits a warning naming the directory — an empty scoped directory
  looks, from the user's side, identical to "my lanes vanished".
- Reversed the built-in collision rule: a custom lane (from either the
  catalogue or the project directory) whose `id` matches a built-in now
  **overrides** it — full content and prompt — instead of being refused. The
  override always produces a warning naming the file and the id.
- Added the gate-weakening warning (T-5-01): if the override drops
  `requires-human-exit`, `terminal`, `requires-outcome` or
  `requires-nonmodel-signal` that the displaced built-in carried, a **second,
  separate** warning line names every dropped protection. Example text:
  `lane <path>: overriding "signoff" drops requires-human-exit — an agent
  could now accept its own work here undetected`.
- Added a defensive warning (called out explicitly in Task 1's `<decided>`
  block as a Task 2 responsibility, though not itemized in Task 2's
  `<behavior>` list): if, after resolution, no installed lane declares
  `requires-specified`, the loader warns that a ticket can never be promoted
  into work. This only fires if an override strips the property from the
  `todo` built-in, since built-ins are always part of the base set.
- Updated the package doc comment and `Load`'s doc comment to explain why
  overriding is allowed and why it is loud, replacing the old "must never
  shadow a built-in" language.
- Updated all **seventeen** `lane.Load()` call sites (the plan said sixteen;
  `core/gate/review_test.go` was an additional site not listed in the task's
  `<files>` but needed the same fix to keep the package compiling — see
  Deviations):
  - `internal/cli`: `root.go` (`loadEnv`, via `s.Root`), `tickets.go` (two
    sites via `s.Root`, one via a new `bestEffortRoot()` helper for
    `newLanesCmd`, which has no open store), `checklist.go`, `resume.go`,
    `validate.go` (all via `s.Root`), `mergedriver.go` (via
    `bestEffortRoot()`, since the merge driver has no `Store` either — see
    Deviations).
  - `internal/tui`: `model.go` (`m.store.Root`), `home.go` (`""`, with a
    comment explaining the launcher spans many boards).
  - Test call sites: `core/gate/plan_test.go`, `core/gate/review_test.go`,
    `core/validate/validate_test.go`, `core/merge/merge_test.go`,
    `internal/tui/signoff_test.go`, `internal/cli/json_test.go` — all pass
    `""`.
- Corrected `.planning/REQUIREMENTS.md` LANE-06 to say a colliding lane
  overrides with a warning rather than being rejected. **This edit was made
  but is intentionally left uncommitted** — see Deviations.
- `core/lane/lane_test.go` is new: ten tests covering built-ins alone,
  catalogue loading and anchor ordering, project-directory authority
  (including the catalogue-lane-does-not-leak case and the
  no-project-directory-still-gets-catalogue case), the empty-project-directory
  warning, override + exact warning text, single- and multi-protection
  dropping, duplicate ids within one directory, the requires-specified
  safety net, and an unknown-anchor warning-and-ordering case.

## Verify output (verbatim)

```
$ export PATH=$PATH:$HOME/.local/go/bin && go build ./... && go test ./core/lane/... ./core/gate/... ./core/validate/... ./core/merge/... ./internal/... 2>&1 | tail -20
ok  	github.com/BeMuCa/jaira/core/lane	0.151s
ok  	github.com/BeMuCa/jaira/core/gate	0.168s
ok  	github.com/BeMuCa/jaira/core/validate	0.156s
ok  	github.com/BeMuCa/jaira/core/merge	0.107s
ok  	github.com/BeMuCa/jaira/internal/cli	0.108s
ok  	github.com/BeMuCa/jaira/internal/tui	1.654s

$ export PATH=$PATH:$HOME/.local/go/bin && grep -rn "lane.Load()" --include=*.go . | grep -v "^./.git" | wc -l | grep -qx 0 && echo "no zero-arg Load remains"
no zero-arg Load remains
```

Full suite after the commit:

```
$ export PATH=$PATH:$HOME/.local/go/bin && go build ./... && go vet ./... && go test ./... -count=1
?   	github.com/BeMuCa/jaira/cmd/jaira	[no test files]
ok  	github.com/BeMuCa/jaira/core/board	0.011s
ok  	github.com/BeMuCa/jaira/core/gate	0.169s
?   	github.com/BeMuCa/jaira/core/gitrepo	[no test files]
ok  	github.com/BeMuCa/jaira/core/lane	0.152s
ok  	github.com/BeMuCa/jaira/core/merge	0.134s
ok  	github.com/BeMuCa/jaira/core/project	0.033s
ok  	github.com/BeMuCa/jaira/core/release	0.015s
?   	github.com/BeMuCa/jaira/core/session	[no test files]
ok  	github.com/BeMuCa/jaira/core/ticket	0.031s
ok  	github.com/BeMuCa/jaira/core/validate	0.132s
ok  	github.com/BeMuCa/jaira/internal/cli	0.093s
ok  	github.com/BeMuCa/jaira/internal/tui	1.552s
?   	github.com/BeMuCa/jaira/scripts/iconpreview	[no test files]
?   	github.com/BeMuCa/jaira/scripts/shotgen	[no test files]
```

Zero-warnings proof, run against a freshly `git init`'d + `jaira init`'d repo
outside this checkout:

```
$ go run ./cmd/jaira -C "$TMPD" lanes 2>&1 1>/dev/null
(nothing printed to stderr)
```

## Deviations from plan

**1. [Rule 3 — blocking issue] `core/gate/review_test.go` was an uncounted call site.**
The plan's `<files>` list for Task 2 and its "sixteen call sites" claim both
omit `core/gate/review_test.go`, which also calls `lane.Load()`. Leaving it
alone would have left `core/gate` unbuildable. Updated it the same way as
every other test call site (`lane.Load("")`). Included in the same commit
since it's mechanically identical to the other listed test fixes.

**2. [Design-conflict resolution, not a bug] Skipped the D-02 "drift" test/warning in `Load()`.**
Task 2's `<behavior>` block lists "Drift, per D-02: same id in project and
catalogue with different bytes warns" as one of `lane_test.go`'s cases. But
Task 1's `<decided>` block — which the instructions told me to treat as
settling the design and to read before Task 2 — explicitly states: *"Rather
than checking on every command, the lane settings screen runs the check when
it opens."* `Load()` runs on every command (it's called by `loadEnv`, `jaira
lanes`, the TUI reload, the merge driver, etc.), so implementing a drift
warning there directly contradicts the decided design, which places this
check in Task 6's lane settings screen instead. I followed the decided text
(more recent, explicitly settling the question) over the earlier, unrevised
acceptance bullet. **The drift check is not implemented anywhere by this
commit** — it is entirely Task 6's responsibility. Flagging this explicitly
so Task 6's author does not assume any part of it already exists.

**3. [Adaptation — plan assumed a `Store` that isn't there] `mergedriver.go` and `newLanesCmd` had no open `Store`.**
The plan says to pass `s.Root` in `mergedriver.go`, but `newMergeDriverCmd`'s
`RunE` never opens a `ticket.Store` — it only reads three file paths git gives
it. Likewise `newLanesCmd` (`jaira lanes`) has never required a board to run
in (arguably by design — it should work to discover where to write a lane
file before a board is fully set up). Added a small `bestEffortRoot()` helper
in `internal/cli/root.go` (`ticket.Discover(g.dir)`, `""` on any error) and
used it in both places instead. This preserves "jaira lanes works outside a
project" and gives the merge driver the correct root when it does run inside
one.

**4. [Explicit instruction override — not committed] `.planning/REQUIREMENTS.md` LANE-06 edit made but left uncommitted.**
Task 2's `<action>` explicitly asks for this correction, and I made it (see
above — LANE-06 now reads "overrides it... reported as a warning, never
silent"). However, my direct execution instructions for this run said
explicitly: *"Do NOT commit .planning/ artifacts"* and *"Commit it
atomically, code changes only."* I followed the more specific, more recent
instruction and left `.planning/REQUIREMENTS.md` as an uncommitted local
change rather than including it in the commit. **Whoever picks this up next
should either commit that file separately or fold it into a later commit** —
it is not lost, just not part of `dee99a7`.

**5. Interpretation call: built-ins are always part of the resolved set, in both catalogue and project modes.**
Task 2's action text says "Resolution order is built-ins, then the source
chosen by D-03" — read literally, built-ins are an unconditional base layer
and the project directory (when authoritative) only adds to or overrides
that base by id; it does not let a project drop a built-in it doesn't
mention. I implemented it this way, and it is the reading consistent with
the task's own behavior-bullet for the empty-directory case ("warns *rather
than silently yielding a built-ins-only board*" — implying the empty-case
outcome, sans the warning, would already have been built-ins-only, i.e.
built-ins are always there). Real scoping of *which built-ins* a project
shows (e.g. a default board that drops `signoff`) is not something `Load()`
does in this task — that is `Materialise`'s job in Task 8, which writes
literal file copies for the lanes a project wants and relies on the override
mechanism built here, not on `Load()` excluding un-mentioned built-ins.
Flagging this so Task 8's author checks this assumption before building on
it — if Task 8 needs `Load()` to actually exclude built-ins under
project-directory mode, that is a change to this task's code, not something
already supported.

## What the next task's author needs to know

- `lane.Load(root string)` is the new signature everywhere. Pass a real root
  whenever a `*ticket.Store` is in hand (`s.Root` / `m.store.Root`); use the
  new `bestEffortRoot()` helper in `internal/cli` when no store is open but a
  best-effort root is still useful; pass `""` only when truly no project is
  in scope (the launcher).
- `lane.ProjectLanesDir(root)` exists and is exported for Task 6/8 to write
  into (`Export`, `Materialise`).
- The override mechanism lives entirely inside `Load()`: any future writer
  (Export, Materialise, Adopt) just needs to write a valid lane file to the
  right directory — the loader already handles override detection, the two
  distinct warnings, and duplicate-id resolution.
- The D-02 drift check does not exist yet anywhere in the codebase. Task 6
  needs to implement it from scratch in the lane settings screen, per the
  decided text.
- `.planning/REQUIREMENTS.md` has an uncommitted LANE-06 fix sitting in the
  working tree (see Deviation 4) — worth committing alongside whatever task
  runs next, or explicitly carrying forward.

## Known stubs

None — this task only touches `core/lane` and lane-loading call sites, no UI
surface with placeholder data.

## Self-Check

```
$ [ -f core/lane/lane.go ] && echo FOUND || echo MISSING
FOUND
$ [ -f core/lane/lane_test.go ] && echo FOUND || echo MISSING
FOUND
$ git log --oneline --all | grep -q dee99a7 && echo FOUND || echo MISSING
FOUND
```

## Self-Check: PASSED

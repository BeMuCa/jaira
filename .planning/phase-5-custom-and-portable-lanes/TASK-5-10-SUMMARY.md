# Phase 5 Tasks 5-10: Ordering proof, lane sharing, the default board, the agent surface, and the merge-rank relabel

## Commits

| Task | Commit | Message |
|---|---|---|
| 5 | `b643e31` | test: prove ordering, tier aliasing and unknown-lane handling for lanes |
| 6 | `0a15748` | feat: a lane settings screen to read, use and publish a lane |
| 7 | `152a122` | feat: adopt a teammate's shared lane |
| 8 | `0c168cb` | feat: a per-user default board decides a fresh board's lanes and options |
| 9 | `20a6d6b` | feat: an agent-facing surface for the default board |
| 10 | `271f27c` | feat: check lane input contracts against load order, and stop calling precedence a position |

(One unrelated commit, `e03a549` "docs: plan what a progress note contains and when one is due", landed between tasks 9 and 10 from a concurrent session on this same branch — no worktree isolation was in effect. It touches only `.planning/quick/`, nothing this run modified, and its presence in `git log` is expected, not a mistake on this run's part.)

## What each task did

**Task 5** — wrote the tests the task asked for (chained anchors, cyclic anchors, no-anchor placement, `ModelTier` round-tripping, a source-scanning guard against ever comparing a tier to a model name, `Set.Columns`/`Set.Precedence` for unknown lanes, and `gate.Decide` refusing movement both into and out of an unrecognized lane). **Every test passed on the first run — no fix was needed in `core/lane/lane.go`.** Task 5 predates the ordering rescope recorded in `lane-system-design.md`; nothing in its behavior list contradicted the rescoped Task 10, so no adaptation was needed there.

**Task 6** — `core/lane/share.go` gained `Bytes`, `Export`, `Publish` and `Adopt`, all copying a lane's bytes verbatim and naming the destination from the lane's own `ID`, never a source filename. `Drift`/`RefreshDrift` implement D-02: a project lane whose bytes differ from its catalogue copy of the same id is flagged, and pulled back in on request, never automatically. `internal/tui/lanes.go` is the lane settings screen (`L` from the board): lists every loaded lane with its source and any drift, shows the selected one's full prompt, `u` exports it into the project, `p` publishes it to `.jaira/shared/<slug>/`, `R` pulls a drifted lane's catalogue copy back in.

**Task 7** — `lane.Shared` walks `.jaira/shared/*/*.md`, skipping and warning about anything unparseable rather than hiding everyone else's lanes; shared lanes are never loaded by `Load`. The lane settings screen grew a second, `tab`-focused section listing them, with `a` to adopt — a colliding catalogue id is held for a second confirming press rather than overwritten on the first. `jaira lanes shared` is the same list without the TUI.

**Task 8** — `core/lane/defaultboard.go`: `LoadDefaultBoard`/`SaveDefaultBoard` read and write `~/.jaira/default-board.md` (`$JAIRA_DEFAULT_BOARD` overrides the path); `Differs`/`Materialise` write `.jaira/lanes/` only when the selection differs from the built-in set, so an unmodified board carries no lane files. `core/ticket/body.go`'s `NewBody` replaces the CLI's hardcoded `newTicketBody`, resolving Options from `Set.Options()` (the installed lanes' `requires-option` values) with the default board's choices pre-ticked; the TUI's own ticket creation, which previously wrote an empty body, now goes through the same function. `jaira init` loads and materialises the default board, but only when the project does not already have its own `.jaira/lanes/` (see Deviations). The default board screen (`d` from the home screen) ticks lanes and options and opens a lane's file in `$EDITOR`; it does not reorder.

**Task 9** — `lane.Validate` reports a default board naming an uninstalled lane or an unclaimed option as warnings (folding in `LoadDefaultBoard`'s own parse-failure warning too), through the one warning channel `jaira lanes` already surfaces everywhere. `jaira lanes path` now names the default board file and whether it exists, in both output modes and in and out of a project; `jaira lanes template --board` prints its shape. `jaira lanes --help` documents the four-step read-write-verify loop.

**Task 10** — was rescoped before this run (see `lane-system-design.md`'s STOP section): no `precedence` value changed, `order()` is untouched, and `core/lane/builtin/*.md` was not touched (confirmed: `git diff --stat` across every commit in this run against `core/lane/builtin/` is empty). What changed: a new `checkContracts` walks lanes in display order and warns when a lane's `input-requires` names a field nothing before it — the ticket itself (`ticket.SuppliedFields`) or an earlier lane's `output-produces` — supplies; and `jaira lanes`/`lanes show` relabel the number `RANK`/`Merge rank` instead of `PREC`/`Precedence`, since it decides a merge, not a position. `docs/AGENTS.md` gained a short note saying the same.

## Deviations from the plan

**1. [Rule 1 — bug, pre-existing] `internal/tui` carried its own duplicate `identity()` function.** `core/identity` was created in an earlier commit (`0180214`, before this run) and `internal/cli/root.go` was updated to call through to it, but `internal/tui/view.go` still had a byte-identical local copy that `model.go` and `signoff.go` called instead. Task 6 explicitly asks for "the TUI cannot import `internal/cli` without a cycle" reasoning to justify `core/identity`'s existence — the TUI simply hadn't been migrated onto it. Fixed as part of Task 6: `model.go` and `signoff.go` now call `identity.Current(...)`, and the dead duplicate in `view.go` is deleted. Files: `internal/tui/model.go`, `internal/tui/signoff.go`, `internal/tui/view.go`. Commit: `0a15748`.

**2. [Interpretation — plan's own text contradicts its literal signatures] `Export`/`Publish`/`Adopt` gained an `overwrite bool` parameter.** The plan lists `Export(l *Lane, dstDir string) (string, error)` etc. with no `overwrite` argument, but the same paragraph says "Refuse to overwrite an existing file unless an explicit `overwrite bool` is passed; the caller asks," and Task 7 requires `Adopt`'s "refusing if that id already exists there unless confirmed" — which needs exactly this parameter. Implemented all three with a trailing `overwrite bool`. `Publish` additionally takes an explicit `who string` for the creator stamp, rather than inferring it from the destination directory's basename, since that seemed clearer and directly testable. Files: `core/lane/share.go`. Commit: `0a15748`.

**3. [Interpretation — same class of gap] `lane.Shared` returns warnings as a third value.** The plan's signature is `Shared(root string) ([]SharedLane, error)`, but its own behavior text requires a file that fails to parse to be "skipped with a warning rather than failing the walk," and Task 7's test list requires asserting that warning. Implemented as `Shared(root string) ([]SharedLane, []string, error)`. Files: `core/lane/share.go`. Commit: `152a122`.

**4. [Rule 2 — missing critical functionality] `jaira init` does not re-materialise the default board once a project has its own `.jaira/lanes/`.** The plan does not say this explicitly, but `newInitCmd`'s own doc string promises "Safe to run more than once," and unconditionally re-running `Materialise` on every `init` would silently overwrite a hand-edited project lane file with whatever the default board currently says — a second `init` after someone edited a materialised lane file would destroy that edit. Guarded with `lane.ProjectLanesActive(s.Root)`: once a project's lane directory exists and holds a file, a later `init` reports "This project already scopes its own lanes; the default board was not applied" instead of touching it. Tested in `TestInitTwiceDoesNotReapplyDefaultBoardOverProjectChoices`. Files: `internal/cli/tickets.go`. Commit: `0c168cb`.

**5. [Interpretation — plan's own text contradicts its literal signature] `lane.Validate` takes the board, not only the set.** The plan's signature is `Validate(set *Set) []string`, but checking "a default board naming a lane that is not installed" is not possible without also having the board's selection in hand. Implemented as `Validate(board *DefaultBoard, set *Set) []string`. Files: `core/lane/defaultboard.go`. Commit: `20a6d6b`.

**6. [Resolution of an internal contradiction in the plan] `"plan"` is not in `ticket.SuppliedFields`.** Task 10's action text says the ticket-supplied exemption set covers, "at minimum, title, goal, context, definition-of-done, assignee, plan and diff" — but its own required test says "a custom lane requiring `plan` and ordered before `pre-process` warns." Those two statements cannot both be true: if `"plan"` is unconditionally exempt, a lane requiring it can never warn regardless of order, which is exactly the case the test demands. I kept `"plan"` out of the exempt set (the Plan heading a new ticket is seeded with holds no items, so it does not actually satisfy a lane's `input-requires` — see `flow.go`'s `showForLane`), and kept `"diff"` in it (nothing produces `diff` via `output-produces`, so excluding it would put the shipped `review` lane above zero warnings, which the plan's own success criteria forbid). This resolves the contradiction in the direction that keeps every stated test and the zero-warnings bar true simultaneously. Documented in `ticket.SuppliedFields`'s doc comment. Files: `core/ticket/schema.go`. Commit: `271f27c`.

**7. [Cosmetic, not functional] Task 8's commit message contains one inaccurate line.** It says "internal/tui now calls through to core/identity... (unrelated fix noticed while wiring ticket creation through NewBody)" — that fix was actually made and committed in Task 6 (`0a15748`), not Task 8 (`0c168cb`). Copy-paste error while drafting the message; the `0c168cb` diff itself does not touch `identity` anywhere (confirmed via `git show --stat`). No code or test is affected.

## Constraint compliance

- No lane's `precedence` value was changed. `core/lane/builtin/*.md` was not modified by any commit in this run (`git diff --stat` against it across the whole range is empty).
- `order()` in `core/lane/lane.go` is untouched; the new `checkContracts` only reads the list `order()` already produced.
- Every task's `<verify>` block was run and its output is below; no commit was made on a red tree.

## Verbatim final `go build && go vet && go test ./...`

```
$ export PATH=$PATH:$HOME/.local/go/bin && go build ./... && go vet ./... && go test ./... -count=1
?   	github.com/BeMuCa/jaira/cmd/jaira	[no test files]
ok  	github.com/BeMuCa/jaira/core/board	0.022s
ok  	github.com/BeMuCa/jaira/core/gate	0.248s
?   	github.com/BeMuCa/jaira/core/gitrepo	[no test files]
ok  	github.com/BeMuCa/jaira/core/identity	0.014s
ok  	github.com/BeMuCa/jaira/core/lane	0.518s
ok  	github.com/BeMuCa/jaira/core/merge	0.125s
ok  	github.com/BeMuCa/jaira/core/project	0.025s
ok  	github.com/BeMuCa/jaira/core/release	0.007s
?   	github.com/BeMuCa/jaira/core/session	[no test files]
ok  	github.com/BeMuCa/jaira/core/ticket	0.027s
ok  	github.com/BeMuCa/jaira/core/validate	0.109s
ok  	github.com/BeMuCa/jaira/internal/cli	0.234s
ok  	github.com/BeMuCa/jaira/internal/tui	1.744s
?   	github.com/BeMuCa/jaira/scripts/iconpreview	[no test files]
?   	github.com/BeMuCa/jaira/scripts/shotgen	[no test files]
```

## Zero-warnings proof (re-run after every task)

Reproduced fresh after Task 10, isolated from the real developer machine:

```
$ h=$(mktemp -d) && d=$(mktemp -d) && (cd "$d" && git init -q . \
    && JAIRA_HOME="$h" JAIRA_LANES_DIR="$h/lanes" JAIRA_DEFAULT_BOARD="$h/no-board.md" jaira init >/dev/null 2>&1 \
    && JAIRA_HOME="$h" JAIRA_LANES_DIR="$h/lanes" JAIRA_DEFAULT_BOARD="$h/no-board.md" jaira lanes 1>/dev/null 2>stderr.out; cat stderr.out; wc -l < stderr.out)
0
```

Task 10's own literal verify (`jaira lanes` against the real, unset environment on this machine, which has no `~/.jaira/lanes` or `~/.jaira/default-board.md`):

```
$ jaira lanes 2>&1 1>/dev/null | wc -l
0
```

## Task-by-task verify output (verbatim, condensed to pass/fail — full logs shown during the run)

- Task 5: `go test ./core/lane/... ./core/gate/... -run 'Order|Anchor|Tier|Unknown|Passthrough|Cycle' -v` — all 10 new/relevant tests PASS.
- Task 6: `go test ./core/identity/... ./core/lane/... ./internal/tui/...` — PASS (identity: 8 tests including 3 new path-safety ones; lane: adds Export/Publish/Adopt/Drift tests; tui: adds `lanes_test.go`).
- Task 7: `go test ./core/lane/... ./internal/tui/... ./internal/cli/...` — PASS (adds `Shared`/adopt tests in all three packages).
- Task 8: `go test ./core/lane/... ./core/ticket/... ./internal/cli/... ./internal/tui/...` — PASS. Both of the plan's literal `init` shell verifies reproduced independently:
  - no default board → no `.jaira/lanes` — confirmed.
  - a default board naming `[backlog, todo, done]` → exactly those three files, `review.md` absent — confirmed.
- Task 9: `go test ./internal/cli/... ./core/lane/...` — PASS. Both of the plan's literal shell verifies reproduced:
  - `lanes path --json` names `default_board`; `lanes template --board` emits a file starting `lanes:` — confirmed.
  - a default board naming `nosuchlane`/`nosuchoption` produces both warnings on `jaira lanes` stderr — confirmed.
- Task 10: `go build ./... && go vet ./... && go test ./... -count=1` — PASS (shown above). `jaira lanes` warning-line count on a clean install — `0` (shown above).

## Human-check items left unverified

These are TUI visual/interactive checks the constraints explicitly say to name rather than claim:

- **Task 6's human-check**: opening the board, pressing `L`, selecting the Review lane, confirming its prompt is readable and scrolls, publishing with `p`, confirming the message names a path under `.jaira/shared/`, that the file exists and is `git status`-untracked-but-not-ignored, and that a second `p` refuses. I verified the file-system and message-content halves of this at the model level (`internal/tui/lanes_test.go`'s `TestLaneScreenPublishWritesUnderIdentitySlug`, `TestLaneScreenPublishRefusesSecondPublish`) and confirmed via `core/board`'s existing gitignore tests that `.jaira/shared/` is never covered by `/.jaira/lanes/`'s ignore rule (Task 4, already merged). I did **not** open an actual terminal and visually confirm the prompt pane renders and scrolls correctly — there is no terminal available to this run to do that in.
- **Task 7's human-check**: opening the lane settings screen with a lane already published under `.jaira/shared/someone-else/`, confirming it appears with its creator, reading the prompt, pressing `a`, and confirming the message and `jaira lanes` from a different project. I verified the adopt mechanics, the confirm-then-overwrite flow, and the cross-catalogue visibility at the model/function level (`TestLaneScreenAdoptWritesIntoCatalogue`, `TestLaneScreenAdoptRefusesThenConfirms`, `TestAdoptThenLoadFindsIt`). I did **not** visually confirm the rendered shared-section listing in a real terminal.
- **Task 8's human-check**: launching the launcher, pressing `d`, unticking a lane and ticking `brainstorm`, saving, confirming `~/.jaira/default-board.md` holds that selection, pressing `e` on a lane and confirming it opens in `$EDITOR` with the board redrawing on exit, then running `jaira init` in a fresh repo and confirming the materialised lanes match. I verified the toggle/save round trip and the built-in "no file to edit" refusal at the model level (`internal/tui/defaultboard_test.go`), and verified `jaira init`'s materialisation behavior directly via the CLI (`internal/cli/init_test.go` and the plan's own shell verifies, reproduced above). I did **not** launch a real `$EDITOR` process from the TUI and watch the screen redraw on exit — `openSelectedInEditor`'s `tea.ExecProcess` wiring is exercised by code inspection and by the "refuses for a built-in" test, not by an actual editor round trip.

## Known stubs

None. Every screen and command added in tasks 5-10 is wired to real data (`lane.Load`, `lane.Shared`, `lane.LoadDefaultBoard`) with no hardcoded or placeholder values.

## Threat model

No new trust boundaries. Tasks 6/7 implement exactly the mitigations the phase's threat register already assigned to them (T-5-02: shared lanes never auto-load; T-5-03: filenames always derive from the validated `ID`/slug, tested with traversal inputs; T-5-07: `Materialise` names files from the resolved lane's `ID`, never the default board's string, tested in `TestMaterialiseWarnsOnUnknownLaneAndContinues`). Task 10 adds no new surface — `checkContracts` only reads already-loaded, already-validated `Lane` structs.

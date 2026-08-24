# Next tasks — 2026-08-24

**Read `HANDOFF.md` first** for what landed and what the user's own board looks
like. This file is the execution list: seven items, ordered, each with the
decision already made, where the code is, and what proves it done. A session
that starts with "go" works this list top to bottom.

## Ground rules for this list

- **Every item goes through a GSD entry point** (`/gsd:quick`), per the project
  CLAUDE.md. One quick task per item, atomic commits, a SUMMARY, a STATE row.
- **Verify on the running binary.** `go build -o ~/.local/bin/jaira ./cmd/jaira`
  after every change — a stale binary already made a working feature look broken
  once. Never `go install`.
- **Gate on `gofmt -l core internal` listing exactly `core/gate/gate.go` and
  `internal/cli/tickets.go`.** Both carry pre-existing alignment groups.
- **The user's real board is a different repo** and must not be probed with a
  live `jaira move`: a probe that passes actually moves his ticket. Read it with
  `list --json`, `show --json`, `next --per-lane`. See the memory note
  `berks-req-board`.
- `go test ./... -race` with the cache cleared, before claiming anything.

## 1. The critique prompt contradicts itself — DONE (`260824-pgj`)

**The problem, in his words:** "finde standardmäßig etwas" against "loope bis die
Kritisier-Lane nix mehr zu meckern hat". If the lane always finds something the
loop never ends. He is right, and it is my wording.

**Where:** `lanes/critique.md`, the paragraph beginning "Default to finding
problems" (last section of the prompt). The escape clause is already there — "If
you genuinely find nothing on the fourth pass, say so plainly" — but the headline
overshadows it.

**Decision:** keep the bias against lazy approval, make termination explicit.
Rewrite so that:
- a finding must name a file and a concrete alternative, or it is not a finding;
- "nothing left" is a legitimate and expected outcome, not a failure of nerve;
- the loop's end condition is stated in the prompt itself.

Then refresh the copies: `jaira lanes adopt lanes/critique.md --force` into the
catalogue and `jaira lanes use critique --force` on his board.

**Done when:** the prompt cannot be read as "always find something", and
`core/lane/shipped_test.go` still passes.

## 2. Show whether a lane has been run — DONE (`260824-piz`)

**The evidence:** another session found nine tickets sitting in `critique`
uncritiqued, rendered identically to the two that had been through it. To know
the difference today you must read `review-summary` and check it for emptiness.

**Decision: derive it, no new field.** The lane's `output-produces` already says
what leaving it requires. A card in a lane whose output contract is unsatisfied
is "not worked yet".

**Where:**
- `core/gate` — `fieldFilled` (~`gate.go:406`) is private, and the enforcement
  loop is at ~`gate.go:362`. Export something small, e.g.
  `func OutputOwed(l *lane.Lane, t *ticket.Ticket) []string`, and have the
  existing loop use it so there is one implementation.
- Card flags: `internal/tui/view.go:305-311` (board) and
  `internal/tui/lanefocus.go:207-213` (lane focus). Both already switch on lane
  flags; add the "owes its lane's output" case there.
- `internal/cli/flow.go` — `emitPerLane` should prefer unworked tickets, and
  `next --lane <id>` likewise, so a loop can drain a lane.

**Careful:** do not mark a ticket in a non-agentic lane as unworked — `todo` and
`backlog` produce nothing, and every card there would light up.

**Done when:** a critique ticket with an empty `review-summary` is visibly
distinct from one that has it, on the board and in lane focus; `next --lane
critique` hands back an unworked one first; tests cover both states.

## 3. `dod --text` and a `[-]` state for superseded — DONE (`260824-pow`)

**The evidence:** there is no way to reword a checklist item. When two items were
redefined, the only route was ticking the old ones with the proof "obsolete,
replaced by item 6" and appending new ones — which leaves ticked items that were
never met and destroys the meaning of the ticks.

**Where:**
- States live in `core/ticket/schema.go:185-205` (`StateTodo`, `StateDoing`,
  `StateDone`, and `Marker()`).
- `SetItemState` is `core/ticket/checklist.go:75`, `AddItem` at `:189`,
  `replaceMarker` at `:173`. The writer is surgical by design: one marker
  character changes and the item's line is otherwise byte-identical. Keep that.
- CLI: `internal/cli/checklist.go` (the `dod` command and its flags).

**Decisions:**
- `jaira dod <id> <n> --text "..."` rewrites one item's text, leaving its state
  and its proof alone.
- A fourth state `[-]` for superseded, distinct from done: it did not happen and
  will not. `DoDComplete` must treat it as **not** blocking completion (that is
  the whole point) while `Marker()` and the parser round-trip it.
- Both work on `--plan` too, since the plan checklist shares the machinery.

**Done when:** an item can be reworded without touching its state, a superseded
item does not read as achieved, `dod --done` on all remaining items still
completes the ticket, and the round-trip tests cover `[-]`.

## 4. The small package — DONE (`260824-pvd`)

Four separate frictions, one quick task, one commit each.

- **`move --dry-run`.** Run the gates, print the violations, write nothing.
  `internal/cli/flow.go`, `newMoveCmd` — the gate call is at ~`:169-190`, and the
  mutation right after it. An agent already created and deleted a throwaway
  ticket just to test behaviour, and I moved one of his tickets by accident
  probing a gate. Both would have been a `--dry-run`.
- **`show --notes-last N`.** One ticket is 263 lines / 56 KB because notes append
  without bound. `printDetail` (`internal/cli/tickets.go:~505-560`) dumps the
  whole body last. Decide the default: show all (today) or count them and show
  the last few. The user's register argues for the latter.
- **A DoD counter in `jaira list`.** The listing already prints handle, title and
  assignee per lane; `dod_complete` and the item counts exist in the JSON
  already. Add `DoD n/m` to the human line.
- **A claim expiry warning.** `ClaimActive` (`internal/cli/claim.go:25-33`)
  silently treats a claim older than `ClaimTTL` as abandoned. Nothing says so
  when it happens. Warn where a stale claim is stepped over.

**Do not "fix" the lane warning going to stdout — it already goes to stderr**
(`internal/cli/root.go:261-263`). It only looks mixed in a terminal. The other
session's feedback was wrong on this, and on two more points recorded in
HANDOFF.md.

## 5. Delete a ticket, with a double check — DONE (`260824-q2i`)

**Today:** `Store` has `Archive` (`core/ticket/store.go:120`) and `Restore`
(`:145`) and no delete at all. Archiving is right for finished work; a ticket
created by mistake or a throwaway probe leaves a file that only `rm` removes.

**Decisions:**
- `jaira delete <id>` removes the file. It asks first, and the confirmation is
  the ticket's handle typed back, not a y/n — deletion is the one irreversible
  operation in a tool whose whole promise is that nothing is lost.
- `--force` skips the prompt for scripts, and says what it deleted.
- In the TUI it belongs on the ticket, not the board, and needs the same second
  step. `x` is archive and must not change meaning.
- Say in the help that archive is almost always what you want.

**Done when:** a mistyped confirmation deletes nothing, a correct one removes the
file, `jaira validate` is clean afterwards, and no other command starts treating
a missing ticket as an error.

## 6. Assign yourself when creating — DONE (`260824-q7x`)

**Careful — invariant 12:** capture belongs to nobody, and the pull claims. That
is deliberate and must not be reversed: `create` leaves `assignee` empty, and
moving an unassigned ticket assigns the mover. So this is an **opt-in**, not a
default.

**Decisions:**
- `jaira create --mine` sets `assignee` to `identity()` at creation.
  `--assignee <name>` keeps precedence over it.
- The TUI already claims on the pull into `todo`, so nothing is missing there for
  the normal flow. If a key is wanted at creation time, it goes in the field
  editor and must not become the default.

**Done when:** `--mine` assigns you, plain `create` still leaves it empty (a test
asserts that, because it is the invariant), and `--mine --assignee x` is x.

## 7. Keep the selected lane centred on the board — DONE (`260824-q9k`)

**His words:** the selected lane should always be in the middle, so the lanes
after it stay visible while the window allows.

**Where:** `internal/tui/renderBoard`, `internal/tui/view.go:~122-130`. Today:

    start := 0
    if m.laneIdx >= perScreen { start = m.laneIdx - perScreen + 1 }

which keeps the focused lane at the **right** edge — so what follows it is never
visible. Centre it instead: `start := m.laneIdx - perScreen/2`, clamped to
`[0, len(m.cols)-perScreen]` so neither end shows blank columns.

**Careful:** the board's own scroll state is shared with lane focus
(`m.scroll[laneID]`, see `renderLaneFocus`), and the render has a ~10-line floor.
Existing tests assert board rendering at 40 columns and at 170; both must stay
green.

**Done when:** with more lanes than fit, the focused lane sits mid-row with its
successors visible; at either end the row is still full, not padded with blanks;
a test covers first lane, middle lane, last lane.

## Not on this list, and why

- **Orchestration / a `run` command.** Another session called the lanes "a
  contract without an enforcer". True, and an explicit non-goal: the README
  excludes board-spawned background agents on purpose. Item 2 is the part of that
  complaint that is real.
- **`blocked_by` sitting empty while dependencies live in prose.** Real, no
  decision yet. Candidate: `validate` warns when a note names a handle that is
  not in `blocked_by`.
- **`docs/img/board.png` is stale** (shows the removed "N lane(s) off-screen"
  notice, lacks `v compact`). The other three screenshots were checked against
  the code and are current. Regenerating needs the demo world built first, which
  was never scripted — see the recipe in the previous handoff at `9cf71cc`.
- **v0.1.1 is uncut.** Everything after `62989f1` is unreleased.

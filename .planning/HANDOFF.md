# Handoff — 2026-08-24

State before a context clear; this file is the memory. The previous handoff
(2026-08-18, the release and the team flow) is at `9cf71cc`. Longer-lived
decisions also live in the agent memory file `jaira-design-invariants`.

**Updated after the session of 2026-08-24 that worked the list of seven.** All
seven are done (`96462ea`..`3067c44`, eleven commits). What that session landed
is the first section below; what it did *not* touch is `NEXT-STEPS.md`, which is
now a second, shorter list.

## Where things are

- **Public: https://github.com/BeMuCa/jaira** — clean tree, pushed, CI green.
- **The first outside contributions are in** (2026-08-25): PRs #1, #2, #5 merged
  whole, #3 reduced to its `z` half. Issue #4 is answered by #5. Every one of
  them was merged with my review fixes as a commit on top, named
  `fix(#N review): ...`, so what I changed about someone else's work is legible
  as its own diff rather than folded into the merge.
- The user's binary: `~/.local/bin/jaira`, built from `3067c44`.
  `go build -o ~/.local/bin/jaira ./cmd/jaira` — never `go install`. **Rebuild
  it after every change**: a stale binary made `jaira whoami` look broken.
- `gofmt -l core internal` lists `core/gate/gate.go` and
  `internal/cli/tickets.go` and always has (pre-existing alignment groups,
  verified with `git stash`). Gate on those two being the only entries.
- Still unreleased: everything after `62989f1`. Cutting v0.1.1 means a `## 0.1.1`
  block in `core/release/NOTES.md`, tag, push the tag.

## What landed since the last handoff

TUI, three fixes about not losing your place
- **A single-lane view holds its lane** (`06a1433`). `rebuild()` re-selected by
  ID and `selectByID` sets `m.laneIdx` too, but in the compact view and lane
  focus that index *is* the page. `holdsLane()`/`holdLane()` keep it; the
  multi-column board still follows the ticket.
- **A dialog goes back to the page it was opened from** (`95a04de`). New
  `returnTo` beside `detailFrom`; `modeHelp` still goes to the board.
  `finishMove` reloads an open ticket so it shows the lane it landed in.
- **`f` overrides a refused move** (`391f642`), then asks `y`. Same mutation and
  wording as the CLI's `--force`, every refusal code, nothing written to the
  ticket. Any other exit drops the offer.
- **`n` writes a follow-up beside its ticket** (`ae6cee3`): split view, the
  follow-up is a draft in memory until `ctrl+s` (esc discards, the board never
  saw it), then it reads as a ticket. `n` again chains. `tab` swaps panes once
  saved; while typing the editor owns tab, so the left pane scrolls with
  shift+arrows. No split below 80 columns or 20 rows.
- **The sign-off marker sits on the lane a person must empty** (`cfddf9d`).
  Three places keyed it off the literal id `"review"`, so Review wore
  "◆ sign off" and Human Review wore nothing. Now `RequiresHumanExit`, which is
  what the compact view always used.

Lanes and routing
- **`rejects-to:`** (`f6015d1`) declares a loop's back edge. Validated (must
  resolve, may not be itself), counted as drift, in `lanes show` and `--json`.
  No built-in declares one, so existing project copies stay equivalent.
- **`lanes/critique.md` and `lanes/optimize.md`** (`9ea9a64`, `90ff38d`) —
  catalogue lanes, deliberately **not** built-ins, since built-ins are injected
  into every board. `core/lane/shipped_test.go` parses everything in `lanes/` in
  CI; it immediately caught an unquoted colon in a description that had cost the
  lane its `id`.
- **`review-check`** (`0994fc0`): the fourth review output, the steps a person
  follows to check the change, as a flow. Enforced by the review lane's
  contract, printed last on the sign-off screen. `optimize` also produces it
  (`95d46d7`) — whichever lane is last before the human owes it.
  **Discovered while wiring it up:** the review fields lived only in the TUI;
  `ticketJSON` and `jaira show` carried none of them, so an agent could not hand
  a reader a review at all. Both carry the whole block now.
- **`next_lane`** on every `--json` ticket (`0994fc0`), and **`next --per-lane`**
  (`95d46d7`) mapping every lane with work, in pipeline order, with its
  `agentic` flag.

Identity
- **One person, several names** (`e51f9ae`). The ownership rail compared one
  string to one string. On the real board: 61 tickets assigned to a work email,
  15 to a git user.name, 1 to another email, all the same person — so every move
  needed `--force`, which trains the override. `identity.Aliases` collects
  JAIRA_USER/user.name, git user.email, and `~/.jaira/identity` (one per line);
  `gate.Request.ActorAliases` carries them. `jaira whoami` prints the list.

Docs and contracts
- The writing register now says **"and knows none of what you know"**, short
  concrete lines with one point each, no jargon — in all five shipped copies
  (`announce.go`, `create` help, `note` help, the pre-process prompt,
  `docs/AGENTS.md`).
- **An accepted ticket leaves the board once its work is pushed**, in that order
  (the done lane's description, SKILL.md, AGENTS.md, README).
- **The ticket rides in the same commit as the code** — in all three agent
  contracts, together with the reason it must be a rule: jaira never commits.
  `core/gitrepo` only reads; the sole git write anywhere is `git config --local`
  for the merge driver.
- **`--force` is reported in the command's output, not recorded on the ticket.**
  AGENTS.md and SKILL.md claimed the opposite in four places. Corrected.
- `docs/img/board.png` is **stale**: it shows "4 lane(s) off-screen" and no
  "v compact", both changed in `c02e49c`. The other three screenshots were
  checked against the code and are current. The demo world was never scripted,
  so regenerating means building one first.

## The user's own board (a different repo)

`/home/berk/git/requirementsgenie-feature-requirements-coverage-elicitation`,
branch `feature/requirements-coverage-elicitation`. Confirmed pipeline:

    backlog → todo → in-progress → critique ⇄ in-progress → optimize → review (he checks) → done

- `.jaira/lanes/review.md` was rewritten there (on 2026-08-21) as a **human** lane
  (`agentic: false`, no contract, `requires-human-exit: true`, after optimize),
  and `signoff` was removed from that board. Reason, measured: 28 tickets sat in
  `review` with **zero** review fields and **zero** commits, and nothing had ever
  reached Human Review. He uses `review` as his own inbox while the tool held it
  to be the model's lane, whose four outputs could never be produced without a
  diff. The board now warns `id "review" overrides the built-in lane` — correct,
  it genuinely does.
- `critique` and `optimize` are installed there and committed by him.
- `~/.jaira/identity` holds `berk.calabakan@partner.bmw.de` and
  `berk.calabakan@esprit-engineering.com`.
- **I moved GJB36F to done by accident** while measuring the identity fix (a real
  `move`, expected to be refused, went through once ownership no longer blocked).
  Moved back and `git checkout`ed; verified identical to HEAD. Lesson: there is
  no `--dry-run`, so never probe a gate with a live `move` on his board.

## The agent block now describes the board (`0e7db74`)

Worth knowing before touching `core/board`: the generated block in `CLAUDE.md`
and `AGENTS.md` used to be one `const string`. It taught the commands and said
nothing about *this* board's lanes, so an agent had to discover the route by
running `jaira lanes` — and would often not. That is the mechanism behind the
twelve uncritiqued tickets, not a lazy model.

It now renders the board's own lanes: order, any declared loop, and per lane
whose it is. `core/board` still imports only the stdlib — it takes plain
`LaneFact` values and the two packages that already load lanes do the
translation. Empty facts render nothing, so the lane-less note is still exactly
reproducible.

**Nothing runs by itself, and that has not changed.** Board-spawned agents are a
non-goal in the README, and the only `exec.Command` calls in the whole codebase
launch the user's `$EDITOR`. What was missing was never a runner; it was the
board telling the session what to do next, which is what this writes down.

## What to do next

`NEXT-STEPS.md` is the execution list: dependencies that live in prose rather
than in `blocked-by`, the stale `board.png`, the uncut v0.1.1, and the two
threads left open with the contributor. A session starting from "go" works it
top to bottom, one `/gsd:quick` per item.

## The list of seven, all done (2026-08-24)

Each went through `/gsd:quick`; each has a SUMMARY under `.planning/quick/` and a
row in `.planning/STATE.md` with the reasoning. Short version:

1. **`96462ea` — the critique prompt can terminate.** "Default to finding
   something" sat above a lane whose back edge is `rejects-to: in-progress`. Now
   a finding must name a file and a concrete alternative, and the end condition
   is written in the prompt. Copies refreshed in the catalogue and on his board.
2. **`610e33a` — a lane that has not been run says so.** `gate.OutputOwed` is the
   gate's own private computation, exported; `◇ unworked` on the card in both
   renderers, agentic lanes only; `next` prefers an unworked ticket inside a
   lane. **Measured: all 12 tickets in his `critique` lane are unworked.**
3. **`0620c9e` — `dod --text` and `[-]`.** Reword an item without touching its
   state or proof; retire one that will not happen. `Checked()` still means only
   `[x]`; the new `Settled()` is what completion, the gate and the counters use.
   Behaviour change: a hand-written `[-]` used to parse as todo.
4. **`a942ad0`, `8c195a5`, `a49c6fe`, `0043cb8` — the small package.**
   `move --dry-run` (stages in memory, writes nothing, same exit code),
   `show --notes-last N` (default deliberately unchanged), `[DoD n/m]` on every
   `list` row, and a warning when an abandoned claim is stepped over.
5. **`1c8884f` — `jaira delete`,** confirmed by typing the handle back, `X` on
   the open ticket in the TUI, refused while another ticket still points at it
   (a dangling `blocked-by` is a `validate` **error**).
6. **`f529c68` — `create --mine`,** an opt-in that must never become the default;
   a test guards the invariant.
7. **`3067c44` — the focused lane sits in the middle of the board,** clamped so
   the row stays full at both ends. The window is now `laneWindow()`.

Two things worth carrying forward from that session:

- **`--dry-run` removes the last reason to probe a gate on his board.** The rule
  in `berks-req-board` stands, but the workaround it asked for is gone.
- **`archive`'s help said "nothing here removes anything"**, which `delete` made
  false. Corrected. Watch for the same in any doc that predates `1c8884f`.

## Feedback from another session, triaged

Two of its points were wrong and worth not re-fixing: missing fields **do** come
as one list (measured: 5 violations in one message), and `jaira list` **does**
print a compact per-lane listing (what it lacks is a DoD counter). Its structural
point stands but not as "a missing daemon" — orchestration is an explicit
non-goal (README) — the real gap is item 2 above. `blocked_by` sitting empty
while dependencies live in prose is real and unaddressed.

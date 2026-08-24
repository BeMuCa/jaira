# Handoff — 2026-08-24

State before a context clear; this file is the memory. The previous handoff
(2026-08-18, the release and the team flow) is at `9cf71cc`. Everything below
landed after it and is pushed; CI green at `e51f9ae`. Longer-lived decisions
also live in the agent memory file `jaira-design-invariants`.

## Where things are

- **Public: https://github.com/BeMuCa/jaira** — `master` = origin, clean tree,
  CI green.
- The user's binary: `~/.local/bin/jaira`, built from `e51f9ae`.
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

## What to do next

`NEXT-STEPS.md` is the execution list: the seven items below, each with the
decision already made, the file and line to change, the traps, and what proves it
done. A session starting from "go" works it top to bottom, one `/gsd:quick` per
item.

## Queue, agreed with him and not started

1. **Critique prompt contradiction.** "Default to finding something" against
   "loop until nothing left to say" — his objection, and he is right. The escape
   clause exists but the headline overshadows it. Sharpen so termination is
   explicit and a finding must name a file and a concrete alternative.
2. **Show whether a lane has been run.** Nine tickets sat in critique
   uncritiqued and looked identical to the two that had been through. Derive it
   from the lane's `output-produces` (no new field) and surface it on the card;
   then `next --lane critique` can prefer unworked ones.
3. **`dod --text` and a `[-]` state for superseded.** There is no way to reword a
   checklist item, so items got ticked with "obsolete, replaced by 6" as proof —
   which destroys the meaning of the ticks.
4. **Small package:** `move --dry-run`, `show --notes-last N` (one ticket is 263
   lines), a DoD counter in `jaira list`, a warning when a claim expires, and the
   lane warning to stderr so it stops polluting listings.
5. **Delete a ticket**, with a double confirmation.
6. **Set yourself as assignee when creating** a ticket.
7. **Keep the selected lane centred** on the board so the following lanes stay
   visible.

## Feedback from another session, triaged

Two of its points were wrong and worth not re-fixing: missing fields **do** come
as one list (measured: 5 violations in one message), and `jaira list` **does**
print a compact per-lane listing (what it lacks is a DoD counter). Its structural
point stands but not as "a missing daemon" — orchestration is an explicit
non-goal (README) — the real gap is item 2 above. `blocked_by` sitting empty
while dependencies live in prose is real and unaddressed.

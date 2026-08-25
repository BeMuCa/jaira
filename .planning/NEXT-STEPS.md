# Next tasks — 2026-08-25 (third list, revised after his review)

**Read `HANDOFF.md` first.** The list of seven is finished, and so is the first
round of outside contributions: PRs #1, #2, #3 (reduced) and #5 are merged, and
issue #4 is answered by #5. This is what is left.

## Ground rules (unchanged)

- **Every item goes through a GSD entry point** (`/gsd:quick`). One quick task
  per item, atomic commits, a SUMMARY, a STATE row.
- **Verify on the running binary.** `go build -o ~/.local/bin/jaira ./cmd/jaira`
  after every change. Never `go install`.
- **Gate on `gofmt -l core internal` listing exactly `core/gate/gate.go` and
  `internal/cli/tickets.go`.** Both carry pre-existing alignment groups.
- **The user's real board is a different repo** and must not be probed with a
  live `jaira move`. It now has `move --dry-run`, so there is no longer any
  reason to. See the memory note `berks-req-board`.
- `go test ./... -race` with the cache cleared, before claiming anything.

## 1. `blocked_by` sits empty while dependencies live in prose

**Real, and still undecided.** Tickets say "waits on the auth work" in a note or
in `context`, and `blocked-by` — the field the gate actually reads — is empty. So
the board reports work as actionable that a person knows is not.

**Candidate, from the previous handoff:** `jaira validate` warns when a note or a
context names a handle that is not in `blocked-by`. A warning, not an error:
mentioning a ticket is not the same as depending on it, and a validator that
cries wolf gets `--strict` turned off.

**Where to look:** `core/validate/validate.go:111-120` already walks `BlockedBy`
per ticket; the handle grammar is `ticket.Handle`, six characters, so a scan of
`t.Body` and `t.Context` for handle-shaped tokens is the cheap version. Decide
what a match against an **archived** ticket means before writing it.

## 2. `docs/img/board.png` is stale

It shows "4 lane(s) off-screen" and no "v compact", both changed in `c02e49c`,
and the board now centres its focused lane (`3067c44`), so the framing is wrong
too. The other three screenshots were checked against the code and are current.

Regenerating needs a demo world, which was never scripted — the recipe is in the
handoff at `9cf71cc`. Script it this time, so the next change to the board does
not leave the same problem behind.

## 3. v0.1.1 is uncut

Everything after `62989f1` is unreleased, which is now a large amount: eleven
commits landed on 2026-08-24 alone. Cutting it means a `## 0.1.1` block in
`core/release/NOTES.md`, a tag, and pushing the tag.

**Worth deciding first:** `[-]` changes what an existing hand-written `[-]` means
(it used to parse as todo and now stops blocking completion). That belongs in the
release notes as a behaviour change, not buried in a feature list.

## 4. The contributor's two open threads

Both were merged with my fixes on top rather than waiting for a reply, so the
next session should expect a response and read it as review, not as noise.

- **#3's second commit was deliberately not taken.** It makes the board window
  state and gives it a `scrolloff` margin — a different answer to the problem
  `3067c44` solved by centring. If he makes the case that the margin still buys
  something once he sees centring running, that is its own PR. `z` and the
  hidden-lane notice are in.
- **#2's `--version` pin can still downgrade**, and `versionLine` compares
  versions by string equality, so a yanked release or one older than the running
  binary reads as "available" on a real release. Called out in the review,
  deliberately not fixed: ordering versions needs a comparison this project does
  not have yet.

## 5. Asked for and not built

Each of these he raised and I did not build, either because it was a question
rather than an instruction or because it changes something already merged.

- **A thin empty lane instead of a hidden one.** His words: hiding can confuse.
  `z` currently drops empty lanes from `boardFit`; the alternative is to keep
  them and render them at a minimal width, so "empty" is visible rather than
  absent. Touches `boardFit`, `renderColumn` and the contributor's tests for
  `z`. I asked and got no answer; it is his call, not a defect.
- **`overrides: <id>` on a lane file.** His board's `review.md` deliberately
  replaces the built-in lane, and the warning about it fires on every command.
  A declared override would silence it and make the remaining warnings mean
  something again.
- **A second back edge.** `rejects-to` names one lane and validates that it is
  installed and not itself, so `rejects-to: human` already works — but a lane
  cannot declare both "back to implementing" and "to the human for a decision".
  Today the prompt carries that (critique's does). Two edges would be a schema
  change.
- **Hiding the version line when it has nothing to say.** On a source build the
  footer reads `jaira dev`, which is decoration. One line to drop it.

## Also open, no decision

- **Nothing tells an agent a lane has been run.** The board and `next` now do
  (`610e33a`), but `ticketJSON` still carries no `output_owed`. Deliberate — the
  ordering answers the question an agent actually asks — and worth revisiting
  only if an agent is seen re-working a ticket its lane already produced.
- **`checklistJSON` reports `"done": false` for a superseded item**, with
  `"state": "superseded"` beside it. Correct as written; revisit only if an agent
  is seen trying to satisfy a retired criterion.

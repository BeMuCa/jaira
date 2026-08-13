---
phase: quick/260813-1as
plan: 01
subsystem: core/ticket, core/gate, core/lane/builtin, internal/cli, internal/tui, docs
tags: [traceability, review, definition-of-done, checklist, proof]
requires: []
provides:
  - "review-summary and review-gaps are mandatory lane output alongside review-verdict, enforced by the existing gate.fieldFilled/OutputProduces mechanism"
  - "DoDItem.Proof: a proof: sub-line under a checklist item, parsed, written and never counted as its own criterion"
  - "SetItemProof writer; AddItem skips a trailing proof line when appending"
  - "jaira dod --proof, usable alone or combined with a state flag, and its state-vs-proof at-most-one/exactly-one usage rules"
  - "proof rendered under its item in the detail pane and the sign-off screen; review summary/gaps rendered on the sign-off screen and the detail pane's Review block"
affects: [core/ticket, core/gate, core/lane/builtin, internal/cli, internal/tui, .claude/skills/jaira, docs/AGENTS.md]
tech-stack:
  added: []
  patterns:
    - "a checklist sub-line (proof:) rides directly beneath its item rather than in its own section, so SetItemState/AddItem's index-addressed edits leave it alone by construction"
    - "new lane-output fields follow the review-verdict five-place pattern: const, struct field, Decode line, gate.fieldFilled case, lane output-produces"
key-files:
  created:
    - core/gate/review_test.go
    - core/ticket/proof_test.go
    - internal/cli/checklist_test.go
  modified:
    - core/ticket/schema.go
    - core/ticket/checklist.go
    - core/gate/gate.go
    - core/lane/builtin/40-review.md
    - core/lane/builtin/20-in-progress.md
    - internal/cli/checklist.go
    - internal/cli/tickets.go
    - internal/cli/json_test.go
    - internal/tui/view.go
    - internal/tui/signoff.go
    - internal/tui/signoff_test.go
    - internal/tui/checklist_test.go
    - .claude/skills/jaira/SKILL.md
    - docs/AGENTS.md
decisions:
  - "proofPrefix ('proof:') defined once in core/ticket/checklist.go and shared by the parser (schema.go's checklistUnder) and the writer (SetItemProof), so the two can never drift into two spellings of the same marker"
  - "SetItemProof collapses whitespace runs in the proof text to single spaces before writing, so a newline in --proof text cannot inject a line that re-parses as a checkbox or heading (threat T-1as-01)"
  - "AddItem advances past a trailing proof line before computing its insertion point, so a new item never lands between the last item and its own evidence"
  - "jaira dod's exactly-one-of-doing/done/todo rule only relaxes to at-most-one when --proof is present; with no --proof, a bare 'dod <id> <n>' is still a usage error, unchanged from before"
metrics:
  duration: "~1h"
  completed: "2026-08-13"
---

# Quick Task 260813-1as: Traceability — review summary/gaps and proof Summary

The review lane now requires a plain-language summary of what shipped and an
explicit answer on what's missing (empty is refused, "none" is accepted)
before a ticket can leave it, and every definition-of-done or plan item can
carry a `proof:` evidence line written in the same call that ticks it,
surviving later edits to other items and rendering wherever a person judges
the work.

## What Was Built

**Task 1 — review-summary and review-gaps, end to end.** Added
`FieldReviewSummary`/`FieldReviewGaps` next to `FieldReviewVerdict` in
`core/ticket/schema.go` (const, struct field, `Decode` line — the same five-ish
places the interface notes named, no others). Added the matching
`fieldFilled` cases in `core/gate/gate.go`; no other gate change was needed
since the output-produces enforcement loop at `gate.go:264-276` already covers
any field a lane declares. `core/lane/builtin/40-review.md`'s
`output-produces` now lists `[review-summary, review-gaps, review-verdict]`.
`core/gate/review_test.go` covers: empty summary blocks the move (naming the
field), empty gaps blocks it independently, `review-gaps: none` counts as
filled, all three filled produces no `CodeMissingProduces` violation, and the
loaded review lane declares all three fields. No existing test moved a ticket
out of review, so there was no fixture fallout to fix.

**Task 2 — proof per item.** `DoDItem` gained a `Proof string` field.
`checklistUnder` (schema.go) now recognizes a trimmed line starting
case-insensitively with `proof:` as evidence for the item parsed immediately
before it, continuing rather than emitting it as a checklist item of its own —
last one wins if a hand-edited file somehow has two. `SetItemProof` (new, in
`checklist.go`) writes or replaces that line, deriving its indent from the
item's own leading whitespace plus two spaces, collapsing whitespace runs in
the proof text to guard against an injected checkbox/heading, and rejecting
empty text. `AddItem`'s last-item insertion point now skips past a trailing
proof line so a new checkbox lands below it rather than between the item and
its evidence. `jaira dod` gained `--proof`: rejected alongside `--add`/
`--option` (usage error), and it relaxes the "exactly one state flag" rule to
"at most one" only when present, so `--proof` alone records evidence without
touching the marker. The human-readable output now prints an item's proof
indented beneath it, and `checklistJSON` emits `"proof"` for every item.
Tests: `core/ticket/proof_test.go` covers parsing, section independence,
insert/replace, survival across `SetItemState` on a different item, `AddItem`
appending after a proof-carrying last item, the full round trip, and the
whitespace-collapse/error cases. `internal/cli/checklist_test.go` (new — no
CLI-command-execution test pattern existed in this codebase before, built one
using `newRoot`/`ticket.At`/`s.Init` following the fixture style already used
in `core/ticket/probe_test.go`) drives the real `dod` command end to end:
tick-and-record in one call, `--proof` alone leaving the marker alone, a
second `--proof` replacing rather than stacking, and the `--add`/`--option`
usage rejection. `internal/cli/json_test.go`'s existing checklist test was
extended with a proof line and assertions on the `"proof"` key.

**Task 3 — get both written and read.** Rewrote the review lane prompt
(`40-review.md`) to ask for `review-summary`, then `review-gaps` (explicitly
"none" if nothing is found), then `review-verdict`, in that order, with the
concrete `jaira set` calls shown; kept the existing spine (judge the diff not
the account, default to finding problems, not the last word) verbatim in
meaning. Extended the in-progress prompt (`20-in-progress.md`) to tie
`--proof` to the sentence already there about not claiming completion without
pointing to what satisfies it. `internal/tui/view.go`: `renderChecklist` now
prints a wrapped, styled proof line under an item (covers both Plan and
Definition of Done in the detail pane); `renderDetail` gained a Review block
(summary/gaps/verdict rows) placed after Outcome. `internal/tui/signoff.go`:
added "What the reviewer says it does" and "What the reviewer found missing"
sections between the implementer's account and the verdict, updated the
function's doc comment (previously "four questions"), and the DoD item loop
now prints an item's proof underneath it. `.claude/skills/jaira/SKILL.md` and
`docs/AGENTS.md` were updated to show `--proof` in the dod examples and to
describe the three review fields in the new order. Extended
`internal/tui/checklist_test.go` and `internal/tui/signoff_test.go` with
assertions that a proof renders under its item (both panes) and that the
review summary/gaps appear on the sign-off screen alongside the verdict.

## Deviations from Plan

**1. [Rule 1 — stale/incorrect doc] Fixed `.claude/skills/jaira/SKILL.md`'s
review-lane bullet, which said an agent "cannot move a ticket out of" review
and that attempting the move "exits 3".** This is contrary to how the lanes
are actually wired: `40-review.md` has no `requires-human-exit`, only
`45-signoff.md` does, and the review lane's own prompt (both before and after
this plan) instructs the reviewer to write its verdict and move the ticket to
sign-off itself. `docs/AGENTS.md`'s equivalent paragraph was already
consistent with the real behavior ("a review agent writes its verdict... then
a person accepts... in the board"). Since I was already rewriting this exact
bullet to name the three fields, I corrected the claim rather than leave
actively misleading text next to it. No code changed; this is a doc-only
correction.
- **Found during:** Task 3
- **Files modified:** `.claude/skills/jaira/SKILL.md`
- **Committed in:** a7527dc (Task 3 commit)

No architectural changes, no new gate on missing proofs, `PromotionFields`
unchanged, and no unrelated frontmatter field was reordered or reformatted.

## Manual End-to-End Verification (scratch repo, per plan's `<verification>`)

Built a binary from this branch and ran the plan's scenario in a scratch git
repo (adjusting for `move`'s actual flag shape, `--to <lane>` rather than a
positional target, which differs from the plan's illustrative shorthand):

1. Created a ticket, moved it through `todo` → `in-progress` → `review`
   (with `--what`/`--why`/`--resolves`/`--commits` on the move into review).
2. `jaira move <id> --to signoff` → refused (exit 3), naming all three:
   `review-summary`, `review-gaps`, `review-verdict` still empty.
3. `jaira set <id> review-summary=... review-gaps=none review-verdict=...` →
   the same move then succeeded.
4. `jaira dod <id> 1 --done --proof "internal/x.go:12; TestX"`, then
   `jaira dod <id> 1 --proof "different"` → the file has one proof line,
   replaced (confirmed by reading the ticket file directly).
5. `jaira dod <id> --add "another criterion"` → the new checkbox landed
   below the previous item's `proof:` line.
6. `jaira show <id> --json | ...` → each `dod_items` entry carries `"proof"`
   (empty string when none recorded).

All six steps behaved as specified.

## Self-Check: PASSED

- FOUND: core/ticket/schema.go (FieldReviewSummary, FieldReviewGaps,
  ReviewSummary/ReviewGaps struct fields, Decode lines, DoDItem.Proof,
  checklistUnder proof handling)
- FOUND: core/ticket/checklist.go (proofPrefix, isProofLine, proofText,
  leadingWhitespace, SetItemProof, AddItem's proof-skip)
- FOUND: core/ticket/proof_test.go
- FOUND: core/gate/gate.go (FieldReviewSummary/FieldReviewGaps cases)
- FOUND: core/gate/review_test.go
- FOUND: core/lane/builtin/40-review.md (output-produces, rewritten prompt)
- FOUND: core/lane/builtin/20-in-progress.md (--proof tie-in)
- FOUND: internal/cli/checklist.go (--proof flag, usage rules, output)
- FOUND: internal/cli/checklist_test.go
- FOUND: internal/cli/tickets.go (checklistJSON "proof" key)
- FOUND: internal/cli/json_test.go (proof assertions)
- FOUND: internal/tui/view.go (renderChecklist proof line, Review block)
- FOUND: internal/tui/signoff.go (summary/gaps sections, proof under DoD item)
- FOUND: internal/tui/signoff_test.go, internal/tui/checklist_test.go (new assertions)
- FOUND: .claude/skills/jaira/SKILL.md, docs/AGENTS.md (updated)
- Commit 0713069: found in `git log --oneline --all`
- Commit e7577d9: found in `git log --oneline --all`
- Commit a7527dc: found in `git log --oneline --all`

## Final `go test ./...` output

```
?   	github.com/BeMuCa/jaira/cmd/jaira	[no test files]
ok  	github.com/BeMuCa/jaira/core/board	0.009s
ok  	github.com/BeMuCa/jaira/core/gate	0.065s
?   	github.com/BeMuCa/jaira/core/gitrepo	[no test files]
?   	github.com/BeMuCa/jaira/core/lane	[no test files]
ok  	github.com/BeMuCa/jaira/core/merge	0.064s
ok  	github.com/BeMuCa/jaira/core/project	0.019s
?   	github.com/BeMuCa/jaira/core/session	[no test files]
ok  	github.com/BeMuCa/jaira/core/ticket	0.031s
ok  	github.com/BeMuCa/jaira/core/validate	0.072s
ok  	github.com/BeMuCa/jaira/internal/cli	0.090s
ok  	github.com/BeMuCa/jaira/internal/tui	1.510s
?   	github.com/BeMuCa/jaira/scripts/iconpreview	[no test files]
?   	github.com/BeMuCa/jaira/scripts/shotgen	[no test files]
```

(ran as `go build ./... && go vet ./... && go test ./... -count=1`; build and
vet both produced no output, i.e. clean)

## Commits

| Task | Commit | Message |
|------|--------|---------|
| 1 | 0713069 | feat: require review-summary and review-gaps to leave review |
| 2 | e7577d9 | feat: proof sub-line for definition-of-done and plan items |
| 3 | a7527dc | docs: get review-summary, review-gaps and proof written and read |

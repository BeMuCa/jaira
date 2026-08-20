---
quick_id: 260820-t7q
status: complete
commits: 90ff38d, 0994fc0
---

# Summary: a review says how to check it, and a ticket names its next lane

Three things in one pass: a rename the user caught, and items 2 and 5 from the
open list.

## 90ff38d — the lane is optimize, not optimise

File, id and name. The catalogue copy and the req board's copy were replaced in
the same pass rather than left as a dead `optimise`. Removing the `order` file the
second `lanes add` wrote restored anchor ordering, so the board still reads
`Implementing → Critique → Optimize → … → Review`.

## 0994fc0 — review-check

`FieldReviewCheck` / `review-check`: the steps a person follows to check the
change themselves, as a flow. Numbered steps, one action each, exact commands and
exact paths, and what they should see — the register the rest of the project
already uses, spelled out in the review lane's prompt.

- The review lane declares it in `output-produces`, so a review that skips it
  cannot leave the lane. `core/gate` gained the `fieldFilled` case; the existing
  gate tests were extended rather than relaxed (one of them failed exactly as it
  should have when the contract grew).
- The sign-off screen prints it **last**, after the verdict: you read the account,
  then you go and look. A test asserts that order.
- The detail pane's Review block now also appears for a ticket carrying only a
  check — its condition named three fields.

**Found while wiring it up:** the review fields were reachable only from the TUI.
`ticketJSON` carried none of them and `printDetail` printed none of them, so an
agent asked to hand a reader the review could not read one. Both now carry the
whole block (`review.summary`, `.gaps`, `.verdict`, `.check`). Without this the
feature the user asked for was impossible, not merely undocumented.

## 0994fc0 — next_lane

`Set.Next` (`core/lane/next.go`): forward through the display order, skipping
steps this ticket opted out of via its Options, never landing on a parking or
question lane. Exposed as `next_lane` on every `--json` ticket.

**Verified against the user's real board before trusting it**, which caught the
case worth having: a ticket sitting in HITL reported `next_lane: "review"` and
would have sent work straight past critique and optimize. A ticket that is itself
parked or waiting on an answer resumes where it stopped, and nothing on the board
records where that was, so the honest answer is empty. Fixed, with a test naming
the reason.

Live spread on the req board afterwards: `backlog→todo`, `todo→in-progress`,
`review→signoff`, `human→""`, `blocked→""`, `done→""`.

## Consequence handled

Editing the built-in review lane made the req board's byte-identical copy differ,
and the board warned exactly as predicted:

    lane .../.jaira/lanes/review.md: id "review" overrides the built-in lane of the same name

`jaira lanes use review --force` refreshed it and the warning went. This is the
trap the design invariants already record: there is no migration path for a
built-in's text, only re-copy or accept the warning.

## Verification

`go build ./...`, `go vet ./...`, `go test ./... -race` green with the cache
cleared. `gofmt -l core internal` lists only the two pre-existing files. New
tests: four in `core/lane/next_test.go` plus the HITL case, three in
`internal/tui/reviewcheck_test.go`, two in `internal/cli/json_test.go`, one added
to `core/gate/review_test.go`, and `TestReviewLaneDeclaresAllThreeFields` renamed
to `…EveryReviewField` with the fourth field added.

Docs: `docs/COMMANDS.md` and `docs/AGENTS.md` describe `next_lane` and the review
block; `SKILL.md` lists `review-check` as the fourth required review field.

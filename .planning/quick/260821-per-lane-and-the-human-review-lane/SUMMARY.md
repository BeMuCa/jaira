---
quick_id: 260821-perlane
status: complete
commits: 95d46d7
---

# Summary: why critique looked dead, and what it actually was

## The report

"critique and optimize do not seem to get passed on automatically by the agent —
maybe because I had not restarted jaira?"

## Three answers, all measured

**jaira needed no restart.** Lane files are read on every CLI invocation and
every TUI reload. What does need restarting is the *agent session*: the
announce.go block in CLAUDE.md and SKILL.md are read at session start, so a
session older than the lanes does not know they exist.

**Nothing was broken. The lane was starved.** `jaira next` returns one ticket,
the furthest along (`sortByProgress`, deliberate: finish in-flight work first).
On the real board that meant 28 tickets in review outranking 2 in critique, every
single time. An agent driving `jaira next` would never have been handed critique
work.

**No bug in the hand-off.** I briefly suspected the critique lane got no diff,
because `input` had four keys and no `diff`. Wrong: `diff` is a separate
top-level key in the `--for-lane` payload. For the real ticket: `complete: true`,
`missing: null`, 53,800 bytes of diff. The agent gets everything the prompt asks
for.

## The bigger finding underneath

Measured on the board:

| | |
|---|---|
| tickets in `review` | 28 |
| of those with any review field filled | **0** |
| of those with commits (needed for the diff) | **0** |
| tickets that ever reached `signoff` (Human Review) | **0** |

And the gate's own words when trying to move one out of review: it demands
`review-summary`, `review-gaps`, `review-verdict` and `review-check` — the four
outputs of a *model* review lane, which could never run because no commits means
no diff.

So the board used `review` as the user's inbox while the tool held it to be the
model's lane. Nothing had ever reached Human Review, which is why the user asked
to delete that lane earlier: right instinct, one lane off.

## What changed

`next --per-lane` (`internal/cli/flow.go`, `emitPerLane`): every lane with
actionable work, pipeline order, with the lane's `agentic` flag and one ticket.
Four tests in `internal/cli/perlane_test.go`, including the 28-hides-2 shape.

`lanes/optimize.md` also produces `review-check`, with a prompt section for it.
It is the last agent before a person, so the handover belongs there. The built-in
review lane keeps producing it too: whichever lane is last before the human owes
it.

Docs: `docs/COMMANDS.md` and `docs/AGENTS.md` describe both questions, and
AGENTS.md now states the drive-forward rule — carry a started ticket to its
`next_lane` until a human step, a gate refusal, or a blocker.

## Done on the user's board only (not in this repo)

Confirmed pipeline: `backlog → todo → in-progress → critique ⇄ in-progress →
optimize → review (the user checks) → done`.

- `.jaira/lanes/review.md` rewritten as a human lane: `agentic: false`, no
  input/output contract, `requires-human-exit: true`, anchored after optimize.
- `signoff` removed from that board (built-in and catalogue untouched;
  `jaira lanes add signoff` brings it back).
- `optimize` copy refreshed so it carries `review-check`.

The board now warns `id "review" overrides the built-in lane of the same name` —
correct and expected, it genuinely does.

Verified afterwards: moving a review ticket no longer trips the four review
fields. What remains is `done`'s own contract (outcome + commits) and "review is
a human checkpoint — open the board and sign this off". For the 28 legacy tickets
that means accepting them with `f` (force), since they carry neither outcome nor
commits.

## Verification

`go build ./...`, `go vet ./...`, `go test ./... -race` green, cache cleared.
`gofmt -l core internal` lists only the two pre-existing files. `next --per-lane`
exercised live against the real board in both human and JSON form.

## Not verified

The two lanes have still never been executed by an agent. Everything checked is
contract, routing and rendering.

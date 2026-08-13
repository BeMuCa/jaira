# Quick task 260813-nzq: Say what a progress note must contain and when one is due — Summary

`jaira note` and `jaira resume` already worked; nothing told an agent to use
them. This adds the contract (three parts: what you were doing, what you
found — especially what did not work — and what the next step is, written as
if the reader has mild ADHD) and the occasions for writing one (before you
stop, when something did not work, when the plan or context turns out wrong,
before starting something long) to the surfaces an agent actually reads: the
`note` command's own help, the working/pre-process/review lane prompts, the
per-session agent note, and `docs/AGENTS.md`.

## Commits

1. `7e9f68b` — docs: state the three parts a progress note needs in its own help
2. `768723f` — docs: tell the working, planning and review lanes when a note is due
3. `f3e704f` — docs: tighten the note help and say to write notes ADHD-style (follow-up to commit 1, requested mid-task)
4. `a8d67c7` — docs: put the note occasion in front of every agent, even one that skips lane prompts

## What changed, per file

- `internal/cli/resume.go` — `note --help`'s Long text now states the three
  parts as a tight list, says why skipping the third (the next step) leaves
  the reader to re-derive a decision already made, and tells the agent to
  write the note itself as if the reader has mild ADHD (matching the literal
  phrasing already used for the `context` field in `announce.go`,
  `15-pre-process.md`, and `jaira create --help`).
- `core/lane/builtin/20-in-progress.md` — the working lane, which previously
  never mentioned `jaira note`, now names it and lists the four occasions
  (stopping, something not working, a plan/context contradiction, before
  something long), closing with the borrowed `checkpoint` register: "the
  value is a readable trail, not a transcript."
- `core/lane/builtin/15-pre-process.md` — one sentence: reasoning behind a
  plan decision goes in a note so a later reader does not have to re-derive
  it.
- `core/lane/builtin/40-review.md` — one sentence: a send-back reason goes in
  a note too, not only in `review-gaps`, since the next review overwrites
  that field.
- `core/board/announce.go` — the agent note's existing `jaira note` bullet now
  carries an occasion ("before you stop, or when something did not work")
  instead of listing the command with no trigger. One line added, per the
  budget stated in the task ("a line earns its place or it goes") — I was
  tempted to also add the three-part shape here and did not; the full
  contract lives in `note --help` and `docs/AGENTS.md` instead.
- `docs/AGENTS.md` — new "Leaving a trail" section: when a note is due, the
  three parts, and the mild-ADHD instruction for how to write it.

## Deviations from Plan

### Mid-task addition (not a Rule 1-4 deviation — explicit instruction from the user via the coordinator)

While Task 3 was in progress, the user added a register requirement not in
the original plan text: the note's own required shape should be described as
a tight list rather than a paragraph explaining each part, and any place that
says what to put in a note should also say to write it "as if the reader has
mild ADHD" — the literal, already-established phrasing from `announce.go`,
`15-pre-process.md`, and `jaira create --help`'s `context` field, not a
paraphrase.

Task 1's commit (`7e9f68b`) was already made when this arrived, so per the
coordinator's own instruction ("if Task 1 is already committed, make a small
follow-up commit for it") I amended `internal/cli/resume.go` in a separate
commit (`f3e704f`) rather than rewriting history. Task 3's files
(`announce.go`, `docs/AGENTS.md`) were not yet committed, so the ADHD
phrasing and tight-list shape were folded directly into their first commit
(`a8d67c7`).

I did not add the ADHD phrasing to the three lane prompts
(`20-in-progress.md`, `15-pre-process.md`, `40-review.md`) or to the
`announce.go` bullet, because none of them describe note *content* — they
only state the *occasion* for writing one (a distinction the plan itself
draws: "The triggers are occasions, not a cadence"). The content contract
lives in exactly two places, per the requirement: `note --help` (canonical,
detailed) and `docs/AGENTS.md` (repeated, short) — both now carry it.

### Auto-fixed Issues

None — no bugs, missing functionality, or blocking issues were found.

## Verbatim final `go test ./...` output

```
?   	github.com/BeMuCa/jaira/cmd/jaira	[no test files]
ok  	github.com/BeMuCa/jaira/core/board	0.012s
ok  	github.com/BeMuCa/jaira/core/gate	0.220s
?   	github.com/BeMuCa/jaira/core/gitrepo	[no test files]
ok  	github.com/BeMuCa/jaira/core/identity	0.014s
ok  	github.com/BeMuCa/jaira/core/lane	0.465s
ok  	github.com/BeMuCa/jaira/core/merge	0.123s
ok  	github.com/BeMuCa/jaira/core/project	0.020s
ok  	github.com/BeMuCa/jaira/core/release	0.008s
?   	github.com/BeMuCa/jaira/core/session	[no test files]
ok  	github.com/BeMuCa/jaira/core/ticket	0.039s
ok  	github.com/BeMuCa/jaira/core/validate	0.087s
ok  	github.com/BeMuCa/jaira/internal/cli	0.203s
ok  	github.com/BeMuCa/jaira/internal/tui	1.812s
?   	github.com/BeMuCa/jaira/scripts/iconpreview	[no test files]
?   	github.com/BeMuCa/jaira/scripts/shotgen	[no test files]
```

`go vet ./...` and `go build ./...` both passed silently (no output).

## Command names checked against a built binary

Ran `/tmp/jaira-nzq <cmd> --help` for every command named in the new text —
`note`, `resume`, `dod`, `set`, `move` — and confirmed each exists as written.
`jaira note --help` and `jaira resume --help` output are quoted above under
verification in the plan; both match what the new prose promises.

## Self-Check

- `internal/cli/resume.go` — FOUND, modified as described
- `core/lane/builtin/20-in-progress.md` — FOUND, modified as described
- `core/lane/builtin/15-pre-process.md` — FOUND, modified as described
- `core/lane/builtin/40-review.md` — FOUND, modified as described
- `core/board/announce.go` — FOUND, modified as described
- `docs/AGENTS.md` — FOUND, modified as described
- Commit `7e9f68b` — FOUND in `git log`
- Commit `768723f` — FOUND in `git log`
- Commit `f3e704f` — FOUND in `git log`
- Commit `a8d67c7` — FOUND in `git log`

## Self-Check: PASSED

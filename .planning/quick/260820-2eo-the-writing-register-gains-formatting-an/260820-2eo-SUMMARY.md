---
quick_id: 260820-2eo
status: complete
---

# Summary: the register says short, it never said tidy, and it assumed too much

## The question that was asked, answered

| Asked for | Was it in? |
|---|---|
| mild ADHD | yes, all five places |
| precise wording | half: "short and concrete, names and paths rather than adjectives" |
| cleanly formatted | **no, nothing anywhere** |
| written for a beginner | **no** — the reader was only "someone not in the conversation", which assumes they know everything else the writer knows |

## The one sentence, now everywhere

> Write it as if the reader has mild ADHD **and knows none of what you know** —
> {what} first, **short concrete lines with one point each**, names and paths
> rather than adjectives, **no jargon** and no preamble.

Five sites, same wording, each keeping its own lead-in and punctuation:

- `core/board/announce.go` — the block `jaira init` / `jaira update` writes into a
  repository's `CLAUDE.md`
- `internal/cli/tickets.go` — `jaira create --help`
- `internal/cli/resume.go` — `jaira note --help`
- `core/lane/builtin/15-pre-process.md` — the pre-process lane's prompt
- `docs/AGENTS.md` — "Leaving a trail"

## The one judgement call

Berk asked for it in the spirit of "written for dummies". The word "dummy" was
kept out of the shipped CLI help: the register is read by strangers who clone the
repo, and "knows none of what you know" plus "no jargon" is the same instruction
an agent can act on. Flagged to him rather than decided silently.

## Verification

- `go build ./...` and `go test ./... -race` green.
- Line-based grep reported 0 hits in two files because the phrase wraps; a
  newline-tolerant check confirms all five carry all three new clauses.
- `jaira create --help`, `jaira note --help` and `jaira lanes show pre-process`
  all print the new sentence.
- `gofmt -l core internal` lists `core/gate/gate.go` and `internal/cli/tickets.go`
  both before and after this change (verified by stashing) — the pre-existing
  alignment diffs HANDOFF records, at line 73 of tickets.go, nowhere near this
  edit.

## Worth knowing

`core/board/announce.go` is the agent-note block, so `jaira update` will now
report a changed block in every repository that has one. That is the intended
effect: it is how the new register reaches existing boards. Projects with their
own copy of the pre-process lane need `jaira lanes use pre-process --force`.

The agent memory `adhd-register-for-replies` was extended to match, so the same
register governs replies in later sessions without being restated.

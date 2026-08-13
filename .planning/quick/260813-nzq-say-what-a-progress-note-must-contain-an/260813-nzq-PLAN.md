---
phase: quick/260813-nzq
plan: 01
type: execute
wave: 1
depends_on: []
files_modified:
  - core/lane/builtin/20-in-progress.md
  - core/lane/builtin/15-pre-process.md
  - core/lane/builtin/40-review.md
  - core/board/announce.go
  - internal/cli/resume.go
  - docs/AGENTS.md
autonomous: true
requirements: [QUICK-260813-nzq]

must_haves:
  truths:
    - "The working lane's prompt names `jaira note` and says when a note is due, not only that notes exist"
    - "A note's required shape is stated in one place and repeated in the prompts that need it: what you were doing, what you found, what the next step is"
    - "The triggers are occasions, not a cadence — a note per stopping point, not a note per turn"
    - "Every command named in the new text exists, with the flags as written"
  artifacts:
    - path: "core/lane/builtin/20-in-progress.md"
      provides: "the working lane finally tells an agent to leave a trail"
    - path: "internal/cli/resume.go"
      provides: "the `note` command's own help states the three parts"
---

<objective>
`jaira note` writes a timestamped, attributed line into a ticket's `## Progress`
section, and `jaira resume` gathers those notes for whatever was left mid-flight.
The machinery works. What is missing is that nothing tells an agent to use it.

Verified: the `in-progress` lane prompt — the one lane where work happens over
time — never mentions `jaira note` at all. Only `05-brainstorm.md` names the
command, and `core/board/announce.go:47` lists it as a line without an occasion.
The `note` command's own Long text describes the intent well, but it is read by
someone running `--help`, not by an agent in the middle of a task.

Purpose: a ticket that stops — a usage limit, a crash, a closed laptop, a
handover — must be resumable from what is written in it. Output: a stated
contract for what a note contains and when one is due, placed where an agent
actually reads it.
</objective>

<context>
@CLAUDE.md
@core/lane/builtin/20-in-progress.md
@core/lane/builtin/15-pre-process.md
@core/board/announce.go
@internal/cli/resume.go

<interfaces>
`jaira note <id> <text>` appends to the `## Progress` heading
(`noteHeading` at internal/cli/resume.go:19). The rendered line carries a
timestamp and the author:

    - **2026-08-12 23:21 · BeMuCa** — erste Idee: writer-Interface statt []byte

`jaira resume` lists tickets left in progress together with these notes.
`jaira checkpoint --focus "…" --ticket <id>` records the current session focus
outside the repo; its help already sets the right register — "Call it when the
topic changes rather than on every turn — the value is a readable trail, not a
transcript." Match that tone; do not ask for more.

Go is not on the default PATH:
  export PATH=$PATH:$HOME/.local/go/bin
</interfaces>
</context>

<tasks>

<task type="auto">
  <name>Task 1: State what a note contains, in the note command itself</name>
  <files>internal/cli/resume.go</files>
  <action>
Extend the `note` command's Long text with the three parts a note needs to make a
ticket resumable, keeping the existing prose rather than replacing it:

1. what you were doing
2. what you found — especially what did **not** work, which is the part that
   saves the next session the most time and is the part everyone omits
3. what the next concrete step would be

Say plainly why the third matters: a note that describes only the past leaves the
reader to re-derive the decision that was already made. "Next: try X because Y is
ruled out" is what makes a ticket resumable; "tried X, did not work" alone is not.

Keep it short. This is a help text, not an essay, and the repo's register is
plain sentences that explain why.
  </action>
  <verify>
    <automated>export PATH=$PATH:$HOME/.local/go/bin && go build ./... && go run ./cmd/jaira note --help 2>&1 | head -20</automated>
  </verify>
  <done>`jaira note --help` states the three parts and why the third one matters.</done>
</task>

<task type="auto">
  <name>Task 2: Say when a note is due, in the lanes where work happens</name>
  <files>core/lane/builtin/20-in-progress.md, core/lane/builtin/15-pre-process.md, core/lane/builtin/40-review.md</files>
  <action>
`20-in-progress.md` is the important one: it currently never names `jaira note`.
Add a short paragraph telling the agent to write one, and when. The occasions,
which are occasions and not a cadence:

- before you stop for any reason — a limit, a question, a handover. The board is
  the only thing that survives the session; anything not written down is lost.
- when something you tried did not work. That is the single most valuable line a
  later session can read, and the one most often skipped.
- when you learn something that contradicts the ticket's context or its plan. Fix
  the plan too, but say what changed and why.
- before starting something long, so an interruption in the middle is recoverable.

Say explicitly that this is not a note per turn. Borrow the register of
`jaira checkpoint`'s help: the value is a readable trail, not a transcript.

`15-pre-process.md`: one sentence — when the plan you are writing rests on
something you had to work out, put that reasoning in a note, so the next reader
knows why the plan looks like this rather than re-deriving it.

`40-review.md`: one sentence — when a review sends work back, the reason goes in a
note as well as in `review-gaps`, so the implementer reading the ticket later sees
it in the trail rather than only in a field that the next review overwrites.

Do not restructure these prompts. Add to them, in their existing voice.
  </action>
  <verify>
    <automated>export PATH=$PATH:$HOME/.local/go/bin && go build ./... && go test ./... -count=1 2>&1 | grep -Ev "no test files" | tail -12</automated>
    <automated>export PATH=$PATH:$HOME/.local/go/bin && grep -c "jaira note" core/lane/builtin/20-in-progress.md core/lane/builtin/15-pre-process.md core/lane/builtin/40-review.md</automated>
  </verify>
  <done>All three prompts name `jaira note` and state the occasion for it; the working lane lists the triggers; the whole suite passes.</done>
</task>

<task type="auto">
  <name>Task 3: Put it in front of an agent that never reads a lane prompt</name>
  <files>core/board/announce.go, docs/AGENTS.md</files>
  <action>
The agent note in `core/board/announce.go` currently lists `jaira note` as a
command without an occasion. Give it one line that says when: before you stop,
and whenever something did not work. That file is read at the start of every
session by agents that may never touch a lane prompt at all, so it is the only
place some readers will ever see this.

`docs/AGENTS.md`: a short paragraph in the surrounding register — what a note
must contain, when one is due, and that `jaira resume` is what reads them back.
State the point plainly: a session that ends abruptly leaves nothing behind
except what was written down.

Keep the agent note's length in check. It is loaded into every session; a line
earns its place or it goes.
  </action>
  <verify>
    <automated>export PATH=$PATH:$HOME/.local/go/bin && go build ./... && go test ./core/board/... -count=1 2>&1 | tail -3</automated>
    <automated>export PATH=$PATH:$HOME/.local/go/bin && go build -o /tmp/jaira-nzq ./cmd/jaira && T=$(mktemp -d) && git -C $T init -q && (cd $T && JAIRA_HOME=$(mktemp -d) /tmp/jaira-nzq init >/dev/null && grep -c "jaira note" CLAUDE.md)</automated>
  </verify>
  <done>A freshly initialised repository's CLAUDE.md tells an agent when to write a note, and docs/AGENTS.md explains the contract.</done>
</task>

</tasks>

<verification>
export PATH=$PATH:$HOME/.local/go/bin && go build ./... && go vet ./... && go test ./... -count=1

Every command named in the new text must exist with the flags as written. Check
each one against `--help` on a built binary rather than trusting the plan.
</verification>

<success_criteria>
- The working lane's prompt tells an agent to write a note and when.
- The three parts of a note are stated in the `note` command's help and echoed where an agent reads them.
- The triggers are occasions, not a cadence.
- A freshly initialised repository's agent note carries the occasion, not just the command.
- `go build ./... && go vet ./... && go test ./... -count=1` passes.
</success_criteria>

<output>
Create `.planning/quick/260813-nzq-say-what-a-progress-note-must-contain-an/260813-nzq-SUMMARY.md` when done
</output>

---
id: in-progress
name: Implementing
after: pre-process
precedence: 30
agentic: true
terminal: false
model-tier: cheap
input-requires: [goal, definition-of-done, context, plan, notes]
output-produces: [outcome-what, outcome-why, outcome-resolves]
description: Carrying out the plan.
---
# Prompt

Carry out the plan on this ticket.

You are given the goal, the context, the definition of done, and the plan worked
out in the previous step. Work through the plan in order. Mark the step you are
on with `jaira dod <id> <n> --doing --plan`, and mark it done when it is done —
that marker is how a person watching the board knows where you are. When you
mark a definition-of-done item done, record the evidence in the same call:
`jaira dod <id> <n> --done --proof "<file:line or test name>"`, naming the
specific code or test that makes it true.

Work only within the plan's scope. If you discover adjacent problems, note them
rather than fixing them. If the plan turns out to be wrong, say so and fix the
plan rather than quietly doing something else.

## Fix the class, not the example

Ticking items is not the job; the change working is. After you, critique,
testing, a bug hunter and a reviewer read this change, and each round they find
the next example of a mistake that was fixed for one example only. So:

- **Know the class before you write code.** The definition of done names the
  class of fault, the input classes and, for a heuristic, the acceptance rule.
  Where the definition of done does not — no planning step ran — name them
  yourself first and write them on the ticket with `jaira dod <id> --add`.
- **Fix the cause, not the example.** A finding or a ticket arrives as one
  example: "3 / 12" is not removed. A fix that adds "3 / 12" to a pattern is
  found out next round by "3/12". Ask why the example fails, and change that —
  read structured data instead of parsing its printed form, validate at the
  one place every caller passes through, not in the caller that was reported.
- **Test the class.** One parametrised test over the input classes, the
  reported example among them — not one test for the example.
- **Search for the siblings.** The same mistake lives wherever the same code
  was copied or the same pattern was used: other write paths, other commands,
  other parsers of the same format. Grep for them and fix them in this change.
  Siblings of this fault are within the plan's scope; a different problem you
  find on the way is still a note.
- **Check the four categories yourself before you move.** Concurrency (two
  writers at once), error paths (a step failing halfway, an error dropped or
  reported as success), unusual input (empty, huge, duplicated, another
  encoding or locale), permissions (what it trusts from outside). Go through
  each for the code you changed, and say in your outcome which you checked and
  what you found. The testing lane and the bug hunt check exactly these; a
  defect they find that you could have found yourself is a round lost.
- **Close the small edges in the same round.** When the ticket comes back,
  the notes that sent it carry blockers and often smaller edges of this change
  beside them: a missing case of the class, a one-line guard, a test that
  covers the example but not its neighbour. Close those too, now — a small
  edge left as a note reaches the person, who sends the ticket back for it
  anyway. A fault that was there before this change and lies outside its scope
  stays a note.
- **A give-back with decisions is the list.** A note `Given back from
  <lane>: beheben:` carries what the person decided per review note, word for
  word: fix exactly what stands under `beheben:`, its `dazu:` line (what
  the person typed) included. What stands under
  `angenommen:` was accepted as it is, and what stands under `Folge-Ticket:`
  has a ticket of its own — leave both alone, small edges or not.

## The third return from the same lane

jaira counts the returns itself. `back` in this lane's payload
(`jaira show <id> --for-lane in-progress --json`) is the last send-back as the
move recorded it:

- `from` — the lane the ticket came back from
- `by` — `person` when a person sent it (from signoff, human or any lane no
  agent works, from the board, or with `jaira move` typed at a terminal),
  `lane` when a loop lane did. An agent whose harness runs its commands in a
  terminal is recorded as a person too, so on such a harness `count` never
  stops a loop: count the returns from the lane's own notes yourself
- `count` — how many returns in a row came from that same lane; a person's
  return, or one from another lane, starts it again at 1, so two loop lanes
  that take turns sending it back never reach 3: count those yourself

`by: person` is the person's word, not a loop: do what they asked, whatever
the count. Their reason is in the newest `Given back from <lane>:` or
`decision:` note.

`by: lane` with `count` 3 or more means the same lane has sent it back three
times or more: do not start another round. Move the ticket to the human lane
with the question instead — the third round's findings, and whether the definition of done is wrong or the
loop reads deeper every round without ending:

```bash
jaira move <id> --to human --question "testing sent this back three times: <findings>. DoD wrong or loop without end? I would <...>"
```

A ticket from before jaira recorded this has an empty `from` and `count` 0:
it counts as no return.

A loop lane that declares `max-rounds` (the shipped critique, optimize and
testing: 3) does not need you to count: jaira itself refuses its next
send-back past the limit (`max_rounds`, exit 3) and names the lane a person
answers in. `rounds_left` in that lane's `--for-lane` payload says how many
remain.

In conversational mode you hand back commit lines instead of committing, so
the bug hunt reads your work uncommitted; once that work is committed the
hunt no longer covers HEAD and `bughunt` in `jaira show <id> --brief` reads
`fail` until the next hunt.

## Progress notes

Write a note with `jaira note <id> <text>`. You are writing to yourself, in a
session that has lost this conversation and can read only the repository and this
ticket.

**The rule: write down what the repository does not already say.**

That test decides both what to write and what to leave out. Applied:

- **a dead end** — "tried X, it does not work because Y". Nothing in the code
  records a thing that was not done. Without it the next session spends its first
  hour rediscovering it. This is the single most valuable line you can leave and
  the one most often skipped.
- **why this and not that** — the code shows the approach you took, never the one
  you rejected or the reason. Without it someone later "fixes" it back.
- **what you had to find out** — "the exporter buffers everything in writeAll()".
  It cost you time to learn; it costs the next reader the same time again.
- **work that is half done** — "callers in Y updated, Z still on the old shape".
- **what you would do next**, when it is not simply the next line of the plan.

And what to skip, because the board already carries it: which step you are on
(the `--doing` marker says so), that a step is finished (the checklist says so),
that a commit exists (git says so), that the definition of done is met (the
`--proof` on the item says so). A note repeating those is a transcript, and a
transcript is not read.

**When to ask yourself the question:** at every natural pause — a plan step
finished, a piece implemented, a test going green, before starting something
long, on handover. Not on every turn.

Do not save it for when you are about to stop. A session killed by a usage limit
or a crash never gets that turn, which is exactly the case the note exists for.

When you are finished, report:

- **what** you changed, concretely
- **why** the change was needed
- **how** the change satisfies the definition of done, referring to it directly
- **which commits** carry it: `--commits "$(git rev-parse HEAD)"` on the move,
  naming every commit you made. Record them now, while you still know which they
  are — at review and at sign-off the diff shown is the diff of exactly these
  commits, and the done lane refuses a ticket that names none, because a change
  with no commits is a change nobody can check, then or months later.

Do not claim the definition of done is met unless you can point to the specific
behaviour that now satisfies it — the `--proof` on each item is where that
pointer goes.

## Items only a person ticks

A definition-of-done item that starts with `seen by a person:` is never yours
to tick, however sure you are: a width table, a run on another platform or the
text of a refusal is a stand-in, not what the person sees. Do the work around
it — the code, the tests a machine can run — name in your outcome what the
person should look at, and leave the item open. The person ticks it at
acceptance.

Where no planning step ran and the change alters what a terminal draws (the
TUI, colours, widths, glyphs, a plugin's panel), behaves differently by
operating system, shell or terminal, or needs a person at the keyboard, add
such an item yourself before you start:
`jaira dod <id> --add "seen by a person: <what to look at, where>"`.
